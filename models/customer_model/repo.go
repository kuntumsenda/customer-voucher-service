package customer_model

import "gorm.io/gorm"

type ICustomerRepo interface {
	CreateCustomer(customer *Customer) error
	ListCustomer() ([]*Customer, error)
	FindCustomerById(id uint) (*Customer, error)
	UpdatePointsCustomer(id uint, newPoints int64) error
}

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) CreateCustomer(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepo) ListCustomer() ([]*Customer, error) {
	var customers []*Customer
	err := r.db.Find(&customers).Where("is_deleted = ?", false).Error
	return customers, err
}

func (r *CustomerRepo) FindCustomerById(id uint) (*Customer, error) {
	var customer Customer
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) UpdatePointsCustomer(id uint, newPoints int64) error {
	return r.db.Model(&Customer{}).Where("id = ? AND is_deleted = ?", id, false).Update("points", newPoints).Error
}
