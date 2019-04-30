package examples

import (
	"github.com/binatify/wechat/client"
	"github.com/binatify/wechat/config"
	"log"
)

func main() {
	c := client.New(&config.Config{
		Duration: 6,
	}).Start()

	c.Listen(func(msg interface{}) {
		content := msg.(map[string]interface{})["Content"].(string)
		log.Println(content)
	})
}
