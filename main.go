package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wayne011872/goSterna/api/mid"
	"github.com/wayne011872/goSterna/db"
	storage "github.com/wayne011872/goSterna/mystorage"
	sternaLog "github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/goSterna/api"
	myapi "github.com/wayne011872/systemMonitorServer/api"
)

var (
	service = flag.String("s","api","service(api)")
	envMode = flag.String("em", "local", "local, container")
)
func main() {
	flag.Parse()
	if *envMode == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("No .env file")
		}
	}
	switch *service {
	case "api":
		runAPI()
	default:
		panic("invalid service")
	}
}

func runAPI() {
	port := os.Getenv("API_PORT")
	ginMode := os.Getenv(("GIN_MODE"))
	serviceName := os.Getenv(("SERVICE"))
	confPath := os.Getenv("CONF_PATH")
	log.Println("run api port: ", port)
	log.Fatal(api.NewGinApiServer(ginMode).Middles(
		mid.NewGinFixDiMid(di,serviceName),
		mid.NewGinDBMid(serviceName),
	).AddAPIs(
		myapi.NewSysInfoAPI(serviceName),
	).Run(port).Error())
}


type di struct {
	*db.MongoConf         `yaml:"mongo,omitempty"`
	*sternaLog.LoggerConf `yaml:"log,omitempty"`
}

func (d *di) IsEmpty() bool {
	if d.MongoConf == nil {
		return true
	}

	if d.LoggerConf == nil {
		return true
	}

	return false
}