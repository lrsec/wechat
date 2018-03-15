package pay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

/*
微信支付退款接口
https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
*/

const (
	REFUND_URL = "https://api.mch.weixin.qq.com/secapi/pay/refund"
)

type RefundParam struct {
	AppId         string `xml:"appid"`
	Mchid         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"` //随机串，必填
	Sign          string `xml:"sign"`
	TransactionId string `xml:"transaction_id"`  // 微信订单号，与商户订单号需要二选一填写
	OutTradeNo    string `xml:"out_trade_no"`    // 商户订单号，与微信订单号需要二选一填写
	OutRefundNo   string `xml:"out_refund_no"`   // 商户退款单号
	TotalFee      int64  `xml:"total_fee"`       // 订单总金额，单位为分，只能为整数
	RefundFee     int64  `xml:"refund_fee"`      //退款总金额，订单总金额，单位为分，只能为整数
	RefundFeeType string `xml:"refund_fee_type"` //退款货币类型，需与支付一致，或者不填。符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	RefundDesc    string `xml:"refund_desc"`     // 若商户传入，会在下发给用户的退款消息中体现退款原因
	NotifyUrl     string `xml:"notify_url"`      // 回调地址
}

type RefundResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	AppId    string `xml:"mch_appid"`
	MchId    string `xml:"mchid"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"Sign"`

	TransactionId       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no"`
	RefundId            string `xml:"refund_id"`             // 微信退款单号
	RefundFee           int64  `xml:"refund_fee"`            // 退款金额
	SettlementRefundFee int64  `xml:"settlement_refund_fee"` // 应结退款金额
	FeeType             string `xml:"fee_type"`
	CashFee             int64  `xml:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type"`
	CashRefundFee       int64  `xml:"cash_refund_fee"`
}

func (self *wechatPay) Refund(transactionId, outTradeNo, outRefundNo string, orderTotalFee, refundFee int64, notifyUrl, refundDesc string) (*RefundResponse, error) {
	param := &RefundParam{
		AppId:         self.AppId,
		Mchid:         self.mchId,
		NonceStr:      randString(self.NonceLen),
		TransactionId: transactionId,
		OutTradeNo:    outTradeNo,
		OutRefundNo:   outRefundNo,
		TotalFee:      orderTotalFee,
		RefundFee:     refundFee,
		RefundFeeType: "CNY",
		RefundDesc:    refundDesc,
		NotifyUrl:     notifyUrl,
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

	request, err := http.NewRequest("POST", REFUND_URL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	result, err := self.secureClient.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	resp := &RefundResponse{}
	if err = xml.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil

}
