package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"be/internal/cache"
	"be/internal/database"
	"be/internal/model"
)

var _ VaultSnapshotsDao = (*vaultSnapshotsDao)(nil)

// VaultSnapshotsDao defining the dao interface
type VaultSnapshotsDao interface {
	Create(ctx context.Context, table *model.VaultSnapshots) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.VaultSnapshots) error
	GetByID(ctx context.Context, id uint64) (*model.VaultSnapshots, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.VaultSnapshots, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.VaultSnapshots) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.VaultSnapshots) error
}

type vaultSnapshotsDao struct {
	db    *gorm.DB
	cache cache.VaultSnapshotsCache // if nil, the cache is not used.
	sfg   *singleflight.Group       // if cache is nil, the sfg is not used.
}

// NewVaultSnapshotsDao creating the dao interface
func NewVaultSnapshotsDao(db *gorm.DB, xCache cache.VaultSnapshotsCache) VaultSnapshotsDao {
	if xCache == nil {
		return &vaultSnapshotsDao{db: db}
	}
	return &vaultSnapshotsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *vaultSnapshotsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new vaultSnapshots, insert the record and the id value is written back to the table
func (d *vaultSnapshotsDao) Create(ctx context.Context, table *model.VaultSnapshots) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a vaultSnapshots by id
func (d *vaultSnapshotsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.VaultSnapshots{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a vaultSnapshots by id, support partial update
func (d *vaultSnapshotsDao) UpdateByID(ctx context.Context, table *model.VaultSnapshots) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *vaultSnapshotsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.VaultSnapshots) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.TotalAssets.IsZero() == false {
		update["total_assets"] = table.TotalAssets
	}
	if table.SharePrice.IsZero() == false {
		update["share_price"] = table.SharePrice
	}
	if table.Tvl.IsZero() == false {
		update["tvl"] = table.Tvl
	}
	if table.DepositorCount != 0 {
		update["depositor_count"] = table.DepositorCount
	}
	if table.DeployedAmount.IsZero() == false {
		update["deployed_amount"] = table.DeployedAmount
	}
	if table.SnapshotAt != nil && table.SnapshotAt.IsZero() == false {
		update["snapshot_at"] = table.SnapshotAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a vaultSnapshots by id
func (d *vaultSnapshotsDao) GetByID(ctx context.Context, id uint64) (*model.VaultSnapshots, error) {
	// no cache
	if d.cache == nil {
		record := &model.VaultSnapshots{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.VaultSnapshots{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.VaultSnapshotsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.VaultSnapshots)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of vaultSnapshotss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *vaultSnapshotsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.VaultSnapshots, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.VaultSnapshotsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.VaultSnapshots{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.VaultSnapshots{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *vaultSnapshotsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.VaultSnapshots) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *vaultSnapshotsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.VaultSnapshots{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *vaultSnapshotsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.VaultSnapshots) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
