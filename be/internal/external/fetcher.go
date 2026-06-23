package external

import (
	"be/internal/proxy"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

// Kline represents a candlestick data point from exchange API.
type Kline struct {
	OpenTime               int64
	Open, High, Low, Close string
	Volume                 string
	CloseTime              int64
}

// TechnicalIndicators holds computed technical analysis values.
type TechnicalIndicators struct {
	RSI14     float64
	MACD      MACDData
	Bollinger BollingerData
	EMA       EMAData
	ATR14     float64
}

// MACDData holds MACD line, signal line and histogram values.
type MACDData struct {
	MACD      float64
	Signal    float64
	Histogram float64
}

// BollingerData holds Bollinger Bands values.
type BollingerData struct {
	Upper  float64
	Middle float64
	Lower  float64
}

// EMAData holds multiple EMA period values.
type EMAData struct {
	EMA7  float64
	EMA25 float64
	EMA99 float64
}

// NewsItem represents a single news article with sentiment.
type NewsItem struct {
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Sentiment   string `json:"sentiment"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

// OnChainData 持有 CryptoQuant 链上指标。
type OnChainData struct {
	ExchangeNetflow24h float64 `json:"exchangeNetflow24h"` // 交易所净流入/流出 24h（负值=流出=看涨）
	ExchangeNetflow7d  float64 `json:"exchangeNetflow7d"`  // 交易所净流入/流出 7d
	ActiveAddresses24h int     `json:"activeAddresses24h"` // 24小时活跃地址数
	MinerOutflow24h    float64 `json:"minerOutflow24h"`    // 矿工流出量(BTC)
	MVRVRatio          float64 `json:"mvrvRatio"`          // MVRV 比率
}

// Fetcher retrieves external market sentiment and news data.
type Fetcher struct {
	client             *http.Client
	gnewsBaseURL       string
	gnewsAPIKey        string
	cryptoquantAPIKey  string
	cryptoquantBaseURL string
	fearGreedBaseURL   string
}

// NewFetcher creates a new Fetcher instance with default settings.
// gnewsAPIKey 是可选的 GNews API Key，为空时 GetNews 会返回 nil。
// cryptoquantAPIKey 和 cryptoquantBaseURL 用于 CryptoQuant 链上数据获取。
// proxyAddr 为代理地址（如 "127.0.0.1:6450"），为空时不使用代理。
func NewFetcher(gnewsBaseURL, gnewsAPIKey, cryptoquantAPIKey, cryptoquantBaseURL, fearGreedBaseURL, proxyAddr string) *Fetcher {
	return &Fetcher{
		client:             proxy.NewHTTPClient(proxyAddr),
		gnewsBaseURL:       gnewsBaseURL,
		gnewsAPIKey:        gnewsAPIKey,
		cryptoquantAPIKey:  cryptoquantAPIKey,
		cryptoquantBaseURL: cryptoquantBaseURL,
		fearGreedBaseURL:   fearGreedBaseURL,
	}
}

// GetFearGreedIndex fetches the Fear and Greed Index from alternative.me API.
// Returns the index value (0-100), classification label (e.g. "Fear", "Greed"), and any error.
func (f *Fetcher) GetFearGreedIndex() (int, string, error) {
	fgi, fgl, _, err := f.GetFearGreedWithTrend(1)
	return fgi, fgl, err
}

// GetFearGreedWithTrend 获取恐惧贪婪指数及多日趋势数组（Python 参考版对应字段）。
// days: 返回的天数（如 7 表示返回最近 7 天的指数数组）。
// 返回: 最新指数值, 分类标签, 趋势数组(从旧到新), error。
func (f *Fetcher) GetFearGreedWithTrend(days int) (int, string, []int, error) {
	// fearGreedBaseURL 已包含完整路径 (如 "https://api.alternative.me/fng")，只追加查询参数
	url := strings.TrimRight(f.fearGreedBaseURL, "/") + fmt.Sprintf("?limit=%d", days)
	resp, err := f.client.Get(url)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to fetch fear greed index: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Data []struct {
			Value          string `json:"value"`
			Classification string `json:"value_classification"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, "", nil, fmt.Errorf("failed to parse fear greed index response: %w", err)
	}

	if len(result.Data) == 0 {
		return 0, "", nil, fmt.Errorf("empty fear greed index data")
	}

	// 解析趋势数组（从旧到新，与 Python 参考版一致）
	trend := make([]int, 0, len(result.Data))
	for i := len(result.Data) - 1; i >= 0; i-- {
		var v int
		if _, err := fmt.Sscanf(result.Data[i].Value, "%d", &v); err == nil {
			trend = append(trend, v)
		}
	}

	var latestValue int
	if _, err := fmt.Sscanf(result.Data[0].Value, "%d", &latestValue); err != nil {
		return 0, "", nil, fmt.Errorf("failed to parse fear greed value: %w", err)
	}

	return latestValue, result.Data[0].Classification, trend, nil
}

