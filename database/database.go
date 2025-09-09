package database

import (
	"fmt"
	"github.com/KQLXK/Family-Finance-System/commen/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.GetConfig()

	// 获取数据库连接字符串
	dsn := cfg.GetDBConnectionString()

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	}

	// 建立数据库连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err = DB.Exec("CREATE DATABASE IF NOT EXISTS family_finance").Error; err != nil {
		log.Fatal("Failed to create database: ", err)
		return
	}

	fmt.Println("MySQL database connection established successfully")
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting SQL DB: %v", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		fmt.Println("Database connection closed")
	}
}
