package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 是应用的顶层配置结构体，包含所有子模块的配置。
type Config struct {
	Server   ServerConfig   // 服务器配置
	Database DatabaseConfig // 数据库配置
	Redis    RedisConfig    // Redis 缓存配置
	MinIO    MinIOConfig    // MinIO 对象存储配置
	Meilis   MeiliConfig    // MeiliSearch 搜索引擎配置
	NATS     NATSConfig     // NATS 消息队列配置
	JWT      JWTConfig      // JWT 认证配置
	Log      LogConfig      // 日志配置
	Upload   UploadConfig   // 本地上传配置
}

// UploadConfig 定义本地文件上传配置。
type UploadConfig struct {
	StorageDir  string // 上传文件保存目录
	PublicPath  string // 上传文件静态访问前缀
	MaxFileSize int64  // 单文件最大大小（字节）
}

// ServerConfig 定义 HTTP 服务器的启动参数。
type ServerConfig struct {
	Port string // 监听端口（如 ":8080"）
	Mode string // 运行模式（"debug" / "release"）
}

// DatabaseConfig 定义 PostgreSQL 数据库连接参数。
type DatabaseConfig struct {
	Host            string // 数据库主机地址
	Port            int    // 数据库端口
	User            string // 数据库用户名
	Password        string // 数据库密码
	DBName          string // 数据库名称
	SSLMode         string // SSL 模式（"disable" / "require"）
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	ConnMaxLifetime int    // 连接最大存活时间（秒）
}

// RedisConfig 定义 Redis 连接参数。
type RedisConfig struct {
	Host     string // Redis 主机地址
	Port     int    // Redis 端口
	Password string // Redis 密码
	DB       int    // 使用的数据库编号
}

// MinIOConfig 定义 MinIO 对象存储连接参数。
type MinIOConfig struct {
	Endpoint  string // MinIO 服务端点
	AccessKey string // 访问密钥
	SecretKey string // 秘密密钥
	Bucket    string // 存储桶名称
	UseSSL    bool   // 是否使用 SSL 连接
}

// MeiliConfig 定义 MeiliSearch 全文搜索引擎连接参数。
type MeiliConfig struct {
	Host   string // MeiliSearch 主机地址
	APIKey string // API 密钥
}

// NATSConfig 定义 NATS 消息队列连接参数。
type NATSConfig struct {
	URL string // NATS 服务器地址
}

// JWTConfig 定义 JWT 令牌的签发参数。
type JWTConfig struct {
	Secret      string // 签名密钥
	ExpireHours int    // 令牌过期时间（小时）
}

// LogConfig 定义日志输出参数，使用 lumberjack 实现日志轮转。
type LogConfig struct {
	Level      string // 日志级别（"debug" / "info" / "warn" / "error"）
	Filename   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}

// Load 读取配置文件 config.yaml，并支持通过 CRUISE_ 前缀的环境变量覆盖配置项。
// root 参数指定配置文件所在目录。
func Load(root string) Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(root))
	v.AddConfigPath(".")

	// 设置环境变量前缀为 CRUISE，并自动读取环境变量
	v.SetEnvPrefix("CRUISE")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 将配置反序列化到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("配置反序列化失败: %w", err))
	}

	applyDatabaseEnvFallbacks(&cfg)
	applyUploadDefaults(&cfg)

	return cfg
}

// applyDatabaseEnvFallbacks 在 CRUISE_DATABASE_* 未提供时回退到 POSTGRES_*，
// 避免仅配置 Docker 变量时后端因空密码无法连接数据库。
func applyDatabaseEnvFallbacks(cfg *Config) {
	if cfg == nil {
		return
	}

	if cfg.Database.Host == "" {
		cfg.Database.Host = firstNonEmptyEnv("CRUISE_DATABASE_HOST", "POSTGRES_HOST")
	}
	if cfg.Database.User == "" {
		cfg.Database.User = firstNonEmptyEnv("CRUISE_DATABASE_USER", "POSTGRES_USER")
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = firstNonEmptyEnv("CRUISE_DATABASE_PASSWORD", "POSTGRES_PASSWORD")
	}
	if cfg.Database.DBName == "" {
		cfg.Database.DBName = firstNonEmptyEnv("CRUISE_DATABASE_DBNAME", "POSTGRES_DB")
	}
	if cfg.Database.Port == 0 {
		if portText := firstNonEmptyEnv("CRUISE_DATABASE_PORT", "POSTGRES_PORT"); portText != "" {
			if port, err := strconv.Atoi(portText); err == nil {
				cfg.Database.Port = port
			}
		}
	}
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func applyUploadDefaults(cfg *Config) {
	if cfg == nil {
		return
	}
	if strings.TrimSpace(cfg.Upload.StorageDir) == "" {
		cfg.Upload.StorageDir = "uploads"
	}
	if strings.TrimSpace(cfg.Upload.PublicPath) == "" {
		cfg.Upload.PublicPath = "/uploads"
	}
	if cfg.Upload.MaxFileSize <= 0 {
		cfg.Upload.MaxFileSize = 10 * 1024 * 1024
	}
}
