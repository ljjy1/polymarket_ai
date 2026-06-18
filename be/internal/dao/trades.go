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

var _ TradesDao = (*tradesDao)(nil)

// TradesDao defining the dao interface
type TradesDao interface {
	Create(ctx context.Context, table *model.Trades) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Trades) error
	GetByID(ctx context.Context, id uint64) (*model.Trades, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Trades, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Trades) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Trades) error
}

type tradesDao struct {
	db    *gorm.DB
	cache cache.TradesCache   // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewTradesDao creating the dao interface
func NewTradesDao(db *gorm.DB, xCache cache.TradesCache) TradesDao {
	if xCache == nil {
		return &tradesDao{db: db}
	}
	return &tradesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *tradesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new trades, insert the record and the id value is written back to the table
func (d *tradesDao) Create(ctx context.Context, table *model.Trades) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a trades by id
func (d *tradesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Trades{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a trades by id, support partial update
func (d *tradesDao) UpdateByID(ctx context.Context, table *model.Trades) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *tradesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Trades) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.StrategyID != 0 {
		update["strategy_id"] = table.StrategyID
	}
	if table.MarketID != 0 {
		update["market_id"] = table.MarketID
	}
	if table.PolymarketOrderID != "" {
		update["polymarket_order_id"] = table.PolymarketOrderID
	}
	if table.Side != "" {
		update["side"] = table.Side
	}
	if table.Action != "" {
		update["action"] = table.Action
	}
	if table.Amount.IsZero() == false {
		update["amount"] = table.Amount
	}
	if table.Price.IsZero() == false {
		update["price"] = table.Price
	}
	if table.Shares.IsZero() == false {
		update["shares"] = table.Shares
	}
	if table.Status != "" {
		update["status"] = table.Status
	}
	if table.Fee.IsZero() == false {
		update["fee"] = table.Fee
	}
	if table.Pnl.IsZero() == false {
		update["pnl"] = table.Pnl
	}
	if table.CloseReason != "" {
		update["close_reason"] = table.CloseReason
	}
	if table.FilledAt != nil && table.FilledAt.IsZero() == false {
		update["filled_at"] = table.FilledAt
	}
	if table.ClosedAt != nil && table.ClosedAt.IsZero() == false {
		update["closed_at"] = table.ClosedAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a trades by id
func (d *tradesDao) GetByID(ctx context.Context, id uint64) (*model.Trades, error) {
	// no cache
	if d.cache == nil {
		record := &model.Trades{}
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
			table := &model.Trades{}
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
			if err = d.cache.Set(ctx, id, table, cache.TradesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Trades)
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

// GetByColumns get a paginated list of tradess by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *tradesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Trades, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.TradesColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Trades{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Trades{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *tradesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Trades) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *tradesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Trades{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *tradesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Trades) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
