package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreatePredictionsRequest request params
type CreatePredictionsRequest struct {
	MarketID             int    `json:"marketID" binding:""`             // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PredictedProbability string `json:"predictedProbability" binding:""` // AI预测的概率值（范围: 0.01 ~ 0.99）
	Confidence           string `json:"confidence" binding:""`           // AI对自己预测的置信度（范围: 0 ~ 1）
	Direction            string `json:"direction" binding:""`            // 预测方向: bullish-看涨, bearish-看跌, neutral-中性
	KeyFactors           string `json:"keyFactors" binding:""`           // 关键技术因素列表（JSON字符串数组）
	RiskFactors          string `json:"riskFactors" binding:""`          // 风险因素列表（JSON字符串数组）
	TechnicalAnalysis    string `json:"technicalAnalysis" binding:""`    // 技术面分析详情
	SentimentAnalysis    string `json:"sentimentAnalysis" binding:""`    // 市场情绪分析详情
	NewsImpact           string `json:"newsImpact" binding:""`           // 新闻影响分析详情
	OnchainAnalysis      string `json:"onchainAnalysis" binding:""`      // 链上数据分析详情
	Reasoning            string `json:"reasoning" binding:""`            // AI推理过程的完整描述
	RecommendedAction    string `json:"recommendedAction" binding:""`    // 推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	MarketProbability    string `json:"marketProbability" binding:""`    // 预测时刻的 Polymarket 市场概率（Yes价格）
	Edge                 string `json:"edge" binding:""`                 // AI概率与市场概率的差值 = predicted_probability - market_probability
	ModelVersion         string `json:"modelVersion" binding:""`         // AI模型版本号
	PromptVersion        string `json:"promptVersion" binding:""`        // 提示词模板版本号
	Seed                 int    `json:"seed" binding:""`                 // 模型推理使用的随机种子（用于复现）
	RawRequest           string `json:"rawRequest" binding:""`           // 发送给AI模型的原始请求体
	RawResponse          string `json:"rawResponse" binding:""`          // AI模型返回的原始响应体
	DataSnapshot         string `json:"dataSnapshot" binding:""`         // 预测时刻的各类数据快照
	TokensUsed           int    `json:"tokensUsed" binding:""`           // 本次预测消耗的Token数量
	LatencyMs            int    `json:"latencyMs" binding:""`            // AI模型推理耗时（毫秒）
}

// UpdatePredictionsByIDRequest request params
type UpdatePredictionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	MarketID             int    `json:"marketID" binding:""`             // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PredictedProbability string `json:"predictedProbability" binding:""` // AI预测的概率值（范围: 0.01 ~ 0.99）
	Confidence           string `json:"confidence" binding:""`           // AI对自己预测的置信度（范围: 0 ~ 1）
	Direction            string `json:"direction" binding:""`            // 预测方向: bullish-看涨, bearish-看跌, neutral-中性
	KeyFactors           string `json:"keyFactors" binding:""`           // 关键技术因素列表（JSON字符串数组）
	RiskFactors          string `json:"riskFactors" binding:""`          // 风险因素列表（JSON字符串数组）
	TechnicalAnalysis    string `json:"technicalAnalysis" binding:""`    // 技术面分析详情
	SentimentAnalysis    string `json:"sentimentAnalysis" binding:""`    // 市场情绪分析详情
	NewsImpact           string `json:"newsImpact" binding:""`           // 新闻影响分析详情
	OnchainAnalysis      string `json:"onchainAnalysis" binding:""`      // 链上数据分析详情
	Reasoning            string `json:"reasoning" binding:""`            // AI推理过程的完整描述
	RecommendedAction    string `json:"recommendedAction" binding:""`    // 推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	MarketProbability    string `json:"marketProbability" binding:""`    // 预测时刻的 Polymarket 市场概率（Yes价格）
	Edge                 string `json:"edge" binding:""`                 // AI概率与市场概率的差值 = predicted_probability - market_probability
	ModelVersion         string `json:"modelVersion" binding:""`         // AI模型版本号
	PromptVersion        string `json:"promptVersion" binding:""`        // 提示词模板版本号
	Seed                 int    `json:"seed" binding:""`                 // 模型推理使用的随机种子（用于复现）
	RawRequest           string `json:"rawRequest" binding:""`           // 发送给AI模型的原始请求体
	RawResponse          string `json:"rawResponse" binding:""`          // AI模型返回的原始响应体
	DataSnapshot         string `json:"dataSnapshot" binding:""`         // 预测时刻的各类数据快照
	TokensUsed           int    `json:"tokensUsed" binding:""`           // 本次预测消耗的Token数量
	LatencyMs            int    `json:"latencyMs" binding:""`            // AI模型推理耗时（毫秒）
}

// PredictionsObjDetail detail
type PredictionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt            *time.Time `json:"createdAt"`            // 创建时间
	UpdatedAt            *time.Time `json:"updatedAt"`            // 更新时间
	MarketID             int        `json:"marketID"`             // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PredictedProbability string     `json:"predictedProbability"` // AI预测的概率值（范围: 0.01 ~ 0.99）
	Confidence           string     `json:"confidence"`           // AI对自己预测的置信度（范围: 0 ~ 1）
	Direction            string     `json:"direction"`            // 预测方向: bullish-看涨, bearish-看跌, neutral-中性
	KeyFactors           string     `json:"keyFactors"`           // 关键技术因素列表（JSON字符串数组）
	RiskFactors          string     `json:"riskFactors"`          // 风险因素列表（JSON字符串数组）
	TechnicalAnalysis    string     `json:"technicalAnalysis"`    // 技术面分析详情
	SentimentAnalysis    string     `json:"sentimentAnalysis"`    // 市场情绪分析详情
	NewsImpact           string     `json:"newsImpact"`           // 新闻影响分析详情
	OnchainAnalysis      string     `json:"onchainAnalysis"`      // 链上数据分析详情
	Reasoning            string     `json:"reasoning"`            // AI推理过程的完整描述
	RecommendedAction    string     `json:"recommendedAction"`    // 推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	MarketProbability    string     `json:"marketProbability"`    // 预测时刻的 Polymarket 市场概率（Yes价格）
	Edge                 string     `json:"edge"`                 // AI概率与市场概率的差值 = predicted_probability - market_probability
	ModelVersion         string     `json:"modelVersion"`         // AI模型版本号
	PromptVersion        string     `json:"promptVersion"`        // 提示词模板版本号
	Seed                 int        `json:"seed"`                 // 模型推理使用的随机种子（用于复现）
	RawRequest           string     `json:"rawRequest"`           // 发送给AI模型的原始请求体
	RawResponse          string     `json:"rawResponse"`          // AI模型返回的原始响应体
	DataSnapshot         string     `json:"dataSnapshot"`         // 预测时刻的各类数据快照
	TokensUsed           int        `json:"tokensUsed"`           // 本次预测消耗的Token数量
	LatencyMs            int        `json:"latencyMs"`            // AI模型推理耗时（毫秒）
}

// CreatePredictionsReply only for api docs
type CreatePredictionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeletePredictionsByIDReply only for api docs
type DeletePredictionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdatePredictionsByIDReply only for api docs
type UpdatePredictionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetPredictionsByIDReply only for api docs
type GetPredictionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Predictions PredictionsObjDetail `json:"predictions"`
	} `json:"data"` // return data
}

// ListPredictionssRequest request params
type ListPredictionssRequest struct {
	query.Params
}

// ListPredictionssReply only for api docs
type ListPredictionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Predictionss []PredictionsObjDetail `json:"predictionss"`
	} `json:"data"` // return data
}
