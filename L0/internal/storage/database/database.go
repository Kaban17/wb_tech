package database

import (
	"context"
	"database/sql"
	"fmt"
	"orders/models"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Connect устанавливает соединение с PostgreSQL
func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	// Устанавливаем лимиты соединения
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// GetOrder получает полную информацию о заказе по его ID
func GetOrder(db *sql.DB, uid string) (*models.Order, error) {
	// Используем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Начинаем транзакцию для согласованности данных
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Основной запрос для получения заказа
	order, err := getOrderMainData(tx, uid)
	if err != nil {
		return nil, err
	}

	// Получаем данные доставки
	delivery, err := getDeliveryData(tx, uid)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	// Получаем данные платежа
	payment, err := getPaymentData(tx, uid)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	// Получаем товары
	items, err := getItemsData(tx, uid)
	if err != nil {
		return nil, err
	}
	order.Items = items

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return order, nil
}

// CreateOrder создает новый заказ в базе данных
func CreateOrder(db *sql.DB, order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Вставляем основные данные заказа
	if err := insertOrderMainData(tx, order); err != nil {
		return err
	}

	// Вставляем данные доставки
	if err := insertDeliveryData(tx, order); err != nil {
		return err
	}

	// Вставляем данные платежа
	if err := insertPaymentData(tx, order); err != nil {
		return err
	}

	// Вставляем товары
	if err := insertItemsData(tx, order); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
func GetAllOrders(db *sql.DB) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query to get all order IDs first (more efficient than loading everything at once)
	orderIDs, err := getAllOrderIDs(db, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order IDs: %v", err)
	}

	var orders []models.Order
	for _, id := range orderIDs {
		order, err := GetOrder(db, id)
		if err != nil {
			// You might want to log the error and continue, or return it
			return nil, fmt.Errorf("failed to get order %s: %v", id, err)
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func getAllOrderIDs(db *sql.DB, ctx context.Context) ([]string, error) {
	query := `SELECT order_uid FROM orders`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query order IDs: %v", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan order ID: %v", err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning order IDs: %v", err)
	}

	return ids, nil
}
func getOrderMainData(tx *sql.Tx, uid string) (*models.Order, error) {
	query := `
    SELECT
        order_uid, track_number, entry, locale,
        internal_signature, customer_id, delivery_service,
        shard_key, sm_id, date_created, oof_shard
    FROM orders
    WHERE order_uid = $1`

	var order models.Order
	err := tx.QueryRow(query, uid).Scan(
		&order.ID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SMID,
		&order.DataCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return &order, nil
}

func getDeliveryData(tx *sql.Tx, uid string) (*models.DeliveryInfo, error) {
	query := `
    SELECT
        name, phone, zip_code, city, address,
        region, email
    FROM delivery
    WHERE order_uid = $1`

	var delivery models.DeliveryInfo
	err := tx.QueryRow(query, uid).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.ZipCode,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery info: %v", err)
	}
	return &delivery, nil
}

func getPaymentData(tx *sql.Tx, uid string) (*models.PaymentInfo, error) {
	query := `
    SELECT
        transaction, request_id, currency, provider,
        amount, payment_dt, bank, delivery_cost,
        goods_total, custom_fee
    FROM payment
    WHERE transaction = $1`

	var payment models.PaymentInfo
	err := tx.QueryRow(query, uid).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment info: %v", err)
	}
	return &payment, nil
}

func getItemsData(tx *sql.Tx, uid string) ([]models.ItemInfo, error) {
	query := `
    SELECT
        chrt_id, track_number, price, rid, name,
        sale, size, total_price, nm_id, brand, status
    FROM items
    WHERE order_uid = $1`

	rows, err := tx.Query(query, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %v", err)
	}
	defer rows.Close()

	var items []models.ItemInfo
	for rows.Next() {
		var item models.ItemInfo
		if err := rows.Scan(
			&item.ChartID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %v", err)
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning items: %v", err)
	}
	return items, nil
}

func insertOrderMainData(tx *sql.Tx, order *models.Order) error {
	query := `
    INSERT INTO orders (
        order_uid, track_number, entry, locale,
        internal_signature, customer_id, delivery_service,
        shard_key, sm_id, date_created, oof_shard
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := tx.Exec(query,
		order.ID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DataCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}
	return nil
}

func insertDeliveryData(tx *sql.Tx, order *models.Order) error {
	query := `
    INSERT INTO delivery (
        order_uid, name, phone, zip_code, city,
        address, region, email
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(query,
		order.ID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.ZipCode,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %v", err)
	}
	return nil
}

func insertPaymentData(tx *sql.Tx, order *models.Order) error {
	query := `
    INSERT INTO payment (
        transaction, request_id, currency, provider,
        amount, payment_dt, bank, delivery_cost,
        goods_total, custom_fee
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := tx.Exec(query,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %v", err)
	}
	return nil
}

func insertItemsData(tx *sql.Tx, order *models.Order) error {
	query := `
    INSERT INTO items (
        order_uid, chrt_id, track_number, price, rid,
        name, sale, size, total_price, nm_id, brand, status
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		_, err := tx.Exec(query,
			order.ID,
			item.ChartID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %v", err)
		}
	}
	return nil
}
