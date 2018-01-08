package client

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
	loginedCode  = "201"
	loginTimeout = "408"
)

func (this *Client) getUUID() (ok bool) {
	url := "https://login.weixin.qq.com/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + util.UnixTimestamp()

	data, err := this.doGet(url)
	if err != nil {
		return
	}

	re := regexp.MustCompile(`"([\S]+)"`)
	found := re.FindStringSubmatch(string(data))

	if len(found) <= 1 {
		return
	}

	this.uuid = found[1]
	return true
}

func (this *Client) qrCode() (ok bool) {
	url := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s?t=webwx&_=%s", this.uuid, util.UnixTimestamp())
	resp, err := this.doGet(url)
	if err != nil {
		log.Fatalf("this.doGet(%s): %v", url, err)
		return
	}

	path := "qrcode.jpg"
	err = ioutil.WriteFile(path, resp, 0755)
	if err != nil {
		log.Fatalf("ioutil.WriteFile(qrcode.jpg): %v", err)
		return
	}

	if runtime.GOOS == "darwin" {
		exec.Command("open", path).Run()
	} else {
		go func() {
			fmt.Printf("please open on web broswer %s/qrcode", this.cfg.Listen)
			http.HandleFunc("/qrcode", func(w http.ResponseWriter, req *http.Request) {
				http.ServeFile(w, req, "qrcode.jpg")
				return
			})
			http.ListenAndServe(this.cfg.Listen, nil)
		}()
	}
	return true
}

func (this *Client) qrCodeConfirm() bool {
	for {
		if !this.doConfirm(1) {
			continue
		}

		log.Println("[*] 请在手机上点击确认 ...")

		if !this.doConfirm(0) {
			continue
		}

		break
	}

	return true
}

func (this *Client) doConfirm(tip int) (ok bool) {
	time.Sleep(time.Duration(tip) * time.Second)

	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	url += "?tip=" + strconv.Itoa(tip) + "&uuid=" + this.uuid + "&_=" + util.UnixTimestamp()

	data, err := this.doGet(url)
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
				this.redirectUri = rUri
				re = regexp.MustCompile(`/`)

				found := re.FindAllStringIndex(rUri, -1)
				this.baseUri = rUri[:found[len(found)-1][0]]
				return true
			}

		case loginedCode:
			return true

		case loginTimeout:
			log.Fatalln("[登陆超时]")

		default:
			log.Fatalln("[登陆异常]")

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

func (this *Client) login() (ok bool) {
	data, err := this.doGet(this.redirectUri)
	if err != nil {
		return
	}

	var v loginResult

	err = xml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalf("error: %v", err)
		return false
	}

	this.skey = v.Skey
	this.sid = v.Wxsid
	this.uin = v.Wxuin
	this.passTicket = v.PassTicket

	this.baseRequest = make(map[string]interface{})
	this.baseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	this.baseRequest["Sid"] = v.Wxsid
	this.baseRequest["Skey"] = v.Skey
	this.baseRequest["DeviceID"] = this.deviceId

	return true
}
