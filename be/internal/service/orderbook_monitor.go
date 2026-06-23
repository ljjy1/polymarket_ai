package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/shopspring/decimal"

	"be/internal/dao"
	bePolymarket "be/internal/polymarket"
)

// OrderBookData 盘口数据缓存
type OrderBookData struct {
	TokenID   string
	Bids      []bePolymarket.OrderBookLevel
	Asks      []bePolymarket.OrderBookLevel
	YesPrice  float64 // 中间价 (best_bid + best_ask) / 2
	Spread    float64 // 价差
	NoPrice   float64 // 1 - yesPrice
	UpdatedAt time.Time
}

// OrderBookMonitor 盘口监控器
// 通过 WebSocket 实时订阅 Polymarket 盘口数据，维护内存中的最新价格
type OrderBookMonitor struct {
	polymarketClient *bePolymarket.Client
	marketDao        dao.MarketsDao

	wsMarketURL string

	mu            sync.RWMutex
	orderBooks    map[string]*OrderBookData // tokenID -> 最新盘口
	subscribedIDs map[string]bool           // 当前已订阅的 tokenID
	wsClient      *bePolymarket.MarketWSClient
	wsCancel      context.CancelFunc

	priceUpdateHooks []func(tokenID string, yesPrice, noPrice float64)

	// marketResolved 并发控制：记录已通过 WS 处理结算的市场 slug
	resolvedMarkets sync.Map
}

// NewOrderBookMonitor 创建盘口监控器
func NewOrderBookMonitor(polyClient *bePolymarket.Client, marketDao dao.MarketsDao, wsMarketURL string) *OrderBookMonitor {
	return &OrderBookMonitor{
		polymarketClient: polyClient,
		marketDao:        marketDao,
		wsMarketURL:      wsMarketURL,
		orderBooks:       make(map[string]*OrderBookData),
		subscribedIDs:    make(map[string]bool),
	}
}

// OnPriceUpdate 注册价格更新回调
func (m *OrderBookMonitor) OnPriceUpdate(hook func(tokenID string, yesPrice, noPrice float64)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.priceUpdateHooks = append(m.priceUpdateHooks, hook)
}

// String 实现 app.IServer 接口
func (m *OrderBookMonitor) String() string {
	return "orderbook-monitor-ws"
}

// Start 启动盘口监控（实现 app.IServer 接口，阻塞运行）
func (m *OrderBookMonitor) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	m.mu.Lock()
	m.wsCancel = cancel
	m.mu.Unlock()

	logger.Info("[盘口监控] 启动 WebSocket 盘口实时监控")

	for {
		select {
		case <-ctx.Done():
			logger.Info("[盘口监控] 上下文取消，停止监控")
			return ctx.Err()
		default:
		}

		// 获取所有活跃市场的 token ID
		tokenIDs, err := m.getActiveTokenIDs(ctx)
		if err != nil {
			logger.Warn("[盘口监控] 查询活跃市场失败", logger.Err(err))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(30 * time.Second):
				continue
			}
		}

		if len(tokenIDs) == 0 {
			logger.Info("[盘口监控] 暂无活跃市场，30 秒后重试")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(30 * time.Second):
				continue
			}
		}

		logger.Info("[盘口监控] 开始订阅盘口数据",
			logger.Int("token_count", len(tokenIDs)),
			logger.String("token_ids", strings.Join(tokenIDs, ",")),
		)

		// 连接 WebSocket 并订阅
		err = m.subscribeAndMonitor(ctx, tokenIDs)
		if err != nil {
			logger.Warn("[盘口监控] WebSocket 连接断开，20 秒后重试",
				logger.Err(err),
			)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(20 * time.Second):
			}
		}
	}
}

// subscribeAndMonitor 订阅指定 token IDs 的盘口并持续监控
func (m *OrderBookMonitor) subscribeAndMonitor(ctx context.Context, tokenIDs []string) error {
	callback := &bePolymarket.WSCallback{
		OnBook: func(snapshot bePolymarket.OrderBookSnapshot) {
			m.handleBook(snapshot)
		},
		OnPriceChange: func(change bePolymarket.PriceChange) {
			m.handlePriceChange(change)
		},
		OnLastTradePrice: func(trade bePolymarket.LastTrade) {
			m.handleLastTradePrice(trade)
		},
		OnBestBidAsk: func(event bePolymarket.BestBidAskEvent) {
			m.handleBestBidAsk(event)
		},
		OnNewMarket: func(event bePolymarket.NewMarketEvent) {
			m.handleNewMarket(event)
		},
		OnMarketResolved: func(event bePolymarket.MarketResolvedEvent) {
			m.handleMarketResolved(event)
		},
		OnTickSizeChange: func(event bePolymarket.TickSizeChangeEvent) {
			m.handleTickSizeChange(event)
		},
	}

	m.mu.Lock()
	m.wsClient = bePolymarket.NewMarketWSClient(tokenIDs, m.wsMarketURL, callback)
	for _, tid := range tokenIDs {
		m.subscribedIDs[tid] = true
	}
	m.mu.Unlock()

	return m.wsClient.Start(ctx)
}

