package supply

import (
	"fmt"
	"math"
)

type Service interface {
	New(supply Supply) error
	Update() error
	Delete(id int) error
	GetAll(createdAt string) ([]InfoDTO, error)
	Products(id int) ([]ProductDTO, error)
}

type service struct {
	repo Storage
}

func NewSupplyService(repo Storage) Service {
	return &service{repo: repo}
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func (s *service) New(supply Supply) error {
	var productsSum float64
	for _, product := range supply.Products {
		sumWithoutTax := product.Amount * product.PriceForUnit
		if sumWithoutTax != product.SumWithoutTax {
			return fmt.Errorf("invalid sum without tax, %f != %f", sumWithoutTax, product.SumWithoutTax)
		}

		if product.Tax > 100 || product.Tax < 0 {
			return fmt.Errorf("invalid tax value")
		}
		tax := 1 + (product.Tax / 100)
		totalSum := sumWithoutTax * tax

		if !almostEqual(totalSum, product.TotalSum) {
			return fmt.Errorf("invalid total sum with tax, %f != %f", totalSum, product.TotalSum)
		}

		productsSum += product.TotalSum
	}

	var paymentsSum float64
	for _, payment := range supply.Payment {
		paymentsSum += payment.PaymentSum
	}

	if productsSum != paymentsSum {
		return fmt.Errorf("all products sum and payment accounts sum doesn't match, %f != %f", productsSum, paymentsSum)
	}

	supply.Info.Sum = productsSum

	err := s.repo.New(supply)
	if err != nil {
		return err
	}

	err = s.repo.UpdateProductsAmount(supply.Products, "+")
	if err != nil {
		return fmt.Errorf("failed to update amount_in_stock, err: %v", err)
	}

	return nil
}

func (s *service) Update() error {
	return nil
}

func (s *service) Delete(id int) error {
	products, err := s.repo.Products(id)
	if err != nil {
		return fmt.Errorf("failed to get supply products from db, %v", err)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete supply from db, %v", err)
	}

	// delete amount from products
	err = s.repo.UpdateProductsAmount(products, "-")
	if err != nil {
		return fmt.Errorf("failed to update product amount after deleting supply, %v", err)
	}

	return err
}

func (s *service) GetAll(createdAt string) ([]InfoDTO, error) {
	return s.repo.GetAll(createdAt)
}

func (s *service) Products(id int) ([]ProductDTO, error) {
	return s.repo.Products(id)
}
