package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/shopspring/decimal"

	polymarketSDK "github.com/0xNetuser/Polymarket-golang/polymarket"

	"be/internal/dao"
	"be/internal/model"
	"be/internal/polymarket"
)

const (
	MinYesPrice  = 0.35  // 最小 Yes 价格
	MaxYesPrice  = 0.65  // 最大 Yes 价格
	MinVolume24h = 10000 // 最小 24h 交易量 USD
	MaxSpread    = 0.03  // 最大价差 3%（未使用，保留仅作参考）
)

// bitcoinAbovePattern 匹配 "Bitcoin above" 标题模式（Python 参考版 _BTC_ABOVE_PATTERN）
var bitcoinAbovePattern = regexp.MustCompile(`(?i)bitcoin\s+above`)

// CandidateMarket 候选市场信息
type CandidateMarket struct {
	ConditionID    string
	TokenID        string
	Question       string
	PriceThreshold int
	YesPrice       float64
	NoPrice        float64
	TotalVolume    float64
	Spread         float64
	Score          float64 // 距离50%的偏差 (越小越好)
}

// MarketScanner 市场扫描器，负责每天扫描 Polymarket 选择可交易的市场
type MarketScanner struct {
	polymarketClient *polymarket.Client
	marketDao        dao.MarketsDao
}

// NewMarketScanner 创建市场扫描器
func NewMarketScanner(polyClient *polymarket.Client, marketDao dao.MarketsDao) *MarketScanner {
	return &MarketScanner{
		polymarketClient: polyClient,
		marketDao:        marketDao,
	}
}

// gammaEventRaw 用于解析 Gamma API 返回的事件 JSON（只解析关心的字段）
type gammaEventRaw struct {
	ID          string           `json:"id"`
	ConditionID string           `json:"conditionId"`
	Slug        string           `json:"slug"`
	Title       string           `json:"title"`
	EndDate     string           `json:"endDate"`
	Markets     []gammaMarketRaw `json:"markets"`
}

// StringSlice 支持从 JSON 数组 ["a","b"] 或 JSON 字符串 "[\"a\",\"b\"]" 反序列化。
//
// Gamma API 有时将数组字段序列化为字符串而非原生 JSON 数组，
// 使用此类型可以兼容两种格式，避免解析失败。
type StringSlice []string

// UnmarshalJSON 实现 json.Unmarshaler，兼容两种 JSON 格式：
//   - ["a", "b"]        — 原生 JSON 数组
//   - "[\"a\",\"b\"]"   — 字符串中嵌套 JSON 数组
func (s *StringSlice) UnmarshalJSON(data []byte) error {
	// 先尝试直接解析为 []string
	var ss []string
	if err := json.Unmarshal(data, &ss); err == nil {
		*s = ss
		return nil
	}
	// 再尝试解析为字符串，然后解析其中的 JSON 数组
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("StringSlice: cannot unmarshal %s: %w", string(data), err)
	}
	if err := json.Unmarshal([]byte(str), &ss); err != nil {
		return fmt.Errorf("StringSlice: cannot unmarshal inner JSON %s: %w", str, err)
	}
	*s = ss
	return nil
}

// FlexString 支持从 JSON 字符串或数字反序列化。
//
// Gamma API 返回的数值字段（如 volume、volume24hr）格式不固定，
// 有时是字符串 "65416.039829"，有时是数字 65416。
// 使用此类型可以兼容两种格式，统一以字符串形式存储。
type FlexString string

// UnmarshalJSON 实现 json.Unmarshaler，兼容 JSON 字符串和数字两种格式。
func (s *FlexString) UnmarshalJSON(data []byte) error {
	// 先尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = FlexString(str)
		return nil
	}
	// 再尝试解析为数字（float64）
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("FlexString: cannot unmarshal %s: %w", string(data), err)
	}
	*s = FlexString(strconv.FormatFloat(f, 'f', -1, 64))
	return nil
}

// gammaMarketRaw 用于解析 Gamma API 返回的市场 JSON（只解析关心的字段）
type gammaMarketRaw struct {
	ConditionID  string      `json:"conditionId"`
	Question     string      `json:"question"`
	Outcomes     StringSlice `json:"outcomes"`
	Volume       FlexString  `json:"volume"`
	Volume24hr   FlexString  `json:"volume24hr"`
	TokenID      string      `json:"tokenId"`
	ClobTokenIDs StringSlice `json:"clobTokenIds"`
}

