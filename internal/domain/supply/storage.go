package supply

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	New(supply Supply) error
	GetAll(createdAt string) ([]InfoDTO, error)
	UpdateProductsAmount(products []ProductDTO, operation string) error
	DeleteAndGetProducts(id int) ([]ProductDTO, error)
}

type storage struct {
	db *sqlx.DB
}

func NewSupplyPostgres(db *sqlx.DB) *storage {
	return &storage{db: db}
}

func (s *storage) GetAll(createdAt string) ([]InfoDTO, error) {
	var supply []InfoDTO

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	querySelectInfo := psql.Select("*").From(postgres.SupplyTable)
	if createdAt != "" {
		querySelectInfo = querySelectInfo.Where(sq.Lt{"created_at": createdAt})
	}
	querySelectInfo = querySelectInfo.OrderBy("created_at DESC").Limit(12)

	querySelectInfoSql, args, err := querySelectInfo.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to get all supply, err: %v", err)
	}

	err = s.db.Select(&supply, querySelectInfoSql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select supply from db, err: %v", err)
	}

	return supply, nil
}

func (s *storage) New(supply Supply) error {
	tx, _ := s.db.Begin()

	var supplyId int
	queryInsetSupplyInfo := fmt.Sprintf(
		"INSERT INTO %s (supplier, supply_time, sum, comment) values ($1, $2, $3, $4) RETURNING id",
		postgres.SupplyTable,
	)
	row := tx.QueryRow(
		queryInsetSupplyInfo,
		supply.Info.Supplier,
		supply.Info.SupplyTime,
		supply.Info.Sum,
		supply.Info.Comment,
	)
	if err := row.Scan(&supplyId); err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, payment := range supply.Payment {
		queryInsertPayment := fmt.Sprintf(
			"INSERT INTO %s (supply_id, payment_account, payment_time, payment_sum) values ($1, $2, $3, $4)",
			postgres.SupplyPaymentTable,
		)

		_, err := tx.Exec(
			queryInsertPayment,
			supplyId,
			payment.PaymentAccount,
			payment.PaymentTime,
			payment.PaymentSum,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	for _, p := range supply.Products {
		queryInsertProduct := fmt.Sprintf(
			"INSERT INTO %s (supply_id, product_id, packaging, amount, price_for_unit, sum_without_tax, tax, total_sum) values ($1, $2, $3, $4, $5, $6, $7, $8)",
			postgres.SupplyProductsTable,
		)

		_, err := tx.Exec(
			queryInsertProduct,
			supplyId,
			p.ProductId,
			p.Packaging,
			p.Amount,
			p.PriceForUnit,
			p.SumWithoutTax,
			p.Tax,
			p.TotalSum,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *storage) UpdateProductsAmount(products []ProductDTO, operation string) error {
	tx, _ := s.db.Begin()

	// TODO check if amount_in_stock is less than amount to delete

	queryUpdateAmount := fmt.Sprintf(
		`
		UPDATE %s 
		SET 
			amount_in_stock = amount_in_stock %s $1,
			current_supply = current_supply %s $2 
		WHERE id = $3
		`,
		postgres.ProductsTable,
		operation,
		operation,
	)

	for _, p := range products {
		_, err := tx.Exec(queryUpdateAmount, p.Amount, p.Amount, p.ProductId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *storage) DeleteAndGetProducts(id int) ([]ProductDTO, error) {
	var products []ProductDTO
	queryGetProducts := fmt.Sprintf("SELECT * FROM %s WHERE supply_id=$1", postgres.SupplyProductsTable)

	err := s.db.Select(&products, queryGetProducts, id)
	if err != nil {
		return nil, err
	}

	queryDeleteSupply := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.SupplyTable)
	_, err = s.db.Exec(queryDeleteSupply, id)
	if err != nil {
		return nil, err
	}

	return products, nil
}
