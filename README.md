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

![screen shot 2018-01-08 at 6 26 04 pm](https://user-images.githubusercontent.com/1459834/34666894-7f92a87a-f4a1-11e7-9dc0-0d49de6d9eb1.png)


### Features

- [x] Auto login with qrcode link. 
- [x] Message linsten.
- [ ] Messge sender.
