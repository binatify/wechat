# wechat

A new WeChat robot with golang.


### Usage

New client with message handler example:

```
package main

import (
	"log"

	"github.com/binatify/wechat/client"
	"github.com/binatify/wechat/config"
)

func main() {
	client.New(&config.Config{
		Listen:   ":9091",
		Duration: 1,
	}).Start().Listen(func(msg interface{}) {
		content := msg.(map[string]interface{})["Content"].(string)
		log.Println(content)
	})
}

```
