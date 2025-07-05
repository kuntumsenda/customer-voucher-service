package transaction_model

import "time"

type Transaction struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID         uint      `gorm:"not null" json:"customer_id"`
	VoucherID          uint      `gorm:"not null" json:"voucher_id"`
	Quantity           int64     `gorm:"not null" json:"quantity"`
	VoucherCostInPoint int64     `gorm:"not null" json:"voucher_cost_in_point"`
	Total              int64     `gorm:"not null" json:"total"`
	Status             int32     `gorm:"not null" json:"status"`
	RedeemDate         time.Time `gorm:"not null" json:"redeem_date"`
	IsDeleted          bool      `gorm:"default:false;not null" json:"is_deleted"`
	CreatedDate        time.Time `gorm:"autoCreateTime" json:"created_date"`
	CreatedBy          string    `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedDate       time.Time `gorm:"autoUpdateTime" json:"modified_date"`
	ModifiedBy         string    `gorm:"type:varchar(255)" json:"modified_by"`
}

func (Transaction) TableName() string {
	return "transaction"
}
