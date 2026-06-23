package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/go-dev-frame/sponge/pkg/logger"

	"be/internal/external"
)

// PredictionResult 是 AI 模型返回的结构化预测结果。
type PredictionResult struct {
	PredictedProbability float64  `json:"predicted_probability"` // 校准后的"答案为 Yes"的概率估计值（范围 0.01~0.99）
	Confidence           float64  `json:"confidence"`            // AI 对自己预测的置信度（范围 0~1），与概率相互独立
	Direction            string   `json:"direction"`             // 市场方向判断: bullish-看涨, bearish-看跌, neutral-中性
	KeyFactors           []string `json:"key_factors"`           // 支持该预测的关键技术/基本面因素列表
	RiskFactors          []string `json:"risk_factors"`          // 可能使预测失效的风险因素列表
	TechnicalAnalysis    string   `json:"technical_analysis"`    // 技术面分析详情，包含 K 线形态、指标解读等
	SentimentAnalysis    string   `json:"sentiment_analysis"`    // 市场情绪分析，含恐惧贪婪指数解读
	NewsImpact           string   `json:"news_impact"`           // 新闻事件对市场的影响评估
	OnchainAnalysis      string   `json:"onchain_analysis"`      // 链上数据分析，含交易所净流量、活跃地址等
	Reasoning            string   `json:"reasoning"`             // AI 完整推理过程，含多因素综合判断的思考链
	RecommendedAction    string   `json:"recommended_action"`    // 推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过（由下游逻辑覆盖）
	RawResponse          string   `json:"raw_response"`          // AI 模型返回的原始响应文本（JSON 解析前的原文）
}

// AIPredictor 使用 Eino ChatModel (DeepSeek) 分析聚合数据，输出结构化预测。
type AIPredictor struct {
	chatModel     *deepseek.ChatModel
	promptVersion string
}

