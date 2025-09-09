package model

import (
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
)

type TransactionTagDao struct{}

var (
	transactionTagOnce sync.Once
	transactionTagDao  *TransactionTagDao
)

// NewTransactionTagDaoInstance 返回 TransactionTagDao 单例实例
func NewTransactionTagDaoInstance() *TransactionTagDao {
	transactionTagOnce.Do(func() {
		transactionTagDao = &TransactionTagDao{}
	})
	return transactionTagDao
}

// CreateTransactionTag 创建交易标签关联
func (TransactionTagDao) CreateTransactionTag(transactionTag *TransactionTag) error {
	if err := database.DB.Create(transactionTag).Error; err != nil {
		log.Printf("创建交易标签关联失败: %v", err)
		return err
	}
	return nil
}

// GetTransactionTagByID 根据ID获取交易标签关联
func (TransactionTagDao) GetTransactionTagByID(id uint) (*TransactionTag, error) {
	var transactionTag TransactionTag
	if err := database.DB.Preload("Transaction").Preload("Tag").First(&transactionTag, id).Error; err != nil {
		log.Printf("获取交易标签关联失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &transactionTag, nil
}

// GetTagsByTransactionID 根据交易ID获取标签列表
func (TransactionTagDao) GetTagsByTransactionID(transactionID uint) ([]Tag, error) {
	var tags []Tag
	if err := database.DB.Joins("JOIN transaction_tags ON transaction_tags.tag_id = tags.id").
		Where("transaction_tags.transaction_id = ?", transactionID).
		Find(&tags).Error; err != nil {
		log.Printf("获取交易标签失败 TransactionID=%d: %v", transactionID, err)
		return nil, err
	}
	return tags, nil
}

// GetTransactionsByTagID 根据标签ID获取交易列表
func (TransactionTagDao) GetTransactionsByTagID(tagID uint) ([]Transaction, error) {
	var transactions []Transaction
	if err := database.DB.Joins("JOIN transaction_tags ON transaction_tags.transaction_id = transactions.id").
		Where("transaction_tags.tag_id = ? AND transactions.status = 'valid'", tagID).
		Preload("Member").Preload("Category").
		Order("transactions.transaction_time DESC").
		Find(&transactions).Error; err != nil {
		log.Printf("获取标签交易失败 TagID=%d: %v", tagID, err)
		return nil, err
	}
	return transactions, nil
}

// DeleteTransactionTag 删除交易标签关联
func (TransactionTagDao) DeleteTransactionTag(id uint) error {
	if err := database.DB.Delete(&TransactionTag{}, id).Error; err != nil {
		log.Printf("删除交易标签关联失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}
