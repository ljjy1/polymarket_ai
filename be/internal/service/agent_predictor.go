package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

// AgentPredictor 使用 DeepAgent（多智能体编排）分析聚合数据，输出结构化预测。
type AgentPredictor struct {
	runner *adk.Runner
	agent  adk.Agent
}

// NewAgentPredictor 创建一个新的 AgentPredictor 实例。
// 它使用 DeepAgent 架构，将技术、情绪、新闻、链上分析分解给子智能体并行处理。
func NewAgentPredictor(ctx context.Context, apiKey, baseURL, modelName string) (*AgentPredictor, error) {
	// 1. 创建 DeepSeek ChatModel
	cfg := &deepseek.ChatModelConfig{
		APIKey:      apiKey,
		Model:       modelName,
		Temperature: 0.2,
		MaxTokens:   4096,
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}

	cm, err := deepseek.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek chat model: %w", err)
	}

	// 2. 创建子智能体
	subAgents, err := NewAnalysisAgents(ctx, cm)
	if err != nil {
		return nil, fmt.Errorf("failed to create analysis sub-agents: %w", err)
	}

	// 3. 创建 DeepAgent（主智能体）
	agent, err := deep.New(ctx, &deep.Config{
		Name:        "market_analyst_orchestrator",
		Description: "比特币预测市场分析主智能体，统筹技术、情绪、新闻、链上分析",
		ChatModel:   cm,
		Instruction: agentPredictorInstruction(),
		SubAgents:   subAgents,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{},
			},
		},
		WithoutGeneralSubAgent: true,
		MaxIteration:           10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepAgent: %w", err)
	}

	// 4. 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	return &AgentPredictor{
		runner: runner,
		agent:  agent,
	}, nil
}

// agentPredictorInstruction 返回主智能体系统提示词。
func agentPredictorInstruction() string {
	return `你是一位比特币预测市场分析师主管，负责统筹多个专业分析智能体完成市场预测任务。

## 工作流程

1. **接收数据**：你将收到完整的市场数据包（MarketDataBundle），包含技术指标、情绪数据、新闻、链上数据等。

2. **制定计划**：使用 write_todos 工具制定分析计划，列出需要调用的子智能体。

3. **分配任务**：通过 task 工具将各维度分析任务分配给对应的子智能体，并将相关数据片段传给它们。

4. **收集结果**：等待所有子智能体返回分析结果。

5. **综合预测**：综合所有分析结果，输出最终的概率预测。

## 子智能体调用规则

- 技术分析子智能体（technical_analyst）：当技术指标数据完整（RSI、MACD、布林带、EMA、ATR 均不为零）时调用
- 情绪分析子智能体（sentiment_analyst）：当恐惧贪婪指数 > 0 时调用
- 新闻分析子智能体（news_analyst）：当新闻头条列表非空时调用
- 链上分析子智能体（onchain_analyst）：当链上数据（交易所净流量、大额交易数等）不为零时调用
- 如果某个维度的数据缺失（如链上数据为空），跳过对应的子智能体调用

## 输出要求（必须严格遵守）

综合所有分析结果后，你必须输出**唯一一个**严格的有效 JSON 对象，匹配以下结构：
{
    "predicted_probability": 0.01-0.99,
    "confidence": 0.0-1.0,
    "direction": "bullish|bearish|neutral",
    "key_factors": ["factor1", "factor2"],
    "risk_factors": ["risk1", "risk2"],
    "technical_analysis": "技术面分析摘要",
    "sentiment_analysis": "情绪分析摘要",
    "news_impact": "新闻影响摘要",
    "onchain_analysis": "链上分析摘要",
    "reasoning": "综合推理过程"
}

### 禁止行为
- ❌ 不允许输出任何 markdown 标记（# ## ** - | 等）
- ❌ 不允许输出表格或列表
- ❌ 不允许输出额外散文、解释或思考过程
- ❌ 不允许使用 json 代码块包裹
- ✅ 只输出一个 JSON 对象，且必须是唯一内容`
}

