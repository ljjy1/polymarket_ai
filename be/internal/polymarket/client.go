package polymarket

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	polymarketSDK "github.com/0xNetuser/Polymarket-golang/polymarket"

	"be/internal/proxy"
)

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

// Client 封装了 Polymarket API 调用，是对 github.com/0xNetuser/Polymarket-golang 的轻量封装。
type Client struct {
	ClobClient *polymarketSDK.ClobClient
	GammaAPI   *polymarketSDK.GammaAPIClient
	DataAPI    *polymarketSDK.DataAPIClient
}

// NewClient 创建 Polymarket 客户端。
//   - privateKey 为空时使用 L0（只读）模式
//   - privateKey 非空时使用 L2 模式（需要同时提供 apiKey、apiSecret、passphrase）
//   - proxyAddr 为代理地址（如 "127.0.0.1:6450"），为空时不使用代理。
//
// 注意：proxyAddr 仅在创建 Client 时生效，且会全局设置 HTTPS_PROXY 环境变量，
// 这会影响进程中所有通过 http.ProxyFromEnvironment 发起的 HTTPS 请求。
func NewClient(clobURL, gammaURL string, chainID int, privateKey, apiKey, apiSecret, passphrase, proxyAddr string) (*Client, error) {
	var creds *polymarketSDK.ApiCreds
	privKey := normalizePrivateKey(privateKey)
	if privKey != "" {
		creds = &polymarketSDK.ApiCreds{
			APIKey:        apiKey,
			APISecret:     apiSecret,
			APIPassphrase: passphrase,
		}
	}

	// ClobClient 内部使用 defaultHTTPTransport，其 Proxy 设置为 http.ProxyFromEnvironment。
	// 设置 HTTPS_PROXY 环境变量，使 ClobClient 的所有请求通过代理。
	if proxyAddr != "" {
		os.Setenv("HTTPS_PROXY", "http://"+proxyAddr)
	}

	clobClient, err := polymarketSDK.NewClobClient(clobURL, chainID, privKey, creds, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create ClobClient: %w", err)
	}

	gammaClient := polymarketSDK.NewGammaAPIClient()
	if gammaURL != "" {
		gammaClient.WithBaseURL(gammaURL)
	}
	// Gamma API 支持注入自定义 HTTP 客户端
	if proxyAddr != "" {
		gammaClient.WithHTTPClient(proxy.NewHTTPClient(proxyAddr))
	}

	dataAPIClient := polymarketSDK.NewDataAPIClient()
	// Data API 支持注入自定义 HTTP 客户端
	if proxyAddr != "" {
		dataAPIClient.WithHTTPClient(proxy.NewHTTPClient(proxyAddr))
	}

	return &Client{
		ClobClient: clobClient,
		GammaAPI:   gammaClient,
		DataAPI:    dataAPIClient,
	}, nil
}

// normalizePrivateKey 规范化私钥：去除 0x 前缀、空格、换行符，并校验 hex 格式。
// 返回规范化后的私钥，如果 key 非空但无效则返回空字符串（调用方视为只读模式）。
func normalizePrivateKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, "0x")
	key = strings.TrimPrefix(key, "0X")

	if key == "" {
		return ""
	}

	if !hexRegex.MatchString(key) {
		// 打印前 4 和后 4 字符，帮助用户排查
		masked := key
		if len(key) > 8 {
			masked = key[:4] + "..." + key[len(key)-4:]
		}
		fmt.Printf("[polymarket] WARNING: private key contains non-hex characters (got=%q len=%d), falling back to read-only mode\n", masked, len(key))
		return ""
	}

	return key
}

// GetGammaEvents 通过 Gamma API 查询事件列表。
func (c *Client) GetGammaEvents(tags string, active bool, endDateMin, endDateMax string) (json.RawMessage, error) {
	query := &polymarketSDK.GammaEventsQuery{
		TagSlug:    tags,
		Active:     &active,
		EndDateMin: endDateMin,
		EndDateMax: endDateMax,
	}
	return c.GammaAPI.ListEvents(query)
}

// GetGammaEventByID 获取单个事件详情（含完整市场信息如 question）。
// 与 Python 参考版一致：先 list events 筛选，再通过 /events/{id} 获取详情。
// id 支持字符串格式（Gamma API 返回的 id 为字符串），内部自动转换为 int。
func (c *Client) GetGammaEventByID(id string) (json.RawMessage, error) {
	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid event id %q: %w", id, err)
	}
	return c.GammaAPI.GetEventByID(n, nil, nil)
}

// GetOrderBook 获取指定 token 的订单簿。
func (c *Client) GetOrderBook(tokenID string) (*polymarketSDK.OrderBookSummary, error) {
	return c.ClobClient.GetOrderBook(tokenID)
}

// GetMarkets 获取指定 conditionID 的市场信息。
func (c *Client) GetMarkets(conditionID string) (json.RawMessage, error) {
	query := &polymarketSDK.GammaMarketsQuery{
		ConditionIDs: []string{conditionID},
	}
	return c.GammaAPI.ListMarkets(query)
}

// CreateOrder 创建并提交 V2 限价订单（GTC 类型）。
func (c *Client) CreateOrder(tokenID string, price, size float64, side string) (*polymarketSDK.PostOrderResultV2, error) {
	args := &polymarketSDK.OrderArgsV2{
		TokenID: tokenID,
		Price:   price,
		Size:    size,
		Side:    side,
	}
	return c.ClobClient.CreateAndPostOrderV2(args, nil, polymarketSDK.OrderTypeGTC, false, false)
}

// CancelOrder 取消指定订单。
func (c *Client) CancelOrder(orderID string) (interface{}, error) {
	return c.ClobClient.Cancel(orderID)
}

// GetPositions 获取当前用户持仓。
// 需要客户端以 L2 模式创建（即提供了 privateKey），否则返回错误。
func (c *Client) GetPositions() ([]polymarketSDK.Position, error) {
	addr := c.ClobClient.GetAddress()
	if addr == "" {
		return nil, fmt.Errorf("no address available, private key is required to query positions")
	}
	query := &polymarketSDK.PositionsQuery{
		User: addr,
	}
	return c.DataAPI.GetPositions(query)
}
