package public

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	IsDebug  bool      `json:"is_debug"`
	WsOrgin  string    `json:"ws_orgin"`
	Mysql    SQLConfig `json:"mysql"`
	UseMysql bool      `json:"use_mysql"`
}

type SQLConfig struct {
	Table    string `json:"table"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

var ConfigData = new(Config)

var configPath = flag.String("c", "data/config.json", "Input config file like data/config.json")

func init() {
	flag.Parse()
	//初始化配置
	println("*configPath - ", *configPath)
	configPtr, _ := os.Open(*configPath)
	defer configPtr.Close()
	decoder := json.NewDecoder(configPtr)
	err := decoder.Decode(&ConfigData)
	if err != nil {
		Printf("配置文件解码失败，", err.Error())
		//给一些默认的值
		ConfigData.IsDebug = false
	}
}

//打印
func Printf(logText ...interface{}) {
	if !ConfigData.IsDebug {
		logPath, line := getCaller(2)
		log.Println("[info] ", logPath, ":", line, " ", fmt.Sprint(logText...))
	} else {
		println(fmt.Sprint(logText...))
	}
}

//debug
func Debug(logText ...interface{}) {
	if ConfigData.IsDebug {
		logPath, line := getCaller(2)
		log.Println("[debug] ", logPath, ":", line, " ", fmt.Sprint(logText...))
		println(fmt.Sprint(logText...))
	}
}

//error
func Error(errText ...interface{}) bool {
	if len(errText) > 0 {
		logPath, line := getCaller(2)
		log.Println("[error] ", logPath, ":", line, " ", fmt.Sprint(errText...))
		return true
	}
	return false
}

//获取行号
func getCaller(skip int) (string, int) {

	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println(file)
	//fmt.Println(line)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}


//两个数之间的随机数
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}


func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}