package order

import (
	"fmt"
	generator "github.com/angelodlfrtr/go-invoice-generator"
	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"log"
)

type Service interface {
	New(order Info) (uuid.UUID, error)
	GetUserOrders(userId int, createdAt string) ([]FullInfo, error)
	GetOrderById(orderId uuid.UUID) (FullInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error)
	DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId uuid.UUID) error
	GetInvoice(id uuid.UUID) (*fpdf.Fpdf, error)
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

func (o *service) New(order Info) (uuid.UUID, error) {
	// TODO
	sum, err := o.OrderSumCount(order.Products)
	if err != nil {
		return [16]byte{}, err
	}
	if sum != order.Order.OrderSumPrice {
		return [16]byte{}, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Order.OrderSumPrice)
	}

	id, err := o.repo.New(order)
	if err != nil {
		return [16]byte{}, err
	}

	return id, nil
}

func (o *service) GetUserOrders(userId int, createdAt string) ([]FullInfo, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *service) GetOrderById(orderId uuid.UUID) (FullInfo, error) {
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

func (o *service) ProcessedOrder(orderId uuid.UUID) error {
	err := o.repo.ProcessedOrder(orderId)
	if err != nil {
		return err
	}

	return nil
}

func (o *service) GetInvoice(id uuid.UUID) (*fpdf.Fpdf, error) {
	order, err := o.GetOrderById(id)
	if err != nil {
		return nil, err
	}

	doc, err := generator.New(generator.Invoice, &generator.Options{
		AutoPrint:       true,
		TextTypeInvoice: "FACTURE",
	})
	if err != nil {
		return nil, err
	}

	doc.SetRef("testref")
	doc.SetVersion("someversion")

	doc.SetDescription("A description")
	doc.SetNotes("I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! ")

	doc.SetDate("02/03/2021")
	doc.SetPaymentTerm("02/04/2021")

	doc.SetCompany(&generator.Contact{
		Name: "Alliance Cup",
		Address: &generator.Address{
			Address:    "Шухевича, 22",
			Address2:   "+380(96) 612-15-16",
			PostalCode: "33018",
			City:       "Рівне",
			Country:    "Україна",
		},
	})

	var orderDelivery string
	for _, d := range order.Delivery {
		orderDelivery += d.DeliveryTitle + " - " + d.DeliveryDescription
	}

	doc.SetCustomer(&generator.Contact{
		Name: fmt.Sprintf("%s %s.%s", order.Info.UserLastName, order.Info.UserFirstName, order.Info.UserMiddleName),
		Address: &generator.Address{
			Address:    orderDelivery,
			PostalCode: order.Info.DeliveryTypeTitle,
		},
	})

	for i := 0; i < 3; i++ {
		doc.AppendItem(&generator.Item{
			Name:        "Cupcake ipsum dolor sit amet bonbon, coucou bonbon lala jojo, mama titi toto",
			Description: "Cupcake ipsum dolor sit amet bonbon, Cupcake ipsum dolor sit amet bonbon, Cupcake ipsum dolor sit amet bonbon",
			UnitCost:    "99876.89",
			Quantity:    "2",
			Tax: &generator.Tax{
				Percent: "20",
			},
		})
	}

	doc.AppendItem(&generator.Item{
		Name:     "Test",
		UnitCost: "99876.89",
		Quantity: "2",
		Tax: &generator.Tax{
			Amount: "89",
		},
		Discount: &generator.Discount{
			Percent: "30",
		},
	})

	doc.AppendItem(&generator.Item{
		Name:     "Test",
		UnitCost: "3576.89",
		Quantity: "2",
		Discount: &generator.Discount{
			Percent: "50",
		},
	})

	doc.AppendItem(&generator.Item{
		Name:     "Test",
		UnitCost: "889.89",
		Quantity: "2",
		Discount: &generator.Discount{
			Amount: "234.67",
		},
	})

	doc.SetDefaultTax(&generator.Tax{
		Percent: "10",
	})

	// doc.SetDiscount(&generator.Discount{
	// Percent: "90",
	// })
	doc.SetDiscount(&generator.Discount{
		Amount: "1340",
	})

	pdf, err := doc.Build()

	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	return pdf, err
}
