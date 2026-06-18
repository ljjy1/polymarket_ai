package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"be/internal/database"
	"be/internal/model"
)

func newPredictionsCache() *gotest.Cache {
	record1 := &model.Predictions{}
	record1.ID = 1
	record2 := &model.Predictions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewPredictionsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_predictionsCache_Set(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Predictions)
	err := c.ICache.(PredictionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(PredictionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_predictionsCache_Get(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Predictions)
	err := c.ICache.(PredictionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PredictionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(PredictionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_predictionsCache_MultiGet(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	var testData []*model.Predictions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Predictions))
	}

	err := c.ICache.(PredictionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PredictionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Predictions))
	}
}

func Test_predictionsCache_MultiSet(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	var testData []*model.Predictions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Predictions))
	}

	err := c.ICache.(PredictionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_predictionsCache_Del(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Predictions)
	err := c.ICache.(PredictionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_predictionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newPredictionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Predictions)
	err := c.ICache.(PredictionsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(PredictionsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewPredictionsCache(t *testing.T) {
	c := NewPredictionsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewPredictionsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewPredictionsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
