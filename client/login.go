package client

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/binatify/wechat/util"
)

const (
	okCode       = "200"
	loggedCode   = "201"
	loginTimeout = "408"
)

func (c *Client) getUUID() (ok bool) {
	url := "https://login.weixin.qq.com/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + util.UnixTimestamp()

	data, err := c.doGet(url)
	if err != nil {
		return
	}

	re := regexp.MustCompile(`"([\S]+)"`)
	found := re.FindStringSubmatch(string(data))

	if len(found) <= 1 {
		return
	}

	c.uuid = found[1]

	return true
}

func (c *Client) qrCode() (ok bool) {
	url := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s?t=webwx&_=%s", c.uuid, util.UnixTimestamp())

	resp, err := c.doGet(url)
	if err != nil {
		c.log.Errorf("this.doGet(%s): %v", url, err)
		return
	}

	path := "qrcode.jpg"
	if err = ioutil.WriteFile(path, resp, 0755); err != nil {
		c.log.Errorf("ioutil.WriteFile(qrcode.jpg): %v", err)
		return
	}

	// hack way for login qrcode display
	if runtime.GOOS == "darwin" {
		exec.Command("open", path).Run()
	} else {
		go func() {
			fmt.Printf("please open on web broswer %s/qrcode", c.cfg.Listen)
			http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
				http.ServeFile(w, req, path)
				return
			})
			http.ListenAndServe(c.cfg.Listen, nil)
		}()
	}

	return true
}

func (c *Client) qrCodeConfirm() (ok bool) {
	for retryTimes := 3; retryTimes > 0; retryTimes-- {

		if !c.doConfirm(1) {
			continue
		}

		c.log.Println("请在手机上点击确认 ...")

		if !c.doConfirm(0) {
			continue
		}

		ok = true

		break
	}

	return
}

func (c *Client) doConfirm(tip int) (ok bool) {
	time.Sleep(time.Duration(tip) * time.Second)

	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	url += "?tip=" + strconv.Itoa(tip) + "&uuid=" + c.uuid + "&_=" + util.UnixTimestamp()

	data, err := c.doGet(url)
	if err != nil {
		return
	}

	re := regexp.MustCompile(`window.code=(\d+);`)
	found := re.FindStringSubmatch(string(data))

	if len(found) > 1 {
		code := found[1]

		switch code {
		case okCode:
			re = regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			found := re.FindStringSubmatch(string(data))

			if len(found) > 1 {
				rUri := found[1] + "&fun=new"
				c.redirectUri = rUri
				re = regexp.MustCompile(`/`)

				found := re.FindAllStringIndex(rUri, -1)
				c.baseUri = rUri[:found[len(found)-1][0]]
				return true
			}

		case loggedCode:
			return true

		case loginTimeout:
			c.log.Error("登陆超时")

		default:
			c.log.Panic("登陆异常")
		}
	}

	return
}

type loginResult struct {
	XMLName xml.Name `xml:"error"`

	Skey       string `xml:"skey"`
	Wxsid      string `xml:"wxsid"`
	Wxuin      string `xml:"wxuin"`
	PassTicket string `xml:"pass_ticket"`
}

func (c *Client) login() (ok bool) {
	data, err := c.doGet(c.redirectUri)
	if err != nil {
		return
	}

	var v loginResult

	if err = xml.Unmarshal(data, &v); err != nil {
		c.log.Errorf("xml.Unmarshal(%#v): %v", v, err)
		return false
	}

	c.sKey = v.Skey
	c.sid = v.Wxsid
	c.uin = v.Wxuin
	c.passTicket = v.PassTicket

	c.baseRequest = make(map[string]interface{})
	c.baseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	c.baseRequest["Sid"] = v.Wxsid
	c.baseRequest["Skey"] = v.Skey
	c.baseRequest["DeviceID"] = c.deviceId

	return true
}
