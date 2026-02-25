package database

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config 定义数据库连接的配置参数。
type Config struct {
	Host            string // 数据库主机地址
	Port            int    // 数据库端口
	User            string // 数据库用户名
	Password        string // 数据库密码
	DBName          string // 数据库名称
	SSLMode         string // SSL 连接模式
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	ConnMaxLifetime int    // 连接最大存活时间（秒）
}

// BuildDSN 根据配置构建 PostgreSQL 的数据源名称（DSN）字符串。
func BuildDSN(cfg Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// Connect 使用 GORM 连接 PostgreSQL 数据库，并配置连接池参数。
// 返回已配置好的 gorm.DB 实例。
func Connect(cfg Config) (*gorm.DB, error) {
	if cfg.Host == "sqlite" || cfg.Host == "memory" {
		return ConnectWithDialector(sqlite.Open("file::memory:?cache=shared"), cfg)
	}
	return ConnectWithDialector(postgres.Open(BuildDSN(cfg)), cfg)
}

// ConnectWithDialector 允许传入自定义的 Dialector（如 sqlite）以便于单元测试。
func ConnectWithDialector(dialector gorm.Dialector, cfg Config) (*gorm.DB, error) {
	// 打开数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return db, nil
}
