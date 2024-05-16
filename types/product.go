package types

import "time"

type Product struct {
	ID        int64      `json:"id" xorm:"'id' pk autoincr"`
	Name      string     `validate:"required" json:"name" xorm:"name"`
	Sku       string     `validate:"required" json:"sku" xorm:"sku"`
	Qty       int64      `validate:"required,min=1" json:"qty" xorm:"qty"`
	CreatedAt time.Time  `json:"createdAt" xorm:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" xorm:"updated_at"`
}

func (*Product) TableName() string {
	return "products"
}

type NewProduct struct {
	Name string `validate:"required" json:"name"`
	Sku  string `validate:"required" json:"sku"`
	Qty  int64  `validate:"required,min=1" json:"qty"`
}

type UpdateProduct struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Sku  *string `json:"sku"`
	Qty  *int64  `json:"qty"`
}
