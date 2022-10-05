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
	New(order Info) (int, error)
	GetUserOrders(userId int, createdAt string) ([]FullInfo, error)
	GetOrderById(orderId int) (FullInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error)
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
		if p.Product.Price != inputPrice {
			return 0, fmt.Errorf("invalid product price")
		}
		if p.Product.Price*float64(item.Quantity) != item.PriceForQuantity {
			return 0, fmt.Errorf("price for quantity mismatch")
		}
		sum += p.Product.Price * float64(item.Quantity)
	}
	return sum, nil
}

func (o *service) New(order Info) (int, error) {
	// TODO
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

func (o *service) GetUserOrders(userId int, createdAt string) ([]FullInfo, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *service) GetOrderById(orderId int) (FullInfo, error) {
	return o.repo.GetOrderById(orderId)
}

func (o *service) GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error) {
	return o.repo.GetAdminOrders(status, lastOrderCreatedAt)
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
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	doc, err := goinvoice.NewWithCyrillic(pwd)
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

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
	doc.SetCustomer(&goinvoice.Customer{
		LastName:     "Жовтанюк",
		FirstName:    "Максим",
		MiddleName:   "В'ячеславович",
		PhoneNumber:  "+380(68) 306-29-75",
		Email:        "zhovtyjshady@gmail.com",
		DeliveryType: "Доставка новою поштою",
		//DeliveryInfo: "м.Рівне, відділення №12",
	})
	for i := 0; i < 15; i++ {
		doc.AppendProductItem(&goinvoice.Product{
			Title:     "Стакан одноразовий Крафт 180мл",
			Price:     8.5,
			Quantity:  100,
			Total:     850,
			Packaging: "шт",
		})
	}

	pdf, err := doc.BuildPdf()
	if err != nil {
		fmt.Println(err.Error())
		return gofpdf.Fpdf{}, err
	}

	return pdf, err
}
