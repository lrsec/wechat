package pay

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

type WechatPay interface {
	GetNonceStr() string
	Sign(param interface{}) (string, error)

	Transfer(openId string, partnerTradeNo string, amount int64, checkName CheckNameMode, receiverName string, desc string) (*TransferResponse, error)
	UnifiedOrder(openId, body, attach, goodsTag, outTradeNo string, totalFee int64, timeStart, timeExpire time.Time, notifyUrl string, tradeType TradeType) (*UnifiedOrderResponse, error)
}

func NewUnSecureWechatPay(mchId, appId, apiSignKey string, nonceLen int, timeout time.Duration) WechatPay {
	if nonceLen > 32 {
		nonceLen = 32
	}

	nonsecureClient := &http.Client{
		Timeout: timeout,
	}

	pay := &wechatPay{
		mchId:           mchId,
		AppId:           appId,
		apiSignKey:      apiSignKey,
		NonceLen:        nonceLen,
		nonSecureClient: nonsecureClient,
	}

	return pay

}

func NewWechatPay(mchId, appId, apiSignKey string, apiKeyFile, apiCertFile string, apiCA []byte, ip, deviceInfo string, nonceLen int, timeout time.Duration) (WechatPay, error) {
	if nonceLen > 32 {
		nonceLen = 32
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(apiCA)
	cliCrt, err := tls.LoadX509KeyPair(apiCertFile, apiKeyFile)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	nonsecureClient := &http.Client{
		Timeout: timeout,
	}

	pay := &wechatPay{
		mchId:           mchId,
		apiSignKey:      apiSignKey,
		ip:              ip,
		deviceInfo:      deviceInfo,
		AppId:           appId,
		NonceLen:        nonceLen,
		secureClient:    client,
		nonSecureClient: nonsecureClient,
	}

	return pay, nil
}

type wechatPay struct {
	mchId      string //商户号
	AppId      string // 应用id, 商户号可以支持多个 appid, 可修改共用
	NonceLen   int    // 随机字符串 nonce_str 长度，最长支持32字符
	apiSignKey string //api 签名用密钥，在后台进行设置

	ip         string //请求 ip 地址
	deviceInfo string // 微信支付分配的终端设备号

	apiPublicKey    string       //api 接口密钥，微信生成，通过后台下载
	secureClient    *http.Client // 要求证书的请求
	nonSecureClient *http.Client // 不要求证书的请求
}

func (pay *wechatPay) GetNonceStr() string {
	return randString(pay.NonceLen)
}

type CheckNameMode string

const (
	NO_CHECK    CheckNameMode = "NO_CHECK"
	FORCE_CHECK               = "FORCE_CHECK"
)
