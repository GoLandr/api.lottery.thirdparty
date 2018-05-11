package utils

import (
	"fmt"
	"mycommon/logs"
	"mycommon/utils"
	"runtime"
	"time"
)

func GetBigSmall(ball int, split int) (int, int) {
	big := 0
	small := 0
	if ball >= split {
		big = 1
	} else {
		small = 1
	}
	return big, small
}
func GetOddEven(ball int) (int, int) {
	odd := 0
	even := 0
	if ball%2 == 0 {
		even = 1
	} else {
		odd = 1
	}
	return odd, even
}

//总和单双大小统计
func GetTotalStat(ball []int, split int) (int, int, int, int) {
	odd := 0
	even := 0
	big := 0
	small := 0
	total := 0
	for _, v := range ball {
		total += v
	}
	if total%2 == 0 {
		even = 1
	} else {
		odd = 1
	}
	if total >= split {
		big = 1
	} else {
		small = 1
	}
	return odd, even, big, small
}

//龙虎统计
func GetPredStat(first_ball int, second_ball int, dragon int, tiger int, draw int) (int, int, int) {
	if first_ball > second_ball {
		dragon += 1
		tiger = 0
		draw = 0
	} else if first_ball < second_ball {
		tiger += 1
		dragon = 0
		draw = 0
	} else {
		draw += 1
		tiger = 0
		dragon = 0
	}
	return dragon, tiger, draw
}

func GetMaxValue(retVal ...int) int {
	maxVal := 0
	for _, v := range retVal {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

//判断场次时间是否满足
func JudgeTime(startTime string, endTime string) bool {
	flag := true

	//判断开始时间和结束时间
	if startTime != "" {
		_at, ok := HMSCompToSys(startTime, false)
		if ok && _at > 0 {
			flag = true
		} else {
			flag = false
		}
	}
	dayFlag := false
	if endTime > startTime {
		dayFlag = true
	}
	//	logs.Debug("startTime_", startTime, "_endTime_", endTime, "_ss_", ss)
	if endTime != "" {
		_at, ok := HMSCompToSys(endTime, dayFlag)
		if ok && _at < 0 {
			flag = true
		} else {
			flag = false
		}
	}
	return flag
}

//传入时间于系统当前时间进行比较 timestr->9:00:00
func HMSCompToSys(timestr string, dayFlag bool) (int, bool) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(error)
			if ok {
				// 日志记录
				for i := 2; i <= 8; i++ {
					_, f, line, ok := runtime.Caller(i)
					if !ok {
						continue
					}
					if i == 2 {
						logs.Error(i, "__err:[", err, "]__fname:[", f, "]__line:[", line, "]")
					} else {
						logs.Error(i, "__fname:[", f, "]__line:[", line, "]")
					}
				}
			}
		}
	}()
	flag := false
	_at := 0
	if timestr != "" {
		timestr = fmt.Sprint(time.Now().Format("2006-01-02"), " ", timestr)
		//		fmt.Println(timestr)
		t := utils.GetTimeFromStr(timestr)
		if dayFlag {
			//日期加一天
			d, _ := time.ParseDuration("24h")
			t = t.Add(d)
			//			logs.Debug("timer_", t.Day())
		}
		flag = true
		if time.Now().Unix() >= t.Unix() {
			_at = 1
		} else {
			_at = -1
		}

	}
	//	fmt.Println("_at_", _at, "flag", flag)
	return _at, flag
}
