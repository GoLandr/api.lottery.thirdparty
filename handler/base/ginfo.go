package base

var GLotteryMgr *LotteryMgr //全局变量
func InitConfigs() {
	if GLotteryMgr == nil {
		mgr := new(LotteryMgr)
		GLotteryMgr = mgr
	}
}