// GetNews 获取 BTC 相关新闻。
// 当配置了 GNews API Key 时调用真实 API，否则返回 nil。
// count 指定最大新闻数量（默认 5，最大 20）。
func (f *Fetcher) GetNews(symbol string, count int) []NewsItem {
	if count <= 0 {
		count = 5
	}
	if count > 20 {
		count = 20
	}

	// 如果有 API Key，调用 GNews 真实接口
	if f.gnewsAPIKey != "" {
		items, err := f.getGNews(symbol, count)
		if err == nil {
			return items
		}
	}
	return nil
}

// getGNews 调用 GNews API 获取实时新闻。
// API 文档: https://gnews.io/docs
func (f *Fetcher) getGNews(symbol string, count int) ([]NewsItem, error) {
	// GNews 要求 q 参数不能直接传 "BTC"，需要加上相关关键词
	query := symbol
	url := fmt.Sprintf("%s?q=%s&lang=en&max=%d&apikey=%s", f.gnewsBaseURL, query, count, f.gnewsAPIKey)

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gnews request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gnews read body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gnews api returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Content     string `json:"content"`
			URL         string `json:"url"`
			Source      struct {
				Name string `json:"name"`
			} `json:"source"`
			PublishedAt string `json:"publishedAt"`
		} `json:"articles"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("gnews parse response failed: %w", err)
	}

	items := make([]NewsItem, 0, len(result.Articles))
	for _, r := range result.Articles {
		// GNews 不提供情感分析，统一设为 neutral
		sentiment := "neutral"

		// 优先使用 description 作为摘要，若为空则使用 title
		summary := r.Description
		if summary == "" {
			summary = r.Title
		}

		// 格式化时间
		pubTime := r.PublishedAt
		if t, err := time.Parse(time.RFC3339, r.PublishedAt); err == nil {
			pubTime = t.Format(time.RFC3339)
		}

		items = append(items, NewsItem{
			Title:       r.Title,
			Summary:     summary,
			Sentiment:   sentiment,
			Source:      r.Source.Name,
			URL:         r.URL,
			PublishedAt: pubTime,
		})
	}

	return items, nil
}

// GetOnChainData 从 CryptoQuant API 获取链上数据。
// 使用 Bearer Token 认证（API Key 作为 Bearer Token）。
// 未配置 API Key 时返回零值，不报错。
// CryptoQuant API 文档: https://docs.cryptoquant.com/api-reference/available-endpoints/endpoints
func (f *Fetcher) GetOnChainData() *OnChainData {
	data := &OnChainData{}

	if f.cryptoquantAPIKey == "" {
		return data // 返回零值，后续使用方自行处理
	}

	base := strings.TrimRight(f.cryptoquantBaseURL, "/")

	type result struct {
		key string
		val float64
		err error
	}
	ch := make(chan result, 5)

	// 1. 交易所净流量（24h）
	// GET /v1/btc/exchange-flows/netflow?exchange=all_exchange&window=day&limit=1
	go func() {
		v, err := f.fetchCryptoQuantMetric(base, "/v1/btc/exchange-flows/netflow?exchange=all_exchange&window=day&limit=1", "netflow_total")
		ch <- result{"exchangeNetflow24h", v, err}
	}()

	// 2. 交易所净流量（7d）
	// 获取 7 天前的数据作为对比基准
	go func() {
		sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("20060102")
		v, err := f.fetchCryptoQuantMetric(base, fmt.Sprintf("/v1/btc/exchange-flows/netflow?exchange=all_exchange&window=day&from=%s&limit=1", sevenDaysAgo), "netflow_total")
		ch <- result{"exchangeNetflow7d", v, err}
	}()

	// 3. 活跃地址数
	// GET /v1/btc/network-data/addresses-count?window=day&limit=1
	go func() {
		v, err := f.fetchCryptoQuantMetric(base, "/v1/btc/network-data/addresses-count?window=day&limit=1", "addresses_count_active")
		ch <- result{"activeAddresses", v, err}
	}()

	// 4. 矿工流出量
	// GET /v1/btc/miner-flows/outflow?miner=all_miner&window=day&limit=1
	go func() {
		v, err := f.fetchCryptoQuantMetric(base, "/v1/btc/miner-flows/outflow?miner=all_miner&window=day&limit=1", "outflow_total")
		ch <- result{"minerOutflow", v, err}
	}()

	// 5. MVRV 比率
	// GET /v1/btc/market-indicator/mvrv?window=day&limit=1
	go func() {
		v, err := f.fetchCryptoQuantMetric(base, "/v1/btc/market-indicator/mvrv?window=day&limit=1", "mvrv")
		ch <- result{"mvrv", v, err}
	}()

	for i := 0; i < 5; i++ {
		r := <-ch
		if r.err != nil {
			continue
		}
		switch r.key {
		case "exchangeNetflow24h":
			data.ExchangeNetflow24h = r.val
		case "exchangeNetflow7d":
			data.ExchangeNetflow7d = r.val
		case "activeAddresses":
			data.ActiveAddresses24h = int(r.val)
		case "minerOutflow":
			data.MinerOutflow24h = r.val
		case "mvrv":
			data.MVRVRatio = r.val
		}
	}

	return data
}

// fetchCryptoQuantMetric 调用 CryptoQuant API 获取指定指标的最新值。
// 使用 Bearer Token 认证，解析 JSON 响应中指定字段的值。
func (f *Fetcher) fetchCryptoQuantMetric(baseURL, endpoint, fieldName string) (float64, error) {
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimLeft(endpoint, "/"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("cryptoquant request creation failed for %s: %w", endpoint, err)
	}
	req.Header.Set("Authorization", "Bearer "+f.cryptoquantAPIKey)

	resp, err := f.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("cryptoquant request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("cryptoquant read body failed for %s: %w", endpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("cryptoquant api returned status %d for %s: %s", resp.StatusCode, endpoint, string(body))
	}

	// CryptoQuant 响应格式: {"status":{...}, "result": {"data": [{"date": "...", fieldName: value}]}}
	var cqResp struct {
		Status struct {
			Code int `json:"code"`
		} `json:"status"`
		Result struct {
			Data []map[string]interface{} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &cqResp); err != nil {
		return 0, fmt.Errorf("cryptoquant parse response failed for %s: %w", endpoint, err)
	}

	if cqResp.Status.Code != 200 {
		return 0, fmt.Errorf("cryptoquant api error for %s: status code %d", endpoint, cqResp.Status.Code)
	}

	if len(cqResp.Result.Data) == 0 {
		return 0, fmt.Errorf("cryptoquant empty response for %s", endpoint)
	}

	// 从最新的数据点中提取目标字段
	latest := cqResp.Result.Data[0]
	val, ok := latest[fieldName]
	if !ok || val == nil {
		return 0, fmt.Errorf("cryptoquant field %s not found in response for %s", fieldName, endpoint)
	}

	// 处理数值类型（可能是 float64 或 int）
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cryptoquant field %s has unexpected type %T for %s", fieldName, val, endpoint)
	}
}

// GetTechnicalIndicators computes RSI, MACD, Bollinger Bands, EMA and ATR
// from the provided kline data. Uses daily klines (klines1d) for calculation;
// klines1h is reserved for potential short-term indicator computation.
// Requires at least 100 data points for EMA99.
func (f *Fetcher) GetTechnicalIndicators(klines1h, klines1d []Kline) (*TechnicalIndicators, error) {
	if len(klines1d) < 100 {
		return nil, fmt.Errorf("insufficient kline data: need at least 100, got %d", len(klines1d))
	}

	// Parse OHLC values from string to float64
	closePrices := make([]float64, len(klines1d))
	highPrices := make([]float64, len(klines1d))
	lowPrices := make([]float64, len(klines1d))
	for i, k := range klines1d {
		c, _ := parseFloat(k.Close)
		h, _ := parseFloat(k.High)
		l, _ := parseFloat(k.Low)
		closePrices[i] = c
		highPrices[i] = h
		lowPrices[i] = l
	}

	// Suppress unused parameter warning
	_ = klines1h

	ti := &TechnicalIndicators{}

	// RSI(14)
	ti.RSI14 = calculateRSI(closePrices, 14)

	// MACD(12, 26, 9)
	ti.MACD = calculateMACD(closePrices, 12, 26, 9)

	// Bollinger Bands(20, 2)
	ti.Bollinger = calculateBollinger(closePrices, 20, 2)

	// EMA(7, 25, 99)
	ti.EMA.EMA7 = calculateEMA(closePrices, 7)
	ti.EMA.EMA25 = calculateEMA(closePrices, 25)
	ti.EMA.EMA99 = calculateEMA(closePrices, 99)

	// ATR(14)
	ti.ATR14 = calculateATR(highPrices, lowPrices, closePrices, 14)

	return ti, nil
}

// GetTechnicalIndicators1h 使用 1h K 线计算技术指标（与 Python 参考版一致）。
// Python 参考版使用 KLINE_LIMIT_1H = 24*7 = 168 根 1h K 线。
// 需要至少 100 个数据点来确保 EMA99 计算的可靠性。
func (f *Fetcher) GetTechnicalIndicators1h(klines1h []Kline) (*TechnicalIndicators, error) {
	if len(klines1h) < 100 {
		return nil, fmt.Errorf("insufficient 1h kline data: need at least 100, got %d", len(klines1h))
	}

	closePrices := make([]float64, len(klines1h))
	highPrices := make([]float64, len(klines1h))
	lowPrices := make([]float64, len(klines1h))
	for i, k := range klines1h {
		c, _ := parseFloat(k.Close)
		h, _ := parseFloat(k.High)
		l, _ := parseFloat(k.Low)
		closePrices[i] = c
		highPrices[i] = h
		lowPrices[i] = l
	}

	ti := &TechnicalIndicators{}

	// RSI(14)
	ti.RSI14 = calculateRSI(closePrices, 14)

	// MACD(12, 26, 9)
	ti.MACD = calculateMACD(closePrices, 12, 26, 9)

	// Bollinger Bands(20, 2)
	ti.Bollinger = calculateBollinger(closePrices, 20, 2)

	// EMA(7, 25, 99)
	ti.EMA.EMA7 = calculateEMA(closePrices, 7)
	ti.EMA.EMA25 = calculateEMA(closePrices, 25)
	ti.EMA.EMA99 = calculateEMA(closePrices, 99)

	// ATR(14)
	ti.ATR14 = calculateATR(highPrices, lowPrices, closePrices, 14)

	return ti, nil
}

// ---------------------------------------------------------------------------
// Technical indicator calculation helpers
// ---------------------------------------------------------------------------

// parseFloat converts a string to float64.
func parseFloat(s string) (float64, error) {
	var v float64
	_, err := fmt.Sscanf(s, "%f", &v)
	return v, err
}

// calculateRSI computes the Relative Strength Index (RSI) for the given period.
func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // neutral default
	}

	tail := prices[len(prices)-period-1:]

	var avgGain, avgLoss float64
	for i := 1; i < len(tail); i++ {
		diff := tail[i] - tail[i-1]
		if diff > 0 {
			avgGain += diff
		} else {
			avgLoss -= diff
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - 100/(1+rs)
}

// calculateEMA computes the Exponential Moving Average for the given period.
// Uses SMA as the initial seed value, then applies the EMA formula recursively.
func calculateEMA(prices []float64, period int) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}
	if period > n {
		period = n
	}

	// Seed: SMA of the first 'period' values
	sum := 0.0
	start := n - period
	for i := start; i < n; i++ {
		sum += prices[i]
	}
	ema := sum / float64(period)

	multiplier := 2.0 / (float64(period) + 1)
	for i := start + 1; i < n; i++ {
		ema = (prices[i]-ema)*multiplier + ema
	}
	return ema
}

// calculateMACD computes MACD line, signal line, and histogram.
func calculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) MACDData {
	n := len(prices)
	if n < slowPeriod+signalPeriod {
		return MACDData{}
	}

	// Build a series of MACD line values
	macdValues := make([]float64, 0, n-slowPeriod+1)
	for i := slowPeriod; i <= n; i++ {
		segment := prices[:i]
		emaFast := calculateEMA(segment, fastPeriod)
		emaSlow := calculateEMA(segment, slowPeriod)
		macdValues = append(macdValues, emaFast-emaSlow)
	}

	if len(macdValues) == 0 {
		return MACDData{}
	}

	macdLine := macdValues[len(macdValues)-1]
	signalLine := calculateEMA(macdValues, signalPeriod)

	return MACDData{
		MACD:      macdLine,
		Signal:    signalLine,
		Histogram: macdLine - signalLine,
	}
}

// calculateBollinger computes Bollinger Bands (middle SMA, upper/lower bands).
func calculateBollinger(prices []float64, period int, numStdDev float64) BollingerData {
	n := len(prices)
	if period > n {
		period = n
	}

	tail := prices[n-period:]

	// SMA (middle band)
	var sum float64
	for _, p := range tail {
		sum += p
	}
	middle := sum / float64(period)

	// Standard deviation
	var variance float64
	for _, p := range tail {
		diff := p - middle
		variance += diff * diff
	}
	variance /= float64(period)
	sd := math.Sqrt(variance)

	return BollingerData{
		Upper:  middle + numStdDev*sd,
		Middle: middle,
		Lower:  middle - numStdDev*sd,
	}
}

// calculateATR computes the Average True Range for the given period.
func calculateATR(high, low, close []float64, period int) float64 {
	n := len(high)
	if n < period+1 {
		return 0
	}

	// Calculate True Range for each bar (except the first)
	trs := make([]float64, n-1)
	for i := 1; i < n; i++ {
		h := high[i]
		l := low[i]
		pc := close[i-1]
		tr := math.Max(h-l, math.Max(math.Abs(h-pc), math.Abs(l-pc)))
		trs[i-1] = tr
	}

	// ATR: simple mean of the last 'period' TRs
	tail := trs[len(trs)-period:]
	var sum float64
	for _, tr := range tail {
		sum += tr
	}
	return sum / float64(period)
}
