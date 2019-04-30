package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/binatify/wechat/util"
)

func (c *Client) webwxinit() bool {
	url := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", c.baseUri, c.passTicket, c.skey, util.UnixTimestamp())
	params := make(map[string]interface{})
	params["BaseRequest"] = c.baseRequest

	res, err := c.doPost(url, params, true)
	if err != nil {
		return false
	}

	ioutil.WriteFile("initdata.txt", res, 777)

	data, ok := util.JsonDecode(string(res)).(map[string]interface{})
	if !ok {
		return false
	}

	c.user = data["User"].(map[string]interface{})
	c.syncKeyMap = data["SyncKey"].(map[string]interface{})
	c.setsynckey()

	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

func (c *Client) setsynckey() {
	keys := []string{}
	for _, keyVal := range c.syncKeyMap["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(int)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(int)))
		keys = append(keys, key+"_"+value)
	}
	c.synckey = strings.Join(keys, "|")
}

func (c *Client) webwxstatusnotify() (ok bool) {
	urlReq := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", c.baseUri, c.passTicket)
	params := make(map[string]interface{})

	params["BaseRequest"] = c.baseRequest
	params["Code"] = 3
	params["FromUserName"] = c.user["UserName"]
	params["ToUserName"] = c.user["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())

	res, err := c.doPost(urlReq, params, true)
	if err != nil {
		return false
	}

	data := util.JsonDecode(string(res)).(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

var (
	syncCheckRegexp = regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
)

func (c *Client) synccheck() (retcode, selector string, ok bool) {
	urlReq := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", c.syncHost)

	v := url.Values{}
	v.Add("r", util.UnixTimestamp())
	v.Add("sid", c.sid)
	v.Add("uin", c.uin)
	v.Add("skey", c.skey)
	v.Add("deviceid", c.deviceId)
	v.Add("synckey", c.synckey)
	v.Add("_", util.UnixTimestamp())

	urlReq = urlReq + "?" + v.Encode()

	data, err := c.doGet(urlReq)
	if err != nil {
		return
	}

	found := syncCheckRegexp.FindStringSubmatch(string(data))
	if len(found) <= 2 {
		return
	}

	return found[1], found[2], true
}

func (c *Client) webwxsync() interface{} {
	urlReq := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&pass_ticket=%s", c.baseUri, c.sid, c.skey, c.passTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = c.baseRequest
	params["SyncKey"] = c.syncKeyMap
	params["rr"] = ^int(time.Now().Unix())
	res, err := c.doPost(urlReq, params, true)
	if err != nil {
		return false
	}

	data := util.JsonDecode(string(res)).(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)

	if retCode == 0 {
		c.syncKeyMap = data["SyncKey"].(map[string]interface{})
		c.setsynckey()
	}

	return data
}