// Stop 停止盘口监控（实现 app.IServer 接口）
func (m *OrderBookMonitor) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.wsCancel != nil {
		m.wsCancel()
	}
	if m.wsClient != nil {
		m.wsClient.Close()
	}
	logger.Info("[盘口监控] 已停止")
	return nil
}

// GetPrice 获取指定 token 的最新价格
func (m *OrderBookMonitor) GetPrice(tokenID string) (yesPrice, noPrice, spread float64, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.orderBooks[tokenID]
	if !exists {
		return 0, 0, 0, false
	}
	return data.YesPrice, data.NoPrice, data.Spread, true
}

// GetOrderBook 获取指定 token 的最新盘口数据
func (m *OrderBookMonitor) GetOrderBook(tokenID string) *OrderBookData {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.orderBooks[tokenID]
	if !exists {
		return nil
	}
	return data
}

// ============================================================
//  事件处理
// ============================================================

// handleBook 处理盘口快照事件
func (m *OrderBookMonitor) handleBook(snapshot bePolymarket.OrderBookSnapshot) {
	yesPrice, spread := calculateYesPrice(snapshot.Bids, snapshot.Asks)
	if yesPrice <= 0 {
		return
	}
	noPrice := math.Round((1-yesPrice)*10000) / 10000
	yesPrice = math.Round(yesPrice*10000) / 10000

	data := &OrderBookData{
		TokenID:   snapshot.AssetID,
		Bids:      snapshot.Bids,
		Asks:      snapshot.Asks,
		YesPrice:  yesPrice,
		Spread:    spread,
		NoPrice:   noPrice,
		UpdatedAt: snapshot.Timestamp,
	}

	m.mu.Lock()
	m.orderBooks[snapshot.AssetID] = data
	hooks := m.priceUpdateHooks
	m.mu.Unlock()

	// 触发回调
	for _, hook := range hooks {
		hook(snapshot.AssetID, yesPrice, noPrice)
	}
}

// handlePriceChange 处理价格变动事件
func (m *OrderBookMonitor) handlePriceChange(change bePolymarket.PriceChange) {
	bestBid, _ := strconv.ParseFloat(change.BestBid, 64)
	bestAsk, _ := strconv.ParseFloat(change.BestAsk, 64)
	if bestBid <= 0 || bestAsk <= 0 || bestAsk < bestBid {
		return
	}

	yesPrice := math.Round(((bestBid+bestAsk)/2)*10000) / 10000
	spread := math.Round((bestAsk-bestBid)*10000) / 10000
	noPrice := math.Round((1-yesPrice)*10000) / 10000

	m.mu.Lock()
	if existing, ok := m.orderBooks[change.AssetID]; ok {
		existing.YesPrice = yesPrice
		existing.NoPrice = noPrice
		existing.Spread = spread
		existing.UpdatedAt = time.Now()
	} else {
		m.orderBooks[change.AssetID] = &OrderBookData{
			TokenID:   change.AssetID,
			YesPrice:  yesPrice,
			NoPrice:   noPrice,
			Spread:    spread,
			UpdatedAt: time.Now(),
		}
	}
	hooks := m.priceUpdateHooks
	m.mu.Unlock()

	// 触发回调
	for _, hook := range hooks {
		hook(change.AssetID, yesPrice, noPrice)
	}
}

