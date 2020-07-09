package mini

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const MediaId = "D5QTs6tockb101m-8yx-BpjTFgvs-3EL7ga6CgMVdP0pFAyXuiUfGubAWpNsdk9e"

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

//DecodeMsg 解密得出消息提
func DecodeMsg(msg EncodedReceiveMsg) (MiniMsg, error) {
	realBytes, err := base64.StdEncoding.DecodeString(msg.Encrypt)
	if err != nil {
		log.Printf("decode msg error %v", err)
		return MiniMsg{}, err
	}
	tpass, err1 := AesDecrypt(realBytes, AesKey)
	if err1 != nil {
		return MiniMsg{}, err1
	}
	buf := bytes.NewBuffer(tpass[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)

	appIdStart := 20 + length

	id := tpass[appIdStart : int(appIdStart)+len(AppId)]
	if string(id) != AppId {
		return MiniMsg{}, fmt.Errorf("invalid appid")
	}
	result := MiniMsg{}
	err = xml.Unmarshal(tpass[20:20+length], &result)
	return result, err
}

func EncodeMsg(msg ImgMsg) (string, error) {
	payload, _ := xml.Marshal(msg)
	payloadStr := fmt.Sprintf("<xml>%s</xml>", string(payload))
	payload = []byte(payloadStr)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(payload)))
	if err != nil {
		return "", err
	}
	bodyLength := buf.Bytes()

	randomBytes := []byte(GenRandomKey(16))

	plainData := bytes.Join([][]byte{randomBytes, bodyLength, payload, []byte(AppId)}, nil)
	cipherData, err := AesEncrypt(plainData, AesKey)
	if err != nil {
		return "", errors.New("AesEncrypt error")
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

//EncodeMiniImgMsg 加密消息
func EncodeMiniImgMsg(toUser, nonce string, timestamp string) (string, error) {
	imgMsg := ImgMsg{
		ToUser:  toUser,
		MsgType: "image",
		Image: ImageMedia{
			MediaId: MediaId,
		},
	}
	payload, _ := json.Marshal(imgMsg)
	return string(payload), nil
}

//SendCustomMsg 发送客服消息
func SendCustomMsg(data []byte) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", AccessToken)

	buf := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Printf("NewRequest error %v", err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("sendCustMsg %+v error", string(data), err)
		return err
	}
	defer resp.Body.Close()
	baseResp := TokenResp{}
	if err := json.NewDecoder(resp.Body).Decode(&baseResp); err != nil {
		log.Printf(" Decode sendCustMsg  %+v resp error", string(data), err)
		return err
	}
	if baseResp.ErrCode > 0 {
		log.Printf(" send sendCustMsg  %+v error %+v", string(data), baseResp)
		return fmt.Errorf("send failed %s", baseResp.Errmsg)
	}
	return nil
}

func EncodeAndSend(msg MiniMsg, nonce string, timestamp string) error {
	UserTextMsg := `我要进群学习\n1.点击下方老师二维码\n2.长按识别二维码添加`
	textMsg := TextMsg{
		ToUser:  msg.FromUserName,
		MsgType: "image",
		Text: TextObject{
			Content: UserTextMsg,
		},
	}
	textPayload, _ := json.Marshal(textMsg)
	SendCustomMsg([]byte(textPayload))

	paylaod, err := EncodeMiniImgMsg(msg.FromUserName, nonce, timestamp)
	if err != nil {
		return err
	}
	log.Printf("send img is %s", string(paylaod))
	return SendCustomMsg([]byte(paylaod))
}
