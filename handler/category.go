// handler/category_handler.go
package handler

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"github.com/KQLXK/Family-Finance-System/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.categoryService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "分类创建成功",
		"data":    category,
	})
}

// GetCategoryByID 根据ID获取分类
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// GetCategoriesByType 根据类型获取分类列表
func (h *CategoryHandler) GetCategoriesByType(c *gin.Context) {
	categoryType := model.CategoryType(c.Query("type"))
	if categoryType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类类型不能为空"})
		return
	}

	categories, err := h.categoryService.GetCategoriesByType(categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetCategoryTreeByType 根据类型获取分类树
func (h *CategoryHandler) GetCategoryTreeByType(c *gin.Context) {
	categoryType := model.CategoryType(c.Query("type"))
	if categoryType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类类型不能为空"})
		return
	}

	categories, err := h.categoryService.GetCategoryTreeByType(categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetCategoriesByParentID 根据父分类ID获取子分类列表
func (h *CategoryHandler) GetCategoriesByParentID(c *gin.Context) {
	parentIDStr := c.Query("parentId")
	var parentID uint = 0
	if parentIDStr != "" {
		id, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的父分类ID"})
			return
		}
		parentID = uint(id)
	}

	categories, err := h.categoryService.GetCategoriesByParentID(parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// GetAllCategories 获取所有分类
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// UpdateCategory 更新分类信息
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
		return
	}

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	category.ID = uint(id)
	if err := h.categoryService.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "分类信息更新成功",
		"data":    category,
	})
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "分类删除成功",
	})
}

// GetFullCategoryPath 获取完整分类路径
func (h *CategoryHandler) GetFullCategoryPath(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
		return
	}

	path, err := h.categoryService.GetFullCategoryPath(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": path,
	})
}
