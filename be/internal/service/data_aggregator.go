package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"

	"be/internal/binance"
	"be/internal/external"
	"be/internal/model"
)

// MarketDataBundle 包含聚合后的多源市场数据，供 Eino Agent 分析使用。
type MarketDataBundle struct {
	Timestamp       time.Time `json:"timestamp"`
	BTCCurrentPrice float64   `json:"btcCurrentPrice"`
	TargetPrice     int       `json:"targetPrice"`
	TargetDatetime  time.Time `json:"targetDatetime"`

	// K-line summary (OHLCV arrays)
	Kline1h7d      []external.Kline `json:"kline1h7d"`
	Kline1d30d     []external.Kline `json:"kline1d30d"`
	PriceChange24h float64          `json:"priceChange24h"`
	PriceChange7d  float64          `json:"priceChange7d"`

	// 24h 行情（来自 Binance 24hr ticker，Python 参考版对应字段）
	High24h   float64 `json:"high24h"`
	Low24h    float64 `json:"low24h"`
	Volume24h float64 `json:"volume24h"`

	// Technical indicators（使用 1h K 线计算，与 Python 参考版一致）
	RSI14     float64                `json:"rsi14"`
	MACD      external.MACDData      `json:"macd"`
	Bollinger external.BollingerData `json:"bollinger"`
	EMA       external.EMAData       `json:"ema"`
	ATR       float64                `json:"atr"`

	// Sentiment
	FearGreedIndex   int     `json:"fearGreedIndex"`
	FearGreedLabel   string  `json:"fearGreedLabel"`
	FearGreedTrend7d []int   `json:"fearGreedTrend7d"` // 7日恐惧贪婪趋势（Python 参考版对应字段）
	FundingRate      float64 `json:"fundingRate"`
	LongShortRatio   float64 `json:"longShortRatio"`

	// News
	NewsHeadlines []external.NewsItem `json:"newsHeadlines"`
	MacroEvents   []MacroEvent        `json:"macroEvents"`

	// On-chain
	ExchangeNetflow24h float64 `json:"exchangeNetflow24h"`
	ExchangeNetflow7d  float64 `json:"exchangeNetflow7d"`
	ActiveAddresses24h int     `json:"activeAddresses24h"`
	MinerOutflow24h    float64 `json:"minerOutflow24h"`
	MVRVRatio          float64 `json:"mvrvRatio"`
}

// MacroEvent 描述宏观事件信息。
type MacroEvent struct {
	Event    string `json:"event"`
	Date     string `json:"date"`
	Expected string `json:"expected"`
	Impact   string `json:"impact"`
}

// DataAggregator 负责收集聚合多源市场数据。
type DataAggregator struct {
	binanceClient   *binance.Client
	externalFetcher *external.Fetcher
}

// NewDataAggregator 创建一个新的 DataAggregator 实例。
func NewDataAggregator(binanceClient *binance.Client, extFetcher *external.Fetcher) *DataAggregator {
	return &DataAggregator{
		binanceClient:   binanceClient,
		externalFetcher: extFetcher,
	}
}

