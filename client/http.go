package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (c *Client) doGet(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("http.NewRequest(%s): %v\n", url, err)
		return
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "login.weixin.qq.com")
	req.Header.Add("Referer", "https://wx.qq.com/")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("httpClient.Do(): %v\n", err)
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll(): %v\n", err)
	}

	return
}

func (c *Client) doPost(reqURL string, params map[string]interface{}, isJson bool) (body []byte, err error) {
	var resp *http.Response

	if isJson {
		payload, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalf("http.NewRequest(%s): %v\n", reqURL, err)
			return nil, err
		}

		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		request.Header.Add("Referer", "https://wx.qq.com/")
		request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")

		resp, err = c.Do(request)
	} else {

		v := url.Values{}

		for key, value := range params {
			v.Add(key, value.(string))
		}

		resp, err = c.PostForm(reqURL, v)
	}

	if err != nil {
		log.Fatalf("doPost(): %v", err)
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll(): %v", err)
	}

	return
}
