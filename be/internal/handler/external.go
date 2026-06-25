package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/response"

	"be/internal/config"
	"be/internal/external"
)

// ---------- 请求体 ----------

type klineReq struct {
	Klines1h []external.Kline `json:"klines1h" binding:"omitempty"`
	Klines1d []external.Kline `json:"klines1d" binding:"omitempty"`
}

type kline1hReq struct {
	Klines1h []external.Kline `json:"klines1h" binding:"required"`
}

// ---------- Handler ----------

type ExternalHandler interface {
	// 恐惧贪婪指数
	GetFearGreedIndex(c *gin.Context)
	GetFearGreedWithTrend(c *gin.Context)
	// 新闻
	GetNews(c *gin.Context)
	// 链上数据
	GetOnChainData(c *gin.Context)
	// 技术指标
	GetTechnicalIndicators(c *gin.Context)
	GetTechnicalIndicators1h(c *gin.Context)
}

type externalHandler struct {
	fetcher *external.Fetcher
}

func NewExternalHandler() ExternalHandler {
	cfg := config.Get()
	f := external.NewFetcher(
		cfg.GNews.BaseURL,
		cfg.GNews.APIKey,
		cfg.CryptoQuant.APIKey,
		cfg.CryptoQuant.BaseURL,
		cfg.FearGreedIndex.URL,
		cfg.Proxy.Addr,
	)
	return &externalHandler{fetcher: f}
}

// GetFearGreedIndex GET /api/v1/external/fear-greed-index
// @Summary 获取恐惧贪婪指数
// @Description 获取当前市场恐惧贪婪指数（0-100）
// @Tags external
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/external/fear-greed-index [get]
func (h *externalHandler) GetFearGreedIndex(c *gin.Context) {
	value, label, err := h.fetcher.GetFearGreedIndex()
	if err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{
		"value": value,
		"label": label,
	})
}

// GetFearGreedWithTrend GET /api/v1/external/fear-greed-index/trend?days=7
// @Summary 获取恐惧贪婪指数及趋势
// @Description 获取恐惧贪婪指数及多日趋势数组
// @Tags external
// @Accept json
// @Produce json
// @Param days query int false "天数（默认7）"
// @Success 200 {object} types.Result
// @Router /api/v1/external/fear-greed-index/trend [get]
func (h *externalHandler) GetFearGreedWithTrend(c *gin.Context) {
	days := 7
	if v := c.Query("days"); v != "" {
		if d, err := parseInt(v); err == nil && d > 0 {
			days = d
		}
	}
	value, label, trend, err := h.fetcher.GetFearGreedWithTrend(days)
	if err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{
		"value": value,
		"label": label,
		"trend": trend,
	})
}

// GetNews GET /api/v1/external/news?symbol=BTC&count=5
// @Summary 获取 BTC 相关新闻
// @Description 获取 BTC 相关新闻（需要配置 GNews API Key）
// @Tags external
// @Accept json
// @Produce json
// @Param symbol query string false "搜索关键词（默认BTC）"
// @Param count query int false "新闻数量（默认5，最大20）"
// @Success 200 {object} types.Result
// @Router /api/v1/external/news [get]
func (h *externalHandler) GetNews(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "BTC")
	count := 5
	if v := c.Query("count"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			count = n
		}
	}
	items := h.fetcher.GetNews(symbol, count)
	response.Success(c, gin.H{"news": items})
}

// GetOnChainData GET /api/v1/external/on-chain
// @Summary 获取链上数据
// @Description 从 CryptoQuant API 获取 BTC 链上数据（需要配置 API Key）
// @Tags external
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/external/on-chain [get]
func (h *externalHandler) GetOnChainData(c *gin.Context) {
	data := h.fetcher.GetOnChainData()
	response.Success(c, data)
}

// GetTechnicalIndicators POST /api/v1/external/technical-indicators
// @Summary 计算技术指标（日线）
// @Description 基于日线 K 线数据计算 RSI、MACD、布林带、EMA、ATR
// @Tags external
// @Accept json
// @Produce json
// @Param data body klineReq true "日线K线数据（至少100根）"
// @Success 200 {object} types.Result
// @Router /api/v1/external/technical-indicators [post]
func (h *externalHandler) GetTechnicalIndicators(c *gin.Context) {
	var req klineReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	ti, err := h.fetcher.GetTechnicalIndicators(req.Klines1h, req.Klines1d)
	if err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, ti)
}

// GetTechnicalIndicators1h POST /api/v1/external/technical-indicators/1h
// @Summary 计算技术指标（1小时线）
// @Description 基于1小时 K 线数据计算 RSI、MACD、布林带、EMA、ATR
// @Tags external
// @Accept json
// @Produce json
// @Param data body kline1hReq true "1小时K线数据（至少100根）"
// @Success 200 {object} types.Result
// @Router /api/v1/external/technical-indicators/1h [post]
func (h *externalHandler) GetTechnicalIndicators1h(c *gin.Context) {
	var req kline1hReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	ti, err := h.fetcher.GetTechnicalIndicators1h(req.Klines1h)
	if err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, ti)
}

// parseInt 辅助函数：将字符串转为 int
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
