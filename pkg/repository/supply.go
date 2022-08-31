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
		return tx.Rollback()
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
			return tx.Rollback()
		}
	}

	for _, p := range supply.Products {
		queryInsertProduct := fmt.Sprintf(
			"INSERT INTO %s (supply_id, product_id, packaging, amount, price_for_unit, sum_without_tax, tax, total_sum) values ($1, $2, $3, $4, $5, $6, $7, $8)",
			productsSupplyTable,
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
			return tx.Rollback()
		}
	}

	return tx.Commit()
}

func (s *SupplyPostgres) UpdateProductsAmount(products []models.ProductSupplyDTO) error {
	tx, _ := s.db.Begin()
	for _, p := range products {
		queryUpdateAmount := fmt.Sprintf(
			"UPDATE %s SET amount_in_stock=amount_in_stock+$1 WHERE id=$2",
			productsTable,
		)
		_, err := tx.Exec(queryUpdateAmount, p.Amount, p.ProductId)
		if err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
