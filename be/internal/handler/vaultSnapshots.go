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

var _ VaultSnapshotsHandler = (*vaultSnapshotsHandler)(nil)

// VaultSnapshotsHandler defining the handler interface
type VaultSnapshotsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type vaultSnapshotsHandler struct {
	iDao dao.VaultSnapshotsDao
}

// NewVaultSnapshotsHandler creating the handler interface
func NewVaultSnapshotsHandler() VaultSnapshotsHandler {
	return &vaultSnapshotsHandler{
		iDao: dao.NewVaultSnapshotsDao(
			database.GetDB(), // db driver is mysql
			cache.NewVaultSnapshotsCache(database.GetCacheType()),
		),
	}
}

// Create a new vaultSnapshots
// @Summary Create a new vaultSnapshots
// @Description Creates a new vaultSnapshots entity using the provided data in the request body.
// @Tags vaultSnapshots
// @Accept json
// @Produce json
// @Param data body types.CreateVaultSnapshotsRequest true "vaultSnapshots information"
// @Success 200 {object} types.CreateVaultSnapshotsReply{}
// @Router /api/v1/vaultSnapshots [post]
// @Security BearerAuth
func (h *vaultSnapshotsHandler) Create(c *gin.Context) {
	form := &types.CreateVaultSnapshotsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	vaultSnapshots := &model.VaultSnapshots{}
	err = copier.Copy(vaultSnapshots, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateVaultSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, vaultSnapshots)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": vaultSnapshots.ID})
}

// DeleteByID delete a vaultSnapshots by id
// @Summary Delete a vaultSnapshots by id
// @Description Deletes a existing vaultSnapshots identified by the given id in the path.
// @Tags vaultSnapshots
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteVaultSnapshotsByIDReply{}
// @Router /api/v1/vaultSnapshots/{id} [delete]
// @Security BearerAuth
func (h *vaultSnapshotsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getVaultSnapshotsIDFromPath(c)
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

// UpdateByID update a vaultSnapshots by id
// @Summary Update a vaultSnapshots by id
// @Description Updates the specified vaultSnapshots by given id in the path, support partial update.
// @Tags vaultSnapshots
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateVaultSnapshotsByIDRequest true "vaultSnapshots information"
// @Success 200 {object} types.UpdateVaultSnapshotsByIDReply{}
// @Router /api/v1/vaultSnapshots/{id} [put]
// @Security BearerAuth
func (h *vaultSnapshotsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getVaultSnapshotsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateVaultSnapshotsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	vaultSnapshots := &model.VaultSnapshots{}
	err = copier.Copy(vaultSnapshots, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDVaultSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, vaultSnapshots)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a vaultSnapshots by id
// @Summary Get a vaultSnapshots by id
// @Description Gets detailed information of a vaultSnapshots specified by the given id in the path.
// @Tags vaultSnapshots
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetVaultSnapshotsByIDReply{}
// @Router /api/v1/vaultSnapshots/{id} [get]
// @Security BearerAuth
func (h *vaultSnapshotsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getVaultSnapshotsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vaultSnapshots, err := h.iDao.GetByID(ctx, id)
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

	data := &types.VaultSnapshotsObjDetail{}
	err = copier.Copy(data, vaultSnapshots)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDVaultSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"vaultSnapshots": data})
}

// List get a paginated list of vaultSnapshotss by custom conditions
// @Summary Get a paginated list of vaultSnapshotss by custom conditions
// @Description Returns a paginated list of vaultSnapshots based on query filters, including page number and size.
// @Tags vaultSnapshots
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListVaultSnapshotssReply{}
// @Router /api/v1/vaultSnapshots/list [post]
// @Security BearerAuth
func (h *vaultSnapshotsHandler) List(c *gin.Context) {
	form := &types.ListVaultSnapshotssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vaultSnapshotss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertVaultSnapshotss(vaultSnapshotss)
	if err != nil {
		response.Error(c, ecode.ErrListVaultSnapshots)
		return
	}

	response.Success(c, gin.H{
		"vaultSnapshotss": data,
		"total":           total,
	})
}

func getVaultSnapshotsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertVaultSnapshots(vaultSnapshots *model.VaultSnapshots) (*types.VaultSnapshotsObjDetail, error) {
	data := &types.VaultSnapshotsObjDetail{}
	err := copier.Copy(data, vaultSnapshots)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertVaultSnapshotss(fromValues []*model.VaultSnapshots) ([]*types.VaultSnapshotsObjDetail, error) {
	toValues := []*types.VaultSnapshotsObjDetail{}
	for _, v := range fromValues {
		data, err := convertVaultSnapshots(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
