package mini

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"sort"
	"strings"
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

var AesKey = DecodeEncodeAesKey(EncryptCode)
var AccessToken = GetToken()

type MiniMsg struct {
	//小程序的原始ID
	ToUserName string `json:"omitempty" xml:"ToUserName,omitempty"`
	//发送者的openid
	FromUserName string `json:"omitempty" xml:"FromUserName,omitempty"`
	//消息创建时间(整型）
	CreateTime int64  `json:"omitempty" xml:"CreateTime,omitempty"`
	MsgType    string `json:"omitempty" xml:"MsgType,omitempty"`
	//消息id，64位整型
	MsgId string `json:"omitempty" xml:"MsgId,omitempty"`
	//文本消息内容
	Content string `json:"omitempty" xml:"Content,omitempty"`
	//图片链接（由系统生成）
	PicUrl string `json:"omitempty" xml:"PicUrl,omitempty"`
	//图片消息媒体id，可以调用[获取临时素材]((getTempMedia)接口拉取数据。
	MediaId      string `json:"omitempty" xml:"MediaId,omitempty"`
	Title        string `json:"omitempty" xml:"Title,omitempty"`
	AppId        string `json:"omitempty" xml:"AppId,omitempty"`
	PagePath     string `json:"omitempty" xml:"PagePath,omitempty"`
	ThumbUrl     string `json:"omitempty" xml:"ThumbUrl,omitempty"`
	ThumbMediaId string `json:"omitempty" xml:"ThumbMediaId,omitempty"`
	Event        string `json:"omitempty" xml:"Event,omitempty"`
	SessionFrom  string `json:"omitempty" xml:"SessionFrom,omitempty"`
}

type TokenResp struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpireIn    int    `json:"expires_in,omitempty"`
	ErrCode     int    `json:"errcode,omitempty"`
	Errmsg      string `json:"errmsg,omitempty"`
}

type ImgMsg struct {
	ToUser  string     `json:"touser"`
	MsgType string     `json:"msgtype"`
	Image   ImageMedia `json:"image"`
}
type ImageMedia struct {
	MediaId string `json:"media_id"`
}

//EncodeMsg 接收来的消息
type EncodedReceiveMsg struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}
type EncodedRespMsg struct {
	Encrypt      string `xml:"Encrypt"`
	MsgSignature string `xml:"MsgSignature"`
	TimeStamp    string `xml:"TimeStamp"`
	Nonce        string `xml:"Nonce"`
}

func Sha1Slice(strs []string) string {
	sort.Slice(strs, func(i, j int) bool {
		return strs[i] < strs[j]
	})
	h := sha1.New()
	h.Write([]byte(strings.Join(strs, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//GenSignature 生成消息推送校验接口的签名
func GenSignature(timestamp, nonce string) string {
	ps := []string{Token, timestamp, nonce}
	return Sha1Slice(ps)
}

//GenEncrpyt  生成消息推送时候的校验签名
func GenEncrpyt(timestamp, nonce, msg_encrypt string) string {
	ps := []string{Token, timestamp, nonce, msg_encrypt}
	return Sha1Slice(ps)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
func DecodeEncodeAesKey(code string) []byte {
	aesKey, _ := base64.StdEncoding.DecodeString(code + "=")
	return aesKey
}

//GenRandomKey 生成指定长度的字符串
func GenRandomKey(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
