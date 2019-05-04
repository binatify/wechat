package client

import (
	"time"
)

func (c *Client) Listen(handler func(interface{})) {
	c.rpcCall("开启状态通知 ...", c.wxStatusNotify)
	c.rpcCall("进行同步线路测试 ...", c.pingSyncCheck)

	c.log.Println("开始接收消息 ...")

	for {
		time.Sleep(c.duration)

		retCode, selector, ok := c.syncCheck()
		if !ok {
			continue
		}

		if retCode == "1100" {
			c.log.Errorf("你在手机上登出了微信，再见")
			return
		}

		if retCode == "1101" {
			c.log.Error("你在其他地方登录了 WEB 版微信，再见")
			return
		}

		if retCode == "0" {
			switch selector {
			case "2":
				r := c.wxSync()

				switch r.(type) {
				case bool:
				default:
					for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
						handler(msg)
					}
				}

			case "6", "4":
				c.wxSync()
			}
		}
	}
}

var (
	syncHosts = []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
	}
)

func (c *Client) pingSyncCheck() (ok bool) {
	for _, host := range syncHosts {
		c.syncHost = host
		if _, _, ok := c.syncCheck(); ok {
			return true
		}
	}

	return false
}
