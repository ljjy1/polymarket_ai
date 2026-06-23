package polymarket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"

	polymarketSDK "github.com/0xNetuser/Polymarket-golang/polymarket"
)

// ============================================================
//  事件类型 —— 对应 Polymarket WebSocket market channel
//  https://docs.polymarket.com/market-data/websocket/market-channel
// ============================================================

// OrderBookSnapshot 盘口快照 (event_type: book)
type OrderBookSnapshot struct {
	AssetID   string
	Bids      []OrderBookLevel
	Asks      []OrderBookLevel
	Timestamp time.Time
}

// OrderBookLevel 盘口价位
type OrderBookLevel struct {
	Price string
	Size  string
}

// PriceChange 价格变动 (event_type: price_change)
type PriceChange struct {
	AssetID string
	Price   string
	Size    string
	Side    string
	BestBid string
	BestAsk string
}

// LastTrade 最新成交 (event_type: last_trade_price)
type LastTrade struct {
	AssetID   string
	Price     string
	Side      string
	Size      string
	Timestamp time.Time
}

// BestBidAskEvent 最优买卖价 (event_type: best_bid_ask, 需要 custom_feature_enabled=true)
type BestBidAskEvent struct {
	AssetID string
	Market  string
	BestBid string
	BestAsk string
	Spread  string
}

// NewMarketEvent 新市场创建 (event_type: new_market, 需要 custom_feature_enabled=true)
type NewMarketEvent struct {
	ID          string   `json:"id"`
	Question    string   `json:"question"`
	Market      string   `json:"market"`
	Slug        string   `json:"slug"`
	AssetsIDs   []string `json:"assets_ids"`
	Outcomes    []string `json:"outcomes"`
	Tags        []string `json:"tags"`
	ConditionID string   `json:"condition_id"`
	Active      bool     `json:"active"`
}

// MarketResolvedEvent 市场结算 (event_type: market_resolved, 需要 custom_feature_enabled=true)
type MarketResolvedEvent struct {
	ID             string   `json:"id"`
	Question       string   `json:"question"`
	Market         string   `json:"market"`
	Slug           string   `json:"slug"`
	AssetsIDs      []string `json:"assets_ids"`
	Outcomes       []string `json:"outcomes"`
	WinningAssetID string   `json:"winning_asset_id"`
	WinningOutcome string   `json:"winning_outcome"`
}

// TickSizeChangeEvent 最小报价单位变动 (event_type: tick_size_change)
type TickSizeChangeEvent struct {
	AssetID     string
	Market      string
	OldTickSize string
	NewTickSize string
}

// ============================================================
//  全局事件回调
// ============================================================

// WSCallback WebSocket 事件回调。
// 各字段可空，不注册 = 丢弃该类型事件。
type WSCallback struct {
	OnBook           func(snapshot OrderBookSnapshot)
	OnPriceChange    func(change PriceChange)
	OnLastTradePrice func(trade LastTrade)
	OnBestBidAsk     func(event BestBidAskEvent)
	OnNewMarket      func(event NewMarketEvent)
	OnMarketResolved func(event MarketResolvedEvent)
	OnTickSizeChange func(event TickSizeChangeEvent)
}

// ============================================================
//  WebSocket 客户端
// ============================================================

// MarketWSClient Polymarket 市场 WebSocket 客户端封装
type MarketWSClient struct {
	tokenIDs []string
	callback *WSCallback
	wsURL    string

	mu      sync.Mutex
	cancel  context.CancelFunc
	running bool
}

// NewMarketWSClient 创建市场 WebSocket 客户端
func NewMarketWSClient(tokenIDs []string, wsURL string, callback *WSCallback) *MarketWSClient {
	return &MarketWSClient{
		tokenIDs: tokenIDs,
		wsURL:    wsURL,
		callback: callback,
	}
}

// Start 启动 WebSocket 连接（阻塞运行，直到 ctx 取消或连接断开）
func (c *MarketWSClient) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("WebSocket 客户端已在运行")
	}
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
	}()

	logger.Info("[WS] 正在连接 Polymarket 市场 WebSocket",
		logger.String("url", c.wsURL),
		logger.Int("token_count", len(c.tokenIDs)),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Info("[WS] 上下文取消，WebSocket 连接已停止")
			return ctx.Err()
		default:
		}

		err := c.run(ctx)
		if err != nil {
			logger.Warn("[WS] WebSocket 连接断开，3 秒后重连",
				logger.Err(err),
			)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
			}
		}
	}
}

