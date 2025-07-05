package voucher_model

import "time"

type Voucher struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BrandID      uint      `gorm:"not null" json:"brand_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	CostInPoint  int       `gorm:"not null" json:"cost_in_point"`
	IsDeleted    bool      `gorm:"default:false;not null" json:"is_deleted"`
	CreatedDate  time.Time `gorm:"autoCreateTime" json:"created_date"`
	CreatedBy    string    `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedDate time.Time `gorm:"autoUpdateTime" json:"modified_date"`
	ModifiedBy   string    `gorm:"type:varchar(255)" json:"modified_by"`
}

func (Voucher) TableName() string {
	return "voucher"
}
