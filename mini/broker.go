package mini

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
接收消息小程序客服消息接口
*/
const (
	Token       = "mmschool"
	EncryptCode = "OxNryzORnx5XVbPDCoKDyMQH5vxbkG4DrkpGWKRasfO"
	AppId       = "wxf0db6bd724116144"
	AppSecret   = "81141bad72ccd3d31541f21863651e62"
)

var AccessToken = GetToken()

type MiniMsg struct {
	//小程序的原始ID
	ToUserName string `xml:"ToUserName"`
	//发送者的openid
	FromUserName string `xml:"FromUserName"`
	//消息创建时间(整型）
	CreateTime int64  `xml:"CreateTime"`
	MsgType    string `xml:"MsgType"`
	//消息id，64位整型
	MsgId string `xml:"MsgId"`
	//文本消息内容
	Content string `json:"omitempty" xml:"Content"`
	//图片链接（由系统生成）
	PicUrl string `json:"omitempty" xml:"PicUrl"`
	//图片消息媒体id，可以调用[获取临时素材]((getTempMedia)接口拉取数据。
	MediaId      string `json:"omitempty" xml:"MediaId"`
	Title        string `json:"omitempty" xml:"Title"`
	AppId        string `json:"omitempty" xml:"AppId"`
	PagePath     string `json:"omitempty" xml:"PagePath"`
	ThumbUrl     string `json:"omitempty" xml:"ThumbUrl"`
	ThumbMediaId string `json:"omitempty" xml:"ThumbMediaId"`
	Event        string `json:"omitempty" xml:"Event"`
	SessionFrom  string `json:"omitempty" xml:"SessionFrom"`
	Encrypt      string `json:"omitempty" xml:"Encrypt"`
}

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type ImgMsg struct {
	ToUser  string     `json:"touser"`
	MsgType string     `json:"msgtype"`
	Image   ImageMedia `json:"image"`
}
type ImageMedia struct {
	MediaId string `json:"media_id"`
}


//http://simg01.gaodunwangxiao.com/v/Uploads/avatar/000/00/00/1536739621__avatar_ori.jpg

//SendCustomMsg 发送客服消息
func SendCustomMsg(msg ImgMsg) {
	//url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", AccessToken)
	//
	//payload, _ := json.Marshal(msg)
	//

}

func GetToken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", AppId, AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("get access_token error %v", err)
		return ""
	}
	defer resp.Body.Close()

	tokenResp := TokenResp{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Printf("decode msg error %v", err)
		return ""
	}
	log.Printf("token is %s", tokenResp.AccessToken)
	return tokenResp.AccessToken
}
