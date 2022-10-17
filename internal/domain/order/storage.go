package order

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"strings"
)

type Storage interface {
	New(order CreateDTO) (int, error)
	GetUserOrders(userId int, createdAt string) ([]SelectDTO, error)
	GetOrderById(orderId int) (SelectDTO, error)
	AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error)
	GetDeliveryTypes() ([]server.DeliveryType, error)
	GetPaymentTypes() ([]server.PaymentType, error)
	ProcessedOrder(orderId int) error
}

type storage struct {
	qb sq.StatementBuilderType
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB, psql sq.StatementBuilderType) *storage {
	return &storage{
		db: db,
		qb: psql,
	}
}

var orderInfoColumnsInsert = []string{
	"user_id",
	"user_lastname",
	"user_firstname",
	"user_middle_name",
	"user_phone_number",
	"user_email",
	"comment",
	"sum_price",
	"delivery_type_id",
	"payment_type_id",
}

func (s *storage) New(order CreateDTO) (int, error) {
	tx, _ := s.db.Begin()

	var deliveryTypeId int
	queryGetDeliveryId := fmt.Sprintf("SELECT id FROM %s WHERE delivery_type_title=$1", postgres.DeliveryTypesTable)
	err := s.db.Get(&deliveryTypeId, queryGetDeliveryId, order.Order.DeliveryTypeTitle)
	if err != nil {
		return 0, fmt.Errorf("failed to create order, delivery type not found %s, error: %v", order.Order.DeliveryTypeTitle, err)
	}

	var paymentTypeId int
	queryGetPaymentTypeId := fmt.Sprintf("SELECT id FROM %s WHERE payment_type_title=$1", postgres.PaymentTypesTable)
	err = s.db.Get(&paymentTypeId, queryGetPaymentTypeId, order.Order.PaymentTypeTitle)
	if err != nil {
		return 0, fmt.Errorf("failed to create order, payment type not found %s, error: %v", order.Order.PaymentTypeTitle, err)
	}

	queryInsertOrder := s.qb.Insert(postgres.OrdersTable).Columns(orderInfoColumnsInsert...)

	queryInsertOrder = queryInsertOrder.Values(
		order.Order.UserId,
		order.Order.UserLastName,
		order.Order.UserFirstName,
		order.Order.UserMiddleName,
		order.Order.UserPhoneNumber,
		order.Order.UserEmail,
		order.Order.Comment,
		order.Order.SumPrice,
		deliveryTypeId, // TODO refactor to sub-query
		paymentTypeId,  // TODO refactor to sub-query
	)

	queryInsertOrderSql, args, err := queryInsertOrder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql query to insert order due to: %v", err)
	}

	var orderId int
	row := tx.QueryRow(queryInsertOrderSql+" RETURNING id", args...)
	if err = row.Scan(&orderId); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("failed to insert new order into table due to: %v", err)
	}

	for _, product := range order.Products {
		queryInsertProducts, args, err := s.qb.Insert(postgres.OrdersProductsTable).
			Columns("order_id", "product_id", "quantity", "price_for_quantity").
			Values(orderId, product.ProductId, product.Quantity, product.PriceForQuantity).
			ToSql()
		if err != nil {
			return 0, err
		}
		_, err = tx.Exec(queryInsertProducts, args...)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	for _, delivery := range order.Delivery {
		queryInsertDelivery, args, err := s.qb.Insert(postgres.OrdersDeliveryTable).
			Columns("order_id", "delivery_title", "delivery_description").
			Values(orderId, delivery.DeliveryTitle, delivery.DeliveryDescription).ToSql()

		_, err = tx.Exec(queryInsertDelivery, args...)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	if order.Order.UserId != nil {
		queryDeleteCartProducts, args, err := s.qb.Delete(postgres.CartsProductsTable).
			Where(sq.Eq{"cart_id": order.Order.UserId}).
			ToSql()
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		_, err = tx.Exec(queryDeleteCartProducts, args...)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return orderId, tx.Commit()
}

func (s *storage) GetUserOrders(userId int, createdAt string) ([]SelectDTO, error) {
	var ordersAmount int
	queryOrdersAmount, args, err := s.qb.Select("count(*)").
		From(postgres.OrdersTable).
		Where(sq.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = s.db.Get(&ordersAmount, queryOrdersAmount, args...); err != nil {
		return nil, err
	}

	ordersLimit := 12
	if ordersAmount <= 12 {
		ordersLimit = ordersAmount
	}

	orders := make([]SelectDTO, ordersLimit)

	query := s.qb.Select(
		"orders.id",
		"orders.user_id",
		"orders.user_lastname",
		"orders.user_firstname",
		"orders.user_middle_name",
		"orders.user_phone_number",
		"orders.user_email",
		"orders.status",
		"orders.comment",
		"orders.sum_price",
		"delivery_types.delivery_type_title",
		"payment_types.payment_type_title",
		"orders.created_at",
		"orders.closed_at",
	).
		From(postgres.OrdersTable).
		LeftJoin(postgres.DeliveryTypesTable + " ON orders.delivery_type_id = delivery_types.id").
		LeftJoin(postgres.PaymentTypesTable + " ON orders.payment_type_id = payment_types.id").
		Where(sq.Eq{"user_id": userId})

	if createdAt != "" {
		query = query.Where(sq.Lt{"orders.created_at": createdAt})
	}

	ordered := query.OrderBy("orders.created_at DESC").Limit(12)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	for i := 0; i < ordersLimit; i++ {
		err = s.db.Get(&orders[i].Info, querySql, args...)
		if err != nil {
			return nil, err
		}
	}

	// TODO "message": "sql: Scan error on column index 1, name \"user_id\": converting NULL to int is unsupported"
	for i := 0; i < ordersLimit; i++ {
		queryOrderProducts, args, err := s.qb.
			Select(
				"id",
				"order_id",
				"article",
				"product_title",
				"img_url",
				"amount_in_stock",
				"price",
				"packaging",
				"created_at",
				"quantity",
				"price_for_quantity",
			).
			From(postgres.OrdersProductsTable).
			LeftJoin(postgres.ProductsTable + " ON orders_products.product_id = products.id").
			Where(sq.Eq{"orders_products.order_id": orders[i].Info.Id}).
			ToSql()
		if err != nil {
			return nil, err
		}

		err = s.db.Select(&orders[i].Products, queryOrderProducts, args...)
		if err != nil {
			return nil, err
		}

		queryOrderDelivery, args, err := s.qb.
			Select("*").
			From(postgres.OrdersDeliveryTable).
			Where(sq.Eq{"order_id": orders[i].Info.Id}).
			ToSql()
		if err != nil {
			return nil, err
		}

		err = s.db.Select(&orders[i].Delivery, queryOrderDelivery, args...)
		if err != nil {
			return nil, err
		}
	}

	return orders, err
}

func (s *storage) GetOrderById(orderId int) (SelectDTO, error) {
	var order SelectDTO

	queryOrderInfoSql, args, err := s.qb.
		Select(
			"orders.id",
			"orders.user_id",
			"orders.user_lastname",
			"orders.user_firstname",
			"orders.user_middle_name",
			"orders.user_phone_number",
			"orders.user_email",
			"orders.status",
			"orders.comment",
			"orders.sum_price",
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

	err = s.db.Get(&order.Info, queryOrderInfoSql, args...)
	if err != nil {
		return SelectDTO{}, fmt.Errorf("failed to get order info due to: %v", err)
	}

	queryProductsSql, args, err := s.qb.
		Select(
			"orders_products.quantity",
			"orders_products.price_for_quantity",
			"products.id",
			"products.article",
			"products.product_title",
			"products.img_url",
			"products.amount_in_stock",
			"products.price",
			"products.packaging",
			"products.created_at",
		).
		From(postgres.OrdersProductsTable).
		LeftJoin(postgres.ProductsTable + " ON products.id=orders_products.product_id").
		Where(sq.Eq{"order_id": orderId}).
		ToSql()

	err = s.db.Select(&order.Products, queryProductsSql, args...)
	if err != nil {
		return SelectDTO{}, fmt.Errorf("failed to get order products due to: %v", err)
	}

	queryDeliverySql, args, err := s.qb.
		Select("order_id, delivery_title, delivery_description").
		From(postgres.OrdersDeliveryTable).
		Where(sq.Eq{"order_id": orderId}).ToSql()

	err = s.db.Select(&order.Delivery, queryDeliverySql, args...)
	if err != nil {
		return SelectDTO{}, fmt.Errorf("failed to get order delivery info due to: %v", err)
	}

	return order, err
}

func (s *storage) AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error) {
	var orders []Order

	queryOrders := s.qb.Select(
		"orders.id",
		"orders.user_lastname",
		"orders.user_firstname",
		"orders.user_middle_name",
		"orders.user_phone_number",
		"orders.user_email",
		"orders.status",
		"orders.comment",
		"orders.sum_price",
		"delivery_types.delivery_type_title",
		"payment_types.payment_type_title",
		"orders.created_at",
		"orders.closed_at",
	).
		From(postgres.OrdersTable).
		LeftJoin(postgres.DeliveryTypesTable + " ON orders.delivery_type_id = delivery_types.id").
		LeftJoin(postgres.PaymentTypesTable + " ON orders.payment_type_id = payment_types.id")

	if status != "" {
		queryOrders = queryOrders.Where(sq.Eq{"orders.order_status": status})
	}

	if lastOrderCreatedAt != "" {
		queryOrders = queryOrders.Where(sq.Lt{"orders.created_at": lastOrderCreatedAt})
	}

	if search != "" {
		searchName := fmt.Sprintf("LOWER(orders.user_lastname)")
		searchNameValue := "%" + strings.ToLower(search) + "%"
		queryOrders = queryOrders.Where(sq.Like{searchName: searchNameValue})
	}

	queryOrdersSql, args, err := queryOrders.OrderBy("orders.created_at DESC").Limit(12).ToSql()

	err = s.db.Select(&orders, queryOrdersSql, args...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *storage) GetDeliveryTypes() (deliveryTypes []server.DeliveryType, err error) {
	queryGetDeliveryTypes := fmt.Sprintf("SELECT * FROM %s", postgres.DeliveryTypesTable)

	err = s.db.Select(&deliveryTypes, queryGetDeliveryTypes)
	if err != nil {
		return nil, err
	}

	return deliveryTypes, err
}

func (s *storage) GetPaymentTypes() (paymentTypes []server.PaymentType, err error) {
	queryGetPaymentTypes := fmt.Sprintf("SELECT * FROM %s", postgres.PaymentTypesTable)

	err = s.db.Select(&paymentTypes, queryGetPaymentTypes)
	if err != nil {
		return nil, err
	}

	return paymentTypes, err
}

func (s *storage) ProcessedOrder(orderId int) error {
	tx, _ := s.db.Begin()

	queryUpdateStatus := fmt.Sprintf("UPDATE %s SET order_status=$1 WHERE id=$2", postgres.OrdersTable)
	_, err := tx.Exec(queryUpdateStatus, "PROCESSED", orderId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update order status in database due to: %v", err)
	}

	order, err := s.GetOrderById(orderId)
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

func (s *storage) ChangeOrderStatus(orderId int, toStatus string) error {
	queryUpdateStatus := fmt.Sprintf("UPDATE %s SET order_status=$1 WHERE id=$2", postgres.OrdersTable)

	_, err := s.db.Exec(queryUpdateStatus, toStatus, orderId)
	if err != nil {
		return fmt.Errorf("failed to update order status in database due to: %v", err)
	}

	return nil
}
