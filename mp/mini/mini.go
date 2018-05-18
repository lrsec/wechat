package mini

import (
	"fmt"
	"net/http"
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
	UnEncryptFromEncryptedData(encryptedData []byte, sessionKey []byte, iv []byte) (string, error)
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

func (mini *wechatMini) UnEncryptFromEncryptedData(encryptedData []byte, sessionKey []byte, iv []byte) (string, error) {
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err := aes.NewCipher(sessionKey)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(encryptedData))
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, encryptedData)

	originInfo, err := base64.StdEncoding.DecodeString(string(decrypted))
	if err != nil {
		return "", err
	}

	return string(originInfo), nil
}
