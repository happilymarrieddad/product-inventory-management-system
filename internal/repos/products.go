package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/utils"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	"xorm.io/xorm"
)

type ProductsFind struct {
	Limit  int
	Offset int
	IDs    []int64
	Names  []string
	Skus   []string
}

//go:generate mockgen -source=./products.go -destination=./mocks/Products.go -package=mock_repos Products
type Products interface {
	Find(ctx context.Context, opts *ProductsFind) ([]*types.Product, int64, error)
	FindTx(ctx context.Context, tx *xorm.Session, opts *ProductsFind) ([]*types.Product, int64, error)
	Get(ctx context.Context, id int64) (*types.Product, bool, error)
	GetTx(ctx context.Context, tx *xorm.Session, id int64) (*types.Product, bool, error)
	Create(ctx context.Context, newProduct types.NewProduct) (*types.Product, error)
	CreateTx(ctx context.Context, tx *xorm.Session, newProduct types.NewProduct) (*types.Product, error)
	Update(ctx context.Context, diff *types.UpdateProduct) (*types.Product, error)
	UpdateTx(ctx context.Context, tx *xorm.Session, diff *types.UpdateProduct) (*types.Product, error)
	Destroy(ctx context.Context, id int64) error
	DestroyTx(ctx context.Context, tx *xorm.Session, id int64) error
}

func NewProducts(db *xorm.Engine) Products {
	return &productsRepo{db}
}

type productsRepo struct {
	db *xorm.Engine
}

func (r *productsRepo) Find(ctx context.Context, opts *ProductsFind) ([]*types.Product, int64, error) {
	var count int64
	res, err := wrapInSession(r.db, func(tx *xorm.Session) (any, error) {
		p, c, e := r.FindTx(ctx, tx, opts)
		if e != nil {
			return nil, e
		}
		count = c
		return p, nil
	})
	if err != nil {
		return nil, 0, err
	}

	return res.([]*types.Product), count, nil
}

func (r *productsRepo) FindTx(ctx context.Context, tx *xorm.Session, opts *ProductsFind) ([]*types.Product, int64, error) {
	if opts == nil {
		opts = &ProductsFind{Limit: 25}
	}

	if opts.Limit > 0 {
		if opts.Offset > 0 {
			tx = tx.Limit(opts.Limit, opts.Offset)
		} else {
			tx = tx.Limit(opts.Limit)
		}
	}

	if len(opts.IDs) > 0 {
		fmt.Println(utils.Int64ArrToInterfaceArr(opts.IDs...))
		tx = tx.In("id", utils.Int64ArrToInterfaceArr(opts.IDs...)...)
	}

	if len(opts.Names) > 0 {
		tx = tx.In("name", utils.AnyArrToInterfaceArr(opts.Names)...)
	}

	if len(opts.Skus) > 0 {
		tx = tx.In("sku", utils.AnyArrToInterfaceArr(opts.Skus)...)
	}

	objs := []*types.Product{}
	count, err := tx.OrderBy("id").FindAndCount(&objs)
	if err != nil {
		return nil, 0, normalizeErr("products", err)
	}

	return objs, count, nil
}

func (r *productsRepo) Get(ctx context.Context, id int64) (*types.Product, bool, error) {
	var exists bool
	res, err := wrapInSession(r.db, func(tx *xorm.Session) (any, error) {
		p, ex, e := r.GetTx(ctx, tx, id)
		if e != nil {
			return nil, e
		}
		exists = ex
		return p, nil
	})
	if err != nil {
		return nil, false, err
	}

	return res.(*types.Product), exists, nil
}

func (r *productsRepo) GetTx(ctx context.Context, tx *xorm.Session, id int64) (*types.Product, bool, error) {
	obj := &types.Product{}
	exists, err := tx.Where("id = ?", id).Get(obj)
	if err != nil {
		return nil, false, normalizeErr("products", err)
	}
	if !exists {
		return nil, exists, nil
	}

	return obj, exists, nil
}

func (r *productsRepo) Create(ctx context.Context, newProduct types.NewProduct) (*types.Product, error) {
	res, err := wrapInSession(r.db, func(tx *xorm.Session) (any, error) {
		return r.CreateTx(ctx, tx, newProduct)
	})
	if err != nil {
		return nil, err
	}

	return res.(*types.Product), nil
}

func (r *productsRepo) CreateTx(ctx context.Context, tx *xorm.Session, newProduct types.NewProduct) (*types.Product, error) {
	obj := &types.Product{
		Name:      newProduct.Name,
		Sku:       newProduct.Sku,
		Qty:       newProduct.Qty,
		CreatedAt: time.Now(),
	}

	if err := types.Validate(obj); err != nil {
		return nil, err
	}

	if _, err := tx.Insert(obj); err != nil {
		return nil, normalizeErr("products", err)
	}

	return obj, nil
}

func (r *productsRepo) Update(ctx context.Context, diff *types.UpdateProduct) (*types.Product, error) {
	res, err := wrapInSession(r.db, func(tx *xorm.Session) (any, error) {
		return r.UpdateTx(ctx, tx, diff)
	})
	if err != nil {
		return nil, err
	}

	return res.(*types.Product), nil
}

func (r *productsRepo) UpdateTx(ctx context.Context, tx *xorm.Session, diff *types.UpdateProduct) (*types.Product, error) {
	obj, exists, err := r.GetTx(ctx, tx, diff.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, types.NewNotFoundError("product not found by id")
	}

	if diff.Name != nil {
		obj.Name = *diff.Name
	}

	if diff.Sku != nil {
		obj.Sku = *diff.Sku
	}

	if diff.Qty != nil {
		obj.Qty = *diff.Qty
	}

	if err := types.Validate(obj); err != nil {
		return nil, err
	}

	obj.UpdatedAt = utils.Ref(time.Now())

	if _, err := tx.ID(diff.ID).Update(obj); err != nil {
		return nil, normalizeErr("products", err)
	}

	return obj, nil
}

func (r *productsRepo) Destroy(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (any, error) {
		return nil, r.DestroyTx(ctx, tx, id)
	})
	return err
}

func (r *productsRepo) DestroyTx(ctx context.Context, tx *xorm.Session, id int64) error {
	count, err := tx.Where("id = ?", id).Delete(&types.Product{})
	if err != nil {
		return normalizeErr("products", err)
	}
	if count == 0 {
		return types.NewBadRequestError("product not found")
	}
	return nil
}
