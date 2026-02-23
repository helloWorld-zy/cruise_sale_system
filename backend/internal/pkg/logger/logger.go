package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New 创建并配置一个生产级的 zap 日志记录器。
// 使用 lumberjack 实现日志文件的自动轮转（按大小切割、保留历史、自动压缩）。
// level 参数指定日志级别（"debug"/"info"/"warn"/"error"），若无效则默认为 info 级别。
func New(level string, filename string) *zap.Logger {
	// 配置 lumberjack 日志轮转
	lj := &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    100,      // 单个日志文件最大大小（MB）
		MaxBackups: 10,       // 保留的旧日志文件数量
		MaxAge:     30,       // 日志文件保留天数
		Compress:   true,     // 是否压缩旧日志
	}

	// 配置日志编码格式
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// 解析日志级别
	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		atomicLevel.SetLevel(zapcore.InfoLevel) // 级别无效时回退到 info
	}

	// 创建日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), // JSON 格式输出
		zapcore.AddSync(lj),                // 写入轮转日志文件
		atomicLevel,                        // 日志级别过滤
	)

	return zap.New(core, zap.AddCaller()) // 添加调用者信息
}
