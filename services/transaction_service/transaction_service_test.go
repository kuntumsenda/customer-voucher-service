package transaction_service

import (
	"context"
	"customer-voucher-service/models/customer_model"
	"customer-voucher-service/models/transaction_model"
	"customer-voucher-service/models/voucher_model"
	pbTransaction "customer-voucher-service/protogen/transaction"
	pbVoucher "customer-voucher-service/protogen/voucher"
	"errors"
	"testing"
	"time"
)

type MockTransactionRepo struct {
	createTransactionFunc func(transaction *transaction_model.Transaction) (*transaction_model.Transaction, error)
	findByIdFunc          func(id uint) (*transaction_model.Transaction, error)
	listTransactionFunc   func(req *pbTransaction.ListTransactionReq) ([]*transaction_model.Transaction, error)
	detailTransactionFunc func(req *pbTransaction.DetailTransactionReq) (*transaction_model.Transaction, error)
}

func (m *MockTransactionRepo) CreateTransaction(transaction *transaction_model.Transaction) (*transaction_model.Transaction, error) {
	if m.createTransactionFunc != nil {
		return m.createTransactionFunc(transaction)
	}
	return transaction, nil
}

func (m *MockTransactionRepo) FindTransactionById(id uint) (*transaction_model.Transaction, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}

func (m *MockTransactionRepo) ListTransaction(req *pbTransaction.ListTransactionReq) ([]*transaction_model.Transaction, error) {
	if m.listTransactionFunc != nil {
		return m.listTransactionFunc(req)
	}
	return []*transaction_model.Transaction{}, nil
}

func (m *MockTransactionRepo) DetailTransaction(req *pbTransaction.DetailTransactionReq) (*transaction_model.Transaction, error) {
	if m.detailTransactionFunc != nil {
		return m.detailTransactionFunc(req)
	}
	return nil, nil
}

type MockVoucherRepo struct {
	createVoucherFunc func(voucher *voucher_model.Voucher) error
	listVoucherFunc   func(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error)
	findByIdFunc      func(id uint) (*voucher_model.Voucher, error)
}

func (m *MockVoucherRepo) CreateVoucher(voucher *voucher_model.Voucher) error {
	if m.createVoucherFunc != nil {
		return m.createVoucherFunc(voucher)
	}
	return nil
}

func (m *MockVoucherRepo) ListVoucher(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error) {
	if m.listVoucherFunc != nil {
		return m.listVoucherFunc(req)
	}
	return []*voucher_model.Voucher{}, nil
}

func (m *MockVoucherRepo) FindVoucherById(id uint) (*voucher_model.Voucher, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}

type MockCustomerRepo struct {
	createCustomerFunc func(customer *customer_model.Customer) error
	listCustomerFunc   func() ([]*customer_model.Customer, error)
	findByIdFunc       func(id uint) (*customer_model.Customer, error)
	updatePointsFunc   func(id uint, newPoints int64) error
}

func (m *MockCustomerRepo) CreateCustomer(customer *customer_model.Customer) error {
	if m.createCustomerFunc != nil {
		return m.createCustomerFunc(customer)
	}
	return nil
}

func (m *MockCustomerRepo) ListCustomer() ([]*customer_model.Customer, error) {
	if m.listCustomerFunc != nil {
		return m.listCustomerFunc()
	}
	return []*customer_model.Customer{}, nil
}

func (m *MockCustomerRepo) FindCustomerById(id uint) (*customer_model.Customer, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}

func (m *MockCustomerRepo) UpdatePointsCustomer(id uint, newPoints int64) error {
	if m.updatePointsFunc != nil {
		return m.updatePointsFunc(id, newPoints)
	}
	return nil
}

func TestTransactionRedeemPoint_Success(t *testing.T) {
	mockCustomer := &customer_model.Customer{
		ID:     1,
		Points: 1000,
	}

	mockVoucher := &voucher_model.Voucher{
		ID:          1,
		CostInPoint: 100,
	}

	expectedTransaction := &transaction_model.Transaction{
		ID:                 1,
		CustomerID:         1,
		VoucherID:          1,
		Quantity:           2,
		VoucherCostInPoint: 100,
		Total:              200,
		Status:             1,
		RedeemDate:         time.Now(),
	}

	mockTransactionRepo := &MockTransactionRepo{
		createTransactionFunc: func(transaction *transaction_model.Transaction) (*transaction_model.Transaction, error) {
			transaction.ID = expectedTransaction.ID
			return transaction, nil
		},
	}

	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}

	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			if id == 1 {
				return mockCustomer, nil
			}
			return nil, errors.New("customer not found")
		},
		updatePointsFunc: func(id uint, newPoints int64) error {
			if id == 1 && newPoints == 800 {
				return nil
			}
			return errors.New("update points failed")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   2,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if !result.IsSuccess {
		t.Error("Expected IsSuccess to be true")
	}
	if result.Data == nil {
		t.Error("Expected Data to not be nil")
	}
	if result.Data.Id != 1 {
		t.Errorf("Expected transaction ID to be 1, got %d", result.Data.Id)
	}
	if result.Data.CustomerId != 1 {
		t.Errorf("Expected customer ID to be 1, got %d", result.Data.CustomerId)
	}
	if result.Data.VoucherId != 1 {
		t.Errorf("Expected voucher ID to be 1, got %d", result.Data.VoucherId)
	}
	if result.Data.Quantity != 2 {
		t.Errorf("Expected quantity to be 2, got %d", result.Data.Quantity)
	}
	if result.Data.Total != 200 {
		t.Errorf("Expected total to be 200, got %d", result.Data.Total)
	}
}