// Aggregate 收集所有数据并返回 MarketDataBundle。
// market: 目标市场对象。
func (a *DataAggregator) Aggregate(ctx context.Context, market *model.Markets) (*MarketDataBundle, error) {
	logger.Info("starting data aggregation",
		logger.String("question", market.Question),
		logger.Int("priceThreshold", market.PriceThreshold),
	)

	bundle := &MarketDataBundle{
		Timestamp:     time.Now(),
		TargetPrice:   market.PriceThreshold,
		MacroEvents:   make([]MacroEvent, 0),
		NewsHeadlines: make([]external.NewsItem, 0),
	}

	if market.TargetDate != nil {
		bundle.TargetDatetime = *market.TargetDate
	}

	symbol := "BTCUSDT"

	// 1. 获取 1h K 线（7 天 = 168 根）
	klines1h, err := a.binanceClient.GetKlines1h(symbol, 168)
	if err != nil {
		logger.Error("failed to fetch 1h klines", logger.Err(err))
	} else {
		bundle.Kline1h7d = toExternalKlines(klines1h)
		logger.Info("fetched 1h klines", logger.Int("count", len(klines1h)))
	}

	// 2. 获取 1d K 线（30 天）
	klines1d, err := a.binanceClient.GetKlines1d(symbol, 30)
	if err != nil {
		logger.Error("failed to fetch 1d klines", logger.Err(err))
	} else {
		bundle.Kline1d30d = toExternalKlines(klines1d)
		logger.Info("fetched 1d klines", logger.Int("count", len(klines1d)))
	}

	// 3. 获取 24h ticker（用于当前价格、24h 涨跌幅、24h 高/低/成交量）
	ticker, err := a.binanceClient.GetTicker24h(symbol)
	if err != nil {
		logger.Error("failed to fetch 24h ticker", logger.Err(err))
	} else {
		if lastPrice, parseErr := strconv.ParseFloat(ticker.LastPrice, 64); parseErr == nil {
			bundle.BTCCurrentPrice = lastPrice
		}
		if priceChange, parseErr := strconv.ParseFloat(ticker.PriceChange, 64); parseErr == nil && bundle.BTCCurrentPrice > 0 {
			prevClose := bundle.BTCCurrentPrice - priceChange
			if prevClose > 0 {
				bundle.PriceChange24h = (priceChange / prevClose) * 100
			}
		}
		// 24h 高/低/成交量（Python 参考版对应字段）
		if high, parseErr := strconv.ParseFloat(ticker.HighPrice, 64); parseErr == nil {
			bundle.High24h = high
		}
		if low, parseErr := strconv.ParseFloat(ticker.LowPrice, 64); parseErr == nil {
			bundle.Low24h = low
		}
		if vol, parseErr := strconv.ParseFloat(ticker.Volume, 64); parseErr == nil {
			bundle.Volume24h = vol
		}
		logger.Info("fetched 24h ticker",
			logger.Float64("lastPrice", bundle.BTCCurrentPrice),
			logger.Float64("priceChange24h", bundle.PriceChange24h),
			logger.Float64("high24h", bundle.High24h),
			logger.Float64("low24h", bundle.Low24h),
			logger.Float64("volume24h", bundle.Volume24h),
		)
	}

	// 4. 计算 7 日涨跌幅（通过 1d K 线首尾收盘价）
	if len(bundle.Kline1d30d) >= 2 {
		firstClose, parseErr := strconv.ParseFloat(tickerOrEmpty(bundle.Kline1d30d[0].Close, "0"), 64)
		if parseErr != nil {
			firstClose = 0
		}
		lastK := bundle.Kline1d30d[len(bundle.Kline1d30d)-1]
		lastClose, parseErr := strconv.ParseFloat(lastK.Close, 64)
		if parseErr != nil {
			lastClose = 0
		}
		if firstClose > 0 {
			bundle.PriceChange7d = ((lastClose - firstClose) / firstClose) * 100
		}
	}

	// 5. 获取恐惧贪婪指数及 7 日趋势（Python 参考版对应字段）
	fgi, fgl, trend, err := a.externalFetcher.GetFearGreedWithTrend(7)
	if err != nil {
		logger.Error("failed to fetch fear & greed index", logger.Err(err))
	} else {
		bundle.FearGreedIndex = fgi
		bundle.FearGreedLabel = fgl
		bundle.FearGreedTrend7d = trend
		logger.Info("fetched fear & greed index",
			logger.Int("index", fgi),
			logger.String("label", fgl),
			logger.Any("trend", trend),
		)
	}

	// 6. 计算技术指标（使用 1h K 线，与 Python 参考版一致）
	// Python 参考版：KLINE_LIMIT_1H = 24*7=168 根 1h K 线，足够计算 EMA99
	if len(bundle.Kline1h7d) >= 100 {
		ti, err := a.externalFetcher.GetTechnicalIndicators1h(bundle.Kline1h7d)
		if err != nil {
			logger.Error("failed to compute technical indicators from 1h klines", logger.Err(err))
		} else {
			bundle.RSI14 = ti.RSI14
			bundle.MACD = ti.MACD
			bundle.Bollinger = ti.Bollinger
			bundle.EMA = ti.EMA
			bundle.ATR = ti.ATR14
			logger.Info("computed technical indicators from 1h klines",
				logger.Float64("rsi14", ti.RSI14),
				logger.Float64("atr", ti.ATR14),
				logger.Int("candle_count", len(bundle.Kline1h7d)),
			)
		}
	} else {
		logger.Warn("insufficient 1h kline data for technical indicators",
			logger.Int("count", len(bundle.Kline1h7d)),
			logger.Int("required", 100),
		)
	}

	// 7. 获取新闻
	news := a.externalFetcher.GetNews("BTC", 20)
	bundle.NewsHeadlines = news
	logger.Info("fetched news headlines", logger.Int("count", len(news)))

	// 8. 获取链上数据 (CryptoQuant)
	onchainData := a.externalFetcher.GetOnChainData()
	if onchainData != nil {
		bundle.ExchangeNetflow24h = onchainData.ExchangeNetflow24h
		bundle.ExchangeNetflow7d = onchainData.ExchangeNetflow7d
		bundle.ActiveAddresses24h = onchainData.ActiveAddresses24h
		bundle.MinerOutflow24h = onchainData.MinerOutflow24h
		bundle.MVRVRatio = onchainData.MVRVRatio
		logger.Info("fetched on-chain data",
			logger.Float64("exchangeNetflow24h", onchainData.ExchangeNetflow24h),
			logger.Float64("mvrv", onchainData.MVRVRatio),
		)
	}

	// 9. FundingRate / LongShortRatio — 当前无可用 API，保留零值
	bundle.FundingRate = 0
	bundle.LongShortRatio = 0

	logger.Info("data aggregation completed",
		logger.Float64("btcPrice", bundle.BTCCurrentPrice),
		logger.Int("fearGreedIndex", bundle.FearGreedIndex),
		logger.Int("newsCount", len(bundle.NewsHeadlines)),
	)

	return bundle, nil
}

// toExternalKlines 将 binance.Kline 切片转换为 external.Kline 切片。
func toExternalKlines(klines []binance.Kline) []external.Kline {
	if klines == nil {
		return nil
	}
	out := make([]external.Kline, len(klines))
	for i, k := range klines {
		out[i] = external.Kline{
			OpenTime:  k.OpenTime,
			Open:      k.Open,
			High:      k.High,
			Low:       k.Low,
			Close:     k.Close,
			Volume:    k.Volume,
			CloseTime: k.CloseTime,
		}
	}
	return out
}

// tickerOrEmpty 返回 s；如果 s 为空则返回 fallback。
func tickerOrEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
