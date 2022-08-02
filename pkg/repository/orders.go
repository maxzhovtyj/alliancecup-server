package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (o *OrdersPostgres) New(order server.OrderFullInfo) (uuid.UUID, error) {
	// todo

	tx, _ := o.db.Begin()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var deliveryTypeId int
	queryGetDeliveryId := fmt.Sprintf("SELECT id FROM %s WHERE delivery_type_title=$1", deliveryTypesTable)
	err := o.db.Get(&deliveryTypeId, queryGetDeliveryId, order.Info.DeliveryTypeTitle)
	if err != nil {
		return [16]byte{}, err
	}

	var paymentTypeId int
	queryGetPaymentTypeId := fmt.Sprintf("SELECT id FROM %s WHERE payment_type_title=$1", paymentTypesTable)
	err = o.db.Get(&paymentTypeId, queryGetPaymentTypeId, order.Info.PaymentTypeTitle)
	if err != nil {
		return [16]byte{}, err
	}

	queryInsertOrder, args, err := psql.Insert(ordersTable).Columns(
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
	).Values(
		order.Info.UserId,
		order.Info.UserLastName,
		order.Info.UserFirstName,
		order.Info.UserMiddleName,
		order.Info.UserPhoneNumber,
		order.Info.UserEmail,
		order.Info.OrderComment,
		order.Info.OrderSumPrice,
		deliveryTypeId,
		paymentTypeId,
	).ToSql()

	var orderId uuid.UUID
	row := tx.QueryRow(queryInsertOrder+"RETURNING id", args...)
	if err = row.Scan(&orderId); err != nil {
		_ = tx.Rollback()
		return [16]byte{}, err
	}

	for _, product := range order.Products {
		queryInsertProducts, args, err := psql.Insert(ordersProductsTable).
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
		queryInsertDelivery, args, err := psql.Insert(ordersDeliveryTable).
			Columns("order_id", "delivery_title", "delivery_description").
			Values(orderId, delivery.DeliveryTitle, delivery.DeliveryDescription).ToSql()

		_, err = tx.Exec(queryInsertDelivery, args...)
		if err != nil {
			_ = tx.Rollback()
			return [16]byte{}, err
		}
	}

	if order.Info.UserId != 0 {
		queryDeleteCartProducts, args, err := psql.Delete(cartsProductsTable).Where(sq.Eq{"cart_id": order.Info.UserId}).ToSql()
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

func (o *OrdersPostgres) GetUserOrders(userId int, createdAt string) ([]server.Order, error) {
	var orders []server.Order
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(
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
		From(ordersTable).
		LeftJoin(deliveryTypesTable + " ON orders.delivery_type_id=delivery_types.id").
		LeftJoin(paymentTypesTable + " ON orders.payment_type_id=payment_types.id").
		Where(sq.Eq{"user_id": userId})

	if createdAt != "" {
		query = query.Where(sq.Lt{"created_at": createdAt})
	}

	ordered := query.OrderBy("orders.created_at DESC").Limit(12)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	err = o.db.Select(&orders, querySql, args...)
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (o *OrdersPostgres) GetOrderById(orderId uuid.UUID) (server.OrderInfo, error) {
	var order server.OrderInfo

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryOrderInfo := psql.
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
		From(ordersTable).
		LeftJoin(deliveryTypesTable + " ON orders.delivery_type_id=delivery_types.id").
		LeftJoin(paymentTypesTable + " ON orders.payment_type_id=payment_types.id").
		Where(sq.Eq{"orders.id": orderId})

	queryOrderInfoSql, args, err := queryOrderInfo.ToSql()
	if err != nil {
		return server.OrderInfo{}, err
	}

	err = o.db.Get(&order.Info, queryOrderInfoSql, args...)
	if err != nil {
		return server.OrderInfo{}, err
	}

	queryProducts := psql.
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
		From(ordersProductsTable).
		LeftJoin(productsTable + " ON products.id=orders_products.product_id").
		Where(sq.Eq{"order_id": orderId})

	queryProductsSql, args, err := queryProducts.ToSql()
	if err != nil {
		return server.OrderInfo{}, err
	}

	err = o.db.Select(&order.Products, queryProductsSql, args...)
	if err != nil {
		return server.OrderInfo{}, err
	}

	queryDelivery := psql.Select("*").From(ordersDeliveryTable).Where(sq.Eq{"order_id": orderId})
	queryDeliverySql, args, err := queryDelivery.ToSql()
	if err != nil {
		return server.OrderInfo{}, err
	}

	err = o.db.Select(&order.Delivery, queryDeliverySql, args...)
	if err != nil {
		return server.OrderInfo{}, err
	}

	return order, err
}

func (o *OrdersPostgres) GetAdminOrders(status string, lastOrderCreatedAt string) ([]server.Order, error) {
	var orders []server.Order

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryOrders := psql.Select(
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
		From(ordersTable).
		LeftJoin(deliveryTypesTable + " ON orders.delivery_type_id=delivery_types.id").
		LeftJoin(paymentTypesTable + " ON orders.payment_type_id=payment_types.id").
		Where(sq.Eq{"orders.order_status": status})

	if lastOrderCreatedAt != "" {
		queryOrders = queryOrders.Where(sq.Lt{"orders.created_at": lastOrderCreatedAt})
	}

	queryOrders = queryOrders.OrderBy("orders.created_at DESC").Limit(12)

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
