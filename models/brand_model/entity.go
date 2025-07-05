package brand_model

import "time"

type Brand struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Description  string    `gorm:"type:text"`
	IsDeleted    bool      `gorm:"default:false"`
	CreatedDate  time.Time `gorm:"autoCreateTime"`
	CreatedBy    string    `gorm:"type:varchar(255)"`
	ModifiedDate time.Time `gorm:"autoUpdateTime"`
	ModifiedBy   string    `gorm:"type:varchar(255)"`
}

func (Brand) TableName() string {
	return "brand"
}
