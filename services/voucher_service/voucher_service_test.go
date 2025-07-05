package voucher_service

import (
	"context"
	"customer-voucher-service/models/brand_model"
	"customer-voucher-service/models/voucher_model"
	pbVoucher "customer-voucher-service/protogen/voucher"
	"errors"
	"testing"
	"time"
)

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

type MockBrandRepo struct {
	createBrandFunc func(brand *brand_model.Brand) error
	listBrandFunc   func() ([]*brand_model.Brand, error)
	findByIdFunc    func(id uint) (*brand_model.Brand, error)
}

func (m *MockBrandRepo) CreateBrand(brand *brand_model.Brand) error {
	if m.createBrandFunc != nil {
		return m.createBrandFunc(brand)
	}
	return nil
}

func (m *MockBrandRepo) ListBrand() ([]*brand_model.Brand, error) {
	if m.listBrandFunc != nil {
		return m.listBrandFunc()
	}
	return []*brand_model.Brand{}, nil
}

func (m *MockBrandRepo) FindBrandById(id uint) (*brand_model.Brand, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}

func TestCreateVoucher_Success(t *testing.T) {
	mockBrand := &brand_model.Brand{
		ID:   1,
		Name: "Test Brand",
	}

	mockVoucherRepo := &MockVoucherRepo{
		createVoucherFunc: func(voucher *voucher_model.Voucher) error {
			if voucher.BrandID != 1 {
				t.Errorf("Expected brand ID to be 1, got %d", voucher.BrandID)
			}
			if voucher.Name != "Test Voucher" {
				t.Errorf("Expected name to be 'Test Voucher', got '%s'", voucher.Name)
			}
			if voucher.Description != "Test Description" {
				t.Errorf("Expected description to be 'Test Description', got '%s'", voucher.Description)
			}
			if voucher.CostInPoint != 100 {
				t.Errorf("Expected cost in point to be 100, got %d", voucher.CostInPoint)
			}
			if voucher.VoucherCode != "TEST001" {
				t.Errorf("Expected voucher code to be 'TEST001', got '%s'", voucher.VoucherCode)
			}
			return nil
		},
	}

	mockBrandRepo := &MockBrandRepo{
		findByIdFunc: func(id uint) (*brand_model.Brand, error) {
			if id == 1 {
				return mockBrand, nil
			}
			return nil, errors.New("brand not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if !result.IsSuccess {
		t.Error("Expected IsSuccess to be true")
	}
}

func TestCreateVoucher_ValidationError_EmptyBrandId(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     0,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_EmptyName(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_NameTooLong(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	longName := ""
	for i := 0; i < 256; i++ {
		longName += "a"
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        longName,
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_DescriptionTooLong(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	longDescription := ""
	for i := 0; i < 256; i++ {
		longDescription += "a"
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: longDescription,
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_EmptyCostInPoint(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 0,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_EmptyVoucherCode(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_ValidationError_VoucherCodeTooLong(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	longVoucherCode := ""
	for i := 0; i < 256; i++ {
		longVoucherCode += "a"
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: longVoucherCode,
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_BrandNotFound(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{}
	mockBrandRepo := &MockBrandRepo{
		findByIdFunc: func(id uint) (*brand_model.Brand, error) {
			return nil, errors.New("brand not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

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

func TestCreateVoucher_CreateVoucherError(t *testing.T) {
	mockBrand := &brand_model.Brand{
		ID:   1,
		Name: "Test Brand",
	}

	mockVoucherRepo := &MockVoucherRepo{
		createVoucherFunc: func(voucher *voucher_model.Voucher) error {
			return errors.New("database error")
		},
	}

	mockBrandRepo := &MockBrandRepo{
		findByIdFunc: func(id uint) (*brand_model.Brand, error) {
			if id == 1 {
				return mockBrand, nil
			}
			return nil, errors.New("brand not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
		brandRepo:   mockBrandRepo,
	}

	req := &pbVoucher.CreateVoucherReq{
		BrandId:     1,
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 100,
		VoucherCode: "TEST001",
	}

	result, err := service.CreateVoucher(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestListVoucher_Success(t *testing.T) {
	now := time.Now()
	mockVouchers := []*voucher_model.Voucher{
		{
			ID:           1,
			BrandID:      1,
			Name:         "Voucher 1",
			Description:  "Description 1",
			VoucherCode:  "VOUCHER001",
			CostInPoint:  100,
			CreatedDate:  now,
			ModifiedDate: now,
			IsDeleted:    false,
		},
		{
			ID:           2,
			BrandID:      2,
			Name:         "Voucher 2",
			Description:  "Description 2",
			VoucherCode:  "VOUCHER002",
			CostInPoint:  200,
			CreatedDate:  now,
			ModifiedDate: now,
			IsDeleted:    false,
		},
	}

	mockVoucherRepo := &MockVoucherRepo{
		listVoucherFunc: func(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error) {
			return mockVouchers, nil
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.ListVoucherReq{}

	result, err := service.ListVoucher(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 vouchers, got %d", len(result.Data))
	}

	if result.Data[0].Id != 1 {
		t.Errorf("Expected first voucher ID to be 1, got %d", result.Data[0].Id)
	}
	if result.Data[0].BrandId != 1 {
		t.Errorf("Expected first voucher brand ID to be 1, got %d", result.Data[0].BrandId)
	}
	if result.Data[0].Name != "Voucher 1" {
		t.Errorf("Expected first voucher name to be 'Voucher 1', got '%s'", result.Data[0].Name)
	}
	if result.Data[0].Description != "Description 1" {
		t.Errorf("Expected first voucher description to be 'Description 1', got '%s'", result.Data[0].Description)
	}
	if result.Data[0].VoucherCode != "VOUCHER001" {
		t.Errorf("Expected first voucher code to be 'VOUCHER001', got '%s'", result.Data[0].VoucherCode)
	}
	if result.Data[0].CostInPoint != 100 {
		t.Errorf("Expected first voucher cost in point to be 100, got %d", result.Data[0].CostInPoint)
	}
	if *result.Data[0].IsDeleted {
		t.Error("Expected first voucher to not be deleted")
	}

	if result.Data[1].Id != 2 {
		t.Errorf("Expected second voucher ID to be 2, got %d", result.Data[1].Id)
	}
	if result.Data[1].BrandId != 2 {
		t.Errorf("Expected second voucher brand ID to be 2, got %d", result.Data[1].BrandId)
	}
	if result.Data[1].Name != "Voucher 2" {
		t.Errorf("Expected second voucher name to be 'Voucher 2', got '%s'", result.Data[1].Name)
	}
	if result.Data[1].Description != "Description 2" {
		t.Errorf("Expected second voucher description to be 'Description 2', got '%s'", result.Data[1].Description)
	}
	if result.Data[1].VoucherCode != "VOUCHER002" {
		t.Errorf("Expected second voucher code to be 'VOUCHER002', got '%s'", result.Data[1].VoucherCode)
	}
	if result.Data[1].CostInPoint != 200 {
		t.Errorf("Expected second voucher cost in point to be 200, got %d", result.Data[1].CostInPoint)
	}
	if *result.Data[1].IsDeleted {
		t.Error("Expected second voucher to not be deleted")
	}
}

func TestListVoucher_EmptyList(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{
		listVoucherFunc: func(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error) {
			return []*voucher_model.Voucher{}, nil
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.ListVoucherReq{}

	result, err := service.ListVoucher(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected 0 vouchers, got %d", len(result.Data))
	}
}

func TestListVoucher_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")
	mockVoucherRepo := &MockVoucherRepo{
		listVoucherFunc: func(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error) {
			return []*voucher_model.Voucher{}, expectedError
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.ListVoucherReq{}

	result, err := service.ListVoucher(context.Background(), req)

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

func TestListVoucher_WithBrandFilter(t *testing.T) {
	now := time.Now()
	mockVouchers := []*voucher_model.Voucher{
		{
			ID:           1,
			BrandID:      1,
			Name:         "Voucher 1",
			Description:  "Description 1",
			VoucherCode:  "VOUCHER001",
			CostInPoint:  100,
			CreatedDate:  now,
			ModifiedDate: now,
			IsDeleted:    false,
		},
	}

	mockVoucherRepo := &MockVoucherRepo{
		listVoucherFunc: func(req *pbVoucher.ListVoucherReq) ([]*voucher_model.Voucher, error) {
			if req.BrandId != nil && *req.BrandId == 1 {
				return mockVouchers, nil
			}
			return []*voucher_model.Voucher{}, nil
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	brandId := int32(1)
	req := &pbVoucher.ListVoucherReq{
		BrandId: &brandId,
	}

	result, err := service.ListVoucher(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 1 {
		t.Errorf("Expected 1 voucher, got %d", len(result.Data))
	}
	if result.Data[0].BrandId != 1 {
		t.Errorf("Expected voucher brand ID to be 1, got %d", result.Data[0].BrandId)
	}
}

func TestDetailVoucher_Success(t *testing.T) {
	now := time.Now()
	mockVoucher := &voucher_model.Voucher{
		ID:           1,
		BrandID:      1,
		Name:         "Test Voucher",
		Description:  "Test Description",
		VoucherCode:  "TEST001",
		CostInPoint:  100,
		CreatedDate:  now,
		ModifiedDate: now,
		IsDeleted:    false,
	}

	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.DetailVoucherReq{
		Id: 1,
	}

	result, err := service.DetailVoucher(context.Background(), req)

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
		t.Errorf("Expected voucher ID to be 1, got %d", result.Data.Id)
	}
	if result.Data.BrandId != 1 {
		t.Errorf("Expected brand ID to be 1, got %d", result.Data.BrandId)
	}
	if result.Data.Name != "Test Voucher" {
		t.Errorf("Expected name to be 'Test Voucher', got '%s'", result.Data.Name)
	}
	if result.Data.Description != "Test Description" {
		t.Errorf("Expected description to be 'Test Description', got '%s'", result.Data.Description)
	}
	if result.Data.VoucherCode != "TEST001" {
		t.Errorf("Expected voucher code to be 'TEST001', got '%s'", result.Data.VoucherCode)
	}
	if result.Data.CostInPoint != 100 {
		t.Errorf("Expected cost in point to be 100, got %d", result.Data.CostInPoint)
	}
	if *result.Data.IsDeleted {
		t.Error("Expected voucher to not be deleted")
	}
}

func TestDetailVoucher_NotFound(t *testing.T) {
	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			return nil, errors.New("voucher not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.DetailVoucherReq{
		Id: 999,
	}

	result, err := service.DetailVoucher(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected result to be nil")
	}
}

func TestDetailVoucher_WithDeletedVoucher(t *testing.T) {
	now := time.Now()
	mockVoucher := &voucher_model.Voucher{
		ID:           1,
		BrandID:      1,
		Name:         "Deleted Voucher",
		Description:  "Deleted Description",
		VoucherCode:  "DELETED001",
		CostInPoint:  100,
		CreatedDate:  now,
		ModifiedDate: now,
		IsDeleted:    true,
	}

	mockVoucherRepo := &MockVoucherRepo{
		findByIdFunc: func(id uint) (*voucher_model.Voucher, error) {
			if id == 1 {
				return mockVoucher, nil
			}
			return nil, errors.New("voucher not found")
		},
	}

	service := &VoucherService{
		voucherRepo: mockVoucherRepo,
	}

	req := &pbVoucher.DetailVoucherReq{
		Id: 1,
	}

	result, err := service.DetailVoucher(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if result.Data == nil {
		t.Error("Expected Data to not be nil")
	}
	if !*result.Data.IsDeleted {
		t.Error("Expected voucher to be deleted")
	}
}
