/*
 @author: lynn
 @date: 2023/5/3
 @time: 21:43
*/

package BOKA

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zerocodehero/requests"
	"strconv"
)

// 不需要的单据类型

var (
	CDCLASS  = []string{"合计", "员工总合计", "小计"}
	CUTCLASS = []string{"1001", "1002", "1003", "1004", "1005", "1011", "1012", "1013", "1014", "1015", "3008", "3009", "4001"}
)

type UserData struct {
	Raw   gjson.Result
	Emps  []string
	Items []Items
}

func KaRun(params ...string) UserData {
	U := UserData{}
	U.GetUserType2()
	U.GetPerformanceData(params...)
	U.Handler()
	return U
}

type Items struct {
	UserID    string  // 会员卡号
	UserName  string  // 会员名字
	CardType  string  // 卡类型
	BillDate  string  // 消费日期
	ItemName  string  // 项目名字
	ItemID    string  // 项目 ID
	ItemAmt   float64 // 项目总金额
	ItemAmt3  float64 // 实际分得金额
	PayType   string  // 支付方式
	PayCode   string  // 支付方式ID
	EmpName   string  // 员工姓名
	EmpID     string  // 员工ID
	EmpMoney  float64 // 工资
	EmpType   string  // 员工职位
	Times     string  // 次数
	ShareRote string  // 分享比率
	Remark    string  // 备忘
	TypeID    string  // TypeIds["充值"] == 1 课程=2 消疗=3 项目=4 剪发=5
}

// GetPerformanceData
//
//	@Description:
//	@param Params startUser, endUser, startDate endDate
func (U *UserData) GetPerformanceData(Params ...string) {

	startDate, endDate, user1, user2 := paramsHandler(Params...)

	res := requests.Post(requests.Config{
		Url: "https://api.bokao2o.com/s3nos_report/person/v2/empPerformStats",
		Params: map[string]string{
			"sign": "UERLUUcjMDAy",
			"v":    "1",
		},
		Headers: TOKEN.Headers(),
		Body: map[string]interface{}{
			"compid":      "002",
			"compName":    "孔雀宫-迎春路",
			"fromdate":    startDate,
			"todate":      endDate,
			"fromempl":    user1,
			"toempl":      user2,
			"inc_card":    1,
			"inc_service": 1,
			"inc_goods":   1,
			"return_type": 1,
			"paymode":     "",
			"cardtype":    "",
			"recalculate": true,
			"type":        "2",
			"userId":      "ADMIN",
		},
	})

	if res.ERR != nil {
		return
	}
	U.Raw = res.JSON()
}

// Handler
//  @Description: 处理所有的员工业绩
//  @receiver U
//

func (U *UserData) Handler() {
	for _, v := range U.Raw.Get("result").Array() {
		// 员工id
		id := v.Array()[0].Get("person_id")
		// 去除离职员工，博卡有bug
		if U.In(U.Emps, id.String()) {
			// 循环每个员工
			for _, v1 := range v.Array() {
				// 单据类型
				cdclassstr := v1.Get("cdclass").String()

				//  去除合计 员工总合计 小计
				if !U.In(CDCLASS, cdclassstr) {

					code := v1.Get("code").String()
					// 剪发
					typeID := ""
					// 充值
					if cdclassstr == "卡充值" {
						if v1.Get("comboname").String() != "" {
							typeID = "课程"
						} else {
							typeID = "充值"
						}

					} else {
						// 所有非充值单据
						if v1.Get("paycod").String() == "9" {
							typeID = "消疗"
						} else {
							if U.In(CUTCLASS, code) {
								typeID = "剪发"
							} else {
								if code == "7077" || code == "7078" {

									if code == "7077" {
										typeID = "课程"
									} else {
										typeID = "消疗"
									}

								} else {
									typeID = "项目"
								}

							}

						}
					}
					// 所有需要的信息
					U.Items = append(U.Items, Items{
						UserID:    v1.Get("cardid").String(),
						UserName:  v1.Get("memname").String(),
						CardType:  v1.Get("cardtypname").String(),
						BillDate:  v1.Get("billdate").String(),
						ItemName:  v1.Get("name").String(),
						ItemID:    v1.Get("code").String(),
						ItemAmt:   Decimal(v1.Get("amt").Float()),
						ItemAmt3:  Decimal(v1.Get("amt3").Float()),
						PayType:   v1.Get("payway").String(),
						PayCode:   v1.Get("paycod").String(),
						EmpName:   v1.Get("empname").String(),
						EmpID:     v1.Get("person_id").String(),
						EmpType:   v1.Get("zhiname").String(),
						EmpMoney:  Decimal(v1.Get("comm").Float()),
						Times:     v1.Get("quan").String(),
						ShareRote: v1.Get("share_rate").String(),
						Remark:    v1.Get("comboname").String(),
						TypeID:    typeID,
					})
				}
			}
		}

	}

}

// GetUserType2
//
//	@Description: 获取所有的员工
//	@receiver U
func (U *UserData) GetUserType2() {
	res := requests.Get(requests.Config{
		Url: "https://api.bokao2o.com/s3nos_report/person/v2/comp/002/getEmpLsByJob",
		Params: map[string]string{
			"sign":   "UERLUUcjMDAy",
			"status": "2",
		},
		Headers: TOKEN.Headers(),
		Body:    nil,
	})

	if res.ERR != nil {
		return
	}

	for _, v := range res.JSON().Get("result").Array() {
		U.Emps = append(U.Emps, v.Get("empId").String())
	}

}

// In
//
//	@Description:
//	@receiver U
//	@param word
//	@param words
//	@return right
func (U *UserData) In(words []string, word string) (right bool) {
	for _, v := range words {
		if v == word {
			right = true
		}
	}
	return
}

// Decimal
//
//	@Description: 保留两位小数
//	@param num
//	@return float64
func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}

// paramsHandler
//
//	@Description:  处理 GetPerformanceData 的params
//	@param Params
//	@return string
//	@return string
//	@return string
//	@return string
func paramsHandler(Params ...string) (string, string, string, string) {
	startDate, endDate := GetDate()
	user1 := ""
	user2 := ""

	if len(Params) == 1 {

		if len(Params[0]) == 3 {
			user1, user2 = Params[0], Params[0]
		}
		if Params[0] == "上月" {
			startDate, endDate = GetDate("last")
		}

	}

	if len(Params) == 2 {
		if len(Params[0]) == 3 && len(Params[1]) == 3 {
			user1, user2 = Params[0], Params[1]
		}
		if len(Params[0]) == 8 && len(Params[1]) == 8 {
			startDate, endDate = Params[0], Params[1]
		}
	}

	return startDate, endDate, user1, user2
}
