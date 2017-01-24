package main

import (
	_ "github.com/peteryj/medtools/routers"
	"github.com/astaxie/beego"
    "runtime"
    "os/exec"
)

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

func main() {
    go open("http://localhost:8080/")
	beego.Run()
}

