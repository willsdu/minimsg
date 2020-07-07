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
	msgSignature := c.Query("msg_signature")
	if timestamp == "" || nonce == "" || msgSignature == "" {
		log.Printf("invalid params")
		c.String(http.StatusBadRequest, "")
		return
	}
	msg := mini.EncodedReceiveMsg{}
	if err := c.BindXML(&msg); err != nil {
		log.Printf("BindXML error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": err.Error(),
		})
		return
	}
	log.Printf("msg is %+v", msg)
	if msgSignature != mini.GenEncrpyt(timestamp, nonce, msg.Encrypt) {
		log.Printf("GenEncrpyt %+v is not mathed", msg)
		c.String(http.StatusBadRequest, "")
		return
	}
	msgContent, err := mini.DecodeMsg(msg)
	if err != nil {
		log.Printf("DecodeMsg %+v error %+v", msg, err)
		c.String(http.StatusBadRequest, "")
		return
	}
	err = mini.EncodeAndSend(msgContent, nonce, timestamp)
	log.Printf("send msg is %+v, err is %v", msgContent, err)
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
