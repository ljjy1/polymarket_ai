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

var _ PredictionsDao = (*predictionsDao)(nil)

// PredictionsDao defining the dao interface
type PredictionsDao interface {
	Create(ctx context.Context, table *model.Predictions) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Predictions) error
	GetByID(ctx context.Context, id uint64) (*model.Predictions, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Predictions, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Predictions) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Predictions) error
}

type predictionsDao struct {
	db    *gorm.DB
	cache cache.PredictionsCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewPredictionsDao creating the dao interface
func NewPredictionsDao(db *gorm.DB, xCache cache.PredictionsCache) PredictionsDao {
	if xCache == nil {
		return &predictionsDao{db: db}
	}
	return &predictionsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *predictionsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new predictions, insert the record and the id value is written back to the table
func (d *predictionsDao) Create(ctx context.Context, table *model.Predictions) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a predictions by id
func (d *predictionsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Predictions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a predictions by id, support partial update
func (d *predictionsDao) UpdateByID(ctx context.Context, table *model.Predictions) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *predictionsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Predictions) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.MarketID != 0 {
		update["market_id"] = table.MarketID
	}
	if table.PredictedProbability.IsZero() == false {
		update["predicted_probability"] = table.PredictedProbability
	}
	if table.Confidence.IsZero() == false {
		update["confidence"] = table.Confidence
	}
	if table.Direction != "" {
		update["direction"] = table.Direction
	}
	if table.KeyFactors.String() != "" {
		update["key_factors"] = table.KeyFactors
	}
	if table.RiskFactors.String() != "" {
		update["risk_factors"] = table.RiskFactors
	}
	if table.TechnicalAnalysis != "" {
		update["technical_analysis"] = table.TechnicalAnalysis
	}
	if table.SentimentAnalysis != "" {
		update["sentiment_analysis"] = table.SentimentAnalysis
	}
	if table.NewsImpact != "" {
		update["news_impact"] = table.NewsImpact
	}
	if table.OnchainAnalysis != "" {
		update["onchain_analysis"] = table.OnchainAnalysis
	}
	if table.Reasoning != "" {
		update["reasoning"] = table.Reasoning
	}
	if table.RecommendedAction != "" {
		update["recommended_action"] = table.RecommendedAction
	}
	if table.MarketProbability.IsZero() == false {
		update["market_probability"] = table.MarketProbability
	}
	if table.Edge.IsZero() == false {
		update["edge"] = table.Edge
	}
	if table.ModelVersion != "" {
		update["model_version"] = table.ModelVersion
	}
	if table.PromptVersion != "" {
		update["prompt_version"] = table.PromptVersion
	}
	if table.Seed != 0 {
		update["seed"] = table.Seed
	}
	if table.RawRequest.String() != "" {
		update["raw_request"] = table.RawRequest
	}
	if table.RawResponse.String() != "" {
		update["raw_response"] = table.RawResponse
	}
	if table.DataSnapshot.String() != "" {
		update["data_snapshot"] = table.DataSnapshot
	}
	if table.TokensUsed != 0 {
		update["tokens_used"] = table.TokensUsed
	}
	if table.LatencyMs != 0 {
		update["latency_ms"] = table.LatencyMs
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a predictions by id
func (d *predictionsDao) GetByID(ctx context.Context, id uint64) (*model.Predictions, error) {
	// no cache
	if d.cache == nil {
		record := &model.Predictions{}
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
			table := &model.Predictions{}
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
			if err = d.cache.Set(ctx, id, table, cache.PredictionsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Predictions)
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

// GetByColumns get a paginated list of predictionss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *predictionsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Predictions, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.PredictionsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Predictions{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Predictions{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *predictionsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Predictions) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *predictionsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Predictions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *predictionsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Predictions) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
