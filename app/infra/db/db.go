package db

import (
	"log"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// GetDBInstance 返回全局唯一的数据库连接实例
func Init(dbPath string) (*gorm.DB, error) {
	var initErr error
	once.Do(func() {
		var err error
		instance, err = connectDatabase(dbPath)
		if err != nil {
			initErr = err
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

// connectDatabase 创建并配置数据库连接
func connectDatabase(dbPath string) (*gorm.DB, error) {
	// 使用 github.com/glebarez/sqlite 驱动连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取通用数据库对象 sql.DB 以进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// GetDB 获取数据库实例（不带路径参数，适用于已初始化后获取实例）
func GetDB() *gorm.DB {
	if instance == nil {
		log.Fatal("Database instance is not initialized. Call GetDBInstance first.")
	}
	return instance
}
