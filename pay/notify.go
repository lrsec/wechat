package pay

import "encoding/xml"

type NotifyInfo struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"mch_appid"`
	MchId      string `xml:"mchid"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	Openid        string `xml:"openid"`
	TransactionId string `xml:"transaction_id"` // 订单号
	TotalFee      int    `xml:"total_fee"`      // 单位为分
	TradeType     string `xml:"trade_type"`     //JSAPI、NATIVE、APP
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	Attach        string `xml:"attach"`         // 用户透传数据
	TimeEnd       string `xml:"time_end"`       //支付完成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010

	IsSubscribe        string `xml:"is_subscribe"`         // Y-> yes, N-> no,
	BankType           string `xml:"bank_type"`            // like CMC
	SettlementTotalFee int    `xml:"settlement_total_fee"` //应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额
	FeeType            string `xml:"fee_type"`             //货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY
	CashFee            int    `xml:"cash_fee"`             //现金支付金额订单现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`
}

type NotifyReply struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

func (self *wechatPay) ParseNotifyInfo(body []byte) (*NotifyInfo, error) {
	info := &NotifyInfo{}

	if err := xml.Unmarshal(body, info); err != nil {
		return nil, err
	}

	if err := self.VerifySign(info, info.Sign); err != nil {
		return nil, err
	}

	return info, nil
}
