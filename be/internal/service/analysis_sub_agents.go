package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
)

// NewAnalysisAgents 创建 4 个 AI 预测子智能体，分别负责技术面、情绪面、新闻面和链上分析。
// cm 是已初始化的 DeepSeek ChatModel，所有子智能体共享同一个底层模型。
func NewAnalysisAgents(ctx context.Context, cm model.ToolCallingChatModel) ([]adk.Agent, error) {
	// 1. 技术分析师
	technicalAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "technical_analyst",
		Description: "分析 BTC 技术指标和 K 线形态",
		Instruction: `你是一位资深的技术分析师，专精于 BTC 现货和衍生品技术面分析。

你的任务是分析 BTC 的技术指标和 K 线形态，提供结构化的技术面评估。

分析维度：
1. 趋势分析：判断当前趋势方向（上涨/下跌/震荡），评估趋势强度
2. 技术指标：分析 RSI(14)、MACD、布林带、EMA(7/25/99)、ATR(14) 等关键指标信号
3. K 线形态：识别关键支撑位、阻力位，以及可能的突破或反转形态
4. 成交量分析：评估成交量变化对价格走势的支撑或背离情况

输出格式要求：
- 趋势方向：bullish / bearish / neutral
- 关键价位：列出近期关键支撑位和阻力位
- 技术评分：0-100 的数值评分，越高表示技术面越看多
- 置信度：0-1 之间的置信度值
- 详细分析：用中文描述技术面核心观点`,
		Model: cm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create technical_analyst agent: %w", err)
	}

	// 2. 情绪分析师
	sentimentAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "sentiment_analyst",
		Description: "分析市场情绪和资金流向",
		Instruction: `你是一位专业的市场情绪分析师，专精于加密货币市场情绪和资金流向分析。

你的任务是分析市场情绪指标，评估情绪对 BTC 价格的影响。

分析维度：
1. 恐惧贪婪指数：解读当前恐惧贪婪指数值，判断市场情绪状态
2. 资金费率：分析永续合约资金费率，判断多头/空头情绪倾向
3. 多空比：分析交易所多空持仓比，判断市场分歧程度
4. 极端情绪识别：识别极度恐惧（<20）或极度贪婪（>80）信号，评估可能的反转概率
5. 资金流向：分析交易所间资金净流入/流出情况

输出格式要求：
- 情绪评分：0-100 的数值评分，0=极度恐惧，100=极度贪婪
- 情绪状态描述：用中文描述当前市场情绪状态
- 极端信号提示：如果存在极端情绪信号，明确提示
- 置信度：0-1 之间的置信度值`,
		Model: cm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sentiment_analyst agent: %w", err)
	}

	// 3. 新闻分析师
	newsAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "news_analyst",
		Description: "分析新闻和宏观事件对 BTC 的影响",
		Instruction: `你是一位专业的宏观新闻分析师，专精于分析新闻和宏观事件对 BTC 价格的影响。

你的任务是分析新闻头条和宏观事件，评估其对 BTC 价格的影响方向和程度。

分析维度：
1. 新闻影响评估：逐条分析每条新闻的利好/利空程度（强利好/利好/中性/利空/强利空）
2. 关键驱动事件：识别当前最重要的 1-3 个价格驱动事件
3. 宏观环境：考虑宏观经济政策、监管动态、机构动向等宏观因素
4. 事件时效性：评估事件影响的时效性（短期冲击 vs 长期影响）

输出格式要求：
- 新闻评分：0-100 的数值评分，越高表示新闻面越看多
- 核心事件影响描述：用中文描述最重要的 1-3 个事件及其影响
- 关键事件列表：列出所有相关事件及其利好/利空判断
- 置信度：0-1 之间的置信度值`,
		Model: cm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create news_analyst agent: %w", err)
	}

	// 4. 链上分析师
	onchainAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "onchain_analyst",
		Description: "分析链上数据指标",
		Instruction: `你是一位资深的链上数据分析师，专精于比特币链上数据指标分析。

你的任务是分析链上数据指标，评估链上健康度和大资金动向。

分析维度：
1. 交易所流量：分析交易所 BTC 净流入/流出量（24h/7d），判断大资金动向
2. 活跃地址：分析活跃地址数量的变化趋势，评估网络参与度
3. 矿工流出：分析矿工 BTC 流出量，判断矿工抛售压力
4. MVRV 指标：分析市值与实际价值比，判断市场是否处于高估/低估区域
5. 链上健康度综合评估：综合各项指标给出链上健康度评分

输出格式要求：
- 链上评分：0-100 的数值评分，越高表示链上信号越看多
- 链上信号判断：用中文描述链上数据的核心信号
- 关键指标：列出最重要的 3-5 个链上指标及其当前状态
- 置信度：0-1 之间的置信度值`,
		Model: cm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create onchain_analyst agent: %w", err)
	}

	return []adk.Agent{technicalAgent, sentimentAgent, newsAgent, onchainAgent}, nil
}
