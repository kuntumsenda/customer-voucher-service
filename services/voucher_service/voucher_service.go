package voucher_service

import (
	"context"
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/message"
	"customer-voucher-service/db"
	"customer-voucher-service/models/brand_model"
	"customer-voucher-service/models/voucher_model"
	pbVoucher "customer-voucher-service/protogen/voucher"
	"customer-voucher-service/utils/validator"
	"errors"
)

type IVoucherService interface {
	CreateVoucher(ctx context.Context, req *pbVoucher.CreateVoucherReq) (*pbVoucher.CreateVoucherRes, error)
	ListVoucher(ctx context.Context, req *pbVoucher.ListVoucherReq) (*pbVoucher.ListVoucherRes, error)
	DetailVoucher(ctx context.Context, req *pbVoucher.DetailVoucherReq) (*pbVoucher.DetailVoucherRes, error)
}

type VoucherService struct {
	pbVoucher.UnimplementedVoucherServiceServer
	voucherRepo voucher_model.IVoucherRepo
	brandRepo   brand_model.IBrandRepo
}

func NewVoucherService() *VoucherService {
	return &VoucherService{
		voucherRepo: voucher_model.NewVoucherRepo(db.DB),
		brandRepo:   brand_model.NewBrandRepo(db.DB),
	}
}

type createVoucherReqValidate struct {
	BrandId     int32  `validate:"required"`
	Name        string `validate:"required,max=255"`
	Description string `validate:max=255"`
	CostInPoint int64  `validate:"required"`
	VoucherCode string `validate:"required,max=255"`
}

func (s *VoucherService) CreateVoucher(ctx context.Context, req *pbVoucher.CreateVoucherReq) (*pbVoucher.CreateVoucherRes, error) {
	validateReq := createVoucherReqValidate{
		BrandId:     req.BrandId,
		Name:        req.Name,
		Description: req.Description,
		CostInPoint: req.CostInPoint,
		VoucherCode: req.VoucherCode,
	}
	if err := validator.ValidateReqField(validateReq); err != nil {
		return &pbVoucher.CreateVoucherRes{IsSuccess: false}, err
	}
	resBrand, err := s.brandRepo.FindBrandById(uint(req.BrandId))
	if err != nil || resBrand == nil {
		return &pbVoucher.CreateVoucherRes{IsSuccess: false}, errors.New(message.NotFoundMessage("brand"))
	}
	voucher := &voucher_model.Voucher{
		BrandID:     resBrand.ID,
		Name:        req.Name,
		Description: req.Description,
		CostInPoint: req.CostInPoint,
		VoucherCode: req.VoucherCode,
	}
	err = s.voucherRepo.CreateVoucher(voucher)
	if err != nil {
		return nil, err
	}
	return &pbVoucher.CreateVoucherRes{IsSuccess: true}, nil
}

func (s *VoucherService) ListVoucher(ctx context.Context, req *pbVoucher.ListVoucherReq) (*pbVoucher.ListVoucherRes, error) {
	result, err := s.voucherRepo.ListVoucher(req)
	if err != nil {
		return nil, err
	}
	list := []*pbVoucher.Voucher{}

	for _, cust := range result {
		points := int64(cust.CostInPoint)
		data := pbVoucher.Voucher{
			Id:           int32(cust.ID),
			BrandId:      int32(cust.BrandID),
			Name:         cust.Name,
			Description:  cust.Description,
			VoucherCode:  cust.VoucherCode,
			CostInPoint:  points,
			CreatedDate:  cust.CreatedDate.Format(constants.FormatDate),
			ModifiedDate: cust.ModifiedDate.Format(constants.FormatDate),
			IsDeleted:    &cust.IsDeleted,
		}
		list = append(list, &data)
	}
	return &pbVoucher.ListVoucherRes{
		Data: list,
	}, nil
}

func (s *VoucherService) DetailVoucher(ctx context.Context, req *pbVoucher.DetailVoucherReq) (*pbVoucher.DetailVoucherRes, error) {
	result, err := s.voucherRepo.FindVoucherById(uint(req.Id))
	if err != nil {
		return nil, err
	}

	isDeleted := result.IsDeleted
	data := &pbVoucher.Voucher{
		Id:           int32(result.ID),
		BrandId:      int32(result.BrandID),
		Name:         result.Name,
		Description:  result.Description,
		VoucherCode:  result.VoucherCode,
		CostInPoint:  result.CostInPoint,
		CreatedDate:  result.CreatedDate.Format(constants.FormatDate),
		ModifiedDate: result.ModifiedDate.Format(constants.FormatDate),
		IsDeleted:    &isDeleted,
	}

	return &pbVoucher.DetailVoucherRes{
		Data: data,
	}, nil
}
