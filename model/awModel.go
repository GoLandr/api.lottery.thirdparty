package model

// 爱玩model

// "id":"会员ID","openid":"OPENID","nickName":"昵称","realName":实名,"mobile":手机号码,"iconUrl":"头像地址","money":"游卡数量","referralCode":推荐码
// 爱玩用户

type AWFlag struct {
	Flag bool        `json:"flag"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
