package order

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(order Order) (bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// ----------------------------------- Checking order -----------------------------------

	var exstingOrderID int
	checkingQuery := `
		SELECT id 
		FROM orders
		WHERE order_uid = $1
	`
	err = tx.QueryRow(checkingQuery, order.OrderUID).Scan(&exstingOrderID)

	if err == nil {
		return false, nil

	}

	if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	// ----------------------------------- Create Delivery -----------------------------------

	deliverID, err := r.createDelivery(tx, order.Delivery)

	if err != nil {
		return false, err
	}

	// ----------------------------------- Create Payments -----------------------------------

	paymentID, err := r.createPayment(tx, order.Payment)

	if err != nil {
		return false, err
	}

	// ----------------------------------- Create Order -----------------------------------

	uID, err := r.createOrder(tx, order, deliverID, paymentID)
	if err != nil {
		return false, err
	}

	// ----------------------------------- Create Items -----------------------------------

	for _, v := range order.Items {
		err := r.createItem(tx, uID, v)
		if err != nil {
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) createDelivery(tx *sql.Tx, d Delivery) (int, error) {
	var id int

	query := `
		INSERT INTO deliveries (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := tx.QueryRow(query, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email).Scan(&id)

	return id, err
}

func (r *Repository) createPayment(tx *sql.Tx, p Payment) (int, error) {
	var id int

	query := `
		INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := tx.QueryRow(query, p.Transaction, p.RequestID, p.Currency, p.Provider, p.Amount, p.PaymentDt, p.Bank, p.DeliveryCost, p.GoodsTotal, p.CustomFee).Scan(&id)

	return id, err
}

func (r *Repository) createOrder(tx *sql.Tx, o Order, idDelivery int, idPayment int) (string, error) {
	var uid string

	query := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING order_uid
	`

	err := tx.QueryRow(query, o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, o.DateCreated, o.OofShard, idDelivery, idPayment).Scan(&uid)
	return uid, err
}
func (r *Repository) createItem(tx *sql.Tx, uid string, i Item) error {
	query := `
		INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := tx.Exec(query, i.ChrtID, i.TrackNumber, i.Price, i.Rid, i.Name, i.Sale, i.Size, i.TotalPrice, i.NmID, i.Brand, i.Status, uid)
	return err
}

func (r *Repository) GetAll() ([]Order, error) {
	query := `
		SELECT
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email, p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders as o
		JOIN deliveries AS d ON o.delivery_id = d.id
		JOIN payments AS p ON o.payment_id = p.id		
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[string]Order)

	for rows.Next() {
		var o Order

		err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard, &o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email, &o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		o.Items = []Item{}
		ordersMap[o.OrderUID] = o
	}

	itemsQuery := `
		SELECT
			order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM 
			items
	`
	itemRows, err := r.db.Query(itemsQuery)

	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var i Item
		var orderUID string
		err := itemRows.Scan(
			&orderUID, &i.ChrtID, &i.TrackNumber, &i.Price, &i.Rid, &i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NmID, &i.Brand, &i.Status,
		)
		if err != nil {
			return nil, err
		}
		if order, ok := ordersMap[orderUID]; ok {
			order.Items = append(order.Items, i)
			ordersMap[orderUID] = order

		}
	}
	result := make([]Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		result = append(result, order)
	}

	return result, nil
}
