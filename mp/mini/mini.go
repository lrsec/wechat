package mini

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	"crypto/aes"
	"crypto/cipher"

	"encoding/base64"

	"github.com/juju/errors"
)

type WechatMini interface {
	// 小程序登陆接口 https://mp.weixin.qq.com/debug/wxadoc/dev/api/api-login.html#wxloginobject
	GetSessionKeyByCode(jsCode string) (*GetSessionKeyByCodeResponse, error)
	UnEncryptFromEncryptedData(encryptedData string, sessionKey string, iv string) (*UserInfo, error)
}

const (
	get_sessionkey_url = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type GetSessionKeyByCodeResponse struct {
	ErrCode    int64  `json:"errcode"`     // 发生错误时,返回错误码
	ErrMsg     string `json:"errmsg"`      // 发生错误时，返回具体错误信息
	OpenId     string `json:"openid"`      // 登陆用户 openid
	UnionId    string `json:"unionid"`     // 登陆用户 unionid
	SessionKey string `json:"session_key"` // 登陆用户 session_key
}

type UserInfo struct {
	OpenId    string            `json:"openId"`
	NickName  string            `json:"nickName"`
	Gender    int               `json:"gender"` // 0 未知， 1 男， 2 女
	City      string            `json:"city"`
	Province  string            `json:"province"`
	Country   string            `json:"country"`
	AvatarUrl string            `json:"avatarUrl"`
	UnionId   string            `json:"unionId"`
	Watermark UserInfoWaterMark `json:"watermark"`
}

type UserInfoWaterMark struct {
	Appid     string `json:"appid"`
	Timestamp int    `json:"timestamp"`
}

func NewWechatMini(appId, secret string, client *http.Client) WechatMini {

	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &wechatMini{
		appId:  appId,
		secret: secret,
		client: client,
	}
}

type wechatMini struct {
	appId  string
	secret string

	client *http.Client
}

func (mini *wechatMini) GetSessionKeyByCode(jsCode string) (*GetSessionKeyByCodeResponse, error) {
	url := fmt.Sprintf(get_sessionkey_url, mini.appId, mini.secret, jsCode)

	result, err := mini.client.Get(url)
	if err != nil {
		return nil, errors.Trace(err)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	resp := &GetSessionKeyByCodeResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return resp, nil
}

func (mini *wechatMini) UnEncryptFromEncryptedData(encryptedData string, sessionKey string, iv string) (*UserInfo, error) {

	decodedEncryptedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.Trace(err)
	}

	decodedSessionKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.Trace(err)
	}

	decodedIv, err := base64.StdEncoding.DecodeString(iv)

	block, err := aes.NewCipher(decodedSessionKey)
	if err != nil {
		return nil, errors.Trace(err)
	}

	decrypter := cipher.NewCBCDecrypter(block, decodedIv)
	plant := make([]byte, len(decodedEncryptedData))
	decrypter.CryptBlocks(plant, decodedEncryptedData)
	plant = PKCS7UnPadding(plant)
	plant = []byte(strings.Replace(string(plant), "\a", "", -1))

	userInfo := &UserInfo{}
	err = json.Unmarshal(plant, userInfo)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return userInfo, nil

}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}
