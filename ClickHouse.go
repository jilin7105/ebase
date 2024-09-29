package ebase

import (
	std_ck "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"gorm.io/driver/clickhouse"

	"gorm.io/gorm"
	"time"
)

func newCk(ckc config.ClickHouseConfig) (db *gorm.DB, err error) {
	if ckc.ReadTimeout <= 0 {
		ckc.ReadTimeout = 5
	}

	if ckc.WriteTimeout <= 0 {
		ckc.WriteTimeout = 5
	}

	sqlDB := std_ck.OpenDB(&std_ck.Options{
		Addr: ckc.Hosts,
		Auth: std_ck.Auth{
			Database: ckc.Db,
			Username: ckc.User,
			Password: ckc.Pass,
		},
		Settings: std_ck.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Duration(ckc.ReadTimeout) * time.Second,
		ReadTimeout: time.Duration(ckc.WriteTimeout) * time.Second,
		Compression: &std_ck.Compression{
			std_ck.CompressionLZ4,
			3,
		},
		Debug: true,
	})

	db, err = gorm.Open(clickhouse.New(clickhouse.Config{
		Conn: sqlDB, // initialize with existing database conn
	}))

	return
}

func (e *Eb) initCk() {

	for _, ClickHouseConfig := range e.Config.ClickHouse {
		ckClient, err := newCk(ClickHouseConfig)
		if err != nil {
			logger.Error("ClickHouse  初始化失败 %s %s", ClickHouseConfig.Name, err.Error())
		}
		e.ClickHouse[ClickHouseConfig.Name] = ckClient
	}

}
