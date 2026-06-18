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

func newMarketsCache() *gotest.Cache {
	record1 := &model.Markets{}
	record1.ID = 1
	record2 := &model.Markets{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewMarketsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_marketsCache_Set(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Markets)
	err := c.ICache.(MarketsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(MarketsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_marketsCache_Get(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Markets)
	err := c.ICache.(MarketsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MarketsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(MarketsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_marketsCache_MultiGet(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	var testData []*model.Markets
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Markets))
	}

	err := c.ICache.(MarketsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MarketsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Markets))
	}
}

func Test_marketsCache_MultiSet(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	var testData []*model.Markets
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Markets))
	}

	err := c.ICache.(MarketsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_marketsCache_Del(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Markets)
	err := c.ICache.(MarketsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_marketsCache_SetCacheWithNotFound(t *testing.T) {
	c := newMarketsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Markets)
	err := c.ICache.(MarketsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(MarketsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewMarketsCache(t *testing.T) {
	c := NewMarketsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewMarketsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewMarketsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
