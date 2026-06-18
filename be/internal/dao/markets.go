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

var _ MarketsDao = (*marketsDao)(nil)

// MarketsDao defining the dao interface
type MarketsDao interface {
	Create(ctx context.Context, table *model.Markets) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Markets) error
	GetByID(ctx context.Context, id uint64) (*model.Markets, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Markets, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Markets) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Markets) error
}

type marketsDao struct {
	db    *gorm.DB
	cache cache.MarketsCache  // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewMarketsDao creating the dao interface
func NewMarketsDao(db *gorm.DB, xCache cache.MarketsCache) MarketsDao {
	if xCache == nil {
		return &marketsDao{db: db}
	}
	return &marketsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *marketsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new markets, insert the record and the id value is written back to the table
func (d *marketsDao) Create(ctx context.Context, table *model.Markets) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a markets by id
func (d *marketsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Markets{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a markets by id, support partial update
func (d *marketsDao) UpdateByID(ctx context.Context, table *model.Markets) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *marketsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Markets) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.PolymarketConditionID != "" {
		update["polymarket_condition_id"] = table.PolymarketConditionID
	}
	if table.PolymarketTokenID != "" {
		update["polymarket_token_id"] = table.PolymarketTokenID
	}
	if table.EventSlug != "" {
		update["event_slug"] = table.EventSlug
	}
	if table.Question != "" {
		update["question"] = table.Question
	}
	if table.PriceThreshold != 0 {
		update["price_threshold"] = table.PriceThreshold
	}
	if table.ScanDate != nil && table.ScanDate.IsZero() == false {
		update["scan_date"] = table.ScanDate
	}
	if table.TargetDate != nil && table.TargetDate.IsZero() == false {
		update["target_date"] = table.TargetDate
	}
	if table.CurrentYesPrice.IsZero() == false {
		update["current_yes_price"] = table.CurrentYesPrice
	}
	if table.CurrentNoPrice.IsZero() == false {
		update["current_no_price"] = table.CurrentNoPrice
	}
	if table.SelectedAt != nil && table.SelectedAt.IsZero() == false {
		update["selected_at"] = table.SelectedAt
	}
	if table.Status != "" {
		update["status"] = table.Status
	}
	if table.Resolution != "" {
		update["resolution"] = table.Resolution
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a markets by id
func (d *marketsDao) GetByID(ctx context.Context, id uint64) (*model.Markets, error) {
	// no cache
	if d.cache == nil {
		record := &model.Markets{}
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
			table := &model.Markets{}
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
			if err = d.cache.Set(ctx, id, table, cache.MarketsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Markets)
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

// GetByColumns get a paginated list of marketss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *marketsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Markets, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.MarketsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Markets{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Markets{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *marketsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Markets) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *marketsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Markets{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *marketsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Markets) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
