CREATE TABLE deliveries ( 
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    address VARCHAR(50) NOT NULL,
    region VARCHAR(50) NOT NULL,
	email VARCHAR(50) NOT NULL
);


CREATE TABLE payments (
	id SERIAL PRIMARY KEY,
   	"transaction" VARCHAR(50) NOT NULL,
    request_id VARCHAR(50),
    currency VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(30) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT
);

CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
	order_uid VARCHAR(50) UNIQUE NOT NULL,
  	track_number VARCHAR(50) NOT NULL,
  	entry VARCHAR(50) NOT NULL,
	locale VARCHAR(5) NOT NULL,
  	internal_signature VARCHAR(50),
  	customer_id VARCHAR(50) NOT NULL,
  	delivery_service VARCHAR(50) NOT NULL,
  	shardkey VARCHAR(10) NOT NULL,
  	sm_id INT NOT NULL,
  	date_created TIMESTAMP NOT NULL,
  	oof_shard VARCHAR(10) NOT NULL,

	delivery_id INT UNIQUE,
	payment_id INT UNIQUE,

	CONSTRAINT fk_delivery FOREIGN KEY(delivery_id) REFERENCES deliveries(id),
	CONSTRAINT fk_payment FOREIGN KEY(payment_id) REFERENCES payments(id)
);

CREATE TABLE items (
	id SERIAL PRIMARY KEY,
	chrt_id INT NOT NULL,
	track_number VARCHAR(50) NOT NULL,
	price INT NOT NULL,
	rid VARCHAR(50) NOT NULL,
	name VARCHAR(50) NOT NULL,
	sale INT NOT NULL,
	size VARCHAR(10) NOT NULL,
	total_price INT NOT NULL,
	nm_id INT NOT NULL,
	brand VARCHAR(50) NOT NULL,
	status INT NOT NULL,

	order_uid VARCHAR(50) NOT NULL,

	CONSTRAINT fk_order FOREIGN KEY(order_uid) REFERENCES orders(order_uid)
);