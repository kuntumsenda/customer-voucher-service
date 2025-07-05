package customer_model

import "time"

type Customer struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Points       int       `gorm:"default:0;not null" json:"points"`
	IsDeleted    bool      `gorm:"default:false;not null" json:"is_deleted"`
	CreatedDate  time.Time `gorm:"autoCreateTime" json:"created_date"`
	CreatedBy    string    `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedDate time.Time `gorm:"autoUpdateTime" json:"modified_date"`
	ModifiedBy   string    `gorm:"type:varchar(255)" json:"modified_by"`
}

func (Customer) TableName() string {
	return "customer"
}
