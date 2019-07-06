package setting

import (
	"github.com/go-ini/ini"
	"github.com/labstack/gommon/log"
	"time"
)

var (
	Cfg *ini.File
)

type Application struct {
	PageSize        int
	JwtSecret       string
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavaPath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogFileName string
	LogFileExt  string
	TimeFormat   string
}

var AppSetting = &Application{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

func Setup() {
	// 加载配置文件
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("文件【conf/app.ini】解析失败: %v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg MapTo AppSetting err: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg MapTo DatabaseSetting err: %v", err)
	}
}

//func LoadBase() {
//	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
//}
//
//// 【？为啥不是获取内容】
//func LoadServer() {
//	sec, err := Cfg.GetSection("server")
//	if err != nil {
//		log.Fatalf("配置文件获取'server'失败: %v", err)
//	}
//
//	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
//
//	HttpPort = sec.Key("HTTP_PORT").MustInt(9090)
//	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
//	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
//}
//
//func LoadAPP() {
//	sec, err := Cfg.GetSection("app")
//	if err != nil {
//		log.Fatalf("配置文件获取'app'失败: %v", err)
//	}
//
//	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
//	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
//}