// handleLastTradePrice 处理最新成交事件 —— 追踪成交价，发现大额偏离
func (m *OrderBookMonitor) handleLastTradePrice(trade bePolymarket.LastTrade) {
	price, err := strconv.ParseFloat(trade.Price, 64)
	if err != nil {
		return
	}
	size, err := strconv.ParseFloat(trade.Size, 64)
	if err != nil {
		return
	}

	// 对比成交价与盘口中间价的偏离
	m.mu.RLock()
	ob := m.orderBooks[trade.AssetID]
	m.mu.RUnlock()

	deviation := 0.0
	if ob != nil && ob.YesPrice > 0 {
		deviation = math.Abs(price-ob.YesPrice) / ob.YesPrice * 100
	}

	logCtx := []logger.Field{
		logger.String("asset_id", trade.AssetID),
		logger.String("price", trade.Price),
		logger.String("side", trade.Side),
		logger.Float64("size", size),
	}

	// 大额偏离预警（成交价偏离中间价 > 15%）
	if deviation > 15.0 {
		logger.Warn("[盘口监控] 大额偏离成交",
			logger.String("asset_id", trade.AssetID),
			logger.String("price", trade.Price),
			logger.String("side", trade.Side),
			logger.Float64("size", size),
			logger.Float64("deviation_pct", deviation),
		)
	} else {
		logger.Debug("[盘口监控] 最新成交", logCtx...)
	}
}

// handleBestBidAsk 处理最优买卖价事件 —— 直接用 server 算好的 best_bid / best_ask 更新价格
func (m *OrderBookMonitor) handleBestBidAsk(event bePolymarket.BestBidAskEvent) {
	bestBid, err1 := strconv.ParseFloat(event.BestBid, 64)
	bestAsk, err2 := strconv.ParseFloat(event.BestAsk, 64)
	if err1 != nil || err2 != nil || bestBid <= 0 || bestAsk <= 0 || bestAsk < bestBid {
		return
	}

	yesPrice := math.Round(((bestBid+bestAsk)/2)*10000) / 10000
	spread := bestAsk - bestBid
	noPrice := math.Round((1-yesPrice)*10000) / 10000

	m.mu.Lock()
	if existing, ok := m.orderBooks[event.AssetID]; ok {
		existing.YesPrice = yesPrice
		existing.NoPrice = noPrice
		existing.Spread = spread
		existing.UpdatedAt = time.Now()
	} else {
		m.orderBooks[event.AssetID] = &OrderBookData{
			TokenID:   event.AssetID,
			YesPrice:  yesPrice,
			NoPrice:   noPrice,
			Spread:    spread,
			UpdatedAt: time.Now(),
		}
	}
	m.mu.Unlock()
}

// handleNewMarket 处理新市场创建事件 —— 检查是否为 BTC 相关市场
func (m *OrderBookMonitor) handleNewMarket(event bePolymarket.NewMarketEvent) {
	// 检查 tag 是否包含 bitcoin
	hasBitcoinTag := false
	for _, tag := range event.Tags {
		if strings.EqualFold(tag, "bitcoin") || strings.EqualFold(tag, "btc") {
			hasBitcoinTag = true
			break
		}
	}

	logger.Info("[盘口监控] 新市场创建",
		logger.String("slug", event.Slug),
		logger.String("question", event.Question),
		logger.Bool("is_bitcoin_related", hasBitcoinTag),
		logger.Int("tag_count", len(event.Tags)),
	)

	if hasBitcoinTag {
		logger.Info("[盘口监控] 发现 BTC 相关新市场，可触发增量扫描",
			logger.String("slug", event.Slug),
		)
		// 增量扫描逻辑由上层决定，此处仅记录日志
	}
}

// handleMarketResolved 处理市场结算事件
// 并发控制：sync.Map 记录已处理的 slug，防止 WS 重连重复触发。
// 注意：仅更新 DB 状态，PnL 结算由 cron HandleSettlement 处理，定时任务仍然保留。
func (m *OrderBookMonitor) handleMarketResolved(event bePolymarket.MarketResolvedEvent) {
	// 并发控制：每个 slug 只处理一次
	if _, already := m.resolvedMarkets.LoadOrStore(event.Slug, true); already {
		logger.Debug("[盘口监控] 市场已通过 WS 处理过结算，跳过",
			logger.String("slug", event.Slug),
		)
		return
	}

	logger.Info("[盘口监控] 市场已结算（WS 实时通知）",
		logger.String("slug", event.Slug),
		logger.String("question", event.Question),
		logger.String("winning_outcome", event.WinningOutcome),
	)

	// 查找 DB 中对应 market（通过 event_slug）
	ctx := context.Background()
	markets, total, err := m.marketDao.GetByColumns(ctx, &query.Params{
		Page:  1,
		Limit: 1,
		Columns: []query.Column{
			{Name: "event_slug", Value: event.Slug},
		},
	})
	if err != nil || total == 0 {
		logger.Warn("[盘口监控] 未找到对应市场记录",
			logger.String("slug", event.Slug),
			logger.Err(err),
		)
		return
	}

	market := markets[0]
	if market.Status == "resolved" {
		logger.Debug("[盘口监控] 市场已在 DB 中标记为 resolved，跳过",
			logger.String("slug", event.Slug),
		)
		return
	}

	// 更新 DB
	market.Status = "resolved"
	market.Resolution = event.WinningOutcome
	if err := m.marketDao.UpdateByID(ctx, market); err != nil {
		logger.Error("[盘口监控] 更新市场结算状态失败",
			logger.Err(err),
			logger.String("slug", event.Slug),
		)
		return
	}

	logger.Info("[盘口监控] 市场结算状态已更新到 DB（PnL 将由 cron 定时结算处理）",
		logger.String("slug", event.Slug),
		logger.String("resolution", event.WinningOutcome),
	)
}

