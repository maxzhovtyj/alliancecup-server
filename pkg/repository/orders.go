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
		"id",
		"user_id",
		"user_lastname",
		"user_firstname",
		"user_middle_name",
		"user_phone_number",
		"user_email",
		"order_status",
		"order_comment",
		"order_sum_price",
		"delivery_type_title",
		"payment_type_title",
		"created_at",
		"closed_at",
	).From(ordersTable).Where(sq.Eq{"user_id": userId}) // todo

	if createdAt != "" {
		query = query.Where(sq.Lt{"created_at": createdAt})
	}

	ordered := query.OrderBy("products.created_at DESC").Limit(12)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	err = o.db.Select(orders, querySql, args...)
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (o *OrdersPostgres) NewUserOrder() {

}

func (o *OrdersPostgres) GetAllOrders() {

}
