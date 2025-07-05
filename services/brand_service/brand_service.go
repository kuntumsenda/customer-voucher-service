package brand_service

import (
	"context"
	"customer-voucher-service/constants/message"
	"customer-voucher-service/db"
	"customer-voucher-service/models/brand_model"
	pbBrand "customer-voucher-service/protogen/brand"
	"errors"

	"github.com/go-playground/validator/v10"
)

type IBrandService interface {
	CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error)
}

type BrandService struct {
	pbBrand.UnimplementedBrandServiceServer
	brandRepo brand_model.BrandRepo
}

func NewBrandService() *BrandService {
	return &BrandService{brandRepo: *brand_model.NewBrandRepo()}
}

type createBrandReqValidate struct {
	Name        string `validate:"required,max=255"`
	Description string `validate:"max=255"`
}

func validateCreateBrandReq(req *pbBrand.CreateBrandReq) error {
	v := validator.New()
	validateReq := createBrandReqValidate{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
	if err := v.Struct(validateReq); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Tag() {
			case "required":
				return errors.New(message.RequiredMessage(fieldErr.Field()))
			case "max":
				return errors.New(message.MaxLengthMessage(fieldErr.Field(), 255))
			}
		}
		return errors.New("invalid request")
	}
	return nil
}

func (s *BrandService) CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error) {
	if err := validateCreateBrandReq(req); err != nil {
		return &pbBrand.CreateBrandRes{IsSuccess: false}, err
	}
	brand := &brand_model.Brand{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
	err := s.brandRepo.CreateBrand(db.DB, brand)
	if err != nil {
		return &pbBrand.CreateBrandRes{IsSuccess: false}, err
	}
	return &pbBrand.CreateBrandRes{IsSuccess: true}, nil
}
