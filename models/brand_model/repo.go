package brand_model

import "gorm.io/gorm"

type IBrandRepo interface {
	CreateBrand(db *gorm.DB, brand *Brand) error
}

type BrandRepo struct{}

func NewBrandRepo() *BrandRepo {
	return &BrandRepo{}
}

func (r *BrandRepo) CreateBrand(db *gorm.DB, brand *Brand) error {
	return db.Create(brand).Error
}
