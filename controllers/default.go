package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/httplib"
    "strings"
    "strconv"
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

type Entity struct {
    Ti *string `json:"Ti"`
    CC *int64 `json:CC`
    Id *int64 `json:Id`
}

type JSONRet struct {
    Expr string `json:"expr"`
    Entities []Entity `json:"entities"`
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

func doQuery(expr string, attrs []string, count int) (*JSONRet,error) {

    url := fmt.Sprintf("https://%s%s", HOST, URL_FORMAT)

    req := httplib.Get(url)
    // set https insecure
    req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
    // set header
    req.Header("Ocp-Apim-Subscription-Key", KEY)
    // set param
    req.Param("expr", expr)
    req.Param("model", "latest")
    req.Param("attributes", strings.Join(attrs, ","))
    req.Param("count", strconv.Itoa(count))
    req.Param("offset", "0")

    // prepare result
    var result JSONRet
    var err error

    if err = req.ToJSON(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

func queryPaperIdByTitle(title string) int64 {

    expr := fmt.Sprintf("Ti='%s'", title)
    attrs := []string{"Id"}

    result, err := doQuery(expr, attrs, 10)

    if err != nil {
        logs.Error("query failed: %v", err)
        return -1
    }

    if result.Entities == nil || len(result.Entities) == 0 {
        logs.Error("parse result %v failed: null", result)
        return -1
    }

    entity := result.Entities[0]
    if entity.Id == nil {
        logs.Error("parse result %s failed: Id null", result)
        return -1
    }

    return *entity.Id
}

func queryCitationsById(id int64) []string {

    expr := fmt.Sprintf("RId=%v", id)
    attrs := []string{"Ti"}

    result, err := doQuery(expr, attrs, 100000)

    ret := []string{}

    if err != nil {
        logs.Error("query failed: %v", err)
        return nil
    }

    if result.Entities == nil {
        logs.Error("parse result %v failed: null", result)
        return nil
    }

    if len(result.Entities) == 0 {
        return ret
    }

    for _, v := range result.Entities {
        if v.Ti == nil {
            continue
        }

        ret = append(ret, *v.Ti)
    }

    return ret
}

func handleQuery(srcTitle string, outchan chan []string) {

    if outchan == nil {
        return
    }

    var title string
    var ccList []string

    title = convertPaperTitle(srcTitle)
    if title == "" {
        outchan <- ccList
        return
    }

    // query citations 
    paperId := queryPaperIdByTitle(title)

    if paperId < 0 {
        outchan <- ccList
        return
    }

    ccList = queryCitationsById(paperId)
    outchan <- ccList
    return
}

func (c *MainController) Get() {
	c.TplName = "gen.tpl"
}

func (c *MainController) Post() {
    papers := c.GetString("papers")
    paperList := strings.Split(papers, "\r\n")

    outChans := make([]chan []string, len(paperList))

    var result []map[string]interface{}
    for i, p := range paperList {
        logs.Debug("[%d] handling %s ...", i, p)

        outChans[i] = make(chan []string)

        go handleQuery(p, outChans[i])

    }

    for i, c := range outChans {

        var numCC int
        var ccList []string

        p := paperList[i]
        ccList = <- c

        if ccList == nil {
            numCC = -1
        } else {
            numCC = len(ccList)
        }

        tup := map[string]interface{}{
            "title": p,
            "cc": numCC,
            "cclist": ccList,
        }
        result = append(result, tup)

        logs.Debug("[%d] %s : %v, %v", i, p, numCC, ccList)
    }

    c.Data["Result"] = result
    c.TplName = "gen.tpl"
}

