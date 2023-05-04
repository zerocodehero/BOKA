/*
 @author: lynn
 @date: 2023/5/3
 @time: 21:36
*/

package BOKA

import (
	"fmt"
	"github.com/zerocodehero/requests"
	"time"
)

type Token struct {
	Token string
	ID    string
	Expir int64
}

var TOKEN = Token{}

func (T *Token) Get() (string, string) {
	timeNow := time.Now().Unix()

	if T.Token == "" || T.ID == "" || timeNow-T.Expir > CONFIG.BoKa.Sec {
		T.RequestToken()
	}
	return T.Token, T.ID
}

func (T *Token) RequestToken() {

	res := requests.Post(requests.Config{
		Url:    "https://api.bokao2o.com/auth/merchant/v2/user/login",
		Params: nil,
		Headers: map[string]string{
			"referer": "https://s3.boka.vc/",
		},
		Body: map[string]interface{}{
			"custId":   CONFIG.BoKa.CustId,
			"compId":   CONFIG.BoKa.CompId,
			"userName": CONFIG.BoKa.UserName,
			"passWord": CONFIG.BoKa.PassWord,
			"source":   CONFIG.BoKa.Source,
		},
	})

	if res.ERR != nil {
		return
	}

	jsonStr := res.JSON()
	token := jsonStr.Get("result.token").String()
	shopId := jsonStr.Get("result.token").String()

	TOKEN.Token, TOKEN.ID = token, shopId
	TOKEN.Expir = time.Now().Unix()
}

func (T *Token) Headers() map[string]string {
	token, shopId := TOKEN.Get()
	return map[string]string{
		"access_token": token,
		"accesstoken":  token,
		"device_id":    "s3backend",
		"deviceid":     "s3backend",
		"referer":      "https://s3.boka.vc/home",
		"Cookie":       fmt.Sprintf(`subCustType=; token=%s; custCode=PDKQG; custId=PDKQG; compId=002; shopId=%s; empId=admin; empName=%%25E7%%25B3%%25BB%%25E7%%25BB%%259F%%25E7%%25AE%%25A1%%25E7%%2590%%2586%%25E5%%2591%%2598;`, token, shopId),
	}

}
