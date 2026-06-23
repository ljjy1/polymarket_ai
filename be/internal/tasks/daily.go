package tasks

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sasynq"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"

	"be/internal/binance"
	"be/internal/config"
	"be/internal/contract"
	"be/internal/dao"
	"be/internal/database"
	"be/internal/external"
	"be/internal/model"
	"be/internal/polymarket"
	"be/internal/service"
)

// 任务类型常量
const (
	TypeDailyScan     = "daily:scan"
	TypeDailyPredict  = "daily:predict"
	TypeDailyStrategy = "daily:strategy"
	TypeDailyExecute  = "daily:execute"
)

// ---------- 任务 Payload ----------

// ScanPayload 市场扫描任务载荷
type ScanPayload struct{}

// PredictPayload AI 预测任务载荷
type PredictPayload struct {
	MarketID int `json:"marketId"`
}

// StrategyPayload 策略生成任务载荷
type StrategyPayload struct {
	PredictionID int `json:"predictionId"`
	MarketID     int `json:"marketId"`
}

// ExecutePayload 交易执行任务载荷
type ExecutePayload struct {
	MarketID   int `json:"marketId"`
	StrategyID int `json:"strategyId"`
}

// ---------- 依赖注入 ----------

// DailyTaskDeps 日常任务所需的依赖
type DailyTaskDeps struct {
	AsynqClient       *sasynq.Client // 用于任务链入队
	MarketScanner     *service.MarketScanner
	DataAggregator    *service.DataAggregator
	AIPredictor       *service.AIPredictor    // 传统单模型预测器（fallback）
	AgentPredictor    *service.AgentPredictor // DeepAgent 多智能体预测器
	StrategyGenerator *service.StrategyGenerator
	TradeExecutor     *service.TradeExecutor
	MarketsDao        dao.MarketsDao
	PredictionsDao    dao.PredictionsDao
	StrategiesDao     dao.StrategiesDao
	TradesDao         dao.TradesDao
	VaultContract     *contract.VaultContractClient
}

// NewDailyTaskDeps 创建日常任务依赖
func NewDailyTaskDeps() (*DailyTaskDeps, error) {
	cfg := config.Get()

	// DAO 层
	db := database.GetDB()
	marketsDao := dao.NewMarketsDao(db, nil)
	predictionsDao := dao.NewPredictionsDao(db, nil)
	strategiesDao := dao.NewStrategiesDao(db, nil)
	tradesDao := dao.NewTradesDao(db, nil)

	// Polymarket 客户端
	polyClient, err := polymarket.NewClient(
		cfg.Polymarket.ClobAPIURL,
		cfg.Polymarket.GammaAPIURL,
		cfg.Polymarket.ChainID,
		cfg.Polymarket.PrivateKey,
		cfg.Polymarket.APIKey,
		cfg.Polymarket.APISecret,
		cfg.Polymarket.Passphrase,
		cfg.Proxy.Addr,
	)
	if err != nil {
		return nil, err
	}

	// Service 层
	marketScanner := service.NewMarketScanner(polyClient, marketsDao)

	binanceClient := binance.NewClient(cfg.Binance.APIURL, cfg.Proxy.Addr)
	extFetcher := external.NewFetcher(cfg.GNews.BaseURL, cfg.GNews.APIKey, cfg.CryptoQuant.APIKey, cfg.CryptoQuant.BaseURL, cfg.FearGreedIndex.URL, cfg.Proxy.Addr)
	dataAggregator := service.NewDataAggregator(binanceClient, extFetcher)

	// 创建 AI 预测器
	aiPredictor, err := service.NewAIPredictor(cfg.DeepSeek.APIKey, cfg.DeepSeek.BaseURL, cfg.DeepSeek.Model)
	if err != nil {
		return nil, err
	}

	// 创建 DeepAgent 多智能体预测器
	agentPredictor, err := service.NewAgentPredictor(context.Background(), cfg.DeepSeek.APIKey, cfg.DeepSeek.BaseURL, cfg.DeepSeek.Model)
	if err != nil {
		logger.Warn("创建 DeepAgent 预测器失败，将使用传统 AI 预测器作为 fallback", logger.Err(err))
		agentPredictor = nil
	}

	strategyGenerator := service.NewStrategyGenerator(service.DefaultStrategyConfig())

	// 初始化金库合约客户端
	var vaultContract *contract.VaultContractClient
	if cfg.Vault.RPCURL != "" && cfg.Vault.ContractAddress != "" {
		var err error
		vaultContract, err = contract.NewVaultContractClient(cfg.Vault.RPCURL, cfg.Vault.ContractAddress, cfg.Vault.StrategistPrivateKey)
		if err != nil {
			logger.Warn("[tasks] 初始化金库合约客户端失败",
				logger.Err(err),
				logger.String("rpc_url", cfg.Vault.RPCURL),
			)
		} else {
			logger.Info("[tasks] 金库合约客户端初始化成功",
				logger.String("contract", cfg.Vault.ContractAddress),
			)
		}
	} else {
		logger.Warn("[tasks] 金库合约配置为空，策略任务将跳过链上数据读取")
	}

	tradeExecutor := service.NewTradeExecutor(polyClient, tradesDao, strategiesDao, vaultContract)

	// Asynq 客户端（用于任务链入队）
	asynqClient := NewAsynqClient()

	return &DailyTaskDeps{
		AsynqClient:       asynqClient,
		MarketScanner:     marketScanner,
		DataAggregator:    dataAggregator,
		AIPredictor:       aiPredictor,
		AgentPredictor:    agentPredictor,
		StrategyGenerator: strategyGenerator,
		TradeExecutor:     tradeExecutor,
		MarketsDao:        marketsDao,
		PredictionsDao:    predictionsDao,
		StrategiesDao:     strategiesDao,
		TradesDao:         tradesDao,
		VaultContract:     vaultContract,
	}, nil
}