// NewAIPredictor 创建一个新的 AIPredictor 实例。
// apiKey: DeepSeek API Key
// baseURL: DeepSeek API 基础地址（可选，默认 https://api.deepseek.com）
// model: 模型名称（如 deepseek-chat）
func NewAIPredictor(apiKey, baseURL, model string) (*AIPredictor, error) {
	ctx := context.Background()

	cfg := &deepseek.ChatModelConfig{
		APIKey:      apiKey,
		Model:       model,
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

	return &AIPredictor{
		chatModel:     cm,
		promptVersion: "v1.0",
	}, nil
}

// systemPrompt 返回市场分析师的系统提示词。
// 与 Python 参考版 (ai_predictor.py) 保持一致，
// 明确告知模型下游有 25% edge 硬门控，让模型在输出时考虑这一约束。
func (p *AIPredictor) systemPrompt() string {
	return `你是一位拥有10年以上经验的量化加密货币市场分析师，专精于比特币衍生品、链上分析和预测市场。

你的任务是估算BTC/USDT (Binance现货) 在特定结算时间点收盘价是否高于某个价格阈值的概率。你并不是在预测方向或幅度——你是在估算一个校准后的概率，范围在0到1之间。

你必须遵守的关键规则：

1. 输出校准后的概率，而非置信度加权赌注。如果数据确实模糊不清，应输出接近当前市场价格的概率，而不是强烈的主观意见。

2. 仅使用提供的数据。不要编造提示中没有给出的价格、指标或新闻。

3. 区分以下两个概念：
   - "predicted_probability"：你对真实概率的诚实估计
   - "confidence"：考虑到数据质量，你认为自己的估计有多可靠
   这两个是独立的——对50%的概率有高置信度是完全合理的。

4. 考虑距离结算的时间。48小时内结算的市场比7天市场的确定性更高，请相应调整。

5. 仅输出符合所提供JSON Schema的有效JSON。不要在JSON之外包含markdown或任何散文内容。

6. 下游系统仅在满足以下条件时才会交易：
   abs(predicted_probability - market_yes_price) >= 0.25  AND  confidence >= 0.6
   如果数据不支持 >=25% 的价差，请将 recommended_action 设为 "skip"，并报告你诚实的概率估计。`
}

// userPromptTemplate 返回 Go 模板格式的用户提示词模板。
// 与 Python 参考版 (ai_predictor.py) 的结构保持一致，
// 包含: 市场问题、Polymarket状态、BTC价格、技术指标、情绪、新闻、链上数据。
func (p *AIPredictor) userPromptTemplate() string {
	return `# 市场问题
问题: "BTC/USDT (Binance) 在 {{.TargetDatetime}} 时收盘价是否高于 ${{.TargetPrice}}?"
距离结算还有: {{.HoursToResolution}} 小时

# 当前市场状态 (Polymarket)
Yes价格: {{.YesPrice}} (隐含概率 {{.YesPct}}%)
No价格: {{.NoPrice}}
24h交易量: ${{.Volume24h}}
买卖价差: {{.Spread}}%

# 当前BTC价格 (Binance)
现货价格: ${{.BTCCurrentPrice}}
距离阈值: {{.DistancePct}}% ({{.DirectionLabel}})
24h涨跌幅: {{.Change24h}}%
7日涨跌幅: {{.Change7d}}%
24h最高价: ${{.High24h}}
24h最低价: ${{.Low24h}}

# 技术指标 (基于1小时K线)
RSI(14): {{.RSI14}}
MACD: macd={{.MACD}}, signal={{.Signal}}, histogram={{.Histogram}}
布林带: upper={{.BBUpper}}, mid={{.BBMid}}, lower={{.BBLower}}
EMA: ema7={{.EMA7}}, ema25={{.EMA25}}, ema99={{.EMA99}}
ATR(14): {{.ATR}}

# 市场情绪
恐惧贪婪指数: {{.FGI}}/100 ({{.FGILabel}})
恐惧贪婪7日趋势: {{.FGITrend}}

# 新闻头条
{{.NewsHeadlines}}

# 链上信号
交易所净流量24h: ${{.ExchangeNetflow}}
大额交易数24h: {{.LargeTxCount}}
活跃地址变化: {{.ActiveAddrChange}}%

# 你的任务
估算答案为YES的校准概率。输出严格JSON格式，匹配所给的Schema。`
}

// Predict 使用 Eino ChatModel 分析市场数据，返回结构化预测结果。
// bundle: 聚合后的市场数据
// marketYesPrice: Polymarket 上的 Yes 价格（0-1 之间）
func (p *AIPredictor) Predict(ctx context.Context, bundle *MarketDataBundle, marketYesPrice float64) (*PredictionResult, error) {
	if bundle == nil {
		return nil, fmt.Errorf("market data bundle is nil")
	}

	// 1. 构建用户提示词数据
	distance := bundle.BTCCurrentPrice - float64(bundle.TargetPrice)
	distancePct := (distance / float64(bundle.TargetPrice)) * 100
	directionLabel := "above"
	if distance < 0 {
		directionLabel = "below"
	}

	hoursToResolution := time.Until(bundle.TargetDatetime).Hours()
	if hoursToResolution < 0 {
		hoursToResolution = 0
	}

	// 格式化新闻头条
	newsText := formatNewsHeadlines(bundle.NewsHeadlines)

	// 格式恐惧贪婪 7 日趋势
	fgiTrendStr := "n/a"
	if len(bundle.FearGreedTrend7d) > 0 {
		trendStrs := make([]string, len(bundle.FearGreedTrend7d))
		for i, v := range bundle.FearGreedTrend7d {
			trendStrs[i] = fmt.Sprintf("%d", v)
		}
		fgiTrendStr = strings.Join(trendStrs, ",")
	}

	noPrice := 1.0 - marketYesPrice
	spreadPct := 0.0
	if marketYesPrice > 0 {
		spreadPct = 0.0 // 暂无实时价差数据，保留为0
	}

	tplData := map[string]any{
		"TargetDatetime":    bundle.TargetDatetime.Format("2006-01-02 15:04 MST"),
		"TargetPrice":       bundle.TargetPrice,
		"HoursToResolution": fmt.Sprintf("%.1f", hoursToResolution),
		"YesPrice":          fmt.Sprintf("%.4f", marketYesPrice),
		"YesPct":            fmt.Sprintf("%.1f", marketYesPrice*100),
		"NoPrice":           fmt.Sprintf("%.4f", noPrice),
		"Volume24h":         fmt.Sprintf("%.0f", bundle.Volume24h),
		"Spread":            fmt.Sprintf("%.2f", spreadPct),
		"BTCCurrentPrice":   fmt.Sprintf("%.2f", bundle.BTCCurrentPrice),
		"DistancePct":       fmt.Sprintf("%.2f", math.Abs(distancePct)),
		"DirectionLabel":    directionLabel,
		"Change24h":         fmt.Sprintf("%.2f", bundle.PriceChange24h),
		"Change7d":          fmt.Sprintf("%.2f", bundle.PriceChange7d),
		"High24h":           fmt.Sprintf("%.2f", bundle.High24h),
		"Low24h":            fmt.Sprintf("%.2f", bundle.Low24h),
		"RSI14":             fmt.Sprintf("%.2f", bundle.RSI14),
		"MACD":              fmt.Sprintf("%.4f", bundle.MACD.MACD),
		"Signal":            fmt.Sprintf("%.4f", bundle.MACD.Signal),
		"Histogram":         fmt.Sprintf("%.4f", bundle.MACD.Histogram),
		"BBUpper":           fmt.Sprintf("%.2f", bundle.Bollinger.Upper),
		"BBMid":             fmt.Sprintf("%.2f", bundle.Bollinger.Middle),
		"BBLower":           fmt.Sprintf("%.2f", bundle.Bollinger.Lower),
		"EMA7":              fmt.Sprintf("%.2f", bundle.EMA.EMA7),
		"EMA25":             fmt.Sprintf("%.2f", bundle.EMA.EMA25),
		"EMA99":             fmt.Sprintf("%.2f", bundle.EMA.EMA99),
		"ATR":               fmt.Sprintf("%.2f", bundle.ATR),
		"FGI":               bundle.FearGreedIndex,
		"FGILabel":          bundle.FearGreedLabel,
		"FGITrend":          fgiTrendStr,
		"NewsHeadlines":     newsText,
		"ExchangeNetflow":   fmt.Sprintf("%.0f", 0.0),
		"LargeTxCount":      0,
		"ActiveAddrChange":  fmt.Sprintf("%.2f", 0.0),
	}

	// 2. 创建提示词模板
	tpl := prompt.FromMessages(schema.GoTemplate,
		&schema.Message{
			Role:    schema.System,
			Content: p.systemPrompt(),
		},
		&schema.Message{
			Role:    schema.User,
			Content: p.userPromptTemplate(),
		},
	)

	// 3. 格式化消息
	msgs, err := tpl.Format(ctx, tplData)
	if err != nil {
		return nil, fmt.Errorf("failed to format prompt template: %w", err)
	}

	// 4. 调用模型（使用 seed 确保可复现性）
	result, err := p.chatModel.Generate(ctx, msgs,
		deepseek.WithExtraFields(map[string]any{
			"seed": 42,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("AI prediction failed: %w", err)
	}

	// 5. 解析响应 JSON
	prediction, err := parsePredictionResponse(result.Content)
	if err != nil {
		// 重试一次
		logger.Warn("failed to parse AI response, retrying once",
			logger.String("raw_response", result.Content),
			logger.Err(err),
		)
		result, err = p.chatModel.Generate(ctx, msgs,
			deepseek.WithExtraFields(map[string]any{
				"seed": 42,
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("AI prediction retry failed: %w", err)
		}
		prediction, err = parsePredictionResponse(result.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AI prediction response after retry: %w, raw: %s", err, result.Content)
		}
		prediction.RawResponse = result.Content
	} else {
		prediction.RawResponse = result.Content
	}

	// 6. 计算 edge 并设置 recommended_action
	edge := prediction.PredictedProbability - marketYesPrice
	prediction.RecommendedAction = computeRecommendedAction(edge, prediction.Confidence)

	// 7. 记录预测关键数据
	logger.Info("AI prediction completed",
		logger.Float64("predicted_probability", prediction.PredictedProbability),
		logger.Float64("confidence", prediction.Confidence),
		logger.String("direction", prediction.Direction),
		logger.Float64("edge", edge),
		logger.String("recommended_action", prediction.RecommendedAction),
		logger.Float64("market_yes_price", marketYesPrice),
		logger.String("prompt_version", p.promptVersion),
	)

	return prediction, nil
}

// parsePredictionResponse 解析 AI 模型返回的 JSON 响应。
// 处理可能的 markdown 代码块包裹 (```json ... ```) 以及 markdown 表格/散文中有 JSON 对象的情况。
func parsePredictionResponse(content string) (*PredictionResult, error) {
	cleaned := content

	// Step 1: 尝试移除 markdown 代码块标记 (```json ... ```)
	if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimPrefix(cleaned, "```")
		if idx := strings.LastIndex(cleaned, "```"); idx >= 0 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	// Step 2: 直接解析
	var result PredictionResult
	if err := json.Unmarshal([]byte(cleaned), &result); err == nil {
		if validatePredictionResult(&result) {
			return &result, nil
		}
	}

	// Step 3: 兜底 — 从 markdown/散文内容中提取 JSON 对象 { ... }
	if idx := strings.Index(cleaned, "{"); idx >= 0 {
		// 找到第一个 {，然后找与之配对的最后一个 }
		start := idx
		end := strings.LastIndex(cleaned, "}")
		if end > start {
			jsonCandidate := cleaned[start : end+1]
			if err := json.Unmarshal([]byte(jsonCandidate), &result); err == nil {
				if validatePredictionResult(&result) {
					return &result, nil
				}
			}
		}
	}

	// Step 4: 全部失败
	return nil, fmt.Errorf("JSON parse error: invalid content, raw: %s", cleaned)
}

// validatePredictionResult 验证预测结果的关键字段是否合法。
func validatePredictionResult(r *PredictionResult) bool {
	return r.PredictedProbability >= 0.01 && r.PredictedProbability <= 0.99 &&
		r.Confidence >= 0 && r.Confidence <= 1
}

// computeRecommendedAction 根据 edge 和 confidence 计算推荐操作。
func computeRecommendedAction(edge, confidence float64) string {
	const minEdge = 0.25
	const minConfidence = 0.6

	if math.Abs(edge) < minEdge {
		return "skip"
	}
	if confidence < minConfidence {
		return "skip"
	}
	if edge > 0 {
		return "buy_yes"
	}
	return "buy_no"
}

// formatNewsHeadlines 将新闻头条列表格式化为字符串。
func formatNewsHeadlines(news []external.NewsItem) string {
	if len(news) == 0 {
		return "暂无新闻数据"
	}

	var sb strings.Builder
	for i, item := range news {
		if i >= 20 {
			break
		}
		fmt.Fprintf(&sb, "%d. [%s] %s - %s\n", i+1, item.Sentiment, item.Title, item.Source)
	}
	return sb.String()
}
