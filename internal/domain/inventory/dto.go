package inventory

import "time"

type ProductDTO struct {
	InventoryId     int     `json:"-" db:"inventory_id"`
	ProductId       int     `json:"productId" db:"product_id"`
	LastInventoryId *int    `json:"lastInventoryId" db:"last_inventory_id"`
	InitialAmount   float64 `json:"initialAmount" db:"initial_amount"` // amount from the last inventory
	Supply          float64 `json:"supply" db:"supply"`                // from the last inventory
	Spend           float64 `json:"spends" db:"spend"`                 // spending (customers orders) from the last inventory
	WriteOff        float64 `json:"writeOff" db:"write_off"`           // something that wasn't sold
	PlannedAmount   float64 `json:"plannedAmount" db:"planned_amount"` // current amount in stock
	RealAmount      float64 `json:"realAmount" db:"real_amount"`       // inventory input
	//WriteOffPrice   float64 // price for write off amount
	//realAmountPrice float64 // inventory input price
	//differencePrice float64 // difference * product price

	//inventory_id
	//product_id
	//last_inventory
	//initial_amount
	//supply
	//spend
	//write_off
	//write_off_price
	//planned_amount
	//difference
	//difference_price
}

type CurrentProductDTO struct {
	ProductId       int        `json:"productId" db:"id"`
	Title           string     `json:"title" db:"product_title"`
	Price           float64    `json:"price" db:"price"`
	InitialAmount   *float64   `json:"initialAmount" db:"initial_amount"`
	CurrentSupply   float64    `json:"currentSupply" db:"current_supply"`
	CurrentSpend    float64    `json:"currentSpend" db:"current_spend"`
	CurrentWriteOff float64    `json:"currentWriteOff" db:"current_write_off"`
	CurrentAmount   float64    `json:"currentAmount" db:"amount_in_stock"`
	LastInventoryId *int       `json:"lastInventoryId" db:"last_inventory_id"`
	LastInventory   *time.Time `json:"lastInventory" db:"last_inventory"`
}
