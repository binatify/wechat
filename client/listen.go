package client

import (
	"log"
	"time"

	"github.com/binatify/wechat/util"
)

func (c *Client) Listen(handler func(interface{})) {
	util.RpcCall("[*] 开启状态通知 ...", c.webwxstatusnotify)
	util.RpcCall("[*] 进行同步线路测试 ...", c.pingSynccheck)

	log.Println("[*] 开始接收消息 ...")

	for {
		time.Sleep(c.duration)

		retcode, selector, ok := c.synccheck()
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
				r := c.webwxsync()

				switch r.(type) {
				case bool:
				default:
					for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
						handler(msg)
					}
				}

			case "6", "4":
				c.webwxsync()
			}
		}
	}
}

func (c *Client) pingSynccheck() (ok bool) {
	syncHost := []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
	}

	for _, host := range syncHost {
		c.syncHost = host
		_, _, ok := c.synccheck()
		if ok {
			return true
		}
	}

	return false
}
