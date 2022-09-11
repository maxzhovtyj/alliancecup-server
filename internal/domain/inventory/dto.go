package inventory

import "time"

type ProductDTO struct {
	id              int
	productId       int
	inventoryId     int
	productTitle    string
	createdAt       time.Time
	initialAmount   float64 // amount from the last inventory
	supply          float64 // from the last inventory
	spends          float64 // spending (customers orders) from the last inventory
	writeOffAmount  float64 // something that wasn't sold
	writeOffPrice   float64 // price for write off amount
	plannedAmount   float64 // current amount in stock
	realAmount      float64 // inventory input
	realAmountPrice float64 // inventory input price
	difference      float64 // plannedAmount - realAmount
	differencePrice float64 // difference * product price
}
