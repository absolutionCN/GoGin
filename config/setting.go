package config

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

//var (
//	Cfg          *ini.File
//	RunMode      string
//	HTTPPort     int
//	ReadTimeout  time.Duration
//	WriteTimeout time.Duration
//	PageSize     int
//	YapiHost     string
//	JwtSecret    string
//)
//
//func init() {
//	var err error
//	Cfg, err = ini.Load("config/app.ini")
//	if err != nil {
//		log.Fatalf("Fail to parse 'conf/app.ini':%v", err)
//	}
//	LoadBase()
//	LoadServer()
//	LoadApp()
//	LoadYapi()
//}
//
//func LoadBase() {
//	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
//}
//
//func LoadServer() {
//	sec, err := Cfg.GetSection("server")
//	if err != nil {
//		log.Fatalf("Fail to get section 'server': %v", err)
//	}
//	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
//	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
//	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
//}
//
//func LoadApp() {
//	sec, err := Cfg.GetSection("app")
//	if err != nil {
//		log.Fatalf("Fail to get section 'app': %v", err)
//	}
//	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
//}
//
//func LoadYapi() {
//	sec, err := Cfg.GetSection("yapi")
//	if err != nil {
//		log.Fatalf("Fail to get section 'yapi': %v", err)
//	}
//	YapiHost = sec.Key("HOST").MustString("")
//}

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string
	ImagePrefixUrl  string
	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowExts  []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type DataBase struct {
	Type        string
	User        string
	PassWord    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &DataBase{}

func Setup() {
	Cfg, err := ini.Load("config/app.ini")
	if err != nil {
		log.Fatal("Fail to parse 'conf/app.ini': %v", err)
	}
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatal("Cfg.MapTo AppSetting err: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatal("Cfg.MapTo ServerSetting err: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatal("Cfg.MapTo DataSetting err: %v", err)
	}
}
