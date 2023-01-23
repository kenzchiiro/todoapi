package todo

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) New(todo *Todo) error {
	return g.db.Create(todo).Error
}

func (g *GormStore) Find(todos *[]Todo) error {
	return g.db.Find(&todos).Error
}

func (g *GormStore) Delete(todo *Todo, id int) error {
	return g.db.Where("id = ?", id).Delete(todo).Error
}
