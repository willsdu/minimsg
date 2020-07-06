package main

import (
	"github.com/gin-gonic/gin"
	"log"
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
	log.Printf("signature %s, timestamp %s, nonce %s, echostr %s", signature, timestamp, nonce, echostr)
	localSign := mini.GenSignature(timestamp, nonce)
	if localSign == signature {
		c.String(200, echostr)
		return
	}
	c.String(http.StatusBadRequest, "")
}
func PostMsgHandle(c *gin.Context) {
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	if timestamp == "" || nonce == "" {
		c.String(http.StatusBadRequest, "")
		return
	}
	msg := mini.MiniMsg{}
	if err := c.BindXML(&msg); err != nil {
		log.Printf("err %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}
	if msg.Encrypt != mini.CheckEncrpyt(timestamp, nonce, msg.Encrypt) {
		c.String(http.StatusBadRequest, "")
		return
	}
	msgContent, err := mini.DecodeMsg(msg)
	log.Printf("msg is %+v, err is %v", msgContent, err)
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
	r.POST("/mini/msg", Broker)
	r.Run(":3305")
}