// ---------- 任务处理函数 ----------

// HandleScanTask 扫描 Polymarket 市场，发现可交易的市场
// 扫描成功后自动入队 Predict 任务（延迟 10 分钟等待数据就绪，Predict 内部会调用 Aggregate 获取市场数据）
func HandleScanTask(ctx context.Context, p *ScanPayload, deps *DailyTaskDeps) error {
	logger.Info("[tasks] 开始执行市场扫描任务", logger.String("type", TypeDailyScan))

	market, err := deps.MarketScanner.Scan(ctx)
	if err != nil {
		logger.Error("[tasks] 市场扫描失败", logger.Err(err))
		return err
	}

	if market == nil {
		logger.Info("[tasks] 市场扫描完成，未发现符合条件的市场")
		return nil
	}

	logger.Info("[tasks] 市场扫描完成",
		logger.Int("market_id", int(market.ID)),
		logger.String("question", market.Question),
	)

	if market.ID <= 0 {
		logger.Error("[tasks] 扫描返回的市场 ID 无效（<=0），跳过入队 Predict 任务", logger.Uint64("market_id", market.ID))
		return nil
	}

	// 入队 Predict 任务（延迟 10 分钟执行，Predict 内部会调用 Aggregate 获取实时市场数据）
	_, _, err = deps.AsynqClient.EnqueueIn(10*time.Minute, TypeDailyPredict, PredictPayload{MarketID: int(market.ID)})
	if err != nil {
		logger.Warn("[tasks] 入队 Predict 任务失败", logger.Err(err))
	} else {
		logger.Info("[tasks] 已入队 Predict 任务",
			logger.Int("market_id", int(market.ID)),
			logger.String("delay", "10m"),
		)
	}
	return nil
}