// formatMarketDataBundle 将 MarketDataBundle 格式化为完整的文本输入。
func formatMarketDataBundle(bundle *MarketDataBundle, marketYesPrice float64) string {
	var sb strings.Builder

	hoursToResolution := time.Until(bundle.TargetDatetime).Hours()
	if hoursToResolution < 0 {
		hoursToResolution = 0
	}

	distance := bundle.BTCCurrentPrice - float64(bundle.TargetPrice)
	distancePct := (distance / float64(bundle.TargetPrice)) * 100
	directionLabel := "above"
	if distance < 0 {
		directionLabel = "below"
	}

	sb.WriteString("# 市场问题\n")
	sb.WriteString(fmt.Sprintf("问题: BTC/USDT (Binance) 在 %s 时收盘价是否高于 $%d?\n",
		bundle.TargetDatetime.Format("2006-01-02 15:04 MST"), bundle.TargetPrice))
	sb.WriteString(fmt.Sprintf("距离结算还有: %.1f 小时\n\n", hoursToResolution))

	sb.WriteString("# 当前市场状态 (Polymarket)\n")
	sb.WriteString(fmt.Sprintf("Yes 价格: %.4f (隐含概率 %.1f%%)\n\n", marketYesPrice, marketYesPrice*100))

	sb.WriteString("# 当前 BTC 价格\n")
	sb.WriteString(fmt.Sprintf("现货价格 (Binance): $%.2f\n", bundle.BTCCurrentPrice))
	sb.WriteString(fmt.Sprintf("距离阈值: %.2f%% (%s)\n", distancePct, directionLabel))
	sb.WriteString(fmt.Sprintf("24h 涨跌幅: %.2f%%\n", bundle.PriceChange24h))
	sb.WriteString(fmt.Sprintf("7日 涨跌幅: %.2f%%\n\n", bundle.PriceChange7d))

	sb.WriteString("# 技术指标\n")
	sb.WriteString(fmt.Sprintf("RSI(14): %.2f\n", bundle.RSI14))
	sb.WriteString(fmt.Sprintf("MACD: macd=%.4f, signal=%.4f, histogram=%.4f\n",
		bundle.MACD.MACD, bundle.MACD.Signal, bundle.MACD.Histogram))
	sb.WriteString(fmt.Sprintf("布林带: upper=%.2f, mid=%.2f, lower=%.2f\n",
		bundle.Bollinger.Upper, bundle.Bollinger.Middle, bundle.Bollinger.Lower))
	sb.WriteString(fmt.Sprintf("EMA: ema7=%.2f, ema25=%.2f, ema99=%.2f\n",
		bundle.EMA.EMA7, bundle.EMA.EMA25, bundle.EMA.EMA99))
	sb.WriteString(fmt.Sprintf("ATR(14): %.2f\n\n", bundle.ATR))

	sb.WriteString("# 市场情绪\n")
	sb.WriteString(fmt.Sprintf("恐惧贪婪指数: %d/100 (%s)\n", bundle.FearGreedIndex, bundle.FearGreedLabel))
	sb.WriteString(fmt.Sprintf("资金费率: %.6f\n", bundle.FundingRate))
	sb.WriteString(fmt.Sprintf("多空比: %.4f\n\n", bundle.LongShortRatio))

	sb.WriteString("# 新闻头条\n")
	if len(bundle.NewsHeadlines) > 0 {
		sb.WriteString(formatNewsHeadlines(bundle.NewsHeadlines))
	} else {
		sb.WriteString("暂无新闻数据\n")
	}
	sb.WriteString("\n")

	sb.WriteString("# 链上数据 (CryptoQuant)\n")
	sb.WriteString(fmt.Sprintf("交易所净流量 24h: $%.0f\n", bundle.ExchangeNetflow24h))
	sb.WriteString(fmt.Sprintf("交易所净流量 7d: $%.0f\n", bundle.ExchangeNetflow7d))
	sb.WriteString(fmt.Sprintf("活跃地址 24h: %d\n", bundle.ActiveAddresses24h))
	sb.WriteString(fmt.Sprintf("矿工流出 24h: $%.0f\n", bundle.MinerOutflow24h))
	sb.WriteString(fmt.Sprintf("MVRV 比率: %.4f\n\n", bundle.MVRVRatio))

	sb.WriteString("# 宏观事件\n")
	if len(bundle.MacroEvents) > 0 {
		for i, ev := range bundle.MacroEvents {
			sb.WriteString(fmt.Sprintf("%d. [%s] %s - 预期: %s, 影响: %s\n",
				i+1, ev.Date, ev.Event, ev.Expected, ev.Impact))
		}
	} else {
		sb.WriteString("暂无宏观事件数据\n")
	}
	sb.WriteString("\n")

	sb.WriteString("# 你的任务\n")
	sb.WriteString("请使用 write_todos 制定分析计划，然后通过 task 工具调用合适的子智能体进行分析。\n")
	sb.WriteString("\n")
	sb.WriteString("【严格遵守】最终输出必须是严格的有效 JSON 格式，匹配以下结构：\n")
	sb.WriteString("{\n")
	sb.WriteString("    \"predicted_probability\": 0.42,\n")
	sb.WriteString("    \"confidence\": 0.65,\n")
	sb.WriteString("    \"direction\": \"bearish\",\n")
	sb.WriteString("    \"key_factors\": [\"因子1\", \"因子2\"],\n")
	sb.WriteString("    \"risk_factors\": [\"风险1\", \"风险2\"],\n")
	sb.WriteString("    \"technical_analysis\": \"技术面分析摘要\",\n")
	sb.WriteString("    \"sentiment_analysis\": \"情绪分析摘要\",\n")
	sb.WriteString("    \"news_impact\": \"新闻影响摘要\",\n")
	sb.WriteString("    \"onchain_analysis\": \"链上分析摘要\",\n")
	sb.WriteString("    \"reasoning\": \"综合推理过程\"\n")
	sb.WriteString("}\n")
	sb.WriteString("不允许输出任何 markdown 标记、表格、标题（# ## **等）、额外散文或思考过程。\n")
	sb.WriteString("只输出一个 JSON 对象，不要用 ```json 包裹。\n")

	return sb.String()
}