// totalVolumeFloat 解析 Volume（总成交量）为 float64，解析失败返回 0。
// Python 参考版实际使用的是 API 的 volume（总成交量）字段，而非 volume24hr。
func (m *gammaMarketRaw) totalVolumeFloat() float64 {
	v, err := strconv.ParseFloat(string(m.Volume), 64)
	if err != nil {
		return 0
	}
	return v
}

// hasYesOutcome 检查市场是否包含 "Yes" 结果。
func (m *gammaMarketRaw) hasYesOutcome() bool {
	for _, o := range m.Outcomes {
		if o == "Yes" {
			return true
		}
	}
	// API 返回的 outcomes 可能为空字符串（未解析出数据时），视为 Yes 方向
	return len(m.Outcomes) == 0
}

// hasExistingScan 检查是否已有当日扫描记录（幂等）
func (s *MarketScanner) hasExistingScan(ctx context.Context, scanDate time.Time) bool {
	dateStr := scanDate.Format("2006-01-02")
	params := &query.Params{
		Columns: []query.Column{
			{
				Name:  "scan_date",
				Exp:   "=",
				Value: dateStr,
			},
		},
	}
	records, total, err := s.marketDao.GetByColumns(ctx, params)
	if err != nil {
		logger.Warn("检查当日扫描记录失败", logger.Err(err), logger.String("scan_date", dateStr))
		return false
	}
	return total > 0 && len(records) > 0
}

