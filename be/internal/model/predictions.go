package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// Predictions AI预测记录表
type Predictions struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	MarketID             int              `gorm:"column:market_id;type:int(11);not null;comment:关联市场ID（逻辑外键→markets.id, 无物理约束）" json:"marketID"`
	PredictedProbability *decimal.Decimal `gorm:"column:predicted_probability;type:decimal(38,18);not null;comment:AI预测的概率值（范围: 0.01 ~ 0.99）" json:"predictedProbability"`
	Confidence           *decimal.Decimal `gorm:"column:confidence;type:decimal(38,18);not null;comment:AI对自己预测的置信度（范围: 0 ~ 1）" json:"confidence"`
	Direction            string           `gorm:"column:direction;type:varchar(16);not null;comment:预测方向: bullish-看涨, bearish-看跌, neutral-中性" json:"direction"`
	KeyFactors           *datatypes.JSON  `gorm:"column:key_factors;type:json;not null;comment:关键技术因素列表（JSON字符串数组）" json:"keyFactors"`
	RiskFactors          *datatypes.JSON  `gorm:"column:risk_factors;type:json;not null;comment:风险因素列表（JSON字符串数组）" json:"riskFactors"`
	TechnicalAnalysis    string           `gorm:"column:technical_analysis;type:text;not null;comment:技术面分析详情" json:"technicalAnalysis"`
	SentimentAnalysis    string           `gorm:"column:sentiment_analysis;type:text;not null;comment:市场情绪分析详情" json:"sentimentAnalysis"`
	NewsImpact           string           `gorm:"column:news_impact;type:text;not null;comment:新闻影响分析详情" json:"newsImpact"`
	OnchainAnalysis      string           `gorm:"column:onchain_analysis;type:text;not null;comment:链上数据分析详情" json:"onchainAnalysis"`
	Reasoning            string           `gorm:"column:reasoning;type:text;not null;comment:AI推理过程的完整描述" json:"reasoning"`
	RecommendedAction    string           `gorm:"column:recommended_action;type:varchar(16);not null;comment:推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过" json:"recommendedAction"`
	MarketProbability    *decimal.Decimal `gorm:"column:market_probability;type:decimal(38,18);not null;comment:预测时刻的 Polymarket 市场概率（Yes价格）" json:"marketProbability"`
	Edge                 *decimal.Decimal `gorm:"column:edge;type:decimal(38,18);not null;comment:AI概率与市场概率的差值 = predicted_probability - market_probability" json:"edge"`
	ModelVersion         string           `gorm:"column:model_version;type:varchar(32);not null;comment:AI模型版本号" json:"modelVersion"`
	PromptVersion        string           `gorm:"column:prompt_version;type:varchar(16);not null;comment:提示词模板版本号" json:"promptVersion"`
	Seed                 int              `gorm:"column:seed;type:int(11);not null;comment:模型推理使用的随机种子（用于复现）" json:"seed"`
	RawRequest           *datatypes.JSON  `gorm:"column:raw_request;type:json;comment:发送给AI模型的原始请求体（可选，暂未捕获）" json:"rawRequest"`
	RawResponse          *datatypes.JSON  `gorm:"column:raw_response;type:json;comment:AI模型返回的原始响应体（可选，暂未捕获）" json:"rawResponse"`
	DataSnapshot         *datatypes.JSON  `gorm:"column:data_snapshot;type:json;comment:预测时刻的各类数据快照（可选，暂未捕获）" json:"dataSnapshot"`
	TokensUsed           int              `gorm:"column:tokens_used;type:int(11);default:0;not null;comment:本次预测消耗的Token数量" json:"tokensUsed"`
	LatencyMs            int              `gorm:"column:latency_ms;type:int(11);default:0;not null;comment:AI模型推理耗时（毫秒）" json:"latencyMs"`
}

// PredictionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var PredictionsColumnNames = map[string]bool{
	"id":                    true,
	"created_at":            true,
	"updated_at":            true,
	"deleted_at":            true,
	"market_id":             true,
	"predicted_probability": true,
	"confidence":            true,
	"direction":             true,
	"key_factors":           true,
	"risk_factors":          true,
	"technical_analysis":    true,
	"sentiment_analysis":    true,
	"news_impact":           true,
	"onchain_analysis":      true,
	"reasoning":             true,
	"recommended_action":    true,
	"market_probability":    true,
	"edge":                  true,
	"model_version":         true,
	"prompt_version":        true,
	"seed":                  true,
	"raw_request":           true,
	"raw_response":          true,
	"data_snapshot":         true,
	"tokens_used":           true,
	"latency_ms":            true,
}
