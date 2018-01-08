package client

import (
	"log"
	"time"

	"github.com/binatify/wechat/util"
)

func (this *Client) Listen(handler func(interface{})) {
	util.RpcCall("[*] 开启状态通知 ...", this.webwxstatusnotify)
	util.RpcCall("[*] 进行同步线路测试 ...", this.pingSynccheck)

	log.Println("[*] 开始接收消息 ...")

	for {
		time.Sleep(this.duration)

		retcode, selector, ok := this.synccheck()
		if !ok {
			continue
		}

		if retcode == "1100" {
			log.Println("[*] 你在手机上登出了微信，再见")
			return
		}

		if retcode == "1101" {
			log.Println("[*] 你在其他地方登录了 WEB 版微信，再见")
			return
		}

		if retcode == "0" {
			switch selector {
			case "2":
				r := this.webwxsync()

				switch r.(type) {
				case bool:
				default:
					for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
						handler(msg)
					}
				}

			case "6", "4":
				this.webwxsync()
			}
		}
	}
}

func (this *Client) pingSynccheck() (ok bool) {
	syncHost := []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
	}

	for _, host := range syncHost {
		this.syncHost = host
		_, _, ok := this.synccheck()
		if ok {
			return true
		}
	}

	return false
}
