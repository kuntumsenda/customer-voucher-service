package brand_service

import (
	"context"
	"customer-voucher-service/models/brand_model"
	pbBrand "customer-voucher-service/protogen/brand"
	"errors"
	"testing"
	"time"
)

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

func TestCreateBrand_Success(t *testing.T) {
	mockRepo := &MockBrandRepo{
		createBrandFunc: func(brand *brand_model.Brand) error {
			if brand.Name != "Test Brand" || brand.Description != "Test Description" {
				t.Errorf("Expected brand with name 'Test Brand' and description 'Test Description', got name '%s' and description '%s'", brand.Name, brand.Description)
			}
			return nil
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.CreateBrandReq{
		Name:        "Test Brand",
		Description: "Test Description",
	}

	result, err := service.CreateBrand(context.Background(), req)

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

func TestCreateBrand_ValidationError_EmptyName(t *testing.T) {
	mockRepo := &MockBrandRepo{}
	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.CreateBrandReq{
		Name:        "",
		Description: "Test Description",
	}

	result, err := service.CreateBrand(context.Background(), req)

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

func TestCreateBrand_ValidationError_NameTooLong(t *testing.T) {
	mockRepo := &MockBrandRepo{}
	service := &BrandService{
		brandRepo: mockRepo,
	}

	longName := ""
	for i := 0; i < 256; i++ {
		longName += "a"
	}

	req := &pbBrand.CreateBrandReq{
		Name:        longName,
		Description: "Test Description",
	}

	result, err := service.CreateBrand(context.Background(), req)

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

func TestCreateBrand_ValidationError_DescriptionTooLong(t *testing.T) {
	mockRepo := &MockBrandRepo{}
	service := &BrandService{
		brandRepo: mockRepo,
	}

	longDescription := ""
	for i := 0; i < 256; i++ {
		longDescription += "a"
	}

	req := &pbBrand.CreateBrandReq{
		Name:        "Test Brand",
		Description: longDescription,
	}

	result, err := service.CreateBrand(context.Background(), req)

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

func TestCreateBrand_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")
	mockRepo := &MockBrandRepo{
		createBrandFunc: func(brand *brand_model.Brand) error {
			return expectedError
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.CreateBrandReq{
		Name:        "Test Brand",
		Description: "Test Description",
	}

	result, err := service.CreateBrand(context.Background(), req)

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

func TestListBrand_Success(t *testing.T) {
	now := time.Now()
	mockBrands := []*brand_model.Brand{
		{
			ID:           1,
			Name:         "Brand 1",
			Description:  "Description 1",
			IsDeleted:    false,
			CreatedDate:  now,
			ModifiedDate: now,
		},
		{
			ID:           2,
			Name:         "Brand 2",
			Description:  "Description 2",
			IsDeleted:    false,
			CreatedDate:  now,
			ModifiedDate: now,
		},
	}

	mockRepo := &MockBrandRepo{
		listBrandFunc: func() ([]*brand_model.Brand, error) {
			return mockBrands, nil
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.ListBrandReq{}

	result, err := service.ListBrand(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 brands, got %d", len(result.Data))
	}

	if result.Data[0].Id != 1 {
		t.Errorf("Expected first brand ID to be 1, got %d", result.Data[0].Id)
	}
	if result.Data[0].Name != "Brand 1" {
		t.Errorf("Expected first brand name to be 'Brand 1', got '%s'", result.Data[0].Name)
	}
	if result.Data[0].Description != "Description 1" {
		t.Errorf("Expected first brand description to be 'Description 1', got '%s'", result.Data[0].Description)
	}
	if *result.Data[0].IsDeleted {
		t.Error("Expected first brand to not be deleted")
	}

	if result.Data[1].Id != 2 {
		t.Errorf("Expected second brand ID to be 2, got %d", result.Data[1].Id)
	}
	if result.Data[1].Name != "Brand 2" {
		t.Errorf("Expected second brand name to be 'Brand 2', got '%s'", result.Data[1].Name)
	}
	if result.Data[1].Description != "Description 2" {
		t.Errorf("Expected second brand description to be 'Description 2', got '%s'", result.Data[1].Description)
	}
	if *result.Data[1].IsDeleted {
		t.Error("Expected second brand to not be deleted")
	}
}

func TestListBrand_EmptyList(t *testing.T) {
	mockRepo := &MockBrandRepo{
		listBrandFunc: func() ([]*brand_model.Brand, error) {
			return []*brand_model.Brand{}, nil
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.ListBrandReq{}

	result, err := service.ListBrand(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 0 {
		t.Errorf("Expected 0 brands, got %d", len(result.Data))
	}
}

func TestListBrand_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")
	mockRepo := &MockBrandRepo{
		listBrandFunc: func() ([]*brand_model.Brand, error) {
			return []*brand_model.Brand{}, expectedError
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.ListBrandReq{}

	result, err := service.ListBrand(context.Background(), req)

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

func TestListBrand_WithDeletedBrands(t *testing.T) {
	now := time.Now()
	mockBrands := []*brand_model.Brand{
		{
			ID:           1,
			Name:         "Active Brand",
			Description:  "Active Description",
			IsDeleted:    false,
			CreatedDate:  now,
			ModifiedDate: now,
		},
		{
			ID:           2,
			Name:         "Deleted Brand",
			Description:  "Deleted Description",
			IsDeleted:    true,
			CreatedDate:  now,
			ModifiedDate: now,
		},
	}

	mockRepo := &MockBrandRepo{
		listBrandFunc: func() ([]*brand_model.Brand, error) {
			return mockBrands, nil
		},
	}

	service := &BrandService{
		brandRepo: mockRepo,
	}

	req := &pbBrand.ListBrandReq{}

	result, err := service.ListBrand(context.Background(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
	}
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 brands, got %d", len(result.Data))
	}

	if *result.Data[0].IsDeleted {
		t.Error("Expected first brand to not be deleted")
	}
	if !*result.Data[1].IsDeleted {
		t.Error("Expected second brand to be deleted")
	}
}
