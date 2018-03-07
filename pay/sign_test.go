package pay

import (
	"strings"
	"testing"

	"code.inke.cn/boc/financial/fin-conversion/util/md5"
)

var (
	pay     *wechatPay
	param   testParam
	pParam  *testParam
	wantStr string
	wantMd5 string
)

func init() {
	pay = &wechatPay{
		apiSignKey: "test-sign-key",
	}

	i := 10
	f := 11.1
	b := true
	es := ""
	s := "test"

	param = testParam{
		S:           s,
		I:           i,
		F:           f,
		EmptyString: es,
		B:           b,

		PS:           &s,
		PI:           &i,
		PF:           &f,
		PEmptyString: &es,
		PB:           &b,
	}

	pParam = &param

	wantStr = "b=true&f=11.1&i=10&pb=true&pf=11.1&pi=10&ps=test&s=test&key=test-sign-key"

	wantMd5 = strings.ToUpper(md5.Md5(wantStr))

}

type testParam struct {
	S           string  `xml:"s"`
	I           int     `xml:"i"`
	F           float64 `xml:"f"`
	EmptyString string  `xml:"es"`
	B           bool    `xml:"b"`

	PS           *string  `xml:"ps"`
	PI           *int     `xml:"pi"`
	PF           *float64 `xml:"pf"`
	PEmptyString *string  `xml:"pes"`
	PB           *bool    `xml:"pb"`
}

func Test_wechatPay_sign(t *testing.T) {
	result, err := pay.sign(param)
	if err != nil {
		t.Errorf("sign return err: %v", err)
	}

	if result != wantMd5 {
		t.Errorf("sign fail for struct. want: %v. get: %v", wantMd5, result)
	}

	result, err = pay.sign(pParam)
	if err != nil {
		t.Errorf("sign return err: %v", err)
	}

	if result != wantMd5 {
		t.Errorf("sign fail for pointer. want: %v. get: %v", wantMd5, result)
	}
}

func Test_wechatPay_genContentStr(t *testing.T) {
	result, err := pay.genContentStr(param)
	if err != nil {
		t.Errorf("genContentStr return err: %v", err)
	}

	if result != wantStr {
		t.Errorf("genContentStr fail for struct. want: %v. get: %v", wantStr, result)
	}

	result, err = pay.genContentStr(pParam)
	if err != nil {
		t.Errorf("genContentStr return err: %v", err)
	}

	if result != wantStr {
		t.Errorf("genContentStr fail for pointer. want: %v. get: %v", wantStr, result)
	}
}
