package order

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	New(order Info) (uuid.UUID, error)
	GetUserOrders(userId int, createdAt string) ([]FullInfo, error)
	GetOrderById(orderId uuid.UUID) (FullInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error)
	GetDeliveryTypes() ([]server.DeliveryType, error)
	GetPaymentTypes() ([]server.PaymentType, error)
	ProcessedOrder(orderId uuid.UUID) error
}

type storage struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *storage {
	return &storage{db: db}
}

var orderInfoColumnsInsert = []string{
	"user_id",
	"user_lastname",
	"user_firstname",
	"user_middle_name",
	"user_phone_number",
	"user_email",
	"order_comment",
	"order_sum_price",
	"delivery_type_id",
	"payment_type_id",
}

func (o *storage) New(order Info) (uuid.UUID, error) {
	tx, _ := o.db.Begin()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var deliveryTypeId int
	queryGetDeliveryId := fmt.Sprintf("SELECT id FROM %s WHERE delivery_type_title=$1", postgres.DeliveryTypesTable)
	err := o.db.Get(&deliveryTypeId, queryGetDeliveryId, order.Order.DeliveryTypeTitle)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed to create order, delivery type not found %s, error: %v", order.Order.DeliveryTypeTitle, err)
	}

	var paymentTypeId int
	queryGetPaymentTypeId := fmt.Sprintf("SELECT id FROM %s WHERE payment_type_title=$1", postgres.PaymentTypesTable)
	err = o.db.Get(&paymentTypeId, queryGetPaymentTypeId, order.Order.PaymentTypeTitle)
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed to create order, payment type not found %s, error: %v", order.Order.PaymentTypeTitle, err)
	}

	queryInsertOrder := psql.Insert(postgres.OrdersTable).Columns(orderInfoColumnsInsert...)

	queryInsertOrder = queryInsertOrder.Values(
		order.Order.UserId,
		order.Order.UserLastName,
		order.Order.UserFirstName,
		order.Order.UserMiddleName,
		order.Order.UserPhoneNumber,
		order.Order.UserEmail,
		order.Order.OrderComment,
		order.Order.OrderSumPrice,
		deliveryTypeId,
		paymentTypeId,
	)

	queryInsertOrderSql, args, err := queryInsertOrder.ToSql()
	if err != nil {
		return [16]byte{}, fmt.Errorf("failed to build sql query to insert order due to: %v", err)
	}

	var orderId uuid.UUID
	row := tx.QueryRow(queryInsertOrderSql+" RETURNING id", args...)
	if err = row.Scan(&orderId); err != nil {
		_ = tx.Rollback()
		return [16]byte{}, fmt.Errorf("failed to insert new order into table due to: %v", err)
	}

	for _, product := range order.Products {
		queryInsertProducts, args, err := psql.Insert(postgres.OrdersProductsTable).
			Columns("order_id", "product_id", "quantity", "price_for_quantity").
			Values(orderId, product.ProductId, product.Quantity, product.PriceForQuantity).
			ToSql()
		if err != nil {
			return [16]byte{}, err
		}
		_, err = tx.Exec(queryInsertProducts, args...)
		if err != nil {
			_ = tx.Rollback()
			return [16]byte{}, err
		}
	}

	for _, delivery := range order.Delivery {
		queryInsertDelivery, args, err := psql.Insert(postgres.OrdersDeliveryTable).
			Columns("order_id", "delivery_title", "delivery_description").
			Values(orderId, delivery.DeliveryTitle, delivery.DeliveryDescription).ToSql()

		_, err = tx.Exec(queryInsertDelivery, args...)
		if err != nil {
			_ = tx.Rollback()
			return [16]byte{}, err
		}
	}

	if order.Order.UserId != nil {
		queryDeleteCartProducts, args, err := psql.Delete(postgres.CartsProductsTable).Where(sq.Eq{"cart_id": order.Order.UserId}).ToSql()
		if err != nil {
			_ = tx.Rollback()
			return [16]byte{}, err
		}
		_, err = tx.Exec(queryDeleteCartProducts, args...)
		if err != nil {
			_ = tx.Rollback()
			return [16]byte{}, err
		}
	}

	return orderId, tx.Commit()
}