// HandlePredictTask AI 预测任务
// 根据配置使用 DeepAgent 多智能体架构或传统单模型预测
// 预测完成后保存结果到 predictions 表，并自动入队 Strategy 任务
func HandlePredictTask(ctx context.Context, p *PredictPayload, deps *DailyTaskDeps) error {
	logger.Info("[tasks] 开始执行 AI 预测任务",
		logger.String("type", TypeDailyPredict),
		logger.Int("market_id", p.MarketID),
	)

	if p.MarketID <= 0 {
		logger.Warn("[tasks] 预测任务 market_id 无效（<=0），跳过执行", logger.Int("market_id", p.MarketID))
		return nil
	}

	market, err := deps.MarketsDao.GetByID(ctx, uint64(p.MarketID))
	if err != nil {
		logger.Error("[tasks] 获取市场信息失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	bundle, err := deps.DataAggregator.Aggregate(ctx, market)
	if err != nil {
		logger.Error("[tasks] 预测阶段数据聚合失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	var yesPrice float64
	if market.CurrentYesPrice != nil {
		yesPrice = market.CurrentYesPrice.InexactFloat64()
	}

	cfg := config.Get()
	useAgent := cfg.DeepSeek.UseAgentPredictor && deps.AgentPredictor != nil

	var result *service.PredictionResult
	if useAgent {
		logger.Info("[tasks] 使用 DeepAgent 多智能体架构进行预测")
		result, err = deps.AgentPredictor.Predict(ctx, bundle, yesPrice)
	} else {
		logger.Info("[tasks] 使用传统 AI 预测器进行预测")
		result, err = deps.AIPredictor.Predict(ctx, bundle, yesPrice)
	}
	if err != nil {
		logger.Error("[tasks] AI 预测失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	logger.Info("[tasks] AI 预测完成",
		logger.Int("market_id", p.MarketID),
		logger.Bool("use_agent", useAgent),
		logger.Float64("probability", result.PredictedProbability),
		logger.Float64("confidence", result.Confidence),
	)

	// 计算 Edge
	edge := result.PredictedProbability - yesPrice
	edgeDec := decimal.NewFromFloat(edge)

	// 保存预测结果到 predictions 表
	keyFactorsJSON, _ := json.Marshal(result.KeyFactors)
	riskFactorsJSON, _ := json.Marshal(result.RiskFactors)

	probDec := decimal.NewFromFloat(result.PredictedProbability)
	confDec := decimal.NewFromFloat(result.Confidence)
	mpDec := decimal.NewFromFloat(yesPrice)

	// 构造 raw_request（发送给 AI 的请求概要）
	rawReq := map[string]any{
		"market_id":        p.MarketID,
		"market_question":  market.Question,
		"yes_price":        yesPrice,
		"bundle_timestamp": bundle.Timestamp,
		"btc_price":        bundle.BTCCurrentPrice,
		"target_price":     bundle.TargetPrice,
		"target_datetime":  bundle.TargetDatetime,
		"fgi":              bundle.FearGreedIndex,
		"news_count":       len(bundle.NewsHeadlines),
	}
	rawReqJSON, _ := json.Marshal(rawReq)

	// 构造 data_snapshot（预测时刻的完整聚合数据快照）
	snapshotJSON, _ := json.Marshal(bundle)

	// 构造 raw_response
	rawRespJSON := json.RawMessage(result.RawResponse)

	prediction := &model.Predictions{
		MarketID:             p.MarketID,
		PredictedProbability: &probDec,
		Confidence:           &confDec,
		Direction:            result.Direction,
		KeyFactors:           (*datatypes.JSON)(&keyFactorsJSON),
		RiskFactors:          (*datatypes.JSON)(&riskFactorsJSON),
		TechnicalAnalysis:    result.TechnicalAnalysis,
		SentimentAnalysis:    result.SentimentAnalysis,
		NewsImpact:           result.NewsImpact,
		OnchainAnalysis:      result.OnchainAnalysis,
		Reasoning:            result.Reasoning,
		RecommendedAction:    result.RecommendedAction,
		MarketProbability:    &mpDec,
		Edge:                 &edgeDec,
		RawRequest:           (*datatypes.JSON)(&rawReqJSON),
		RawResponse:          (*datatypes.JSON)(&rawRespJSON),
		DataSnapshot:         (*datatypes.JSON)(&snapshotJSON),
	}
	if err := deps.PredictionsDao.Create(ctx, prediction); err != nil {
		logger.Error("[tasks] 保存预测记录失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	logger.Info("[tasks] 预测记录已保存到数据库",
		logger.Uint64("prediction_id", prediction.ID),
		logger.Int("market_id", p.MarketID),
		logger.Float64("edge", edge),
	)

	// 入队 Strategy 任务
	_, _, err = deps.AsynqClient.EnqueueNow(TypeDailyStrategy, StrategyPayload{
		PredictionID: int(prediction.ID),
		MarketID:     p.MarketID,
	})
	if err != nil {
		logger.Warn("[tasks] 入队 Strategy 任务失败", logger.Err(err))
	} else {
		logger.Info("[tasks] 已入队 Strategy 任务",
			logger.Int("prediction_id", int(prediction.ID)),
			logger.Int("market_id", p.MarketID),
		)
	}
	return nil
}

// HandleStrategyTask 策略生成任务
// 策略生成后保存到 strategies 表，非 skip 时自动入队 Execute 任务
func HandleStrategyTask(ctx context.Context, p *StrategyPayload, deps *DailyTaskDeps) error {
	logger.Info("[tasks] 开始执行策略生成任务",
		logger.String("type", TypeDailyStrategy),
		logger.Int("prediction_id", p.PredictionID),
		logger.Int("market_id", p.MarketID),
	)

	if p.MarketID <= 0 || p.PredictionID <= 0 {
		logger.Warn("[tasks] 策略任务参数无效（market_id 或 prediction_id <=0），跳过执行",
			logger.Int("market_id", p.MarketID),
			logger.Int("prediction_id", p.PredictionID),
		)
		return nil
	}

	prediction, err := deps.PredictionsDao.GetByID(ctx, uint64(p.PredictionID))
	if err != nil {
		logger.Error("[tasks] 获取预测记录失败", logger.Err(err), logger.Int("prediction_id", p.PredictionID))
		return err
	}

	market, err := deps.MarketsDao.GetByID(ctx, uint64(p.MarketID))
	if err != nil {
		logger.Error("[tasks] 获取市场信息失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	predResult := &service.PredictionResult{
		PredictedProbability: prediction.PredictedProbability.InexactFloat64(),
		Confidence:           prediction.Confidence.InexactFloat64(),
		Direction:            prediction.Direction,
	}

	// 从金库读取可用余额
	var vaultBalance float64
	if deps.VaultContract != nil {
		balance, err := deps.VaultContract.AvailableBalance(ctx)
		if err != nil {
			logger.Warn("[tasks] 读取金库余额失败，跳过策略生成",
				logger.Err(err),
				logger.Int("market_id", p.MarketID),
			)
			// 跳过策略生成，不入队 Execute
			return nil
		}
		vaultBalance = float64(balance.Int64()) / 1e6 // USDC 6位精度
	} else {
		logger.Warn("[tasks] 金库合约未配置，跳过策略生成",
			logger.Int("market_id", p.MarketID),
		)
		return nil
	}

	var yesPrice float64
	if market.CurrentYesPrice != nil {
		yesPrice = market.CurrentYesPrice.InexactFloat64()
	}

	strategyResult := deps.StrategyGenerator.Generate(predResult, yesPrice, vaultBalance)

	logger.Info("[tasks] 策略生成完成",
		logger.Int("market_id", p.MarketID),
		logger.String("action", strategyResult.Action),
		logger.Float64("position_size", strategyResult.PositionSize),
	)

	// 保存策略到 strategies 表
	posSizeDec := decimal.NewFromFloat(strategyResult.PositionSize)
	entryPriceDec := decimal.NewFromFloat(strategyResult.EntryPrice)
	tpDec := decimal.NewFromFloat(strategyResult.TakeProfit)
	slDec := decimal.NewFromFloat(strategyResult.StopLoss)
	kellyDec := decimal.NewFromFloat(strategyResult.KellyFraction)
	edgeDec := decimal.NewFromFloat(strategyResult.Edge)

	status := "pending"
	if strategyResult.Action == "skip" {
		status = "skipped"
	}

	strategy := &model.Strategies{
		PredictionID:  p.PredictionID,
		MarketID:      p.MarketID,
		Action:        strategyResult.Action,
		Side:          strategyResult.Side,
		PositionSize:  &posSizeDec,
		EntryPrice:    &entryPriceDec,
		TakeProfit:    &tpDec,
		StopLoss:      &slDec,
		KellyFraction: &kellyDec,
		Edge:          &edgeDec,
		SkipReason:    strategyResult.SkipReason,
		Status:        status,
	}
	if err := deps.StrategiesDao.Create(ctx, strategy); err != nil {
		logger.Error("[tasks] 保存策略记录失败", logger.Err(err),
			logger.Int("market_id", p.MarketID),
			logger.Int("prediction_id", p.PredictionID),
		)
		return err
	}

	logger.Info("[tasks] 策略记录已保存到数据库",
		logger.Uint64("strategy_id", strategy.ID),
		logger.String("action", strategyResult.Action),
		logger.Float64("position_size", strategyResult.PositionSize),
	)

	// 非 skip 策略才入队 Execute
	if strategyResult.Action != "skip" {
		_, _, err = deps.AsynqClient.EnqueueNow(TypeDailyExecute, ExecutePayload{
			MarketID:   p.MarketID,
			StrategyID: int(strategy.ID),
		})
		if err != nil {
			logger.Warn("[tasks] 入队 Execute 任务失败", logger.Err(err))
		} else {
			logger.Info("[tasks] 已入队 Execute 任务",
				logger.Int("market_id", p.MarketID),
				logger.Int("strategy_id", int(strategy.ID)),
			)
		}
	} else {
		logger.Info("[tasks] 策略为 skip，不触发 Execute",
			logger.String("skip_reason", strategyResult.SkipReason),
		)
	}
	return nil
}

// HandleExecuteTask 交易执行任务
func HandleExecuteTask(ctx context.Context, p *ExecutePayload, deps *DailyTaskDeps) error {
	logger.Info("[tasks] 开始执行交易任务",
		logger.String("type", TypeDailyExecute),
		logger.Int("market_id", p.MarketID),
		logger.Int("strategy_id", p.StrategyID),
	)

	if p.MarketID <= 0 || p.StrategyID <= 0 {
		logger.Warn("[tasks] 交易任务参数无效（market_id 或 strategy_id <=0），跳过执行",
			logger.Int("market_id", p.MarketID),
			logger.Int("strategy_id", p.StrategyID),
		)
		return nil
	}

	market, err := deps.MarketsDao.GetByID(ctx, uint64(p.MarketID))
	if err != nil {
		logger.Error("[tasks] 获取市场信息失败", logger.Err(err), logger.Int("market_id", p.MarketID))
		return err
	}

	strategy, err := deps.StrategiesDao.GetByID(ctx, uint64(p.StrategyID))
	if err != nil {
		logger.Error("[tasks] 获取策略记录失败", logger.Err(err), logger.Int("strategy_id", p.StrategyID))
		return err
	}

	strategyResult := &service.StrategyResult{
		ID:     p.StrategyID,
		Action: strategy.Action,
		Side:   strategy.Side,
	}

	// 非 skip 策略：先从金库提取策略资金
	if strategyResult.Action != "skip" && deps.VaultContract != nil {
		if strategy.PositionSize != nil {
			// 将 positionSize（USDC，6位精度）转换为 *big.Int
			posFloat, _ := strategy.PositionSize.Float64()
			amount := new(big.Int).SetUint64(uint64(posFloat * 1e6))
			if err := deps.VaultContract.WithdrawToStrategy(ctx, amount); err != nil {
				logger.Error("[tasks] 从金库提取策略资金失败",
					logger.Err(err),
					logger.Int("strategy_id", p.StrategyID),
					logger.String("amount", amount.String()),
				)
				// 更新策略状态为 failed
				strategy.Status = "failed"
				_ = deps.StrategiesDao.UpdateByID(ctx, strategy)
				return err
			}
			logger.Info("[tasks] 金库资金提取成功",
				logger.Int("strategy_id", p.StrategyID),
				logger.String("amount", amount.String()),
			)
		}
	}

	trade, err := deps.TradeExecutor.Execute(ctx, strategyResult, market)
	if err != nil {
		logger.Error("[tasks] 交易执行失败", logger.Err(err),
			logger.Int("market_id", p.MarketID),
			logger.Int("strategy_id", p.StrategyID),
		)
		return err
	}

	if trade != nil {
		logger.Info("[tasks] 交易执行完成", logger.Int("trade_id", int(trade.ID)))
	} else {
		logger.Info("[tasks] 交易执行完成（策略跳过）")
	}
	return nil
}

// ---------- 任务注册 ----------

// RegisterDailyTasks 注册日常流水线任务处理函数
func RegisterDailyTasks(srv *sasynq.Server, deps *DailyTaskDeps) {
	mux := srv.Mux()

	// 使用 HandleFunc 包装依赖注入
	sasynq.RegisterTaskHandler(mux, TypeDailyScan, sasynq.HandleFunc(func(ctx context.Context, p ScanPayload) error {
		return HandleScanTask(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeDailyPredict, sasynq.HandleFunc(func(ctx context.Context, p PredictPayload) error {
		return HandlePredictTask(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeDailyStrategy, sasynq.HandleFunc(func(ctx context.Context, p StrategyPayload) error {
		return HandleStrategyTask(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeDailyExecute, sasynq.HandleFunc(func(ctx context.Context, p ExecutePayload) error {
		return HandleExecuteTask(ctx, &p, deps)
	}))
}

// RegisterDailyCron 注册日常流水线定时任务
func RegisterDailyCron(scheduler *sasynq.Scheduler) {
	// 每天早上 8:00 执行市场扫描（流水线起点）
	scheduler.RegisterTask("0 8 * * *", TypeDailyScan, ScanPayload{})
	// Aggregate / Predict / Strategy / Execute 由任务链自动触发，无需独立 cron

	logger.Info("[tasks] 日常流水线定时任务注册完成")
}
