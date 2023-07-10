package service

import (
	"gorm.io/gorm"
)

type IBaseService[T any] interface {
	FindById(id int, dbs ...*gorm.DB) (*T, error)
	FindByIds(ids []int, dbs ...*gorm.DB) (result []*T, err error)
	ListByModel(model *T, dbs ...*gorm.DB) (result []*T, err error)

	DeleteById(model *T, id int, dbs ...*gorm.DB) error
	DeleteByIds(model *T, ids []int, dbs ...*gorm.DB) error
	DeleteByCond(model *T, where string, v []any, dbs ...*gorm.DB) error

	//UpdateById 更据主键更新，只更新非零值，建议数据库对象的属性建议使用*int而不是int
	UpdateById(target *T, dbs ...*gorm.DB) error
	UpdateByCond(where string, v []any, target *T, dbs ...*gorm.DB) error
	Save(model *T, dbs ...*gorm.DB) (*T, error)
	SaveBatch(models []*T, dbs ...*gorm.DB) error

	EndTx(db *gorm.DB, err error) (e error)
	// todo implement
}

type BaseService[T any] struct {
	DB *gorm.DB
}

func NewBaseService[T any](db *gorm.DB) BaseService[T] {
	return BaseService[T]{DB: db}
}

func (bService BaseService[T]) getDB(dbs ...*gorm.DB) *gorm.DB {
	if dbs == nil || len(dbs) == 0 {
		return bService.DB
	}
	return dbs[0]
}

func (bService BaseService[T]) FindById(id int, dbs ...*gorm.DB) (result *T, err error) {
	// SELECT * FROM tab WHERE id = 10;
	err = bService.getDB(dbs...).First(&result, id).Error
	return
}

func (bService BaseService[T]) FindByIds(ids []int, dbs ...*gorm.DB) (result []*T, err error) {
	if len(ids) == 0 {
		ids = []int{-1}
	}
	// SELECT * FROM tab WHERE id IN (1,2,3);
	err = bService.getDB(dbs...).Find(&result, ids).Error
	return
}

func (bService BaseService[T]) ListByModel(model *T, dbs ...*gorm.DB) (result []*T, err error) {
	err = bService.getDB(dbs...).Where(model).Find(&result).Error
	return
}

func (bService BaseService[T]) DeleteById(model *T, id int, dbs ...*gorm.DB) error {
	err := bService.getDB(dbs...).Delete(model, id).Error
	return err
}

func (bService BaseService[T]) DeleteByIds(model *T, ids []int, dbs ...*gorm.DB) error {
	err := bService.getDB(dbs...).Delete(model, ids).Error
	return err
}

func (bService BaseService[T]) DeleteByCond(model *T, where string, v []any, dbs ...*gorm.DB) error {
	// where->name LIKE ? v->string[]{"%jinzhu%"}
	// where->name LIKE ? v->string[]{"%jinzhu%",18}
	err := bService.getDB(dbs...).Where(where, v...).Delete(model).Error
	return err
}

func (bService BaseService[T]) UpdateById(target *T, dbs ...*gorm.DB) error {
	//这里执行Updates后是不会把值填充到target的，所以只返回err就好了，返回target没意义，和传进来的一样
	err := bService.getDB(dbs...).Updates(target).Error
	return err
}

func (bService BaseService[T]) UpdateByCond(where string, v []any, target *T, dbs ...*gorm.DB) error {
	//这里执行Updates后是不会把值填充到target的，所以只返回err就好了，返回target没意义，和传进来的一样
	err := bService.getDB(dbs...).Where(where, v...).Updates(target).Error
	return err
}

func (bService BaseService[T]) Save(model *T, dbs ...*gorm.DB) (*T, error) {
	err := bService.getDB(dbs...).Save(model).Error
	return model, err
}

func (bService BaseService[T]) SaveBatch(models []*T, dbs ...*gorm.DB) error {
	err := bService.getDB(dbs...).Save(&models).Error
	return err
}

func (bService BaseService[T]) EndTx(db *gorm.DB, err error) (e error) {
	if err != nil {
		e = db.Rollback().Error
	} else {
		e = db.Commit().Error
	}
	return e
}