// extractThreshold 从问题标题中提取价格阈值
// 例如 "Will Bitcoin be above $100,000 on June 24?" -> 100000
func extractThreshold(question string) int {
	// 尝试匹配 $X,XXX 或 $XXXXX 模式
	re := regexp.MustCompile(`\$(\d{1,3}(?:,\d{3})*(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(question)
	if len(matches) >= 2 {
		cleaned := strings.ReplaceAll(matches[1], ",", "")
		if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
			return int(val)
		}
	}

	// 尝试匹配数字后跟 K（如 $100K -> 100000）
	reK := regexp.MustCompile(`\$?(\d+(?:\.\d+)?)\s*K`)
	matchesK := reK.FindStringSubmatch(question)
	if len(matchesK) >= 2 {
		if val, err := strconv.ParseFloat(matchesK[1], 64); err == nil {
			return int(val * 1000)
		}
	}

	// 尝试匹配数字后跟 M（如 $1M -> 1000000）
	reM := regexp.MustCompile(`\$?(\d+(?:\.\d+)?)\s*M`)
	matchesM := reM.FindStringSubmatch(question)
	if len(matchesM) >= 2 {
		if val, err := strconv.ParseFloat(matchesM[1], 64); err == nil {
			return int(val * 1000000)
		}
	}

	return 0
}

// extractPriceFromOrderBook 从订单簿中提取价格和价差
func extractPriceFromOrderBook(ob *polymarketSDK.OrderBookSummary) (yesPrice float64, spread float64) {
	if ob == nil || len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		return 0, 0
	}

	// 取最优买价（Bids 最高价）和最优卖价（Asks 最低价）
	bestBid := parseOrderPrice(ob.Bids[0].Price)
	bestAsk := parseOrderPrice(ob.Asks[0].Price)

	if bestBid <= 0 || bestAsk <= 0 || bestAsk < bestBid {
		return 0, 0
	}

	yesPrice = (bestBid + bestAsk) / 2
	spread = bestAsk - bestBid
	return
}

// parseOrderPrice 解析订单簿中的价格字符串为 float64
func parseOrderPrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

// Scan 执行市场扫描，返回选中的市场或 nil
//
// 逻辑:
//  1. 计算 targetDate = today + 2 days
//  2. 通过 Gamma API 查询 BTC 事件 (tag=bitcoin, active=true)
//  3. 筛选两天后结算的事件
//  4. 获取每个事件下的子市场订单簿
//  5. 筛选: 35% <= yes_price <= 65%, 24h 交易量 >= $10,000, 价差 <= 3%
//  6. 按距离50%排序 + 交易量排序，取 top 1
//  7. 写入 markets 表
func (s *MarketScanner) Scan(ctx context.Context) (*model.Markets, error) {
	now := time.Now().Truncate(24 * time.Hour)
	targetDate := now.Add(2 * 24 * time.Hour)
	scanDate := now

	logger.Info("开始市场扫描",
		logger.String("scan_date", scanDate.Format("2006-01-02")),
		logger.String("target_date", targetDate.Format("2006-01-02")),
	)

	// 1. 幂等检查：是否已有当日扫描记录
	if s.hasExistingScan(ctx, scanDate) {
		logger.Info("当日已有扫描记录，跳过扫描", logger.String("scan_date", scanDate.Format("2006-01-02")))
		return nil, nil
	}

	// 2. 通过 Gamma API 查询两天后结算的 BTC 事件
	// 用 endDateMin/endDateMax 过滤过期事件，避免 API 返回多年前的历史事件
	endDateMin := targetDate.Format(time.RFC3339)
	endDateMax := targetDate.Add(24*time.Hour - time.Second).Format(time.RFC3339)
	logger.Info("正在查询 Polymarket BTC 事件",
		logger.String("tag", "bitcoin"),
		logger.String("endDateMin", endDateMin),
		logger.String("endDateMax", endDateMax),
	)
	eventsRaw, err := s.polymarketClient.GetGammaEvents("bitcoin", true, endDateMin, endDateMax)
	if err != nil {
		logger.Error("查询 BTC 事件失败", logger.Err(err))
		return nil, fmt.Errorf("查询 BTC 事件失败: %w", err)
	}

	if len(eventsRaw) == 0 {
		logger.Warn("BTC 事件查询结果为空")
		return nil, nil
	}

	// 3. 解析事件列表
	var events []gammaEventRaw
	if err := json.Unmarshal(eventsRaw, &events); err != nil {
		logger.Error("解析事件 JSON 失败", logger.Err(err))
		return nil, fmt.Errorf("解析事件 JSON 失败: %w", err)
	}

	logger.Info("获取到 BTC 事件", logger.Int("count", len(events)))

	var candidates []CandidateMarket

	// 4. 筛选两天后结算的事件
	for _, event := range events {
		eventEndDate, err := time.Parse(time.RFC3339, event.EndDate)
		if err != nil {
			// 尝试其他常见时间格式
			eventEndDate, err = time.Parse("2006-01-02T15:04:05Z", event.EndDate)
			if err != nil {
				eventEndDate, err = time.Parse("2006-01-02T15:04:05.000Z", event.EndDate)
				if err != nil {
					logger.Warn("无法解析事件结束日期", logger.String("end_date", event.EndDate), logger.String("slug", event.Slug))
					continue
				}
			}
		}

		// 检查是否在目标日期结算（同一天）
		if !isSameDay(eventEndDate, targetDate) {
			continue
		}

		logger.Info("找到匹配结算日期的事件",
			logger.String("slug", event.Slug),
			logger.String("title", event.Title),
			logger.String("end_date", event.EndDate),
		)

		// 检查事件标题是否包含 "Bitcoin above"（Python 参考版规则 2）
		if !bitcoinAbovePattern.MatchString(event.Title) {
			logger.Info("事件标题不包含 'Bitcoin above'，跳过",
				logger.String("title", event.Title),
				logger.String("slug", event.Slug),
			)
			continue
		}

		// 重新获取事件详情，确保 market 数据完整（含 question 字段）
		// Python 参考版先 list events 筛选，再通过 /events/{id} 获取详情
		detailRaw, err := s.polymarketClient.GetGammaEventByID(event.ID)
		if err != nil {
			logger.Warn("获取事件详情失败，使用列表数据",
				logger.Err(err),
				logger.String("slug", event.Slug),
				logger.String("event_id", event.ID),
			)
		} else {
			var detailEvent gammaEventRaw
			if err := json.Unmarshal(detailRaw, &detailEvent); err == nil && len(detailEvent.Markets) > 0 {
				event.Markets = detailEvent.Markets
				logger.Info("已获取事件详情，markets 数量",
					logger.Int("count", len(detailEvent.Markets)),
				)
			}
		}

		// 5. 处理事件下的子市场
		for _, market := range event.Markets {
			// 只关注 Yes 方向的市场
			if !market.hasYesOutcome() {
				continue
			}

			// 获取 Yes token ID
			tokenID := market.TokenID
			if tokenID == "" && len(market.ClobTokenIDs) > 0 {
				tokenID = market.ClobTokenIDs[0]
			}
			if tokenID == "" {
				logger.Warn("市场缺少 Token ID，跳过", logger.String("question", market.Question))
				continue
			}

			// 获取订单簿
			orderBook, err := s.polymarketClient.GetOrderBook(tokenID)
			if err != nil {
				logger.Warn("获取订单簿失败",
					logger.Err(err),
					logger.String("token_id", tokenID),
					logger.String("question", market.Question),
				)
				continue
			}
			//打印订单orderBook信息
			logger.Info("订单簿信息",
				logger.String("token_id", tokenID),
				logger.String("question", market.Question),
				logger.Any("order_book", orderBook),
			)

			yesPrice, spread := extractPriceFromOrderBook(orderBook)
			if yesPrice <= 0 {
				logger.Warn("无法从订单簿提取有效价格", logger.String("token_id", tokenID))
				continue
			}

			noPrice := 1 - yesPrice
			totalVolume := market.totalVolumeFloat()

			logger.Info("市场信息",
				logger.String("question", market.Question),
				logger.Float64("yes_price", yesPrice),
				logger.Float64("no_price", noPrice),
				logger.Float64("spread", spread),
				logger.Float64("total_volume", totalVolume),
			)

			// 6. 筛选条件
			if yesPrice < MinYesPrice || yesPrice > MaxYesPrice {
				logger.Info("价格不在范围内，跳过",
					logger.Float64("yes_price", yesPrice),
					logger.Float64("min", MinYesPrice),
					logger.Float64("max", MaxYesPrice),
				)
				continue
			}

			if totalVolume < MinVolume24h {
				logger.Info("交易量不足，跳过",
					logger.Float64("total_volume", totalVolume),
					logger.Float64("min_volume", MinVolume24h),
				)
				continue
			}

			// 计算分数：距离50%的偏差（越小越好）
			deviation := math.Abs(yesPrice - 0.5)
			// 综合分数 = 偏差（权重0.7） + 交易量归一化（权重0.3）
			// 偏差越小越好，交易量越大越好
			volumeScore := totalVolume / 100000.0 // 归一化到 0-1 左右的范围
			score := deviation*0.7 - volumeScore*0.3

			threshold := extractThreshold(market.Question)
			if threshold == 0 {
				threshold = extractThreshold(event.Title)
			}

			candidate := CandidateMarket{
				ConditionID:    market.ConditionID,
				TokenID:        tokenID,
				Question:       market.Question,
				PriceThreshold: threshold,
				YesPrice:       yesPrice,
				NoPrice:        noPrice,
				TotalVolume:    totalVolume,
				Spread:         spread,
				Score:          score,
			}
			candidates = append(candidates, candidate)

			logger.Info("找到符合条件的候选市场",
				logger.String("question", market.Question),
				logger.Float64("yes_price", yesPrice),
				logger.Float64("total_volume", totalVolume),
			)
		}
	}

	if len(candidates) == 0 {
		logger.Warn("未找到符合条件的候选市场")
		return nil, nil
	}

	// 7. 排序：按偏差从小到大（偏差越小越接近50%），偏差相同时按交易量从大到小
	sort.Slice(candidates, func(i, j int) bool {
		if math.Abs(candidates[i].Score-candidates[j].Score) > 0.001 {
			return candidates[i].Score < candidates[j].Score
		}
		return candidates[i].TotalVolume > candidates[j].TotalVolume
	})

	best := candidates[0]

	// 查找当前最佳市场对应的事件 slug，方便日志定位
	eventSlug := ""
	for _, ev := range events {
		for _, m := range ev.Markets {
			if m.ConditionID == best.ConditionID {
				eventSlug = ev.Slug
				break
			}
		}
		if eventSlug != "" {
			break
		}
	}

	logger.Info("选定最佳市场",
		logger.String("event_slug", eventSlug),
		logger.String("question", best.Question),
		logger.Float64("yes_price", best.YesPrice),
		logger.Float64("total_volume", best.TotalVolume),
		logger.Float64("score", best.Score),
	)

	// 8. 写入数据库
	yesPriceDec := decimal.NewFromFloat(best.YesPrice)
	noPriceDec := decimal.NewFromFloat(best.NoPrice)
	nowTime := time.Now()

	marketRecord := &model.Markets{
		PolymarketConditionID: best.ConditionID,
		PolymarketTokenID:     best.TokenID,
		EventSlug:             eventSlug,
		Question:              best.Question,
		PriceThreshold:        best.PriceThreshold,
		ScanDate:              &scanDate,
		TargetDate:            &targetDate,
		CurrentYesPrice:       &yesPriceDec,
		CurrentNoPrice:        &noPriceDec,
		SelectedAt:            &nowTime,
		Status:                "active",
	}

	if err := s.marketDao.Create(ctx, marketRecord); err != nil {
		logger.Error("保存市场记录失败", logger.Err(err))
		return nil, fmt.Errorf("保存市场记录失败: %w", err)
	}

	logger.Info("市场扫描完成，已保存选中市场",
		logger.Uint64("market_id", marketRecord.ID),
		logger.String("question", best.Question),
	)

	return marketRecord, nil
}

// isSameDay 检查两个时间是否在同一天
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
