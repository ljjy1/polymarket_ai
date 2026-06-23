// Package initial 是启动服务进行初始化的包，包括初始化配置、服务配置、连接数据库以及服务关闭时所需的资源释放。
package initial

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/stat"
	"github.com/go-dev-frame/sponge/pkg/tracer"

	"be/configs"
	"be/internal/config"
	"be/internal/database"
)

var (
	version    string // 服务版本号
	configFile string // 配置文件路径
)

// InitApp 初始化应用配置
func InitApp() {
	// 设置默认时区为中国上海时区 (UTC+8)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic("load timezone location failed: " + err.Error())
	}
	time.Local = loc

	initConfig()        // 初始化配置
	expandEnvConfig()   // 展开环境变量
	cfg := config.Get() // 获取配置对象

	// 初始化日志
	_, err = logger.Init(
		logger.WithLevel(cfg.Logger.Level),   // 设置日志级别
		logger.WithFormat(cfg.Logger.Format), // 设置日志格式
		logger.WithSave(
			cfg.Logger.IsSave, // 是否保存日志到文件
			//logger.WithFileName(cfg.Logger.LogFileConfig.Filename),       // 日志文件名
			//logger.WithFileMaxSize(cfg.Logger.LogFileConfig.MaxSize),     // 日志文件最大大小
			//logger.WithFileMaxBackups(cfg.Logger.LogFileConfig.MaxBackups), // 最大备份数
			//logger.WithFileMaxAge(cfg.Logger.LogFileConfig.MaxAge),       // 日志文件最大保存天数
			//logger.WithFileIsCompression(cfg.Logger.LogFileConfig.IsCompression), // 是否压缩
		),
	)
	if err != nil {
		panic(err) // 日志初始化失败则直接 panic
	}
	logger.Debug(config.Show())             // 打印当前配置信息（调试级别）
	logger.Info("[logger] was initialized") // 记录日志初始化完成

	// 初始化链路追踪
	if cfg.App.EnableTrace {
		tracer.InitWithConfig(
			cfg.App.Name,                       // 应用名称
			cfg.App.Env,                        // 运行环境
			cfg.App.Version,                    // 版本号
			cfg.Jaeger.AgentHost,               // Jaeger Agent 主机地址
			strconv.Itoa(cfg.Jaeger.AgentPort), // Jaeger Agent 端口
			cfg.App.TracingSamplingRate,        // 采样率
		)
		logger.Info("[tracer] was initialized") // 记录追踪初始化完成
	}

	// 初始化系统资源统计与打印
	if cfg.App.EnableStat {
		stat.Init(
			stat.WithLog(logger.Get()), // 使用已初始化的日志实例
			stat.WithAlarm(),           // 开启告警（Windows 下无效），默认 CPU 和内存阈值 0.8，可修改
			stat.WithPrintField(logger.String("service_name", cfg.App.Name), logger.String("host", cfg.App.Host)), // 额外打印字段
		)
		logger.Info("[resource statistics] was initialized") // 记录资源统计初始化完成
	}

	// 初始化数据库
	database.InitDB()
	logger.Infof("[%s] was initialized", cfg.Database.Driver) // 记录数据库驱动初始化完成
	database.InitCache(cfg.App.CacheType)                     // 初始化缓存
	if cfg.App.CacheType != "" {
		logger.Infof("[%s] was initialized", cfg.App.CacheType) // 记录缓存类型初始化完成
	}

	// 自动迁移数据库表结构 暂时不需要 手动管理表结构先
	//models := []interface{}{
	//	&model.User{},
	//	&model.Markets{},
	//	&model.Predictions{},
	//	&model.Strategies{},
	//	&model.Trades{},
	//	&model.VaultSnapshots{},
	//	&model.SystemLogs{},
	//}
	//for _, m := range models {
	//	if err := database.GetDB().AutoMigrate(m); err != nil {
	//		logger.Error("auto migrate table error", logger.Err(err)) // 自动迁移出错记录日志
	//	}
	//}
}

// initConfig 解析命令行参数并加载配置
func initConfig() {
	flag.StringVar(&version, "version", "", "服务版本号")
	flag.StringVar(&configFile, "c", "", "配置文件路径")
	flag.Parse() // 解析命令行参数

	getConfigFromLocal() // 从本地配置文件加载

	if version != "" {
		config.Get().App.Version = version // 若命令行指定版本号则覆盖配置中的版本
	}
}

// getConfigFromLocal 从本地配置文件获取配置
func getConfigFromLocal() {
	if configFile == "" {
		configFile = configs.Location("be.yml") // 未指定配置文件时使用默认路径
	}
	err := config.Init(configFile) // 加载并解析配置文件
	if err != nil {
		panic("init config error: " + err.Error()) // 配置加载失败直接 panic
	}
}

// expandEnvConfig 将配置中所有以 "$" 开头的字符串值替换为实际环境变量值。
// Sponge 的 conf.Parse 不会自动展开环境变量，通过反射递归遍历所有字段。
func expandEnvConfig() {
	cfg := config.Get()
	expandEnvReflect(reflect.ValueOf(cfg))
	config.Set(cfg)
}

// expandEnvReflect 递归遍历结构体，对字符串字段执行 os.ExpandEnv。
func expandEnvReflect(v reflect.Value) {
	// 指针 → 解引用到实际值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			expandEnvReflect(v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Struct || elem.Kind() == reflect.Map {
				expandEnvReflect(elem)
			} else if elem.Kind() == reflect.String {
				s := elem.String()
				if strings.HasPrefix(s, "$") {
					v.SetMapIndex(key, reflect.ValueOf(os.ExpandEnv(s)))
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			expandEnvReflect(v.Index(i))
		}
	case reflect.String:
		s := v.String()
		if strings.HasPrefix(s, "$") {
			v.SetString(os.ExpandEnv(s))
		}
	}
}
