package pay

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

type WechatPay interface {
	// ============通用方法============
	GetNonceStr() string
	Sign(param interface{}) (string, error)          // 生成签名
	VerifySign(param interface{}, sign string) error // 验签

	// ============功能方法============
	// 向用户账户转账接口
	Transfer(openId string, partnerTradeNo string, amount int64, checkName CheckNameMode, receiverName string, desc string, deviceInfo string, ip string) (*TransferResponse, error)

	// 微信支付 - 统一下单接口
	UnifiedOrder(openId, body, attach, goodsTag, outTradeNo string, totalFee int64, timeStart, timeExpire time.Time, notifyUrl string, tradeType TradeType) (*UnifiedOrderResponse, error)
	// 微信支付 - 退款接口
	Refund(transactionId, outTradeNo, outRefundNo string, orderTotalFee, refundFee int64, notifyUrl, refundDesc string) (*RefundResponse, error)
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

func NewWechatPay(mchId, appId, apiSignKey string, apiKeyFile, apiCertFile string, apiCA []byte, nonceLen int, timeout time.Duration) (WechatPay, error) {
	if nonceLen > 32 {
		nonceLen = 32
	}

	cliCrt, err := tls.LoadX509KeyPair(apiCertFile, apiKeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cliCrt},
	}

	if apiCA != nil {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(apiCA)

		tlsConfig.RootCAs = pool
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
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