// handleTickSizeChange 处理最小报价单位变动事件
func (m *OrderBookMonitor) handleTickSizeChange(event bePolymarket.TickSizeChangeEvent) {
	logger.Info("[盘口监控] Tick Size 变动",
		logger.String("asset_id", event.AssetID),
		logger.String("old_tick_size", event.OldTickSize),
		logger.String("new_tick_size", event.NewTickSize),
	)
}

// ============================================================
//  辅助方法
// ============================================================

// getActiveTokenIDs 从数据库获取所有活跃市场的 token ID
func (m *OrderBookMonitor) getActiveTokenIDs(ctx context.Context) ([]string, error) {
	params := &query.Params{
		Page:  0,
		Limit: 100,
		Columns: []query.Column{
			{
				Name:  "status",
				Exp:   "=",
				Value: "active",
			},
		},
	}
	markets, total, err := m.marketDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("查询活跃市场失败: %w", err)
	}
	if total == 0 {
		return nil, nil
	}

	tokenIDs := make([]string, 0, len(markets))
	for _, mkt := range markets {
		if mkt.PolymarketTokenID != "" {
			tokenIDs = append(tokenIDs, mkt.PolymarketTokenID)
		}
	}
	return tokenIDs, nil
}

// UpdateDatabasePrices 将内存中的最新价格批量更新到数据库
func (m *OrderBookMonitor) UpdateDatabasePrices(ctx context.Context) {
	m.mu.RLock()
	snapshots := make([]*OrderBookData, 0, len(m.orderBooks))
	for _, data := range m.orderBooks {
		snapshots = append(snapshots, data)
	}
	m.mu.RUnlock()

	for _, data := range snapshots {
		// 查询对应的市场记录
		params := &query.Params{
			Columns: []query.Column{
				{
					Name:  "polymarket_token_id",
					Exp:   "=",
					Value: data.TokenID,
				},
			},
		}
		markets, total, err := m.marketDao.GetByColumns(ctx, params)
		if err != nil {
			logger.Warn("[盘口监控] 查询市场记录失败",
				logger.Err(err),
				logger.String("token_id", data.TokenID),
			)
			continue
		}
		if total == 0 || len(markets) == 0 {
			continue
		}

		market := markets[0]
		yesPriceDec := decimal.NewFromFloat(data.YesPrice)
		noPriceDec := decimal.NewFromFloat(data.NoPrice)
		market.CurrentYesPrice = &yesPriceDec
		market.CurrentNoPrice = &noPriceDec

		if err := m.marketDao.UpdateByID(ctx, market); err != nil {
			logger.Warn("[盘口监控] 更新市场价格失败",
				logger.Err(err),
				logger.String("token_id", data.TokenID),
			)
		}
	}

	logger.Info("[盘口监控] 数据库价格已更新",
		logger.Int("market_count", len(snapshots)),
	)
}

// calculateYesPrice 从买卖盘计算中间价
func calculateYesPrice(bids, asks []bePolymarket.OrderBookLevel) (float64, float64) {
	if len(bids) == 0 || len(asks) == 0 {
		return 0, 0
	}

	bestBid, err1 := strconv.ParseFloat(bids[0].Price, 64)
	bestAsk, err2 := strconv.ParseFloat(asks[0].Price, 64)

	if err1 != nil || err2 != nil || bestBid <= 0 || bestAsk <= 0 || bestAsk < bestBid {
		return 0, 0
	}

	yesPrice := (bestBid + bestAsk) / 2
	spread := bestAsk - bestBid
	return yesPrice, spread
}
