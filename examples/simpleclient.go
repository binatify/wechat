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
