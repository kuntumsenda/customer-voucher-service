package brand_service

import (
	"context"
	"customer-voucher-service/constants"
	"customer-voucher-service/db"
	"customer-voucher-service/models/brand_model"
	pbBrand "customer-voucher-service/protogen/brand"
	"customer-voucher-service/utils/validator"
)

type IBrandService interface {
	CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error)
	ListBrand(ctx context.Context, req *pbBrand.ListBrandReq) (*pbBrand.ListBrandRes, error)
}

type BrandService struct {
	pbBrand.UnimplementedBrandServiceServer
	brandRepo brand_model.IBrandRepo
}

func NewBrandService() *BrandService {
	return &BrandService{brandRepo: brand_model.NewBrandRepo(db.DB)}
}

type createBrandReqValidate struct {
	Name        string `validate:"required,max=255"`
	Description string `validate:"max=255"`
}

func (s *BrandService) CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error) {
	validateReq := createBrandReqValidate{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := validator.ValidateReqField(validateReq); err != nil {
		return &pbBrand.CreateBrandRes{IsSuccess: false}, err
	}
	brand := &brand_model.Brand{
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.brandRepo.CreateBrand(brand)
	if err != nil {
		return nil, err
	}
	return &pbBrand.CreateBrandRes{IsSuccess: true}, nil
}

func (s *BrandService) ListBrand(ctx context.Context, req *pbBrand.ListBrandReq) (*pbBrand.ListBrandRes, error) {
	result, err := s.brandRepo.ListBrand()
	if err != nil {
		return nil, err
	}
	list := []*pbBrand.Brand{}

	for _, b := range result {
		data := pbBrand.Brand{
			Id:           int32(b.ID),
			Name:         b.Name,
			Description:  b.Description,
			CreatedDate:  b.CreatedDate.Format(constants.FormatDate),
			ModifiedDate: b.ModifiedDate.Format(constants.FormatDate),
			IsDeleted:    &b.IsDeleted,
		}
		list = append(list, &data)
	}
	return &pbBrand.ListBrandRes{
		Data: list,
	}, nil
}
