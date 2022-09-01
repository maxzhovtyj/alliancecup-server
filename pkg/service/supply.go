package service

import (
	"fmt"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
	"math"
)

type SupplyService struct {
	repo repository.Supply
}

func NewSupplyService(repo repository.Supply) *SupplyService {
	return &SupplyService{repo: repo}
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func (s *SupplyService) New(supply models.SupplyDTO) error {
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

func (s *SupplyService) Update() error {
	return nil
}

func (s *SupplyService) Delete(id int) error {
	// delete from supplies
	products, err := s.repo.DeleteAndGetProducts(id)

	// delete amount from products
	err = s.repo.UpdateProductsAmount(products, "-")
	if err != nil {
		return err
	}

	return nil
}

func (s *SupplyService) GetAll(createdAt string) ([]models.SupplyInfoDTO, error) {
	return s.repo.GetAll(createdAt)
}
