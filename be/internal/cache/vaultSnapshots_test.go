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

func newVaultSnapshotsCache() *gotest.Cache {
	record1 := &model.VaultSnapshots{}
	record1.ID = 1
	record2 := &model.VaultSnapshots{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewVaultSnapshotsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_vaultSnapshotsCache_Set(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VaultSnapshots)
	err := c.ICache.(VaultSnapshotsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(VaultSnapshotsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_vaultSnapshotsCache_Get(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VaultSnapshots)
	err := c.ICache.(VaultSnapshotsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VaultSnapshotsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(VaultSnapshotsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_vaultSnapshotsCache_MultiGet(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	var testData []*model.VaultSnapshots
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.VaultSnapshots))
	}

	err := c.ICache.(VaultSnapshotsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VaultSnapshotsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.VaultSnapshots))
	}
}

func Test_vaultSnapshotsCache_MultiSet(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	var testData []*model.VaultSnapshots
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.VaultSnapshots))
	}

	err := c.ICache.(VaultSnapshotsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vaultSnapshotsCache_Del(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VaultSnapshots)
	err := c.ICache.(VaultSnapshotsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vaultSnapshotsCache_SetCacheWithNotFound(t *testing.T) {
	c := newVaultSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VaultSnapshots)
	err := c.ICache.(VaultSnapshotsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(VaultSnapshotsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewVaultSnapshotsCache(t *testing.T) {
	c := NewVaultSnapshotsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewVaultSnapshotsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewVaultSnapshotsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
