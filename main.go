package main

import (
	_ "github.com/peteryj/medtools/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "runtime"
    "os/exec"
    "flag"
    "fmt"
)

const (
    VERSION = "1.1"
)

var (
    RunningPort string = ""
)


func init() {
    runmode := beego.AppConfig.String("runmode")
    if runmode == "dev" {
        logs.SetLevel(logs.LevelDebug)
    } else {
        logs.SetLevel(logs.LevelInfo)
    }

    RunningPort = beego.AppConfig.String("httpport")
}

func open(url string) error{
    var cmd string
    var args[]string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default:// linux, openbsd, netbsd
        cmd = "xdg-open"
    }

    args = append(args, url)
    
    return exec.Command(cmd, args...).Start()
}

func pathMapping() {
    beego.SetStaticPath("/assets", "static/metronic/assets")
    beego.SetStaticPath("/js", "static/js")
}

func main() {

    // cmd args
    flagVersion := flag.Bool("v", false, "Print Version")
    flag.Parse()

    // handling cmd args
    if *flagVersion {
        fmt.Println(VERSION)
        return
    }

    // run web app
    pathMapping()

    go open(fmt.Sprintf("http://localhost:%v/", RunningPort))
	beego.Run()
}

