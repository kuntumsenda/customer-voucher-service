package transaction_service

import (
	"context"
	"customer-voucher-service/constants"
	"customer-voucher-service/constants/message"
	"customer-voucher-service/db"
	"customer-voucher-service/models/customer_model"
	"customer-voucher-service/models/transaction_model"
	"customer-voucher-service/models/voucher_model"
	pbTransaction "customer-voucher-service/protogen/transaction"
	"customer-voucher-service/utils/validator"
	"errors"
	"time"
)

type ITransactionService interface {
	TransactionRedeemPoint(ctx context.Context, req *pbTransaction.TransactionRedeemPointReq) (*pbTransaction.TransactionRedeemPointRes, error)
	ListTransaction(ctx context.Context, req *pbTransaction.ListTransactionReq) (*pbTransaction.ListTransactionRes, error)
	DetailTransaction(ctx context.Context, req *pbTransaction.DetailTransactionReq) (*pbTransaction.DetailTransactionRes, error)
}

type TransactionService struct {
	pbTransaction.UnimplementedTransactionServiceServer
	transactionRepo transaction_model.TransactionRepo
	voucherRepo     voucher_model.VoucherRepo
	customerRepo    customer_model.CustomerRepo
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		transactionRepo: *transaction_model.NewTransactionRepo(db.DB),
		voucherRepo:     *voucher_model.NewVoucherRepo(db.DB),
		customerRepo:    *customer_model.NewCustomerRepo(db.DB),
	}
}

type createTransactionReqValidate struct {
	CustomerId int32 `validate:"required"`
	VoucherId  int32 `validate:"required"`
	Quantity   int64 `validate:"required"`
}

func (s *TransactionService) TransactionRedeemPoint(ctx context.Context, req *pbTransaction.TransactionRedeemPointReq) (*pbTransaction.TransactionRedeemPointRes, error) {
	validateReq := createTransactionReqValidate{
		CustomerId: req.CustomerId,
		VoucherId:  req.VoucherId,
		Quantity:   req.Quantity,
	}
	if err := validator.ValidateReqField(validateReq); err != nil {
		return &pbTransaction.TransactionRedeemPointRes{IsSuccess: false}, err
	}

	// check customer
	resCustomer, err := s.customerRepo.FindCustomerById(uint(req.CustomerId))
	if err != nil || resCustomer == nil {
		return &pbTransaction.TransactionRedeemPointRes{IsSuccess: false}, errors.New(message.NotFoundMessage("customer"))
	}

	// check voucher
	resVoucher, err := s.voucherRepo.FindVoucherById(uint(req.VoucherId))
	if err != nil || resVoucher == nil {
		return &pbTransaction.TransactionRedeemPointRes{IsSuccess: false}, errors.New(message.NotFoundMessage("voucher"))
	}

	totalRedeem := CalculateTotalPointRedeem(resVoucher.CostInPoint, req.Quantity)

	if !IsAbleToRedeem(totalRedeem, resCustomer.Points) {
		return &pbTransaction.TransactionRedeemPointRes{IsSuccess: false}, errors.New("not enough points to redeem")
	}

	transaction := &transaction_model.Transaction{
		CustomerID:         resCustomer.ID,
		VoucherID:          resVoucher.ID,
		Quantity:           req.Quantity,
		VoucherCostInPoint: resVoucher.CostInPoint,
		Total:              totalRedeem,
		Status:             1,
		RedeemDate:         time.Now(),
	}

	result, err := s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = s.customerRepo.UpdatePointsCustomer(resCustomer.ID, RedundantPointsCustomer(result.Total, resCustomer.Points))
	if err != nil {
		return nil, err
	}

	return &pbTransaction.TransactionRedeemPointRes{
		IsSuccess: true,
		Data: &pbTransaction.Transaction{
			Id:         int32(result.ID),
			CustomerId: int32(result.CustomerID),
			VoucherId:  int32(result.VoucherID),
			Quantity:   result.Quantity,
			Total:      result.Total,
			Status:     &result.Status,
			RedeemDate: result.RedeemDate.Format(constants.FormatDate),
		},
	}, nil
}

func CalculateTotalPointRedeem(cip int64, qty int64) int64 {
	total := cip * qty
	return total
}

func RedundantPointsCustomer(total int64, custPoints int64) int64 {
	return custPoints - total
}

func IsAbleToRedeem(total int64, custPoints int64) bool {
	return RedundantPointsCustomer(total, custPoints) > 0
}

func (s *TransactionService) ListTransaction(ctx context.Context, req *pbTransaction.ListTransactionReq) (*pbTransaction.ListTransactionRes, error) {
	result, err := s.transactionRepo.ListTransaction(req)
	if err != nil {
		return nil, err
	}
	list := []*pbTransaction.Transaction{}

	for _, trans := range result {
		status := int32(trans.Status)
		isDeleted := trans.IsDeleted
		data := pbTransaction.Transaction{
			Id:           int32(trans.ID),
			CustomerId:   int32(trans.CustomerID),
			VoucherId:    int32(trans.VoucherID),
			Quantity:     trans.Quantity,
			Total:        int64(trans.Total),
			Status:       &status,
			RedeemDate:   trans.RedeemDate.Format(constants.FormatDate),
			CreatedDate:  trans.CreatedDate.Format(constants.FormatDate),
			ModifiedDate: trans.ModifiedDate.Format(constants.FormatDate),
			IsDeleted:    &isDeleted,
		}
		list = append(list, &data)
	}
	return &pbTransaction.ListTransactionRes{
		Data: list,
	}, nil
}

func (s *TransactionService) DetailTransaction(ctx context.Context, req *pbTransaction.DetailTransactionReq) (*pbTransaction.DetailTransactionRes, error) {
	result, err := s.transactionRepo.FindTransactionById(uint(req.Id))
	if err != nil {
		return nil, err
	}

	status := int32(result.Status)
	isDeleted := result.IsDeleted
	data := &pbTransaction.Transaction{
		Id:           int32(result.ID),
		CustomerId:   int32(result.CustomerID),
		VoucherId:    int32(result.VoucherID),
		Quantity:     result.Quantity,
		Total:        int64(result.Total),
		Status:       &status,
		RedeemDate:   result.RedeemDate.Format(constants.FormatDate),
		CreatedDate:  result.CreatedDate.Format(constants.FormatDate),
		ModifiedDate: result.ModifiedDate.Format(constants.FormatDate),
		IsDeleted:    &isDeleted,
	}

	return &pbTransaction.DetailTransactionRes{
		Data: data,
	}, nil
}
