package main

import (
	"github.com/KQLXK/Family-Finance-System/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	familyHandler := handler.NewFamilyHandler()
	memberHandler := handler.NewMemberHandler()
	categoryHandler := handler.NewCategoryHandler()
	transactionHandler := handler.NewTransactionHandler()

	// 家庭相关路由
	familyGroup := r.Group("/api/families")
	{
		familyGroup.POST("", familyHandler.CreateFamily)
		familyGroup.GET("", familyHandler.GetAllFamilies)

		familyGroup.POST("/:id/members", memberHandler.CreateMember)
		familyGroup.GET("/:id/members", memberHandler.GetMembersByFamilyID)
		familyGroup.GET("/:id/members/active", memberHandler.GetActiveMembersByFamilyID)

		familyGroup.GET("/:id", familyHandler.GetFamilyByID)
		familyGroup.PUT("/:id", familyHandler.UpdateFamily)
		familyGroup.DELETE("/:id", familyHandler.DeleteFamily)

		// 家庭交易相关路由
		familyGroup.POST("/:id/transactions", transactionHandler.CreateTransaction)
		familyGroup.GET("/:id/transactions", transactionHandler.GetTransactionsByFamilyID)
		familyGroup.GET("/:id/transactions/time-range", transactionHandler.GetTransactionsByTimeRange)
		familyGroup.GET("/:id/transactions/summary/category", transactionHandler.GetTransactionSummaryByCategory)
		familyGroup.GET("/:id/transactions/summary/time", transactionHandler.GetTransactionSummaryByTime)
	}

	// 成员相关路由（独立于家庭）
	memberGroup := r.Group("/api/members")
	{
		memberGroup.GET("", memberHandler.GetAllMembers)
		memberGroup.GET("/:id", memberHandler.GetMemberByID)
		memberGroup.PUT("/:id", memberHandler.UpdateMember)
		memberGroup.DELETE("/:id", memberHandler.DeleteMember)
		memberGroup.PUT("/:id/role", memberHandler.ChangeMemberRole)
	}

	categoryGroup := r.Group("/api/categories")
	{
		categoryGroup.POST("", categoryHandler.CreateCategory)
		categoryGroup.GET("", categoryHandler.GetAllCategories)
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
		categoryGroup.GET("/:id/path", categoryHandler.GetFullCategoryPath)
		categoryGroup.GET("/type/list", categoryHandler.GetCategoriesByType)
		categoryGroup.GET("/type/tree", categoryHandler.GetCategoryTreeByType)
		categoryGroup.GET("/parent/children", categoryHandler.GetCategoriesByParentID)
	}

	transactionGroup := r.Group("/api/transactions")
	{
		transactionGroup.GET("/:id", transactionHandler.GetTransactionByID)
		transactionGroup.PUT("/:id", transactionHandler.UpdateTransaction)
		transactionGroup.DELETE("/:id", transactionHandler.DeleteTransaction)
		transactionGroup.POST("/:id/tags", transactionHandler.AddTagToTransaction)
		transactionGroup.DELETE("/:id/tags/:tagId", transactionHandler.RemoveTagFromTransaction)
	}

	return r
}
