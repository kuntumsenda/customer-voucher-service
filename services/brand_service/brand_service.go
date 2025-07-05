package brand_service

import (
	"context"
	pbBrand "customer-voucher-service/protogen/brand"
)

type IBrandService interface {
	CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error)
}

type BrandService struct {
	pbBrand.UnimplementedBrandServiceServer
}

func NewBrandService() *BrandService {
	return &BrandService{}
}

func (s *BrandService) CreateBrand(ctx context.Context, req *pbBrand.CreateBrandReq) (*pbBrand.CreateBrandRes, error) {
	return &pbBrand.CreateBrandRes{
		IsSuccess: true,
	}, nil
}
