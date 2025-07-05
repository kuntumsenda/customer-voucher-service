package voucher_model

import (
	pb "customer-voucher-service/protogen/voucher"

	"gorm.io/gorm"
)

type IVoucherRepo interface {
	CreateVoucher(voucher *Voucher) error
	ListVoucher(*pb.ListVoucherReq) ([]*Voucher, error)
	FindVoucherById(id uint) (*Voucher, error)
}

type VoucherRepo struct {
	db *gorm.DB
}

func NewVoucherRepo(db *gorm.DB) *VoucherRepo {
	return &VoucherRepo{
		db: db,
	}
}

func (r *VoucherRepo) CreateVoucher(voucher *Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *VoucherRepo) ListVoucher(req *pb.ListVoucherReq) ([]*Voucher, error) {
	var vouchers []*Voucher
	query := r.db.Model(&Voucher{}).Where("is_deleted = ?", false)

	if req.BrandId != nil {
		query = query.Where("brand_id = ?", *req.BrandId)
	}

	err := query.Find(&vouchers).Error
	return vouchers, err
}

func (r *VoucherRepo) FindVoucherById(id uint) (*Voucher, error) {
	var voucher Voucher
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}
