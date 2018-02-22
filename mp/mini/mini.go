package mini

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/juju/errors"
)

type WechatMini interface {
	// 小程序登陆接口 https://mp.weixin.qq.com/debug/wxadoc/dev/api/api-login.html#wxloginobject
	GetSessionKeyByCode(jsCode string) (*GetSessionKeyByCodeResponse, error)
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
		AppId:  appId,
		Secret: secret,
		client: client,
	}
}

type wechatMini struct {
	AppId  string
	Secret string

	client *http.Client
}

func (this *wechatMini) GetSessionKeyByCode(jsCode string) (*GetSessionKeyByCodeResponse, error) {
	url := fmt.Sprintf(get_sessionkey_url, this.AppId, this.Secret, jsCode)

	result, err := this.client.Get(url)
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
