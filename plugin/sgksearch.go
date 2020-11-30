package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
)

const (
    DESPET = "-----------------------------------------------\n"
    LAST_DESPET = "==数据太多,只显示前5条==\n"
)

var queryUrl = os.Getenv("CALLBACK_URL")

type queryMessageReqBody struct {
    QueryStr string `json:"query_str"`
    Skip int `json:"skip"`
}

type Hit struct {
    QQ string `json:"qq"`
    Name string `json:"name"`
    Phone string `json:"phone"`
    Password string `json:"password"`
    Area string `json:"area"`
    Email string `json:"email"`
}

type Result struct {
    Name string `json:"name"`
    HitsCount int `json:"hits_count"`
    Hits []Hit `json:"hits"`
    Fields []string `json:"fields"`
}

type queryMessageResponse struct {
    Status int `json:"status"`
    Info string `json:"info"`
    Hits int `json:"hits"`
    Total int `json:"total"`
    Run float32 `json:"run"`
    Results []Result `json:"results"`
}

func querySgkBody(sendMsg string) (body queryMessageResponse, err error) {
    // fmt.Println("request sgk return body response.")
    var strQuery strings.Builder
    strQuery.WriteString("query_str=")
    strQuery.WriteString(sendMsg)
    payload := strings.NewReader(strQuery.String())
    res, err := http.Post(queryUrl, "text/plain", payload)
    if err != nil {
        return body, errors.New("Query sgk web error")
    }
    defer res.Body.Close()
    temp, _ := ioutil.ReadAll(res.Body)
    json.Unmarshal([]byte(temp), &body)
    return
}

func parseMsg(body queryMessageResponse) (returnStr string) {
    var strBuild strings.Builder
    for i := 0; i < len(body.Results); i++ {
        hitCount := body.Results[i].HitsCount
        // fmt.Println(".....hit count:", hitCount)
        if hitCount > 0 {
            strBuild.WriteString("来源: ")
            strBuild.WriteString(body.Results[i].Name)
            strBuild.WriteString("\n")
            strBuild.WriteString("关联")
            strBuild.WriteString(strconv.Itoa(hitCount))
            strBuild.WriteString("次\n")
            strBuild.WriteString(DESPET)
            hits := body.Results[i].Hits
            for j := 0; j < len(hits); j++ {
                for _, v := range body.Results[i].Fields {
                    switch v {
                    case "qq":
                        strBuild.WriteString("QQ: ")
                        strBuild.WriteString(hits[j].QQ)
                        strBuild.WriteString("\n")
                    case "phone":
                        strBuild.WriteString("手机号: ")
                        strBuild.WriteString(hits[j].Phone)
                        strBuild.WriteString("\n")
                    case "password":
                        strBuild.WriteString("密码: ")
                        strBuild.WriteString(hits[j].Password)
                        strBuild.WriteString("\n")
                    case "name":
                        strBuild.WriteString("用户名: ")
                        strBuild.WriteString(hits[j].Name)
                        strBuild.WriteString("\n")
                    case "email":
                        strBuild.WriteString("邮箱: ")
                        strBuild.WriteString(hits[j].Email)
                        strBuild.WriteString("\n")
                    }
                }
                if len(hits[j].Area) > 0 {
                    strBuild.WriteString("归属地: ")
                    strBuild.WriteString(hits[j].Area)
                    strBuild.WriteString("\n")
                }
                strBuild.WriteString(DESPET)
                if j > 3 {
                    strBuild.WriteString(LAST_DESPET)
                    break
                }
            }
        }
    }
    // fmt.Printf("query runtime:%f, hit: %d, total: %d\n", body.Run, body.Hits, body.Total)
    returnStr = strBuild.String()
    return
}

func search(command *bot.PassiveCmd) (msg string, err error) {
    if queryUrl == "" {
        fmt.Println("Undefined web site api!!Selete default http://web:8812/query")
        queryUrl = "http://web:8812/query"
    }
    if len(command.Raw) > 0 {
        resBody, err := querySgkBody(command.Raw)
        if err != nil {
            return "", errors.New(err.Error())
        }
        // if resBody.Status == 400 {
        //     fmt.Println("unlogin!!!pls login!")
        //     return
        // }
        if resBody.Status < 0 {
            return "", errors.New(resBody.Info)
        }
        if resBody.Status == 0 && strings.Contains(resBody.Info, "added") {
            t:=time.Tick(time.Millisecond*200)
            select {
            case <-t:
                fmt.Printf("module sgksearch.go: Wait 200ms,User: {id=%s, nick=%s, realname=%s}, query %s first,query result:%s\n", command.User.ID, command.User.Nick, command.User.RealName, command.Raw, resBody.Info)
                resBody, _ = querySgkBody(command.Raw)
            }
        }
        if resBody.Hits < 1 {
            // return "", errors.New("hit count less then 0")
            return `==祝贺您,你的个人信息没有泄露.==`, nil
        }

        msg = parseMsg(resBody)

    }
    return
}

func goodMorning(channel string) (msg string, err error) {
    msg = fmt.Sprintf("你的积分已重置,今天可使用3次免费查询.")
    return
}

func doping(command *bot.Cmd) (msg string, err error) {
  msg = fmt.Sprintf("Pong success. Hello %s, Your user id is %s.", command.User.RealName, command.User.ID)
  return
}

func init() {
    channels := strings.Split(os.Getenv("CHANNEL_IDS"), ",")
    if len(channels) > 0 {
        config := bot.PeriodicConfig{
            CronSpec: "00 02 * * *",
            Channels: channels,
            CmdFunc:  goodMorning,
        }
        bot.RegisterPeriodicCommand("good_morning", config)
    }

    bot.RegisterPassiveCommand(
        "search",
        search)

    bot.RegisterCommand(
        "ping",
        "Sends a 'ping' message to you on the channel.",
        "",
        doping)
}
