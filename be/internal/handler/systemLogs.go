package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"be/internal/cache"
	"be/internal/dao"
	"be/internal/database"
	"be/internal/ecode"
	"be/internal/model"
	"be/internal/types"
)

var _ SystemLogsHandler = (*systemLogsHandler)(nil)

// SystemLogsHandler defining the handler interface
type SystemLogsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type systemLogsHandler struct {
	iDao dao.SystemLogsDao
}

// NewSystemLogsHandler creating the handler interface
func NewSystemLogsHandler() SystemLogsHandler {
	return &systemLogsHandler{
		iDao: dao.NewSystemLogsDao(
			database.GetDB(), // db driver is mysql
			cache.NewSystemLogsCache(database.GetCacheType()),
		),
	}
}

// Create a new systemLogs
// @Summary Create a new systemLogs
// @Description Creates a new systemLogs entity using the provided data in the request body.
// @Tags systemLogs
// @Accept json
// @Produce json
// @Param data body types.CreateSystemLogsRequest true "systemLogs information"
// @Success 200 {object} types.CreateSystemLogsReply{}
// @Router /api/v1/systemLogs [post]
// @Security BearerAuth
func (h *systemLogsHandler) Create(c *gin.Context) {
	form := &types.CreateSystemLogsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	systemLogs := &model.SystemLogs{}
	err = copier.Copy(systemLogs, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateSystemLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, systemLogs)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": systemLogs.ID})
}

// DeleteByID delete a systemLogs by id
// @Summary Delete a systemLogs by id
// @Description Deletes a existing systemLogs identified by the given id in the path.
// @Tags systemLogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteSystemLogsByIDReply{}
// @Router /api/v1/systemLogs/{id} [delete]
// @Security BearerAuth
func (h *systemLogsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getSystemLogsIDFromPath(c)
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

// UpdateByID update a systemLogs by id
// @Summary Update a systemLogs by id
// @Description Updates the specified systemLogs by given id in the path, support partial update.
// @Tags systemLogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateSystemLogsByIDRequest true "systemLogs information"
// @Success 200 {object} types.UpdateSystemLogsByIDReply{}
// @Router /api/v1/systemLogs/{id} [put]
// @Security BearerAuth
func (h *systemLogsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getSystemLogsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateSystemLogsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	systemLogs := &model.SystemLogs{}
	err = copier.Copy(systemLogs, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDSystemLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, systemLogs)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a systemLogs by id
// @Summary Get a systemLogs by id
// @Description Gets detailed information of a systemLogs specified by the given id in the path.
// @Tags systemLogs
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSystemLogsByIDReply{}
// @Router /api/v1/systemLogs/{id} [get]
// @Security BearerAuth
func (h *systemLogsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getSystemLogsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	systemLogs, err := h.iDao.GetByID(ctx, id)
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

	data := &types.SystemLogsObjDetail{}
	err = copier.Copy(data, systemLogs)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDSystemLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"systemLogs": data})
}

// List get a paginated list of systemLogss by custom conditions
// @Summary Get a paginated list of systemLogss by custom conditions
// @Description Returns a paginated list of systemLogs based on query filters, including page number and size.
// @Tags systemLogs
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListSystemLogssReply{}
// @Router /api/v1/systemLogs/list [post]
// @Security BearerAuth
func (h *systemLogsHandler) List(c *gin.Context) {
	form := &types.ListSystemLogssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	systemLogss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSystemLogss(systemLogss)
	if err != nil {
		response.Error(c, ecode.ErrListSystemLogs)
		return
	}

	response.Success(c, gin.H{
		"systemLogss": data,
		"total":       total,
	})
}

func getSystemLogsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertSystemLogs(systemLogs *model.SystemLogs) (*types.SystemLogsObjDetail, error) {
	data := &types.SystemLogsObjDetail{}
	err := copier.Copy(data, systemLogs)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertSystemLogss(fromValues []*model.SystemLogs) ([]*types.SystemLogsObjDetail, error) {
	toValues := []*types.SystemLogsObjDetail{}
	for _, v := range fromValues {
		data, err := convertSystemLogs(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
