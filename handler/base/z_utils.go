package base

import (
	"mycommon/utils"

	"api.lottery.thirdparty/model"
)

func GetPushMenber(limit map[int][]*model.TLimit, maxVal int) []int {
	var retLst []int
	if lst, bok := limit[maxVal]; bok {
		for _, entry := range lst {
			//判断是否会员
			if entry.Is_valid == 1 && CheckValid(entry.Valid_date) {
				if menber, tok := GMenber[entry.Menber_id]; tok {
					retLst = append(retLst, menber.Code)
				}
			}
		}
	}
	return retLst
}
func CheckValid(valid_data string) bool {
	flag := false
	if utils.Timestamp() <= utils.GetTimeFromStr(valid_data).Unix() {
		flag = true
	}
	return flag
}
