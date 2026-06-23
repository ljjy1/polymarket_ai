package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-dev-frame/sponge/pkg/logger"

	"be/internal/proxy"
)

const (
	defaultBaseURL = "https://api.binance.com"
)

// Kline 表示一条 K 线数据
type Kline struct {
	OpenTime  int64
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	CloseTime int64
}

// Ticker24h 表示 24 小时行情数据
type Ticker24h struct {
	Symbol      string `json:"symbol"`
	LastPrice   string `json:"lastPrice"`
	PriceChange string `json:"priceChange"`
	Volume      string `json:"volume"`
	HighPrice   string `json:"highPrice"`
	LowPrice    string `json:"lowPrice"`
}

// Client 是 Binance API 客户端
type Client struct {
	apiURL     string
	httpClient *http.Client
}

// NewClient 创建一个新的 Binance 客户端
// proxyAddr 为代理地址（如 "127.0.0.1:6450"），为空时不使用代理。
func NewClient(apiURL, proxyAddr string) *Client {
	if apiURL == "" {
		apiURL = defaultBaseURL
	}
	return &Client{
		apiURL:     apiURL,
		httpClient: proxy.NewHTTPClient(proxyAddr),
	}
}

// doRequest 发送 HTTP GET 请求并返回响应体
func (c *Client) doRequest(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("binance request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("binance read body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance request returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// getKlines 获取 K 线数据的通用方法
func (c *Client) getKlines(symbol, interval string, limit int) ([]Kline, error) {
	url := fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", c.apiURL, symbol, interval, limit)

	body, err := c.doRequest(url)
	if err != nil {
		logger.Error("fetch klines failed", logger.String("symbol", symbol), logger.String("interval", interval), logger.Err(err))
		return nil, err
	}

	klines, err := parseKlines(body)
	if err != nil {
		logger.Error("parse klines failed", logger.String("symbol", symbol), logger.String("interval", interval), logger.Err(err))
		return nil, err
	}

	logger.Info("fetched klines successfully",
		logger.String("symbol", symbol),
		logger.String("interval", interval),
		logger.Int("count", len(klines)),
	)
	return klines, nil
}

// parseKlines 解析 Binance K 线 JSON 响应
func parseKlines(data []byte) ([]Kline, error) {
	var rawKlines []interface{}
	if err := json.Unmarshal(data, &rawKlines); err != nil {
		return nil, fmt.Errorf("unmarshal klines failed: %w", err)
	}

	klines := make([]Kline, 0, len(rawKlines))
	for i, raw := range rawKlines {
		arr, ok := raw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid kline format at index %d", i)
		}
		if len(arr) < 7 {
			return nil, fmt.Errorf("kline data too short at index %d: got %d elements", i, len(arr))
		}

		kline, err := parseKlineItem(arr)
		if err != nil {
			return nil, fmt.Errorf("parse kline at index %d: %w", i, err)
		}
		klines = append(klines, kline)
	}

	return klines, nil
}

// parseKlineItem 解析单条 K 线数据
func parseKlineItem(arr []interface{}) (Kline, error) {
	openTime, err := toInt64(arr[0])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid openTime: %w", err)
	}

	open, err := toString(arr[1])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid open: %w", err)
	}

	high, err := toString(arr[2])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid high: %w", err)
	}

	low, err := toString(arr[3])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid low: %w", err)
	}

	close_, err := toString(arr[4])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid close: %w", err)
	}

	volume, err := toString(arr[5])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid volume: %w", err)
	}

	closeTime, err := toInt64(arr[6])
	if err != nil {
		return Kline{}, fmt.Errorf("invalid closeTime: %w", err)
	}

	return Kline{
		OpenTime:  openTime,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close_,
		Volume:    volume,
		CloseTime: closeTime,
	}, nil
}

// toInt64 将 interface{} 转换为 int64
func toInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case float64:
		return int64(val), nil
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			return 0, fmt.Errorf("json.Number to int64: %w", err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unexpected type %T for numeric value", v)
	}
}

// toString 将 interface{} 转换为 string
func toString(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetKlines1h 获取 1 小时 K 线数据
func (c *Client) GetKlines1h(symbol string, limit int) ([]Kline, error) {
	return c.getKlines(symbol, "1h", limit)
}

// GetKlines1d 获取 1 日 K 线数据
func (c *Client) GetKlines1d(symbol string, limit int) ([]Kline, error) {
	return c.getKlines(symbol, "1d", limit)
}

// GetTicker24h 获取 24 小时行情数据
func (c *Client) GetTicker24h(symbol string) (*Ticker24h, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbol=%s", c.apiURL, symbol)

	body, err := c.doRequest(url)
	if err != nil {
		logger.Error("fetch ticker failed", logger.String("symbol", symbol), logger.Err(err))
		return nil, err
	}

	var ticker Ticker24h
	if err := json.Unmarshal(body, &ticker); err != nil {
		logger.Error("parse ticker failed", logger.String("symbol", symbol), logger.Err(err))
		return nil, fmt.Errorf("unmarshal ticker failed: %w", err)
	}

	logger.Info("fetched ticker successfully", logger.String("symbol", symbol))
	return &ticker, nil
}

// GetCurrentPrice 获取当前价格
func (c *Client) GetCurrentPrice(symbol string) (string, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", c.apiURL, symbol)

	body, err := c.doRequest(url)
	if err != nil {
		logger.Error("fetch current price failed", logger.String("symbol", symbol), logger.Err(err))
		return "", err
	}

	var result struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("parse current price failed", logger.String("symbol", symbol), logger.Err(err))
		return "", fmt.Errorf("unmarshal price failed: %w", err)
	}

	return result.Price, nil
}
