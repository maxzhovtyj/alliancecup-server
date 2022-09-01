package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
)

type SupplyPostgres struct {
	db *sqlx.DB
}

func NewSupplyPostgres(db *sqlx.DB) *SupplyPostgres {
	return &SupplyPostgres{db: db}
}

func (s *SupplyPostgres) GetAll(createdAt string) ([]models.SupplyInfoDTO, error) {
	var supply []models.SupplyInfoDTO

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	querySelectInfo := psql.Select("*").From(supplyTable)
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

func (s *SupplyPostgres) New(supply models.SupplyDTO) error {
	tx, _ := s.db.Begin()

	var supplyId int
	queryInsetSupplyInfo := fmt.Sprintf(
		"INSERT INTO %s (supplier, supply_time, comment) values ($1, $2, $3) RETURNING id",
		supplyTable,
	)
	row := tx.QueryRow(
		queryInsetSupplyInfo,
		supply.Info.Supplier,
		supply.Info.SupplyTime,
		supply.Info.Comment,
	)
	if err := row.Scan(&supplyId); err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, payment := range supply.Payment {
		queryInsertPayment := fmt.Sprintf(
			"INSERT INTO %s (supply_id, payment_account, payment_time, payment_sum) values ($1, $2, $3, $4)",
			supplyPaymentTable,
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
			supplyProductsTable,
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

func (s *SupplyPostgres) UpdateProductsAmount(products []models.ProductSupplyDTO, operation string) error {
	tx, _ := s.db.Begin()

	// TODO check if amount_in_stock is less than amount to delete
	//q := `
	//	DO $$
	//		DECLARE
	//			selected_product products%rowtype;
	//		BEGIN
	//		SELECT *
	//		FROM products
	//		INTO selected_product
	//		WHERE product_id=$1;
	//
	//		IF selected_product.amount_in_stock < $2 THEN
	//			UPDATE products SET amount_in_stock = 0;
	//		ELSE
	//			UPDATE products SET amount_in_stock = amount_in_stock-$2;
	//		END IF
	//	END $$
	//`

	for _, p := range products {
		queryUpdateAmount := fmt.Sprintf(
			"UPDATE %s SET amount_in_stock=amount_in_stock%s$1 WHERE id=$2",
			productsTable,
			operation,
		)

		_, err := tx.Exec(queryUpdateAmount, p.Amount, p.ProductId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SupplyPostgres) DeleteAndGetProducts(id int) ([]models.ProductSupplyDTO, error) {
	var products []models.ProductSupplyDTO
	queryGetProducts := fmt.Sprintf("SELECT * FROM %s WHERE supply_id=$1", supplyProductsTable)

	err := s.db.Select(&products, queryGetProducts, id)
	if err != nil {
		return nil, err
	}

	queryDeleteSupply := fmt.Sprintf("DELETE FROM %s WHERE id=$1", supplyTable)
	_, err = s.db.Exec(queryDeleteSupply, id)
	if err != nil {
		return nil, err
	}

	return products, nil
}
