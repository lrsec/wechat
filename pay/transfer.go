package pay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/kataras/iris/core/errors"
)

/*
企业付款相关功能
*/

const (
	TRANSFER_URL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

type TransferResponse struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	AppId          string `xml:"mch_appid"`
	MchId          string `xml:"mchid"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	PaymentNo      string `xml:"payment_no"`
	PaymentTime    string `xml:"payment_time"`
}

type transferParam struct {
	AppId          string `xml:"mch_appid"`
	Mchid          string `xml:"mchid"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	Openid         string `xml:"openid"`
	CheckName      string `xml:"check_name"`
	ReUserName     string `xml:"re_user_name"`
	Amount         int64  `xml:"amount"`
	Desc           string `xml:"desc"`
	SPBillCreateIP string `xml:"spbill_create_ip"`
}

// 企业付款到用户零钱账户
func (self *wechatPay) Transfer(openId string, partnerTradeNo string, amount int64, checkName CheckNameMode, receiverName string, desc string, deviceInfo string, ip string) (*TransferResponse, error) {

	if self.secureClient == nil {
		return nil, errors.New("need create secure wechat with CA")
	}

	param := &transferParam{
		AppId:          self.AppId,
		Mchid:          self.mchId,
		DeviceInfo:     deviceInfo,
		NonceStr:       randString(self.NonceLen),
		PartnerTradeNo: partnerTradeNo,
		Openid:         openId,
		CheckName:      string(checkName),
		ReUserName:     receiverName,
		Amount:         amount,
		Desc:           desc,
		SPBillCreateIP: ip,
	}

	sign, err := self.Sign(param)
	if err != nil {
		return nil, err
	}

	param.Sign = sign

	body, err := xml.Marshal(param)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", TRANSFER_URL, bytes.NewReader(body))
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

	resp := &TransferResponse{}
	if err = xml.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
