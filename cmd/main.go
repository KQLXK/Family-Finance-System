package main

import (
	"github.com/KQLXK/Family-Finance-System/model"
	"gorm.io/gorm"
	"log"

	"github.com/KQLXK/Family-Finance-System/database"
)

func main() {

	//// 设置Gin运行模式
	//if cfg.Server.Mode == "release" {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	// 初始化数据库
	database.InitDB()

	// 确保程序退出时关闭数据库连接
	defer database.CloseDB(database.DB)

	// 运行数据库迁移
	if err := AutoMigrate(database.DB); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migration completed successfully")

	//设置路由
	r := SetupRouter()

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Family{},
		&model.Member{},
		&model.Category{},
		&model.Tag{},
		&model.Transaction{},
		&model.TransactionTag{},
	)
}
