package brand_model

import "gorm.io/gorm"

type IBrandRepo interface {
	CreateBrand(brand *Brand) error
	ListBrand() (*[]Brand, error)
	FindBrandById(id uint) (*Brand, error)
}

type BrandRepo struct {
	db *gorm.DB
}

func NewBrandRepo(db *gorm.DB) *BrandRepo {
	return &BrandRepo{
		db: db,
	}
}

func (r *BrandRepo) CreateBrand(brand *Brand) error {
	return r.db.Create(brand).Error
}

func (r *BrandRepo) ListBrand() ([]*Brand, error) {
	var brand []*Brand
	err := r.db.Find(&brand).Where("is_deleted = ?", false).Error
	return brand, err
}

func (r *BrandRepo) FindBrandById(id uint) (*Brand, error) {
	var brand Brand
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}