// run 建立单次 WebSocket 连接并处理消息
func (c *MarketWSClient) run(ctx context.Context) error {
	handler := &polymarketSDK.WSHandler{
		OnBook: func(e polymarketSDK.WSBookEvent) {
			c.handleBook(e)
		},
		OnPriceChange: func(e polymarketSDK.WSPriceChangeEvent) {
			c.handlePriceChange(e)
		},
		OnLastTradePrice: func(e polymarketSDK.WSLastTradePriceEvent) {
			c.handleLastTradePrice(e)
		},
		OnTickSizeChange: func(e polymarketSDK.WSTickSizeChangeEvent) {
			c.handleTickSizeChange(e)
		},
		// best_bid_ask, new_market, market_resolved 通过 HandleUnknown 兜底解析
		HandleUnknown: func(eventType string, raw json.RawMessage) {
			c.handleUnknown(eventType, raw)
		},
	}

	client, err := polymarketSDK.NewMarketWSClient(c.tokenIDs, true, handler)
	if err != nil {
		return fmt.Errorf("创建市场 WebSocket 客户端失败: %w", err)
	}

	logger.Info("[WS] 市场 WebSocket 已连接", logger.Int("token_count", len(c.tokenIDs)))
	return client.Run(ctx)
}

// Close 关闭 WebSocket 连接
func (c *MarketWSClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
}

// ============================================================
//  SDK 事件 → 内部类型 转换 & 分发
// ============================================================

func (c *MarketWSClient) handleBook(e polymarketSDK.WSBookEvent) {
	if c.callback == nil || c.callback.OnBook == nil {
		return
	}

	snapshot := OrderBookSnapshot{
		AssetID:   e.AssetID,
		Timestamp: time.Now(),
	}

	for _, bid := range e.Bids {
		snapshot.Bids = append(snapshot.Bids, OrderBookLevel{
			Price: bid.Price,
			Size:  bid.Size,
		})
	}
	for _, ask := range e.Asks {
		snapshot.Asks = append(snapshot.Asks, OrderBookLevel{
			Price: ask.Price,
			Size:  ask.Size,
		})
	}

	c.callback.OnBook(snapshot)
}

func (c *MarketWSClient) handlePriceChange(e polymarketSDK.WSPriceChangeEvent) {
	if c.callback == nil || c.callback.OnPriceChange == nil {
		return
	}

	for _, pc := range e.PriceChanges {
		change := PriceChange{
			AssetID: pc.AssetID,
			Price:   pc.Price,
			Size:    pc.Size,
			Side:    pc.Side,
			BestBid: pc.BestBid,
			BestAsk: pc.BestAsk,
		}
		c.callback.OnPriceChange(change)
	}
}

func (c *MarketWSClient) handleLastTradePrice(e polymarketSDK.WSLastTradePriceEvent) {
	if c.callback == nil || c.callback.OnLastTradePrice == nil {
		return
	}

	trade := LastTrade{
		AssetID:   e.AssetID,
		Price:     e.Price,
		Side:      e.Side,
		Size:      e.Size,
		Timestamp: time.Now(),
	}
	c.callback.OnLastTradePrice(trade)
}

func (c *MarketWSClient) handleTickSizeChange(e polymarketSDK.WSTickSizeChangeEvent) {
	if c.callback == nil || c.callback.OnTickSizeChange == nil {
		return
	}

	event := TickSizeChangeEvent{
		AssetID:     e.AssetID,
		Market:      e.Market,
		OldTickSize: e.OldTickSize,
		NewTickSize: e.NewTickSize,
	}
	c.callback.OnTickSizeChange(event)
}

// handleUnknown 解析 HandleUnknown 兜底事件 (best_bid_ask / new_market / market_resolved)
func (c *MarketWSClient) handleUnknown(eventType string, raw json.RawMessage) {
	if c.callback == nil {
		return
	}

	switch eventType {
	case "best_bid_ask":
		if c.callback.OnBestBidAsk == nil {
			return
		}
		var rawEvent struct {
			AssetID string `json:"asset_id"`
			Market  string `json:"market"`
			BestBid string `json:"best_bid"`
			BestAsk string `json:"best_ask"`
			Spread  string `json:"spread"`
		}
		if err := json.Unmarshal(raw, &rawEvent); err != nil {
			logger.Warn("[WS] 解析 best_bid_ask 失败", logger.Err(err))
			return
		}
		c.callback.OnBestBidAsk(BestBidAskEvent{
			AssetID: rawEvent.AssetID,
			Market:  rawEvent.Market,
			BestBid: rawEvent.BestBid,
			BestAsk: rawEvent.BestAsk,
			Spread:  rawEvent.Spread,
		})

	case "new_market":
		if c.callback.OnNewMarket == nil {
			return
		}
		var ev NewMarketEvent
		if err := json.Unmarshal(raw, &ev); err != nil {
			logger.Warn("[WS] 解析 new_market 失败", logger.Err(err))
			return
		}
		c.callback.OnNewMarket(ev)

	case "market_resolved":
		if c.callback.OnMarketResolved == nil {
			return
		}
		var ev MarketResolvedEvent
		if err := json.Unmarshal(raw, &ev); err != nil {
			logger.Warn("[WS] 解析 market_resolved 失败", logger.Err(err))
			return
		}
		c.callback.OnMarketResolved(ev)

	default:
		logger.Debug("[WS] 未识别的事件类型", logger.String("event_type", eventType))
	}
}
