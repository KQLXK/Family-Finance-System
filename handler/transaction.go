// handler/transaction_handler.go
package handler

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"github.com/KQLXK/Family-Finance-System/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TransactionHandler 交易处理器
type TransactionHandler struct {
	transactionService service.TransactionService
}

// NewTransactionHandler 创建交易处理器
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		transactionService: service.NewTransactionService(),
	}
}

// CreateTransaction 创建交易
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction model.Transaction
	familyIdstr := c.Param("id")
	familyId, _ := strconv.ParseUint(familyIdstr, 10, 32)
	transaction.FamilyID = uint(familyId)
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.transactionService.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "交易创建成功",
		"data":    transaction,
	})
}

// GetTransactionByID 根据ID获取交易
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

// GetTransactionsByFamilyID 根据家庭ID获取交易列表
func (h *TransactionHandler) GetTransactionsByFamilyID(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取过滤参数
	filters := make(map[string]interface{})
	if typeStr := c.Query("type"); typeStr != "" {
		filters["type"] = model.TransactionType(typeStr)
	}
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			filters["category_id"] = uint(categoryID)
		}
	}
	if memberIDStr := c.Query("memberId"); memberIDStr != "" {
		memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
		if err == nil {
			filters["member_id"] = uint(memberID)
		}
	}
	if paymentMethod := c.Query("paymentMethod"); paymentMethod != "" {
		filters["payment_method"] = paymentMethod
	}

	transactions, total, err := h.transactionService.GetTransactionsByFamilyID(uint(familyID), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetTransactionsByTimeRange 根据时间范围获取交易列表
func (h *TransactionHandler) GetTransactionsByTimeRange(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	// 获取时间范围参数
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认开始时间为30天前
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认结束时间为当前时间
		endTime = time.Now()
	}

	// 获取过滤参数
	filters := make(map[string]interface{})
	if typeStr := c.Query("type"); typeStr != "" {
		filters["type"] = model.TransactionType(typeStr)
	}
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			filters["category_id"] = uint(categoryID)
		}
	}
	if memberIDStr := c.Query("memberId"); memberIDStr != "" {
		memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
		if err == nil {
			filters["member_id"] = uint(memberID)
		}
	}

	transactions, err := h.transactionService.GetTransactionsByTimeRange(uint(familyID), startTime, endTime, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// UpdateTransaction 更新交易信息
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	var transaction model.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	transaction.ID = uint(id)
	if err := h.transactionService.UpdateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "交易信息更新成功",
		"data":    transaction,
	})
}

// DeleteTransaction 删除交易
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	if err := h.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "交易删除成功",
	})
}

// AddTagToTransaction 为交易添加标签
func (h *TransactionHandler) AddTagToTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	var request struct {
		TagID uint `json:"tag_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.transactionService.AddTagToTransaction(uint(transactionID), request.TagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签添加成功",
	})
}

// RemoveTagFromTransaction 从交易移除标签
func (h *TransactionHandler) RemoveTagFromTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	tagIDStr := c.Param("tagId")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	if err := h.transactionService.RemoveTagFromTransaction(uint(transactionID), uint(tagID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "标签移除成功",
	})
}

// GetTransactionSummaryByCategory 按分类统计交易金额
func (h *TransactionHandler) GetTransactionSummaryByCategory(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	// 获取时间范围参数
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认开始时间为30天前
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认结束时间为当前时间
		endTime = time.Now()
	}

	// 获取交易类型
	transactionType := model.TransactionType(c.Query("type"))
	if transactionType == "" {
		transactionType = model.Expense // 默认统计支出
	}

	summary, err := h.transactionService.GetTransactionSummaryByCategory(uint(familyID), startTime, endTime, transactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

// GetTransactionSummaryByTime 按时间统计交易金额
func (h *TransactionHandler) GetTransactionSummaryByTime(c *gin.Context) {
	familyIDStr := c.Param("id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的家庭ID"})
		return
	}

	// 获取时间范围参数
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认开始时间为30天前
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式，请使用RFC3339格式"})
			return
		}
	} else {
		// 默认结束时间为当前时间
		endTime = time.Now()
	}

	// 获取分组方式
	groupBy := c.DefaultQuery("groupBy", "month")

	summary, err := h.transactionService.GetTransactionSummaryByTime(uint(familyID), startTime, endTime, groupBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}
