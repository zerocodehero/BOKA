/*
 @author: lynn
 @date: 2023/5/3
 @time: 21:42
*/

package main

import "time"

// Async
//
//	@Description: 定时器 开始执行一次，后根据延迟传入的秒数循环执行
//	@param times
//	@param FuncName
func Async(times int, FuncName func()) {
	// 先执行一次
	FuncName()
	//创建Ticker，设置多长时间触发一次
	ticker := time.NewTicker(time.Duration(int64(times)) * time.Second)
	go func() {
		//遍历ticker.C，如果有值，则会执行do something，否则阻塞
		for range ticker.C {
			FuncName()
		}
	}()

}

// GetDate
//
//	@Description: 获取当月和上月日期  20230101 20230131
//	@param params 有值代表上月
//	@return startDate
//	@return endDate
func GetDate(params ...string) (startDate string, endDate string) {

	timeNow := time.Now()

	if len(params) > 0 {
		//上个月
		LastMonth := timeNow.AddDate(0, -1, 0)
		currentYear, currentMonth, _ := LastMonth.Date()
		currentLocation := LastMonth.Location()
		startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		startDate = startOfMonth.Format("20060102")
		endDate = endOfMonth.Format("20060102")

	} else {
		// 当月
		currentYear, currentMonth, _ := timeNow.Date()
		currentLocation := timeNow.Location()
		startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		startDate = startOfMonth.Format("20060102")
		endDate = endOfMonth.Format("20060102")

	}
	return
}
