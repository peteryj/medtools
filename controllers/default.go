package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/httplib"
    "strings"
    "fmt"
    "crypto/tls"
)

const (
    HOST = "api.projectoxford.ai"
    URL_FORMAT = "/academic/v1.0/evaluate"
    KEY = "d43df72e86cf4a81bb6a594d519d1b78"
)

type MainController struct {
	beego.Controller
}

type JSONRet struct {
    Expr string `json:"expr"`
    Entities []map[string]*float64 `json:"entities"`

}

func convertPaperTitle(title string) string {
    var ti string = title
    ti = strings.ToLower(ti)
    ti = strings.Trim(ti, " \n\r")
    ti = strings.Replace(ti, "/", " ", -1)
    ti = strings.Replace(ti, "-", " ", -1)
    ti = strings.Replace(ti, "\"", " ", -1)
    ti = strings.Replace(ti, ": ", " ", -1)
    ti = strings.Replace(ti, " :", " ", -1)
    ti = strings.Replace(ti, " (", " ", -1)
    ti = strings.Replace(ti, ") ", " ", -1)

    // check again
    ti = strings.Replace(ti, "(", " ", -1)
    ti = strings.Replace(ti, ")", " ", -1)
    ti = strings.Replace(ti, ":", " ", -1)

    // strip ' ', '.'
    ti = strings.Trim(ti, " ")
    ti = strings.TrimRight(ti, ".")

    return ti
}

func queryCC(title string) float64 {

    url := fmt.Sprintf("https://%s%s", HOST, URL_FORMAT)

    req := httplib.Get(url)
    // set https insecure
    req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
    // set header
    req.Header("Ocp-Apim-Subscription-Key", KEY)
    // set param
    req.Param("expr", fmt.Sprintf("Ti='%s'", title))
    req.Param("model", "latest")
    req.Param("attributes", "CC")
    req.Param("count", "10")
    req.Param("offset", "0")

    // prepare result
    var result JSONRet

    if err := req.ToJSON(&result); err != nil {
        logs.Error("query failed: %v", err)
        return -1
    }

    if result.Entities == nil || len(result.Entities) == 0 {
        logs.Error("parse result %v failed: null", result)
        return -1
    }

    ss := result.Entities[0]
    if ss["CC"] == nil {
        logs.Error("parse result %s failed: CC null", result)
        return -1
    }

    return *ss["CC"]
}

func (c *MainController) Get() {
	c.TplName = "gen.tpl"
}

func (c *MainController) Post() {
    papers := c.GetString("papers")
    paperList := strings.Split(papers, "\r\n")

    var result []map[string]string
    for i, p := range paperList {
        logs.Debug("[%d] handling %s ...", i, p)

        var title string
        var cc float64

        title = convertPaperTitle(p)
        if title == "" {
            continue
        }

        cc = queryCC(title)

        tup := map[string]string{
            "title": title,
            "cc": fmt.Sprintf("%v", cc),
        }
        result = append(result, tup)

        logs.Debug("[%d] %s : %v", i, p, cc)
    }

    c.Data["Result"] = result
    c.TplName = "gen.tpl"
}


