package main

import (
	"github.com/gin-gonic/gin"
	"minimsg/mini"
	"net/http"
)

func Check(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	if signature == "" {
		c.String(http.StatusBadRequest, "")
		return
	}
	localSign := mini.GenSignature(timestamp, nonce)
	if localSign == signature {
		c.String(200, echostr)
		return
	}
	c.String(http.StatusBadRequest, "")
}
func PostMsgHandle(c *gin.Context) {
	msg := mini.MiniMsg{}
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}
	go mini.HandleMsg(msg)
	c.String(200, "success")
}

func Broker(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		Check(c)
		return
	}
	PostMsgHandle(c)
}

func main() {
	r := gin.Default()
	r.GET("/mini/msg", Broker)
	r.Run(":3305")
}
