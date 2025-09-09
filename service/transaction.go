// service/transaction_service.go
package service

import (
	"errors"
	"fmt"
	"github.com/KQLXK/Family-Finance-System/model"
	"time"
)

// TransactionService 交易服务接口
type TransactionService interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransactionByID(id uint) (*model.Transaction, error)
	GetTransactionsByFamilyID(familyID uint, page, pageSize int, filters map[string]interface{}) ([]model.Transaction, int64, error)
	GetTransactionsByTimeRange(familyID uint, startTime, endTime time.Time, filters map[string]interface{}) ([]model.Transaction, error)
	UpdateTransaction(transaction *model.Transaction) error
	DeleteTransaction(id uint) error
	AddTagToTransaction(transactionID, tagID uint) error
	RemoveTagFromTransaction(transactionID, tagID uint) error
	GetTransactionSummaryByCategory(familyID uint, startTime, endTime time.Time, transactionType model.TransactionType) (map[string]float64, error)
	GetTransactionSummaryByTime(familyID uint, startTime, endTime time.Time, groupBy string) (map[string]float64, error)
}

// transactionService 交易服务实现
type transactionService struct {
	transactionDao model.TransactionDao
	familyDao      model.FamilyDao
	memberDao      model.MemberDao
	categoryDao    model.CategoryDao
	tagDao         model.TagDao
}

// NewTransactionService 创建交易服务实例
func NewTransactionService() TransactionService {
	return &transactionService{
		transactionDao: *model.NewTransactionDaoInstance(),
		familyDao:      *model.NewFamilyDaoInstance(),
		memberDao:      *model.NewMemberDaoInstance(),
		categoryDao:    *model.NewCategoryDaoInstance(),
		tagDao:         *model.NewTagDaoInstance(),
	}
}

// CreateTransaction 创建交易
func (s *transactionService) CreateTransaction(transaction *model.Transaction) error {
	// 验证交易数据
	if err := s.validateTransaction(transaction); err != nil {
		return err
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(transaction.FamilyID)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return errors.New("关联的家庭不存在")
	}

	// 检查成员是否存在且属于该家庭
	memberExists, err := s.memberExistsInFamily(transaction.MemberID, transaction.FamilyID)
	if err != nil {
		return fmt.Errorf("检查成员是否存在时出错: %v", err)
	}
	if !memberExists {
		return errors.New("成员不存在或不属于该家庭")
	}

	// 检查分类是否存在且类型匹配
	categoryExists, err := s.categoryExistsAndMatchesType(transaction.CategoryID, transaction.Type)
	if err != nil {
		return fmt.Errorf("检查分类是否存在时出错: %v", err)
	}
	if !categoryExists {
		return errors.New("分类不存在或类型不匹配")
	}

	// 设置默认状态
	transaction.Status = model.Valid
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	// 创建交易
	if err := s.transactionDao.CreateTransaction(transaction); err != nil {
		return fmt.Errorf("创建交易失败: %v", err)
	}

	return nil
}

// GetTransactionByID 根据ID获取交易
func (s *transactionService) GetTransactionByID(id uint) (*model.Transaction, error) {
	// 验证ID
	if id == 0 {
		return nil, errors.New("无效的交易ID")
	}

	// 获取交易
	transaction, err := s.transactionDao.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取交易失败: %v", err)
	}

	if transaction == nil || transaction.Status == model.Deleted {
		return nil, errors.New("交易不存在或已被删除")
	}

	return transaction, nil
}

// GetTransactionsByFamilyID 根据家庭ID获取交易列表
func (s *transactionService) GetTransactionsByFamilyID(familyID uint, page, pageSize int, filters map[string]interface{}) ([]model.Transaction, int64, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, 0, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, 0, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, 0, errors.New("家庭不存在")
	}

	// 构建查询条件
	query := map[string]interface{}{
		"family_id": familyID,
		"status":    model.Valid,
	}

	// 添加过滤条件
	for key, value := range filters {
		query[key] = value
	}

	// 获取交易列表
	transactions, total, err := s.transactionDao.GetTransactionsByFamilyID(familyID, page, pageSize, query)
	if err != nil {
		return nil, 0, fmt.Errorf("获取交易列表失败: %v", err)
	}

	return transactions, total, nil
}

// GetTransactionsByTimeRange 根据时间范围获取交易列表
func (s *transactionService) GetTransactionsByTimeRange(familyID uint, startTime, endTime time.Time, filters map[string]interface{}) ([]model.Transaction, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 构建查询条件
	query := map[string]interface{}{
		"family_id": familyID,
		"status":    model.Valid,
	}

	// 添加过滤条件
	for key, value := range filters {
		query[key] = value
	}

	// 获取交易列表
	transactions, err := s.transactionDao.GetTransactionsByTimeRange(familyID, startTime, endTime, query)
	if err != nil {
		return nil, fmt.Errorf("获取时间段交易失败: %v", err)
	}

	return transactions, nil
}

