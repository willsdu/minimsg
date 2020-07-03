package mini

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

/*
接收消息小程序客服消息接口
*/
const (
	token       = "mmschool"
	EncryptCode = "OxNryzORnx5XVbPDCoKDyMQH5vxbkG4DrkpGWKRasfO"
	AppId       = "wxf0db6bd724116144"
	AppSecret   = "81141bad72ccd3d31541f21863651e62"
)

var AccessToken = GetToken()

type MiniMsg struct {
	//小程序的原始ID
	ToUserName string
	//发送者的openid
	FromUserName string
	//消息创建时间(整型）
	CreateTime int64
	MsgType    string
	//消息id，64位整型
	MsgId string
	//文本消息内容
	Content string `json:"omitempty"`
	//图片链接（由系统生成）
	PicUrl string `json:"omitempty"`
	//图片消息媒体id，可以调用[获取临时素材]((getTempMedia)接口拉取数据。
	MediaId      string `json:"omitempty"`
	Title        string `json:"omitempty"`
	AppId        string `json:"omitempty"`
	PagePath     string `json:"omitempty"`
	ThumbUrl     string `json:"omitempty"`
	ThumbMediaId string `json:"omitempty"`
	Event        string `json:"omitempty"`
	SessionFrom  string `json:"omitempty"`
}

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

func HandleMsg(msg MiniMsg) {

}

//GenSignature 生成消息接收的时候的签名
func GenSignature(timestamp, nonce string) string {
	ps := []string{AccessToken, timestamp, nonce}
	sort.Slice(ps, func(i, j int) bool {
		return ps[i] < ps[j]
	})
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(ps, ""))))
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
	return tokenResp.AccessToken
}
