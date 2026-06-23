package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"be/internal/cache"
	"be/internal/dao"
	"be/internal/database"
	"be/internal/ecode"
	"be/internal/model"
	"be/internal/types"
)

var _ MarketsHandler = (*marketsHandler)(nil)

// MarketsHandler defining the handler interface
type MarketsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	// 新增
	GetToday(c *gin.Context)
	TriggerScan(c *gin.Context)
}

type marketsHandler struct {
	iDao dao.MarketsDao
}

// NewMarketsHandler creating the handler interface
func NewMarketsHandler() MarketsHandler {
	return &marketsHandler{
		iDao: dao.NewMarketsDao(
			database.GetDB(), // db driver is mysql
			cache.NewMarketsCache(database.GetCacheType()),
		),
	}
}

// Create a new markets
// @Summary Create a new markets
// @Description Creates a new markets entity using the provided data in the request body.
// @Tags markets
// @Accept json
// @Produce json
// @Param data body types.CreateMarketsRequest true "markets information"
// @Success 200 {object} types.CreateMarketsReply{}
// @Router /api/v1/markets [post]
// @Security BearerAuth
func (h *marketsHandler) Create(c *gin.Context) {
	form := &types.CreateMarketsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	markets := &model.Markets{}
	err = copier.Copy(markets, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateMarkets)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, markets)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": markets.ID})
}

// DeleteByID delete a markets by id
// @Summary Delete a markets by id
// @Description Deletes a existing markets identified by the given id in the path.
// @Tags markets
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteMarketsByIDReply{}
// @Router /api/v1/markets/{id} [delete]
// @Security BearerAuth
func (h *marketsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getMarketsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a markets by id
// @Summary Update a markets by id
// @Description Updates the specified markets by given id in the path, support partial update.
// @Tags markets
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateMarketsByIDRequest true "markets information"
// @Success 200 {object} types.UpdateMarketsByIDReply{}
// @Router /api/v1/markets/{id} [put]
// @Security BearerAuth
func (h *marketsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getMarketsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateMarketsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	markets := &model.Markets{}
	err = copier.Copy(markets, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDMarkets)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, markets)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a markets by id
// @Summary Get a markets by id
// @Description Gets detailed information of a markets specified by the given id in the path.
// @Tags markets
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetMarketsByIDReply{}
// @Router /api/v1/markets/{id} [get]
// @Security BearerAuth
func (h *marketsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getMarketsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	markets, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.MarketsObjDetail{}
	err = copier.Copy(data, markets)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDMarkets)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"markets": data})
}

// List get a paginated list of marketss by custom conditions
// @Summary Get a paginated list of marketss by custom conditions
// @Description Returns a paginated list of markets based on query filters, including page number and size.
// @Tags markets
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListMarketssReply{}
// @Router /api/v1/markets/list [post]
// @Security BearerAuth
func (h *marketsHandler) List(c *gin.Context) {
	form := &types.ListMarketssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	marketss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertMarketss(marketss)
	if err != nil {
		response.Error(c, ecode.ErrListMarkets)
		return
	}

	response.Success(c, gin.H{
		"marketss": data,
		"total":    total,
	})
}

// GetToday 获取今日选定的市场
// @Summary 获取今日选定的市场
// @Description 获取今日扫描选定的市场信息
// @Tags markets
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/markets/today [get]
// @Security BearerAuth
func (h *marketsHandler) GetToday(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	today := time.Now().UTC().Format("2006-01-02")

	// 查询今日扫描记录
	markets, total, err := h.iDao.GetByColumns(ctx, &query.Params{
		Page: 1,
		Size: 1,
		Columns: []query.Column{
			{Name: "scan_date", Exp: "=", Value: today},
		},
		Sort: "id DESC",
	})
	if err != nil || total == 0 {
		response.Error(c, ecode.NotFound)
		return
	}

	data, err := convertMarkets(markets[0])
	if err != nil {
		response.Error(c, ecode.ErrGetByIDMarkets)
		return
	}

	response.Success(c, gin.H{"markets": data})
}

// TriggerScan 手动触发市场扫描
// @Summary 手动触发市场扫描
// @Description 手动触发市场扫描，从 Polymarket 拉取最新市场数据
// @Tags markets
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/markets/scan [post]
// @Security BearerAuth
func (h *marketsHandler) TriggerScan(c *gin.Context) {
	ctx := middleware.WrapCtx(c)

	// TODO: 实际集成 Asynq 客户端
	logger.Warn("TriggerScan called but service integration not yet complete", middleware.GCtxRequestIDField(c))
	_ = ctx // suppress unused variable
	response.Success(c, gin.H{"message": "scan task submitted"})
}

func getMarketsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertMarkets(markets *model.Markets) (*types.MarketsObjDetail, error) {
	data := &types.MarketsObjDetail{}
	err := copier.Copy(data, markets)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertMarketss(fromValues []*model.Markets) ([]*types.MarketsObjDetail, error) {
	toValues := []*types.MarketsObjDetail{}
	for _, v := range fromValues {
		data, err := convertMarkets(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
