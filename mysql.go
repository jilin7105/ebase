package ebase

import (
	"fmt"
	"github.com/jilin7105/ebase/logger"
	"gorm.io/driver/postgres"

	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (e *Eb) initMysql() {
	// 初始化所有数据库连接
	for _, dbConfig := range e.Config.Databases {
		if dbConfig.Type == "" {
			dbConfig.Type = "mysql"
		}

		if dbConfig.Type == "mysql" {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				dbConfig.Username,
				dbConfig.Password,
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.Dbname,
			)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				logger.Error("Failed to connect to database %s: %v", dbConfig.Name, err)
				continue
			}
			// 设置最大空闲连接数和最大打开连接数
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
			sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
			// 将数据库连接保存到Eb结构体中
			e.DBs[dbConfig.Name] = db
		}

		if dbConfig.Type == "pgsql" {
			dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
				dbConfig.Username,
				dbConfig.Password,
				dbConfig.Dbname,
				dbConfig.Host,
				dbConfig.Port,
			)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				logger.Error("Failed to connect to database %s: %v", dbConfig.Name, err)
				continue
			}
			e.DBs[dbConfig.Name] = db
		}
	}

}
