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

var _ PredictionsHandler = (*predictionsHandler)(nil)

// PredictionsHandler defining the handler interface
type PredictionsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type predictionsHandler struct {
	iDao dao.PredictionsDao
}

// NewPredictionsHandler creating the handler interface
func NewPredictionsHandler() PredictionsHandler {
	return &predictionsHandler{
		iDao: dao.NewPredictionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewPredictionsCache(database.GetCacheType()),
		),
	}
}

// Create a new predictions
// @Summary Create a new predictions
// @Description Creates a new predictions entity using the provided data in the request body.
// @Tags predictions
// @Accept json
// @Produce json
// @Param data body types.CreatePredictionsRequest true "predictions information"
// @Success 200 {object} types.CreatePredictionsReply{}
// @Router /api/v1/predictions [post]
// @Security BearerAuth
func (h *predictionsHandler) Create(c *gin.Context) {
	form := &types.CreatePredictionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	predictions := &model.Predictions{}
	err = copier.Copy(predictions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreatePredictions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, predictions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": predictions.ID})
}

// DeleteByID delete a predictions by id
// @Summary Delete a predictions by id
// @Description Deletes a existing predictions identified by the given id in the path.
// @Tags predictions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePredictionsByIDReply{}
// @Router /api/v1/predictions/{id} [delete]
// @Security BearerAuth
func (h *predictionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getPredictionsIDFromPath(c)
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

// UpdateByID update a predictions by id
// @Summary Update a predictions by id
// @Description Updates the specified predictions by given id in the path, support partial update.
// @Tags predictions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePredictionsByIDRequest true "predictions information"
// @Success 200 {object} types.UpdatePredictionsByIDReply{}
// @Router /api/v1/predictions/{id} [put]
// @Security BearerAuth
func (h *predictionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getPredictionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePredictionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	predictions := &model.Predictions{}
	err = copier.Copy(predictions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDPredictions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, predictions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a predictions by id
// @Summary Get a predictions by id
// @Description Gets detailed information of a predictions specified by the given id in the path.
// @Tags predictions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPredictionsByIDReply{}
// @Router /api/v1/predictions/{id} [get]
// @Security BearerAuth
func (h *predictionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getPredictionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	predictions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.PredictionsObjDetail{}
	err = copier.Copy(data, predictions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDPredictions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"predictions": data})
}

// List get a paginated list of predictionss by custom conditions
// @Summary Get a paginated list of predictionss by custom conditions
// @Description Returns a paginated list of predictions based on query filters, including page number and size.
// @Tags predictions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListPredictionssReply{}
// @Router /api/v1/predictions/list [post]
// @Security BearerAuth
func (h *predictionsHandler) List(c *gin.Context) {
	form := &types.ListPredictionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	predictionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertPredictionss(predictionss)
	if err != nil {
		response.Error(c, ecode.ErrListPredictions)
		return
	}

	response.Success(c, gin.H{
		"predictionss": data,
		"total":        total,
	})
}

func getPredictionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertPredictions(predictions *model.Predictions) (*types.PredictionsObjDetail, error) {
	data := &types.PredictionsObjDetail{}
	err := copier.Copy(data, predictions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertPredictionss(fromValues []*model.Predictions) ([]*types.PredictionsObjDetail, error) {
	toValues := []*types.PredictionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertPredictions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
