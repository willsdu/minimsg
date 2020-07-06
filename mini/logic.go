package mini

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"strings"
)

func Sha1Slice(strs []string) string {
	sort.Slice(strs, func(i, j int) bool {
		return strs[i] < strs[j]
	})
	h := sha1.New()
	h.Write([]byte(strings.Join(strs, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//GenSignature 生成消息接收的时候的签名
func GenSignature(timestamp, nonce string) string {
	ps := []string{Token, timestamp, nonce}
	return Sha1Slice(ps)
}

func CheckEncrpyt(timestamp, nonce, msg_encrypt string) string {
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

//DecodeMsg 解密得出消息提
func DecodeMsg(msg MiniMsg) (MiniMsg, error) {
	realBytes := []byte{}
	_, err := base64.StdEncoding.Decode(realBytes, []byte(msg.Encrypt))
	if err != nil {
		log.Printf("decode msg error %v", err)
		return MiniMsg{}, err
	}
	tpass, err1 := AesDecrypt(realBytes, []byte(EncryptCode))
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
	log.Printf("msg content is %+v", result)
	return result, err
}
