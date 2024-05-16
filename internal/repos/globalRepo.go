package repos

import (
	"sync"

	"xorm.io/xorm"
)

var singular GlobalRepo

//go:generate mockgen -source=./globalRepo.go -destination=./mocks/GlobalRepo.go -package=mock_repos GlobalRepo
type GlobalRepo interface {
	DB() *xorm.Engine
	Products() Products
}

func NewGlobalRepo(db *xorm.Engine) (GlobalRepo, error) {
	if singular == nil {
		singular = &globalRepo{
			db:    db,
			mutex: &sync.RWMutex{},
			repos: make(map[string]interface{}),
		}
	}

	return singular, nil
}

type globalRepo struct {
	db    *xorm.Engine
	repos map[string]interface{}
	mutex *sync.RWMutex
}

func (gr *globalRepo) DB() *xorm.Engine {
	return gr.db
}

func (gr *globalRepo) factory(key string, fn func(db *xorm.Engine) interface{}) interface{} {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	val, exists := gr.repos[key]
	if exists {
		return val
	}

	nFac := fn(gr.db)
	gr.repos[key] = nFac

	return nFac
}

func (gr *globalRepo) Products() Products {
	return gr.factory("Products", func(db *xorm.Engine) interface{} { return NewProducts(db) }).(Products)
}
