package model

import (
	"fmt"
	"github.com/KQLXK/Family-Finance-System/database"
	"log"
	"sync"
	"time"
)

type TransactionDao struct{}

var (
	transactionOnce sync.Once
	transactionDao  *TransactionDao
)

// NewTransactionDaoInstance 返回 TransactionDao 单例实例
func NewTransactionDaoInstance() *TransactionDao {
	transactionOnce.Do(func() {
		transactionDao = &TransactionDao{}
	})
	return transactionDao
}

// CreateTransaction 创建交易
func (TransactionDao) CreateTransaction(transaction *Transaction) error {
	if err := database.DB.Create(transaction).Error; err != nil {
		log.Printf("创建交易失败: %v", err)
		return err
	}
	return nil
}

// GetTransactionByID 根据ID获取交易
func (TransactionDao) GetTransactionByID(id uint) (*Transaction, error) {
	var transaction Transaction
	if err := database.DB.Preload("Family").Preload("Member").Preload("Category").
		Preload("Labels").First(&transaction, id).Error; err != nil {
		log.Printf("获取交易失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &transaction, nil
}

func (TransactionDao) GetTransactionsByFamilyID(familyID uint, page, pageSize int, filters map[string]interface{}) ([]Transaction, int64, error) {
	var transactions []Transaction
	var total int64

	// 构建查询
	query := database.DB.Model(&Transaction{}).Where("family_id = ? AND status = ?", familyID, Valid)

	// 添加过滤条件
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取交易总数失败 FamilyID=%d: %v", familyID, err)
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("Member").Preload("Category").Preload("Labels").
		Order("transaction_time DESC").
		Offset(offset).Limit(pageSize).
		Find(&transactions).Error; err != nil {
		log.Printf("获取交易列表失败 FamilyID=%d: %v", familyID, err)
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionsByTimeRange 根据时间范围获取交易列表
func (TransactionDao) GetTransactionsByTimeRange(familyID uint, startTime, endTime time.Time, filters map[string]interface{}) ([]Transaction, error) {
	var transactions []Transaction

	// 构建查询
	query := database.DB.Where("family_id = ? AND status = ? AND transaction_time BETWEEN ? AND ?",
		familyID, Valid, startTime, endTime)

	// 添加过滤条件
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 获取数据
	if err := query.Preload("Member").Preload("Category").Preload("Labels").
		Order("transaction_time DESC").
		Find(&transactions).Error; err != nil {
		log.Printf("获取时间段交易失败 FamilyID=%d: %v", familyID, err)
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction 更新交易信息
func (TransactionDao) UpdateTransaction(transaction *Transaction) error {
	if err := database.DB.Save(transaction).Error; err != nil {
		log.Printf("更新交易失败 ID=%d: %v", transaction.ID, err)
		return err
	}
	return nil
}

// DeleteTransaction 软删除交易
func (TransactionDao) DeleteTransaction(id uint) error {
	if err := database.DB.Model(&Transaction{}).Where("id = ?", id).Update("status", "deleted").Error; err != nil {
		log.Printf("删除交易失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}

// AddTagToTransaction 为交易添加标签
func (TransactionDao) AddTagToTransaction(transactionID, tagID uint) error {
	transactionTag := TransactionTag{
		TransactionID: transactionID,
		TagID:         tagID,
	}
	if err := database.DB.Create(&transactionTag).Error; err != nil {
		log.Printf("为交易添加标签失败 TransactionID=%d, TagID=%d: %v", transactionID, tagID, err)
		return err
	}
	return nil
}

// RemoveTagFromTransaction 从交易移除标签
func (TransactionDao) RemoveTagFromTransaction(transactionID, tagID uint) error {
	if err := database.DB.Where("transaction_id = ? AND tag_id = ?", transactionID, tagID).
		Delete(&TransactionTag{}).Error; err != nil {
		log.Printf("从交易移除标签失败 TransactionID=%d, TagID=%d: %v", transactionID, tagID, err)
		return err
	}
	return nil
}

func (TransactionDao) GetTransactionSummaryByCategory(familyID uint, startTime, endTime time.Time, transactionType TransactionType) (map[string]float64, error) {
	summary := make(map[string]float64)

	// 执行SQL查询
	rows, err := database.DB.Table("transactions").
		Select("categories.name, SUM(transactions.amount)").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.family_id = ? AND transactions.status = ? AND transactions.type = ? AND transactions.transaction_time BETWEEN ? AND ?",
			familyID, Valid, transactionType, startTime, endTime).
		Group("categories.name").
		Rows()
	if err != nil {
		log.Printf("按分类统计交易金额失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var categoryName string
		var totalAmount float64
		if err := rows.Scan(&categoryName, &totalAmount); err != nil {
			log.Printf("扫描统计结果失败: %v", err)
			continue
		}
		summary[categoryName] = totalAmount
	}

	return summary, nil
}

// GetTransactionSummaryByTime 按时间统计交易金额
func (TransactionDao) GetTransactionSummaryByTime(familyID uint, startTime, endTime time.Time, groupBy string) (map[string]float64, error) {
	summary := make(map[string]float64)

	// 根据分组方式构建SQL
	var timeFormat string
	switch groupBy {
	case "day":
		timeFormat = "%Y-%m-%d"
	case "month":
		timeFormat = "%Y-%m"
	case "year":
		timeFormat = "%Y"
	default:
		timeFormat = "%Y-%m"
	}

	// 执行SQL查询
	rows, err := database.DB.Table("transactions").
		Select(fmt.Sprintf("DATE_FORMAT(transaction_time, '%s') as time_period, SUM(amount)"), timeFormat).
		Where("family_id = ? AND status = ? AND transaction_time BETWEEN ? AND ?",
			familyID, Valid, startTime, endTime).
		Group("time_period").
		Rows()
	if err != nil {
		log.Printf("按时间统计交易金额失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var timePeriod string
		var totalAmount float64
		if err := rows.Scan(&timePeriod, &totalAmount); err != nil {
			log.Printf("扫描统计结果失败: %v", err)
			continue
		}
		summary[timePeriod] = totalAmount
	}

	return summary, nil
}

// TagExistsInTransaction 检查标签是否存在于交易
func (TransactionDao) TagExistsInTransaction(transactionID, tagID uint) (bool, error) {
	var count int64
	if err := database.DB.Model(&TransactionTag{}).
		Where("transaction_id = ? AND tag_id = ?", transactionID, tagID).
		Count(&count).Error; err != nil {
		log.Printf("检查标签是否存在于交易失败 TransactionID=%d, TagID=%d: %v", transactionID, tagID, err)
		return false, err
	}
	return count > 0, nil
}
