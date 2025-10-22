package order

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"gt=0"`
	PaymentDt    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"gt=0"`
	GoodsTotal   int    `json:"goods_total" validate:"gt=0"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"gt=0"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"gt=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"gt=0"`
	NmID        int    `json:"nm_id" validate:"gt=0"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"gt=0"`

	OrderUID string `json:"-"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"email"`
}

type Order struct {
	OrderUID    string `json:"order_uid" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`

	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []Item   `json:"items" validate:"required,min=1,dive"`

	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

func (o *Order) Validate() error {
	// validate.Struct(o) проверяет структуру `o` на основе тегов `validate`
	return validate.Struct(o)
}
