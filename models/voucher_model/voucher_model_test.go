package voucher_model

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn: db,
		DSN:  "sqlmock_db_0",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}
	return gormDB, mock, func() { db.Close() }
}

func TestCreateVoucher(t *testing.T) {
	db, mock, closeFn := setupMockDB(t)
	defer closeFn()
	repo := NewVoucherRepo(db)

	voucher := &Voucher{
		BrandID:     1,
		Name:        "Voucher Test",
		Description: "Desc",
		CostInPoint: 100,
		VoucherCode: "CODE123",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "voucher"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateVoucher(voucher)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), voucher.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListVoucher(t *testing.T) {
	db, mock, closeFn := setupMockDB(t)
	defer closeFn()
	repo := NewVoucherRepo(db)

	mockRows := sqlmock.NewRows([]string{"id", "brand_id", "name", "description", "cost_in_point", "voucher_code", "created_date", "modified_date", "is_deleted"}).
		AddRow(1, 1, "Voucher 1", "Desc 1", 100, "CODE1", time.Now(), time.Now(), false).
		AddRow(2, 2, "Voucher 2", "Desc 2", 200, "CODE2", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "voucher" WHERE is_deleted = $1`)).
		WithArgs(false).
		WillReturnRows(mockRows)

	vouchers, err := repo.ListVoucher(&struct{ BrandId *int32 }{})
	assert.NoError(t, err)
	assert.Len(t, vouchers, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindVoucherById(t *testing.T) {
	db, mock, closeFn := setupMockDB(t)
	defer closeFn()
	repo := NewVoucherRepo(db)

	now := time.Now()
	mockRows := sqlmock.NewRows([]string{"id", "brand_id", "name", "description", "cost_in_point", "voucher_code", "created_date", "modified_date", "is_deleted"}).
		AddRow(1, 1, "Voucher 1", "Desc 1", 100, "CODE1", now, now, false)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "voucher" WHERE id = $1 AND is_deleted = $2 ORDER BY "voucher"."id" LIMIT 1`)).
		WithArgs(1, false).
		WillReturnRows(mockRows)

	voucher, err := repo.FindVoucherById(1)
	assert.NoError(t, err)
	assert.NotNil(t, voucher)
	assert.Equal(t, uint(1), voucher.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