func (o *storage) GetUserOrders(userId int, createdAt string) ([]FullInfo, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var ordersAmount int
	queryOrdersAmount, args, err := psql.Select("count(*)").From(postgres.OrdersTable).Where(sq.Eq{"user_id": userId}).ToSql()
	if err != nil {
		return nil, err
	}

	if err = o.db.Get(&ordersAmount, queryOrdersAmount, args...); err != nil {
		return nil, err
	}

	ordersLimit := 12
	if ordersAmount <= 12 {
		ordersLimit = ordersAmount
	}

	orders := make([]FullInfo, ordersLimit)

	query := psql.Select(
		"order.id",
		"order.user_id",
		"order.user_lastname",
		"order.user_firstname",
		"order.user_middle_name",
		"order.user_phone_number",
		"order.user_email",
		"order.order_status",
		"order.order_comment",
		"order.order_sum_price",
		"delivery_types.delivery_type_title",
		"payment_types.payment_type_title",
		"order.created_at",
		"order.closed_at",
	).
		From(postgres.OrdersTable).
		LeftJoin(postgres.DeliveryTypesTable + " ON order.delivery_type_id=delivery_types.id").
		LeftJoin(postgres.PaymentTypesTable + " ON order.payment_type_id=payment_types.id").
		Where(sq.Eq{"user_id": userId})

	if createdAt != "" {
		query = query.Where(sq.Lt{"created_at": createdAt})
	}

	ordered := query.OrderBy("order.created_at DESC").Limit(12)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	for i := 0; i < ordersLimit; i++ {
		err = o.db.Get(&orders[i].Info, querySql, args...)
		if err != nil {
			return nil, err
		}
	}

	// TODO "message": "sql: Scan error on column index 1, name \"user_id\": converting NULL to int is unsupported"
	for i := 0; i < ordersLimit; i++ {
		queryOrderProducts, args, err := psql.
			Select(
				"id",
				"order_id",
				"article",
				"product_title",
				"img_url",
				"amount_in_stock",
				"price",
				"units_in_package",
				"packages_in_box",
				"created_at",
				"quantity",
				"price_for_quantity",
			).
			From(postgres.OrdersProductsTable).
			LeftJoin(postgres.ProductsTable + " ON orders_products.product_id=products.id").
			Where(sq.Eq{"orders_products.order_id": orders[i].Info.Id}).
			ToSql()
		if err != nil {
			return nil, err
		}

		err = o.db.Select(&orders[i].Products, queryOrderProducts, args...)
		if err != nil {
			return nil, err
		}

		queryOrderDelivery, args, err := psql.
			Select("*").
			From(postgres.OrdersDeliveryTable).
			Where(sq.Eq{"order_id": orders[i].Info.Id}).
			ToSql()
		if err != nil {
			return nil, err
		}

		err = o.db.Select(&orders[i].Delivery, queryOrderDelivery, args...)
		if err != nil {
			return nil, err
		}
	}

	return orders, err
}

func (o *storage) GetOrderById(orderId uuid.UUID) (FullInfo, error) {
	var order FullInfo

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryOrderInfoSql, args, err := psql.
		Select(
			"orders.id",
			"orders.user_id",
			"orders.user_lastname",
			"orders.user_firstname",
			"orders.user_middle_name",
			"orders.user_phone_number",
			"orders.user_email",
			"orders.order_status",
			"orders.order_comment",
			"orders.order_sum_price",
			"delivery_types.delivery_type_title",
			"payment_types.payment_type_title",
			"orders.created_at",
			"orders.closed_at",
		).
		From(postgres.OrdersTable).
		LeftJoin(postgres.DeliveryTypesTable + " ON orders.delivery_type_id=delivery_types.id").
		LeftJoin(postgres.PaymentTypesTable + " ON orders.payment_type_id=payment_types.id").
		Where(sq.Eq{"orders.id": orderId}).
		ToSql()

	err = o.db.Get(&order.Info, queryOrderInfoSql, args...)
	if err != nil {
		return FullInfo{}, fmt.Errorf("failed to get order info due to: %v", err)
	}

	queryProductsSql, args, err := psql.
		Select(
			"orders_products.quantity",
			"orders_products.price_for_quantity",
			"products.id",
			"products.article",
			"products.product_title",
			"products.img_url",
			"products.amount_in_stock",
			"products.price",
			"products.units_in_package",
			"products.packages_in_box",
			"products.created_at",
		).
		From(postgres.OrdersProductsTable).
		LeftJoin(postgres.ProductsTable + " ON products.id=orders_products.product_id").
		Where(sq.Eq{"order_id": orderId}).
		ToSql()

	err = o.db.Select(&order.Products, queryProductsSql, args...)
	if err != nil {
		return FullInfo{}, fmt.Errorf("failed to get order products due to: %v", err)
	}

	queryDeliverySql, args, err := psql.
		Select("order_id, delivery_title, delivery_description").
		From(postgres.OrdersDeliveryTable).
		Where(sq.Eq{"order_id": orderId}).ToSql()

	err = o.db.Select(&order.Delivery, queryDeliverySql, args...)
	if err != nil {
		return FullInfo{}, fmt.Errorf("failed to get order delivery info due to: %v", err)
	}

	return order, err
}

