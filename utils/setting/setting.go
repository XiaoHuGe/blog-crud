package setting

import (
	"github.com/go-ini/ini"
	"github.com/labstack/gommon/log"
	"time"
)

var (
	Cfg *ini.File

	RunMode string

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

func init() {
	// 加载配置文件
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("文件【conf/app.ini】解析失败: %v", err)
	}

	LoadBase()
	LoadServer()
	LoadAPP()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

// 【？为啥不是获取内容】
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("配置文件获取'server'失败: %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	HttpPort = sec.Key("HTTP_PORT").MustInt(9090)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadAPP() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("配置文件获取'app'失败: %v", err)
	}

	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
}