// UpdateTransaction 更新交易信息
func (s *transactionService) UpdateTransaction(transaction *model.Transaction) error {
	// 验证交易ID
	if transaction.ID == 0 {
		return errors.New("无效的交易ID")
	}

	// 验证交易数据
	if err := s.validateTransaction(transaction); err != nil {
		return err
	}

	// 检查交易是否存在
	existingTransaction, err := s.transactionDao.GetTransactionByID(transaction.ID)
	if err != nil {
		return fmt.Errorf("检查交易是否存在时出错: %v", err)
	}
	if existingTransaction == nil || existingTransaction.Status == model.Deleted {
		return errors.New("交易不存在或已被删除")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(transaction.FamilyID)
	if err != nil {
		return fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return errors.New("关联的家庭不存在")
	}

	// 检查成员是否存在且属于该家庭
	memberExists, err := s.memberExistsInFamily(transaction.MemberID, transaction.FamilyID)
	if err != nil {
		return fmt.Errorf("检查成员是否存在时出错: %v", err)
	}
	if !memberExists {
		return errors.New("成员不存在或不属于该家庭")
	}

	// 检查分类是否存在且类型匹配
	categoryExists, err := s.categoryExistsAndMatchesType(transaction.CategoryID, transaction.Type)
	if err != nil {
		return fmt.Errorf("检查分类是否存在时出错: %v", err)
	}
	if !categoryExists {
		return errors.New("分类不存在或类型不匹配")
	}

	// 更新交易信息
	transaction.UpdatedAt = time.Now()
	if err := s.transactionDao.UpdateTransaction(transaction); err != nil {
		return fmt.Errorf("更新交易信息失败: %v", err)
	}

	return nil
}

// DeleteTransaction 删除交易（软删除）
func (s *transactionService) DeleteTransaction(id uint) error {
	// 验证ID
	if id == 0 {
		return errors.New("无效的交易ID")
	}

	// 检查交易是否存在
	transaction, err := s.transactionDao.GetTransactionByID(id)
	if err != nil {
		return fmt.Errorf("检查交易是否存在时出错: %v", err)
	}
	if transaction == nil || transaction.Status == model.Deleted {
		return errors.New("交易不存在或已被删除")
	}

	// 软删除交易（设置状态为deleted）
	if err := s.transactionDao.DeleteTransaction(id); err != nil {
		return fmt.Errorf("删除交易失败: %v", err)
	}

	return nil
}

// AddTagToTransaction 为交易添加标签
func (s *transactionService) AddTagToTransaction(transactionID, tagID uint) error {
	// 验证ID
	if transactionID == 0 || tagID == 0 {
		return errors.New("无效的交易ID或标签ID")
	}

	// 检查交易是否存在
	transaction, err := s.transactionDao.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("检查交易是否存在时出错: %v", err)
	}
	if transaction == nil || transaction.Status == model.Deleted {
		return errors.New("交易不存在或已被删除")
	}

	// 检查标签是否存在
	tag, err := s.tagDao.GetTagByID(tagID)
	if err != nil {
		return fmt.Errorf("检查标签是否存在时出错: %v", err)
	}
	if tag == nil || !tag.IsActive {
		return errors.New("标签不存在或已被禁用")
	}

	// 检查标签是否属于同一个家庭
	if tag.FamilyID != transaction.FamilyID {
		return errors.New("标签不属于该交易的家庭")
	}

	// 检查标签是否已经添加到交易
	exists, err := s.transactionDao.TagExistsInTransaction(transactionID, tagID)
	if err != nil {
		return fmt.Errorf("检查标签是否已存在时出错: %v", err)
	}
	if exists {
		return errors.New("标签已添加到该交易")
	}

	// 添加标签到交易
	if err := s.transactionDao.AddTagToTransaction(transactionID, tagID); err != nil {
		return fmt.Errorf("添加标签到交易失败: %v", err)
	}

	return nil
}