func TestTransactionRedeemPoint_ValidationError_EmptyCustomerId(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{}
	mockCustomerRepo := &MockCustomerRepo{}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 0,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_ValidationError_EmptyVoucherId(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{}
	mockCustomerRepo := &MockCustomerRepo{}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  0,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_ValidationError_EmptyQuantity(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{}
	mockCustomerRepo := &MockCustomerRepo{}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   0,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_CustomerNotFound(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{}
	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			return nil, errors.New("customer not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_VoucherNotFound(t *testing.T) {
	mockCustomer := &customer_model.Customer{
		ID:     1,
		Points: 1000,
	}

	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			return nil, errors.New("voucher not found")
		},
	}
	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			if id == 1 {
				return mockCustomer, nil
			}
			return nil, errors.New("customer not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_NotEnoughPoints(t *testing.T) {
	mockCustomer := &customer_model.Customer{
		ID:     1,
		Points: 50,
	}

	mockVoucher := &voucher_model.Voucher{
		ID:          1,
		CostInPoint: 100,
	}

	mockTransactionRepo := &MockTransactionRepo{}
	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}
	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			if id == 1 {
				return mockCustomer, nil
			}
			return nil, errors.New("customer not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "not enough points to redeem" {
		t.Errorf("Expected error 'not enough points to redeem', got '%s'", err.Error())
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.IsSuccess {
		t.Error("Expected IsSuccess to be false")
	}
}

func TestTransactionRedeemPoint_CreateTransactionError(t *testing.T) {
	mockCustomer := &customer_model.Customer{
		ID:     1,
		Points: 1000,
	}

	mockVoucher := &voucher_model.Voucher{
		ID:          1,
		CostInPoint: 100,
	}

	mockTransactionRepo := &MockTransactionRepo{
		createTransactionFunc: func(transaction *transaction_model.Transaction) (*transaction_model.Transaction, error) {
			return nil, errors.New("database error")
		},
	}
	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}
	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			if id == 1 {
				return mockCustomer, nil
			}
			return nil, errors.New("customer not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestTransactionRedeemPoint_UpdatePointsError(t *testing.T) {
	mockCustomer := &customer_model.Customer{
		ID:     1,
		Points: 1000,
	}

	mockVoucher := &voucher_model.Voucher{
		ID:          1,
		CostInPoint: 100,
	}

	mockTransactionRepo := &MockTransactionRepo{
		createTransactionFunc: func(transaction *transaction_model.Transaction) (*transaction_model.Transaction, error) {
			transaction.ID = 1
			return transaction, nil
		},
	}
	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}
	mockCustomerRepo := &MockCustomerRepo{
		findByIdFunc: func(id uint) (*customer_model.Customer, error) {
			if id == 1 {
				return mockCustomer, nil
			}
			return nil, errors.New("customer not found")
		},
		updatePointsFunc: func(id uint, newPoints int64) error {
			return errors.New("update points failed")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
		voucherRepo:     mockVoucherRepo,
		customerRepo:    mockCustomerRepo,
	}

	req := &pbTransaction.TransactionRedeemPointReq{
		CustomerId: 1,
		VoucherId:  1,
		Quantity:   1,
	}

	result, err := service.TransactionRedeemPoint(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestListTransaction_Success(t *testing.T) {
	now := time.Now()
	mockTransactions := []*transaction_model.Transaction{
		{
			ID:           1,
			CustomerID:   1,
			VoucherID:    1,
			Quantity:     2,
			Total:        200,
			Status:       1,
			RedeemDate:   now,
			CreatedDate:  now,
			ModifiedDate: now,
			IsDeleted:    false,
		},
		{
			ID:           2,
			CustomerID:   2,
			VoucherID:    2,
			Quantity:     1,
			Total:        150,
			Status:       1,
			RedeemDate:   now,
			CreatedDate:  now,
			ModifiedDate: now,
			IsDeleted:    false,
		},
	}

	mockTransactionRepo := &MockTransactionRepo{
		listTransactionFunc: func(req *pbTransaction.ListTransactionReq) ([]*transaction_model.Transaction, error) {
			return mockTransactions, nil
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
	}

	req := &pbTransaction.ListTransactionReq{}

	result, err := service.ListTransaction(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(result.Data))
	}

	if result.Data[0].Id != 1 {
		t.Errorf("Expected first transaction ID to be 1, got %d", result.Data[0].Id)
	}
	if result.Data[0].CustomerId != 1 {
		t.Errorf("Expected first transaction customer ID to be 1, got %d", result.Data[0].CustomerId)
	}
	if result.Data[0].VoucherId != 1 {
		t.Errorf("Expected first transaction voucher ID to be 1, got %d", result.Data[0].VoucherId)
	}
	if result.Data[0].Quantity != 2 {
		t.Errorf("Expected first transaction quantity to be 2, got %d", result.Data[0].Quantity)
	}
	if result.Data[0].Total != 200 {
		t.Errorf("Expected first transaction total to be 200, got %d", result.Data[0].Total)
	}
	if *result.Data[0].Status != 1 {
		t.Errorf("Expected first transaction status to be 1, got %d", *result.Data[0].Status)
	}
	if *result.Data[0].IsDeleted {
		t.Error("Expected first transaction to not be deleted")
	}
}

func TestListTransaction_EmptyList(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{
		listTransactionFunc: func(req *pbTransaction.ListTransactionReq) ([]*transaction_model.Transaction, error) {
			return []*transaction_model.Transaction{}, nil
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
	}

	req := &pbTransaction.ListTransactionReq{}

	result, err := service.ListTransaction(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(result.Data))
	}
}

func TestListTransaction_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")
	mockTransactionRepo := &MockTransactionRepo{
		listTransactionFunc: func(req *pbTransaction.ListTransactionReq) ([]*transaction_model.Transaction, error) {
			return []*transaction_model.Transaction{}, expectedError
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
	}

	req := &pbTransaction.ListTransactionReq{}

	result, err := service.ListTransaction(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestDetailTransaction_Success(t *testing.T) {
	now := time.Now()
	mockTransaction := &transaction_model.Transaction{
		ID:           1,
		CustomerID:   1,
		VoucherID:    1,
		Quantity:     2,
		Total:        200,
		Status:       1,
		RedeemDate:   now,
		CreatedDate:  now,
		ModifiedDate: now,
		IsDeleted:    false,
	}

	mockTransactionRepo := &MockTransactionRepo{
		findByIdFunc: func(id uint) (*transaction_model.Transaction, error) {
			if id == 1 {
				return mockTransaction, nil
			}
			return nil, errors.New("transaction not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
	}

	req := &pbTransaction.DetailTransactionReq{
		Id: 1,
	}

	result, err := service.DetailTransaction(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.Data == nil {
		t.Error("Expected Data to not be nil")
	}
	if result.Data.Id != 1 {
		t.Errorf("Expected transaction ID to be 1, got %d", result.Data.Id)
	}
	if result.Data.CustomerId != 1 {
		t.Errorf("Expected customer ID to be 1, got %d", result.Data.CustomerId)
	}
	if result.Data.VoucherId != 1 {
		t.Errorf("Expected voucher ID to be 1, got %d", result.Data.VoucherId)
	}
	if result.Data.Quantity != 2 {
		t.Errorf("Expected quantity to be 2, got %d", result.Data.Quantity)
	}
	if result.Data.Total != 200 {
		t.Errorf("Expected total to be 200, got %d", result.Data.Total)
	}
	if *result.Data.Status != 1 {
		t.Errorf("Expected status to be 1, got %d", *result.Data.Status)
	}
	if *result.Data.IsDeleted {
		t.Error("Expected transaction to not be deleted")
	}
}

func TestDetailTransaction_NotFound(t *testing.T) {
	mockTransactionRepo := &MockTransactionRepo{
		findByIdFunc: func(id uint) (*transaction_model.Transaction, error) {
			return nil, errors.New("transaction not found")
		},
	}

	service := &TransactionService{
		transactionRepo: mockTransactionRepo,
	}

	req := &pbTransaction.DetailTransactionReq{
		Id: 999,
	}

	result, err := service.DetailTransaction(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestCalculateTotalPointRedeem(t *testing.T) {
	total := CalculateTotalPointRedeem(100, 3)
	if total != 300 {
		t.Errorf("Expected total to be 300, got %d", total)
	}

	total = CalculateTotalPointRedeem(50, 1)
	if total != 50 {
		t.Errorf("Expected total to be 50, got %d", total)
	}

	total = CalculateTotalPointRedeem(0, 5)
	if total != 0 {
		t.Errorf("Expected total to be 0, got %d", total)
	}
}

func TestRedundantPointsCustomer(t *testing.T) {
	redundant := RedundantPointsCustomer(200, 1000)
	if redundant != 800 {
		t.Errorf("Expected redundant points to be 800, got %d", redundant)
	}

	redundant = RedundantPointsCustomer(100, 100)
	if redundant != 0 {
		t.Errorf("Expected redundant points to be 0, got %d", redundant)
	}

	redundant = RedundantPointsCustomer(50, 30)
	if redundant != -20 {
		t.Errorf("Expected redundant points to be -20, got %d", redundant)
	}
}

func TestIsAbleToRedeem(t *testing.T) {
	if !IsAbleToRedeem(200, 1000) {
		t.Error("Expected to be able to redeem 200 points from 1000")
	}

	if IsAbleToRedeem(100, 100) {
		t.Error("Expected to not be able to redeem 100 points from 100")
	}

	if IsAbleToRedeem(50, 30) {
		t.Error("Expected to not be able to redeem 50 points from 30")
	}
}
