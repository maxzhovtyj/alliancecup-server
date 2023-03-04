package inventory

import "time"

type InsertProductDTO struct {
	InventoryId     int      `json:"-" db:"inventory_id"`
	ProductId       int      `json:"productId" db:"product_id" binding:"required"`
	ProductPrice    float64  `json:"productPrice" db:"product_price" binding:"required"`
	LastInventoryId *int     `json:"lastInventoryId" db:"last_inventory_id"`
	InitialAmount   *float64 `json:"initialAmount" db:"initial_amount"`        // amount from the last inventory
	Supply          float64  `json:"currentSupply,omitempty" db:"supply"`      // from the last inventory
	Spend           float64  `json:"currentSpend,omitempty" db:"spend"`        // spending (customers orders) from the last inventory
	WriteOff        float64  `json:"currentWriteOff,omitempty" db:"write_off"` // something that wasn't sold
	PlannedAmount   float64  `json:"currentAmount" db:"planned_amount"`        // current amount in stock
	RealAmount      float64  `json:"realAmount" db:"real_amount"`              // inventory input
}

type SelectProductDTO struct {
	InventoryId     int      `json:"-" db:"inventory_id"`
	ProductId       int      `json:"productId" binding:"required" db:"product_id"`
	ProductTitle    string   `json:"productTitle" binding:"required" db:"product_title"`
	ProductPrice    float64  `json:"productPrice" binding:"required" db:"product_price"`
	LastInventoryId *int     `json:"lastInventoryId" db:"last_inventory_id"`
	InitialAmount   *float64 `json:"initialAmount" db:"initial_amount"`                     // amount from the last inventory
	Supply          float64  `json:"supply" binding:"required" db:"supply"`                 // from the last inventory
	Spend           float64  `json:"spends" binding:"required" db:"spend"`                  // spending (customers orders) from the last inventory
	WriteOff        float64  `json:"writeOff" binding:"required" db:"write_off"`            // something that wasn't sold
	WriteOffPrice   float64  `json:"writeOffPrice" binding:"required" db:"write_off_price"` // something that wasn't sold price
	PlannedAmount   float64  `json:"plannedAmount" binding:"required" db:"planned_amount"`  // current amount in stock
	RealAmount      float64  `json:"realAmount" binding:"required" db:"real_amount"`        // inventory input
	RealAmountPrice float64  `json:"realAmountPrice" db:"real_amount_price"`                // inventory input price
	Difference      float64  `json:"difference" db:"difference"`
	DifferencePrice float64  `json:"differencePrice" db:"difference_price"`
}

type CurrentProductDTO struct {
	ProductId       int        `json:"productId" db:"id"`
	Title           string     `json:"title" db:"product_title"`
	ProductPrice    float64    `json:"productPrice" db:"product_price"`
	InitialAmount   *float64   `json:"initialAmount" db:"initial_amount"`
	CurrentSupply   float64    `json:"currentSupply" db:"current_supply"`
	CurrentSpend    float64    `json:"currentSpend" db:"current_spend"`
	CurrentWriteOff float64    `json:"currentWriteOff" db:"current_write_off"`
	WriteOffPrice   float64    `json:"writeOffPrice" db:"write_off_price"`
	CurrentAmount   float64    `json:"currentAmount" db:"amount_in_stock"`
	RealAmount      float64    `json:"realAmount" db:"current_real_amount"`
	RealAmountPrice float64    `json:"realAmountPrice" db:"real_amount_price"`
	Difference      float64    `json:"difference" db:"difference"`
	DifferencePrice float64    `json:"differencePrice" db:"difference_price"`
	LastInventoryId *int       `json:"lastInventoryId" db:"last_inventory_id"`
	LastInventory   *time.Time `json:"lastInventory" db:"last_inventory"`
}

type DTO struct {
	Id        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
