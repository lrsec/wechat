package pay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

// 微信支付统一下单接口
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1

const (
	UNIFIED_ORDER_URL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type TradeType string

const (
	TRADE_TYPE_JSAPI  = "JSAPI"
	TRADE_TYPE_NATIVE = "NATIVE"
	TRADE_TYPE_APP    = "APP"
)

type unifiedOrderParam struct {
	AppId          string `xml:"appid"`
	Mchid          string `xml:"mch_id"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"Sign"`
	Body           string `xml:"body"`
	Detail         string `xml:"detail"`
	Attach         string `xml:"attach"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int64  `xml:"total_fee"`
	SPBillCreateIP string `xml:"spbill_create_ip"`
	TimeStart      string `xml:"time_start"`
	TimeExpire     string `xml:"time_expire"`
	NotifyUrl      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	Openid         string `xml:"openid"`
	GoodsTag       string `xml:"goods_tag"`
}

type UnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"mch_appid"`
	MchId      string `xml:"mchid"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"Sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	TradeType string `xml:"trade_type"`
	PrePayId  string `xml:"prepay_id"`
	CodeUrl   string `xml:"code_url"`
}

// 统一下单接口
func (self *wechatPay) UnifiedOrder(openId, body, attach, goodsTag, outTradeNo string, totalFee int64, timeStart, timeExpire time.Time, notifyUrl string, tradeType TradeType) (*UnifiedOrderResponse, error) {
	param := &unifiedOrderParam{
		AppId:          self.AppId,
		Mchid:          self.mchId,
		DeviceInfo:     self.deviceInfo,
		NonceStr:       randString(self.NonceLen),
		Body:           body,
		Attach:         attach,
		OutTradeNo:     outTradeNo,
		TotalFee:       totalFee,
		SPBillCreateIP: self.ip,
		TimeStart:      timeStart.Format("20060102150405"),
		TimeExpire:     timeExpire.Format("20060102150405"),
		NotifyUrl:      notifyUrl,
		TradeType:      string(tradeType),
		Openid:         openId,
	}

	sign, err := self.Sign(param)
	if err != nil {
		return nil, err
	}

	param.Sign = sign

	requestBody, err := xml.Marshal(param)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", UNIFIED_ORDER_URL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	result, err := self.nonSecureClient.Do(request)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	resp := &UnifiedOrderResponse{}
	if err = xml.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
