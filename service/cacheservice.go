package service

import (
	"fmt"
	"github.com/linx93/base-service/cache"
	"github.com/ohmountain/memcache"
	"gorm.io/gorm"
)

type CacheService[V any] struct {
	BaseService[V]
	cache.ICache[V]
	cache.KeyGenerator
}

func newCacheService[V any](db *gorm.DB) *CacheService[V] {
	cacheService := CacheService[V]{}
	cacheService.BaseService.DB = db
	return &cacheService
}

func NewDefaultCacheService[V any](db *gorm.DB) *CacheService[V] {
	cacheService := newCacheService[V](db)
	//注入默认缓存
	cacheService.ICache = &cache.DefaultCache[V]{
		Name:  "默认缓存",
		Store: memcache.WithLRU[[]*V](30, false),
	}
	cacheService.KeyGenerator = cache.SimpleKeyGenerator{}
	return cacheService
}

func NewRedisCacheService[V any](db *gorm.DB) *CacheService[V] {
	cacheService := newCacheService[V](db)
	//注入redis缓存，还做具体实现
	cacheService.ICache = cache.NewRedisCache[V]("redis缓存")

	return cacheService
}

// cacheService 从写基本操作方法
func (cs CacheService[V]) getDB(dbs ...*gorm.DB) *gorm.DB {
	if dbs == nil || len(dbs) == 0 {
		return cs.DB
	}
	return dbs[0]
}

func (cs CacheService[V]) FindById(id int, dbs ...*gorm.DB) (result *V, err error) {
	generateKey := cs.Generate(*new(V), "FindByIds", id)
	c := cs.Get(generateKey)
	if c != nil {
		return c[0], nil
	}

	// SELECT * FROM tab WHERE id = 10;
	err = cs.getDB(dbs...).First(&result, id).Error
	if err == nil {
		cs.Put(generateKey, []*V{result})
	}
	return
}

func (cs CacheService[V]) FindByIds(ids []int, dbs ...*gorm.DB) (result []*V, err error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids is empty")
	}

	generateKey := cs.Generate(*new(V), "FindByIds", ids)
	c := cs.Get(generateKey)
	if c != nil {
		return c, nil
	}

	// SELECT * FROM tab WHERE id IN (1,2,3);
	err = cs.getDB(dbs...).Find(&result, ids).Error
	if err == nil {
		cs.Put(generateKey, result)
	}
	return
}

func (cs CacheService[V]) ListByModel(model *V, dbs ...*gorm.DB) (result []*V, err error) {

	generateKey := cs.Generate(*model, "ListByModel", *model)
	c := cs.Get(generateKey)
	if c != nil {
		return c, nil
	}

	err = cs.getDB(dbs...).Where(model).Find(&result).Error
	if err == nil {
		cs.Put(generateKey, result)
	}
	return
}

func (cs CacheService[V]) DeleteById(model *V, id int, dbs ...*gorm.DB) error {
	err := cs.getDB(dbs...).Delete(model, id).Error

	if err == nil {
		//删除成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) DeleteByIds(model *V, ids []int, dbs ...*gorm.DB) error {
	err := cs.getDB(dbs...).Delete(model, ids).Error
	if err == nil {
		//删除成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) DeleteByCond(model *V, where string, v []any, dbs ...*gorm.DB) error {
	// where->name LIKE ? v->string[]{"%jinzhu%"}
	// where->name LIKE ? v->string[]{"%jinzhu%",18}
	err := cs.getDB(dbs...).Where(where, v...).Delete(model).Error
	if err == nil {
		//删除成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) UpdateById(target *V, dbs ...*gorm.DB) error {
	//这里执行Updates后是不会把值填充到target的，所以只返回err就好了，返回target没意义，和传进来的一样
	err := cs.getDB(dbs...).Updates(target).Error
	if err == nil {
		//更新成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) UpdateByCond(where string, v []any, target *V, dbs ...*gorm.DB) error {
	//这里执行Updates后是不会把值填充到target的，所以只返回err就好了，返回target没意义，和传进来的一样
	err := cs.getDB(dbs...).Where(where, v...).Updates(target).Error
	if err == nil {
		//更新成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) Save(model *V, dbs ...*gorm.DB) (*V, error) {
	err := cs.getDB(dbs...).Save(model).Error
	if err == nil {
		//更新或插入成功需要清缓存
		cs.Clear()
	}
	return model, err
}

func (cs CacheService[V]) SaveBatch(models []*V, dbs ...*gorm.DB) error {
	err := cs.getDB(dbs...).Save(&models).Error
	if err == nil {
		//更新或插入成功需要清缓存
		cs.Clear()
	}
	return err
}

func (cs CacheService[V]) EndTx(db *gorm.DB, err error) (e error) {
	if err != nil {
		e = db.Rollback().Error
	} else {
		e = db.Commit().Error
	}
	return e
}
