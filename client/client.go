package client

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/binatify/wechat/config"
	"github.com/binatify/wechat/util"
)

type Client struct {
	cfg *config.Config

	skey       string
	sid        string
	uin        string
	passTicket string

	uuid        string
	deviceId    string
	redirectUri string
	baseUri     string

	user        map[string]interface{}
	syncKeyMap  map[string]interface{}
	synckey     string
	baseRequest map[string]interface{}

	syncHost string
	duration time.Duration

	*http.Client
}

func init() {
	rand.Seed(time.Now().Unix())
}

func New(cfg *config.Config) *Client {
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
		deviceId: deviceId,
		Client:   httpClient,

		duration: time.Duration(cfg.Duration) * time.Second,
	}
}

func (this *Client) Start() *Client {
	log.Println("[*] 微信网页版启动中 ...")

	util.RpcCall("[*] 正在获取 uuid ...", this.getUUID)

	util.RpcCall("[*] 正在获取二维码 ...", this.qrCode)
	util.RpcCall("[*] 请使用微信扫描二维码以登录 ...", this.qrCodeConfirm)
	util.RpcCall("[*] 正在登录 ...", this.login)

	util.RpcCall("[*] 获取微信初始化数据 ...", this.webwxinit)

	return this
}