// RemoveTagFromTransaction 从交易移除标签
func (s *transactionService) RemoveTagFromTransaction(transactionID, tagID uint) error {
	// 验证ID
	if transactionID == 0 || tagID == 0 {
		return errors.New("无效的交易ID或标签ID")
	}

	// 检查交易是否存在
	transaction, err := s.transactionDao.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("检查交易是否存在时出错: %v", err)
	}
	if transaction == nil || transaction.Status == model.Deleted {
		return errors.New("交易不存在或已被删除")
	}

	// 检查标签是否存在于交易
	exists, err := s.transactionDao.TagExistsInTransaction(transactionID, tagID)
	if err != nil {
		return fmt.Errorf("检查标签是否存在于交易时出错: %v", err)
	}
	if !exists {
		return errors.New("标签不存在于该交易")
	}

	// 从交易移除标签
	if err := s.transactionDao.RemoveTagFromTransaction(transactionID, tagID); err != nil {
		return fmt.Errorf("从交易移除标签失败: %v", err)
	}

	return nil
}

// GetTransactionSummaryByCategory 按分类统计交易金额
func (s *transactionService) GetTransactionSummaryByCategory(familyID uint, startTime, endTime time.Time, transactionType model.TransactionType) (map[string]float64, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 获取分类统计
	summary, err := s.transactionDao.GetTransactionSummaryByCategory(familyID, startTime, endTime, transactionType)
	if err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %v", err)
	}

	return summary, nil
}

// GetTransactionSummaryByTime 按时间统计交易金额
func (s *transactionService) GetTransactionSummaryByTime(familyID uint, startTime, endTime time.Time, groupBy string) (map[string]float64, error) {
	// 验证家庭ID
	if familyID == 0 {
		return nil, errors.New("无效的家庭ID")
	}

	// 检查家庭是否存在
	familyExists, err := s.familyExists(familyID)
	if err != nil {
		return nil, fmt.Errorf("检查家庭是否存在时出错: %v", err)
	}
	if !familyExists {
		return nil, errors.New("家庭不存在")
	}

	// 验证分组方式
	if groupBy != "day" && groupBy != "month" && groupBy != "year" {
		return nil, errors.New("无效的分组方式，支持: day, month, year")
	}

	// 获取时间统计
	summary, err := s.transactionDao.GetTransactionSummaryByTime(familyID, startTime, endTime, groupBy)
	if err != nil {
		return nil, fmt.Errorf("获取时间统计失败: %v", err)
	}

	return summary, nil
}

// validateTransaction 验证交易数据
func (s *transactionService) validateTransaction(transaction *model.Transaction) error {
	// 验证金额
	if transaction.Amount <= 0 {
		return errors.New("交易金额必须大于0")
	}

	// 验证类型
	if !s.isValidTransactionType(transaction.Type) {
		return errors.New("无效的交易类型")
	}

	// 验证交易时间
	if transaction.TransactionTime.IsZero() {
		return errors.New("交易时间不能为空")
	}

	// 验证交易时间不能晚于当前时间
	if transaction.TransactionTime.After(time.Now()) {
		return errors.New("交易时间不能晚于当前时间")
	}

	// 验证备注长度
	if len(transaction.Note) > 1000 {
		return errors.New("备注长度不能超过1000个字符")
	}

	// 验证支付方式长度
	if len(transaction.PaymentMethod) > 50 {
		return errors.New("支付方式长度不能超过50个字符")
	}

	return nil
}

// familyExists 检查家庭是否存在
func (s *transactionService) familyExists(familyID uint) (bool, error) {
	if familyID == 0 {
		return false, nil
	}

	family, err := s.familyDao.GetFamilyByID(familyID)
	if err != nil {
		return false, err
	}

	return family != nil, nil
}

// memberExistsInFamily 检查成员是否存在且属于该家庭
func (s *transactionService) memberExistsInFamily(memberID, familyID uint) (bool, error) {
	if memberID == 0 || familyID == 0 {
		return false, nil
	}

	member, err := s.memberDao.GetMemberByID(memberID)
	if err != nil {
		return false, err
	}

	return member != nil && member.FamilyID == familyID && member.Status == 1, nil
}

// categoryExistsAndMatchesType 检查分类是否存在且类型匹配
func (s *transactionService) categoryExistsAndMatchesType(categoryID uint, transactionType model.TransactionType) (bool, error) {
	if categoryID == 0 {
		return false, nil
	}

	category, err := s.categoryDao.GetCategoryByID(categoryID)
	if err != nil {
		return false, err
	}

	if category == nil || category.IsDeleted {
		return false, nil
	}

	// 检查分类类型是否与交易类型匹配
	var expectedCategoryType model.CategoryType
	if transactionType == model.Income {
		expectedCategoryType = model.CategoryIncome
	} else {
		expectedCategoryType = model.CategoryExpense
	}

	return category.Type == expectedCategoryType, nil
}

// isValidTransactionType 验证交易类型是否有效
func (s *transactionService) isValidTransactionType(transactionType model.TransactionType) bool {
	switch transactionType {
	case model.Income, model.Expense:
		return true
	default:
		return false
	}
}
