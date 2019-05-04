# wechat

A new WeChat robot with golang.


### Usage

New client with message handler example:

```
package main

import (
	"github.com/binatify/wechat/client"
	"github.com/binatify/wechat/config"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logger := logrus.WithField("reqID", "reqID")

	c := client.New(&config.Config{
		Duration: 6,
	}, logger).Start()

	c.Listen(func(msg interface{}) {
		content := msg.(map[string]interface{})["Content"].(string)
		logger.Println(content)
	})

}
```

![screen shot 2018-01-08 at 6 26 04 pm](https://user-images.githubusercontent.com/1459834/34666894-7f92a87a-f4a1-11e7-9dc0-0d49de6d9eb1.png)


### Features

- [x] Auto login with qrcode. 
- [x] Message linsten.
- [ ] Messge send.
