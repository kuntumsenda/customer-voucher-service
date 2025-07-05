package transaction_model

import (
	pb "customer-voucher-service/protogen/transaction"

	"gorm.io/gorm"
)

type ITransactionRepo interface {
	CreateTransaction(transaction *Transaction) (*Transaction, error)
	FindTransactionById(id uint) (*Transaction, error)
	ListTransaction(req *pb.ListTransactionReq) ([]*Transaction, error)
	DetailTransaction(req *pb.DetailTransactionReq) (*Transaction, error)
}

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (r *TransactionRepo) CreateTransaction(transaction *Transaction) (*Transaction, error) {
	err := r.db.Create(transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepo) FindTransactionById(id uint) (*Transaction, error) {
	var transaction Transaction
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepo) ListTransaction(req *pb.ListTransactionReq) ([]*Transaction, error) {
	var transactions []*Transaction
	query := r.db.Model(&Transaction{}).Where("is_deleted = ?", false)

	if req.CustomerId != nil {
		query = query.Where("customer_id = ?", *req.CustomerId)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepo) DetailTransaction(req *pb.DetailTransactionReq) (*Transaction, error) {
	return r.FindTransactionById(uint(req.Id))
}
