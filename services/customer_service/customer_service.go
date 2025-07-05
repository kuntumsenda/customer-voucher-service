package customer_service

import (
	"context"
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/message"
	"customer-voucher-service/db"
	"customer-voucher-service/models/customer_model"
	pbCustomer "customer-voucher-service/protogen/customer"
	"customer-voucher-service/utils/validator"
	"errors"
)

type ICustomerService interface {
	CreateCustomer(ctx context.Context, req *pbCustomer.CreateCustomerReq) (*pbCustomer.CreateCustomerRes, error)
	ListCustomer(ctx context.Context, req *pbCustomer.ListCustomerReq) (*pbCustomer.ListCustomerRes, error)
	UpdateCustomerPoints(ctx context.Context, req *pbCustomer.UpdateCustomerPointsReq) (*pbCustomer.UpdateCustomerPointsRes, error)
}

type CustomerService struct {
	pbCustomer.UnimplementedCustomerServiceServer
	customerRepo customer_model.ICustomerRepo
}

func NewCustomerService() *CustomerService {
	return &CustomerService{customerRepo: customer_model.NewCustomerRepo(db.DB)}
}

type createCustomerReqValidate struct {
	FullName string `validate:"required,max=255"`
	Email    string `validate:"required,email,max=255"`
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *pbCustomer.CreateCustomerReq) (*pbCustomer.CreateCustomerRes, error) {
	validateReq := createCustomerReqValidate{
		FullName: req.FullName,
		Email:    req.Email,
	}
	if err := validator.ValidateReqField(validateReq); err != nil {
		return &pbCustomer.CreateCustomerRes{IsSuccess: false}, err
	}
	customer := &customer_model.Customer{
		FullName: req.FullName,
		Email:    req.Email,
		Points:   req.Points,
	}
	err := s.customerRepo.CreateCustomer(customer)
	if err != nil {
		return nil, err
	}
	return &pbCustomer.CreateCustomerRes{IsSuccess: true}, nil
}

func (s *CustomerService) ListCustomer(ctx context.Context, req *pbCustomer.ListCustomerReq) (*pbCustomer.ListCustomerRes, error) {
	result, err := s.customerRepo.ListCustomer()
	if err != nil {
		return nil, err
	}
	list := []*pbCustomer.Customer{}

	for _, cust := range result {
		points := int64(cust.Points)
		data := pbCustomer.Customer{
			Id:           int32(cust.ID),
			FullName:     cust.FullName,
			Email:        cust.Email,
			Points:       &points,
			CreatedDate:  cust.CreatedDate.Format(constants.FormatDate),
			ModifiedDate: cust.ModifiedDate.Format(constants.FormatDate),
			IsDeleted:    &cust.IsDeleted,
		}
		list = append(list, &data)
	}
	return &pbCustomer.ListCustomerRes{
		Data: list,
	}, nil
}

func (s *CustomerService) UpdateCustomerPoints(ctx context.Context, req *pbCustomer.UpdateCustomerPointsReq) (*pbCustomer.UpdateCustomerPointsRes, error) {
	if req.Points < 0 {
		return &pbCustomer.UpdateCustomerPointsRes{IsSuccess: false}, errors.New(message.InvalidFormatMessage("customerPoints"))
	}

	respCustomer, err := s.customerRepo.FindCustomerById(uint(req.Id))
	if err != nil || respCustomer == nil {
		return &pbCustomer.UpdateCustomerPointsRes{IsSuccess: false}, errors.New(message.NotFoundMessage("customer"))
	}

	err = s.customerRepo.UpdatePointsCustomer(respCustomer.ID, req.Points)
	if err != nil {
		return nil, err
	}

	return &pbCustomer.UpdateCustomerPointsRes{IsSuccess: true}, nil
}
