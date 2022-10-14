package order

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	goinvoice "github.com/maxzhovtyj/go-invoice"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"os"
)

type Service interface {
	New(order CreateDTO) (int, error)
	GetUserOrders(userId int, createdAt string) ([]SelectDTO, error)
	GetOrderById(orderId int) (SelectDTO, error)
	AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error)
	DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId int) error
	GetInvoice(orderId int) (gofpdf.Fpdf, error)
}

type service struct {
	repo        Storage
	productRepo product.Storage
}

func NewOrdersService(repo Storage, productRepo product.Storage) Service {
	return &service{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (o *service) OrderSumCount(products []Product) (float64, error) {
	var sum float64
	for _, item := range products {
		p, err := o.productRepo.GetProductById(item.ProductId)
		if err != nil {
			return 0, err
		}
		inputPrice := item.PriceForQuantity / float64(item.Quantity)
		if p.Price != inputPrice {
			return 0, fmt.Errorf("invalid product price")
		}
		if p.Price*float64(item.Quantity) != item.PriceForQuantity {
			return 0, fmt.Errorf("price for quantity mismatch")
		}
		sum += p.Price * float64(item.Quantity)
	}
	return sum, nil
}

func (o *service) New(order CreateDTO) (int, error) {
	sum, err := o.OrderSumCount(order.Products)
	if err != nil {
		return 0, err
	}
	if sum != order.Order.OrderSumPrice {
		return 0, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Order.OrderSumPrice)
	}

	id, err := o.repo.New(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (o *service) GetUserOrders(userId int, createdAt string) ([]SelectDTO, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *service) GetOrderById(orderId int) (SelectDTO, error) {
	return o.repo.GetOrderById(orderId)
}

func (o *service) AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error) {
	return o.repo.AdminGetOrders(status, lastOrderCreatedAt, search)
}

func (o *service) DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error) {
	deliveryTypes, err := o.repo.GetDeliveryTypes()
	if err != nil {
		return server.DeliveryPaymentTypes{}, fmt.Errorf("failed to load delivery types due to: %v", err)
	}

	paymentTypes, err := o.repo.GetPaymentTypes()
	if err != nil {
		return server.DeliveryPaymentTypes{}, fmt.Errorf("failed to load payment types due to: %v", err)
	}

	return server.DeliveryPaymentTypes{
		DeliveryTypes: deliveryTypes,
		PaymentTypes:  paymentTypes,
	}, nil
}

func (o *service) ProcessedOrder(orderId int) error {
	err := o.repo.ProcessedOrder(orderId)
	if err != nil {
		return err
	}

	return nil
}

func (o *service) GetInvoice(orderId int) (gofpdf.Fpdf, error) {
	order, err := o.repo.GetOrderById(orderId)
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	doc, err := goinvoice.NewWithCyrillic(pwd)
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	doc.Language = "UK"

	doc.SetPwd(pwd)
	doc.OrderId = fmt.Sprintf("%d", orderId)
	doc.SetInvoice(&goinvoice.Invoice{
		Logo:  nil,
		Title: "Alliance Cup",
	})
	doc.SetCompany(&goinvoice.Company{
		Name:        "AllianceCup",
		Address:     "Шухевича, 22",
		PhoneNumber: "+380(96) 512-15-16",
		City:        "Рівне",
		Country:     "Україна",
	})

	var orderDeliveryInfo string
	for i, d := range order.Delivery {
		orderDeliveryInfo += d.DeliveryDescription
		if len(order.Delivery)-1 != i {
			orderDeliveryInfo += ", "
		}
	}

	doc.SetCustomer(&goinvoice.Customer{
		LastName:     order.Info.UserLastName,
		FirstName:    order.Info.UserFirstName,
		MiddleName:   order.Info.UserMiddleName,
		PhoneNumber:  order.Info.UserPhoneNumber,
		Email:        order.Info.UserEmail,
		DeliveryType: order.Info.DeliveryTypeTitle,
		DeliveryInfo: orderDeliveryInfo,
	})

	for _, p := range order.Products {
		doc.AppendProductItem(&goinvoice.Product{
			Title:     p.ProductTitle,
			Price:     p.Price,
			Quantity:  float64(p.Quantity),
			Total:     p.PriceForQuantity,
			Packaging: "уп",
		})
	}

	pdf, err := doc.BuildPdf()
	if err != nil {
		fmt.Println(err.Error())
		return gofpdf.Fpdf{}, err
	}

	return pdf, err
}
