package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/binatify/wechat/config"
)

type Client struct {
	cfg *config.Config
	log *logrus.Entry

	sKey       string
	sid        string
	uin        string
	passTicket string

	uuid        string
	deviceId    string
	redirectUri string
	baseUri     string

	syncKey     string
	user        map[string]interface{}
	syncKeyMap  map[string]interface{}
	baseRequest map[string]interface{}

	syncHost string
	duration time.Duration

	*http.Client
}

func init() {
	rand.Seed(time.Now().Unix())
}

func New(cfg *config.Config, log *logrus.Entry) *Client {
	gCookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		CheckRedirect: nil,
		Jar:           gCookieJar,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	str := strconv.Itoa(rand.Int())
	deviceId := "e" + str[2:17]

	return &Client{
		log: log,

		deviceId: deviceId,
		duration: time.Duration(cfg.Duration) * time.Second,

		Client: httpClient,
	}
}

func (c *Client) Start() *Client {

	c.log.Println("微信客户端启动中 ...")

	c.rpcCall("正在获取 uuid ...", c.getUUID)

	c.rpcCall("正在获取二维码 ...", c.qrCode)

	c.rpcCall("请使用微信扫描二维码 ...", c.qrCodeConfirm)

	c.rpcCall("正在登录 ...", c.login)

	c.rpcCall("获取微信初始化数据 ...", c.wxInit)

	return c

}

func (c *Client) rpcCall(description string, f func() bool) {
	c.log.Println(description)

	t1 := time.Now().UnixNano()

	if ok := f(); ok {
		cost := fmt.Sprintf("%.5f", float64(time.Now().UnixNano()-t1)/float64(time.Second))
		c.log.Print("成功, 用时" + cost + "秒")
		return
	}

	c.log.Panic("启动失败，退出程序")
}
