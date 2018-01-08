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

func (this *Client) webwxinit() bool {
	url := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", this.baseUri, this.passTicket, this.skey, util.UnixTimestamp())
	params := make(map[string]interface{})
	params["BaseRequest"] = this.baseRequest

	res, err := this.doPost(url, params, true)
	if err != nil {
		return false
	}

	ioutil.WriteFile("initdata.txt", res, 777)

	data, ok := util.JsonDecode(string(res)).(map[string]interface{})
	if !ok {
		return false
	}

	this.user = data["User"].(map[string]interface{})
	this.syncKeyMap = data["SyncKey"].(map[string]interface{})
	this.setsynckey()

	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

func (this *Client) setsynckey() {
	keys := []string{}
	for _, keyVal := range this.syncKeyMap["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(int)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(int)))
		keys = append(keys, key+"_"+value)
	}
	this.synckey = strings.Join(keys, "|")
}

func (this *Client) webwxstatusnotify() (ok bool) {
	urlReq := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", this.baseUri, this.passTicket)
	params := make(map[string]interface{})

	params["BaseRequest"] = this.baseRequest
	params["Code"] = 3
	params["FromUserName"] = this.user["UserName"]
	params["ToUserName"] = this.user["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())

	res, err := this.doPost(urlReq, params, true)
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

func (this *Client) synccheck() (retcode, selector string, ok bool) {
	urlReq := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", this.syncHost)

	v := url.Values{}
	v.Add("r", util.UnixTimestamp())
	v.Add("sid", this.sid)
	v.Add("uin", this.uin)
	v.Add("skey", this.skey)
	v.Add("deviceid", this.deviceId)
	v.Add("synckey", this.synckey)
	v.Add("_", util.UnixTimestamp())

	urlReq = urlReq + "?" + v.Encode()

	data, err := this.doGet(urlReq)
	if err != nil {
		return
	}

	found := syncCheckRegexp.FindStringSubmatch(string(data))
	if len(found) <= 2 {
		return
	}

	return found[1], found[2], true
}

func (this *Client) webwxsync() interface{} {
	urlReq := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&pass_ticket=%s", this.baseUri, this.sid, this.skey, this.passTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = this.baseRequest
	params["SyncKey"] = this.syncKeyMap
	params["rr"] = ^int(time.Now().Unix())
	res, err := this.doPost(urlReq, params, true)
	if err != nil {
		return false
	}

	data := util.JsonDecode(string(res)).(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)

	if retCode == 0 {
		this.syncKeyMap = data["SyncKey"].(map[string]interface{})
		this.setsynckey()
	}

	return data
}