func (o *storage) GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error) {
	var orders []Order

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryOrders := psql.Select(
		"order.id",
		"order.user_lastname",
		"order.user_firstname",
		"order.user_middle_name",
		"order.user_phone_number",
		"order.user_email",
		"order.order_status",
		"order.order_comment",
		"order.order_sum_price",
		"delivery_types.delivery_type_title",
		"payment_types.payment_type_title",
		"order.created_at",
		"order.closed_at",
	).
		From(postgres.OrdersTable).
		LeftJoin(postgres.DeliveryTypesTable + " ON order.delivery_type_id=delivery_types.id").
		LeftJoin(postgres.PaymentTypesTable + " ON order.payment_type_id=payment_types.id")

	if status != "" {
		queryOrders = queryOrders.Where(sq.Eq{"order.order_status": status})
	}

	if lastOrderCreatedAt != "" {
		queryOrders = queryOrders.Where(sq.Lt{"order.created_at": lastOrderCreatedAt})
	}

	queryOrders = queryOrders.OrderBy("order.created_at DESC").Limit(12)

	queryOrdersSql, args, err := queryOrders.ToSql()
	if err != nil {
		return nil, err
	}

	err = o.db.Select(&orders, queryOrdersSql, args...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *storage) GetDeliveryTypes() (deliveryTypes []server.DeliveryType, err error) {
	queryGetDeliveryTypes := fmt.Sprintf("SELECT * FROM %s", postgres.DeliveryTypesTable)

	err = o.db.Select(&deliveryTypes, queryGetDeliveryTypes)
	if err != nil {
		return nil, err
	}

	return deliveryTypes, err
}

func (o *storage) GetPaymentTypes() (paymentTypes []server.PaymentType, err error) {
	queryGetPaymentTypes := fmt.Sprintf("SELECT * FROM %s", postgres.PaymentTypesTable)

	err = o.db.Select(&paymentTypes, queryGetPaymentTypes)
	if err != nil {
		return nil, err
	}

	return paymentTypes, err
}

func (o *storage) ProcessedOrder(orderId uuid.UUID) error {
	tx, _ := o.db.Begin()

	queryUpdateStatus := fmt.Sprintf("UPDATE %s SET order_status=$1 WHERE id=$2", postgres.OrdersTable)
	_, err := tx.Exec(queryUpdateStatus, "PROCESSED", orderId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update order status in database due to: %v", err)
	}

	order, err := o.GetOrderById(orderId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get order by its id due to: %v", err)
	}

	// TODO case when amount_in_stock is less than needed
	queryUpdateAmount := fmt.Sprintf(
		`
		UPDATE %s
		SET
			amount_in_stock = amount_in_stock - $1,
			current_spend = current_spend + $2
		WHERE id = $3
		`,
		postgres.ProductsTable,
	)

	for _, p := range order.Products {
		_, err = tx.Exec(queryUpdateAmount, p.Quantity, p.Quantity, p.Id)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update product amount due to %v: ", err)
		}
	}

	return tx.Commit()
}

//func (o *storage) ChangeOrderStatus(orderId uuid.UUID, toStatus string) error {
//	queryUpdateStatus := fmt.Sprintf("UPDATE %s SET order_status=$1 WHERE id=$2", postgres.OrdersTable)
//
//	_, err := o.db.Exec(queryUpdateStatus, toStatus, orderId)
//	if err != nil {
//		return fmt.Errorf("failed to update order status in database due to: %v", err)
//	}
//
//	return nil
//}
