package service

import (
	"fmt"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type SupplyService struct {
	repo repository.Supply
}

func NewSupplyService(repo repository.Supply) *SupplyService {
	return &SupplyService{repo: repo}
}

func (s *SupplyService) New(supply models.SupplyDTO) error {
	err := s.repo.New(supply)
	if err != nil {
		return err
	}

	err = s.repo.UpdateProductsAmount(supply.Products)
	if err != nil {
		return fmt.Errorf("failed to update amount_in_stock, err: %v", err)
	}

	return nil
}

func (s *SupplyService) Update() error {
	return nil
}

func (s *SupplyService) Delete(id int) error {
	return nil
}

func (s *SupplyService) GetAll(createdAt string) ([]models.SupplyInfoDTO, error) {
	return s.repo.GetAll(createdAt)
}
