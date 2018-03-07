package pay

import (
	"fmt"
	"reflect"
	"strings"

	"sort"

	"errors"

	"code.inke.cn/boc/financial/fin-conversion/util/md5"
)

/*
微信签名规则，参见: https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=4_3
*/
func (self *wechatPay) sign(param interface{}) (string, error) {

	if content, err := self.genContentStr(param); err != nil {
		return "", err
	} else {
		return strings.ToUpper(md5.Md5(content)), nil
	}
}

func (self *wechatPay) genContentStr(param interface{}) (contentStr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			contentStr = ""
			err = errors.New(fmt.Sprintf("sign panic. msg: %v", r))
		}
	}()

	if param == nil {
		return "", nil
	}

	lt := reflect.TypeOf(param)
	rv := reflect.ValueOf(param)
	if lt.Kind() == reflect.Ptr {
		lt = lt.Elem()
		rv = rv.Elem()
	}

	fieldNum := lt.NumField()

	if fieldNum == 0 {
		return "", nil
	}

	kv := make(map[string]int)
	names := make([]string, 0)
	for i := 0; i < fieldNum; i++ {
		tag := lt.Field(i).Tag.Get("xml")
		if tag != "" && tag != "-" {
			names = append(names, tag)
			kv[tag] = i
		}
	}
	sort.Strings(names)

	for _, name := range names {
		f := rv.Field(kv[name])

		st := lt.Field(kv[name])
		if st.Type.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		valueStr := fmt.Sprintf("%v", f)
		if valueStr != "" {
			contentStr = contentStr + name + "=" + valueStr + "&"
		}
	}

	contentStr = contentStr + "key=" + self.apiSignKey

	return contentStr, nil
}