// Predict 使用 DeepAgent 多智能体架构分析市场数据，返回结构化预测结果。
func (p *AgentPredictor) Predict(ctx context.Context, bundle *MarketDataBundle, marketYesPrice float64) (*PredictionResult, error) {
	if bundle == nil {
		return nil, fmt.Errorf("market data bundle is nil")
	}

	// 1. 格式化输入
	formattedInput := formatMarketDataBundle(bundle, marketYesPrice)

	// 2. 调用 Runner
	iter := p.runner.Query(ctx, formattedInput)

	// 3. 遍历事件流，收集最终输出
	var lastContent string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, fmt.Errorf("agent execution error: %w", event.Err)
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			msgOutput := event.Output.MessageOutput
			if msgOutput.Role == schema.Assistant {
				if !msgOutput.IsStreaming && msgOutput.Message != nil {
					lastContent = msgOutput.Message.Content
				} else if msgOutput.IsStreaming && msgOutput.MessageStream != nil {
					msg, err := schema.ConcatMessageStream(msgOutput.MessageStream)
					if err != nil {
						logger.Warn("failed to concat message stream", logger.Err(err))
						continue
					}
					lastContent = msg.Content
				}
			}
		}
	}

	if lastContent == "" {
		return nil, fmt.Errorf("no valid output from agent")
	}

	// 4. 解析 JSON
	prediction, err := parsePredictionResponse(lastContent)
	if err != nil {
		// 重试一次：重新调用预测
		logger.Warn("failed to parse agent response, retrying once",
			logger.String("raw_response", lastContent),
			logger.Err(err),
		)
		iter2 := p.runner.Query(ctx, formattedInput)
		var retryContent string
		for {
			event, ok := iter2.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				return nil, fmt.Errorf("agent retry execution error: %w", event.Err)
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				msgOutput := event.Output.MessageOutput
				if msgOutput.Role == schema.Assistant {
					if !msgOutput.IsStreaming && msgOutput.Message != nil {
						retryContent = msgOutput.Message.Content
					} else if msgOutput.IsStreaming && msgOutput.MessageStream != nil {
						msg, concatErr := schema.ConcatMessageStream(msgOutput.MessageStream)
						if concatErr != nil {
							logger.Warn("failed to concat message stream on retry", logger.Err(concatErr))
							continue
						}
						retryContent = msg.Content
					}
				}
			}
		}
		if retryContent == "" {
			return nil, fmt.Errorf("no valid output from agent on retry")
		}
		prediction, err = parsePredictionResponse(retryContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse agent response after retry: %w, raw: %s", err, retryContent)
		}
		prediction.RawResponse = retryContent
	} else {
		prediction.RawResponse = lastContent
	}

	// 5. 计算 edge 并设置 recommended_action
	edge := prediction.PredictedProbability - marketYesPrice
	prediction.RecommendedAction = computeRecommendedAction(edge, prediction.Confidence)

	// 6. 记录预测关键数据
	logger.Info("Agent prediction completed",
		logger.Float64("predicted_probability", prediction.PredictedProbability),
		logger.Float64("confidence", prediction.Confidence),
		logger.String("direction", prediction.Direction),
		logger.Float64("edge", edge),
		logger.String("recommended_action", prediction.RecommendedAction),
		logger.Float64("market_yes_price", marketYesPrice),
	)

	return prediction, nil
}
