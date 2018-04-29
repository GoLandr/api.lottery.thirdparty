package webhandler

import (
	"log"
	"mycommon/logs"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"
	"strings"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/sssutils"

	"api.lottery.thirdparty/config"
)

type Statistics struct{}

func (this *Statistics) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "oninfo":
		return this.onInfo(param)

	case "statistics":
		return this.statistics(param)

	case "topagent":
		fallthrough
	case "topAgent":
		return this.topAgent(param)

	case "userInfo":
		return this.userInfo(param)
	}
	log.Println("test__")
	return "找不到方法：" + do, 999
}

func (this *Statistics) onInfo(vals map[string]interface{}) (interface{}, int) {

	shFileName := "oninfo.sh"
	shDir := config.GeShPath()
	resbt, err := sssutils.ExeShell(shDir + shFileName)
	utils.ThrowError(err)

	resstr := string(resbt)
	resstr = strings.TrimSpace(resstr)
	onlineCount := mathstr.Math2intDefault0(resstr)
	onlineCount--

	sqlstr := mathstr.S_SFT(`
		SELECT COUNT(1) playing
		FROM users
		WHERE IFNULL(curroomid,0)<>0
		AND IFNULL(curtableid,0)<>0
		AND IFNULL(curseatid,0)<>0
		`)
	orm := global.GetOrm()
	play, err := orm.QueryOne(sqlstr)
	utils.ThrowError(err)
	playingCount := mathstr.Math2intDefault0(play["playing"])
	idleCount := onlineCount - playingCount

	res := map[string]interface{}{}
	res["online"] = onlineCount
	res["playing"] = playingCount
	res["idle"] = idleCount

	return res, -1
}

func (this *Statistics) statistics(vals map[string]interface{}) (interface{}, int) {
	startDate := fmt.Sprint(vals["start"]) // 2016-01-02
	endDate := fmt.Sprint(vals["end"])     // 2016-01-02
	utils.S_CRPM(vals, "start")

	if "" == endDate {
		endDate = utils.NowDate()
	}
	if startDate > endDate {
		return "起始时间必须小于结束时间", 20
	}

	// 数据格式检测

	wherestr1 := `1=1`
	wherestr2 := `1=1`
	if "" != startDate {
		startTimeObj := utils.GetTimeFromStr(startDate + " 00:00:00")
		start := startTimeObj.Unix()
		wherestr1 += mathstr.S_SFT(` and registertime>={0} `, start)
		wherestr2 += mathstr.S_SFT(` and createtime>={0} `, start)
	}
	if "" != endDate {
		endTimeObj := utils.GetTimeFromStr(endDate + " 23:59:59")
		end := endTimeObj.Unix()
		wherestr1 += mathstr.S_SFT(` and registertime<={0} `, end)
		wherestr2 += mathstr.S_SFT(` and createtime<={0} `, end)
	}

	sqlstr := mathstr.S_SFT(`
		select * from
		(
			SELECT COUNT(1) register FROM users WHERE {0}
		) a,
		(
			SELECT COUNT(1) opened FROM sss_room_records WHERE {1}
		) b,
		(
			select count(1) tables from
			(
				select count(1) cnt,tableid,createtime,FROM_UNIXTIME(createtime,'%Y-%m-%d %h:%i:%s') cdate
				from t_accounts_record where {1}
				group by createtime,tableid
			) xx
		)c
		`, wherestr1, wherestr2)
	logs.Debug("___str:", sqlstr)
	orm := global.GetOrm()
	data1, err := orm.QueryOne(sqlstr)
	utils.ThrowError(err)

	res := map[string]interface{}{}
	res["register"] = mathstr.Math2intDefault0(data1["register"])
	res["opened"] = mathstr.Math2intDefault0(data1["opened"])
	res["tables"] = mathstr.Math2intDefault0(data1["tables"])
	res["animated"] = this.GetAnimatedCount(startDate, endDate)

	return res, -1
}
func (this *Statistics) GetAnimatedCount(startDate string, endDate string) int {

	wherestr := `1=1`
	if "" != startDate {
		wherestr += mathstr.S_SFT(` and c_date>='{0}' `, startDate+" 00:00:00")
	}
	if "" != endDate {
		wherestr += mathstr.S_SFT(` and c_date<='{0}' `, endDate+" 23:59:59")
	}

	// 1 获取期间所有的记录
	sqlrecord := mathstr.S_SFT(`
		select *,
			to_days(c_date) days
		from t_userplay_record
		where {0}
		order by c_date asc
		`, wherestr)
	orm := global.GetOrm()
	logs.Debug("__sqlrecord:", sqlrecord)
	recordLst, err := orm.QueryMap(sqlrecord)
	utils.ThrowError(err)

	if len(recordLst) <= 0 {
		return 0
	}

	// 1.1 加工数据，获取每个人对应的天数
	userIdDaysMap := this.GetUserIdDaysMap(recordLst)
	targetDays := this.GetTargetDays(startDate, endDate)

	res := 0
	for _, e := range userIdDaysMap {
		if e == targetDays {
			res++
		}
	}

	return res
}
func (this *Statistics) GetTargetDays(startDate string, endDate string) string {
	startTimeObj := utils.GetDateFromStr(startDate)
	endTimeObj := utils.GetDateFromStr(endDate)

	startNum := startTimeObj.Unix() / (1 * 24 * 60 * 60)
	endNum := endTimeObj.Unix() / (1 * 24 * 60 * 60)

	res := ""
	for i := startNum; i <= endNum; i++ {
		res += fmt.Sprint(i+719529, ",")
	}

	return res
}
func (this *Statistics) GetUserIdDaysMap(recordLst []map[string]interface{}) map[int]string {

	userIdDaysMap := map[int]string{}
	userIdDaylstMap := map[int][]string{}

	for _, e := range recordLst {
		tuid := mathstr.Math2intDefault0(e["userid"])
		tdays := fmt.Sprint(e["days"])
		userIdDaylstMap[tuid] = append(userIdDaylstMap[tuid], tdays)
	}

	for uid, daysLst := range userIdDaylstMap {
		texistsMap := map[string]bool{}

		daystr := ""
		for _, e := range daysLst {
			if texistsMap[e] {
				continue
			}
			daystr += e + ","
			texistsMap[e] = true
		}
		userIdDaysMap[uid] = daystr
	}

	return userIdDaysMap
}
func (this *Statistics) GetUserIdFromStage(tstage string) []int {
	strlst := strings.Split(tstage, ",")
	var idLst []int
	for _, te := range strlst {
		e := strings.TrimSpace(te)
		if "" == e {
			continue
		}
		elst := strings.Split(e, "_")
		elen := len(elst)
		if elen <= 3 {
			continue
		}
		tid := elst[elen-3] // 取倒数第三个

		idLst = append(idLst, mathstr.Math2intDefault0(tid))
	}
	return idLst
}
func (this *Statistics) addUseridDaylstMap(userIdDaylstMap map[int][]string, userids []int, days string) {
	for _, id := range userids {
		userIdDaylstMap[id] = append(userIdDaylstMap[id], days)
	}
}

func (this *Statistics) topAgent(vals map[string]interface{}) (interface{}, int) {
	startDate := fmt.Sprint(vals["start"]) // 2016-01-02
	endDate := fmt.Sprint(vals["end"])     // 2016-01-02
	num := fmt.Sprint(vals["num"])         // 2016-01-02
	utils.S_CRPM(vals, "start")

	if "" == endDate {
		endDate = utils.NowDate()
	}
	if "" == num {
		num = "50"
	}
	if startDate > endDate {
		return "起始时间必须小于结束时间", 20
	}

	wherestr := `1=1`
	if "" != startDate {
		wherestr += mathstr.S_SFT(` and rec.c_date>='{0}' `, startDate+" 00:00:00")
	}
	if "" != endDate {
		wherestr += mathstr.S_SFT(` and rec.c_date<='{0}' `, endDate+" 23:59:59")
	}

	sqlstr := mathstr.S_SFT(`
		select count(1) round_cnt,ifnull(users.agent_id,0) agent_id
		from t_userplay_record rec
		left join users on rec.userid=users.id
		left join t_accounts_record ar on rec.accounts_id=ar.id
		where ifnull(ar.owner_id,0)<>rec.userid
		and {0}
		group by ifnull(users.agent_id,0)
		order by count(1) desc
		limit 0,{1}
		`, wherestr, num)
	orm := global.GetOrm()
	reslst, err := orm.QueryMap(sqlstr)
	utils.ThrowError(err)

	mathstr.RoundMapLst4(reslst, "round_cnt", "agent_id")

	return reslst, -1
}

func (this *Statistics) userInfo(vals map[string]interface{}) (interface{}, int) {
	uid := mathstr.Math2intDefault0(vals["uid"]) // 2016-01-02
	startDate := fmt.Sprint(vals["start"])       // 2016-01-02
	endDate := fmt.Sprint(vals["end"])           // 2016-01-02
	utils.S_CRPM(vals, "uid")

	if "" == endDate {
		endDate = utils.NowDate()
	}
	if "" == startDate {
		startDate = utils.NowDate()
	}
	if startDate > endDate {
		return "起始时间必须小于结束时间", 20
	}

	wherestr := mathstr.S_SFT(`userid={0}`, uid)
	if "" != startDate {
		wherestr += mathstr.S_SFT(` and rec.c_date>='{0}' `, startDate+" 00:00:00")
	}
	if "" != endDate {
		wherestr += mathstr.S_SFT(` and rec.c_date<='{0}' `, endDate+" 23:59:59")
	}

	sqlstr := mathstr.S_SFT(`
		SELECT ur.score,ar.createtime,ar.tableid,ar.c_date
		FROM t_userplay_record  ur
		LEFT JOIN t_accounts_record ar ON ur.accounts_id=ar.id
		WHERE {0}
		ORDER BY ar.createtime DESC
		`, wherestr)
	orm := global.GetOrm()
	reslst, err := orm.QueryMap(sqlstr)
	utils.ThrowError(err)

	totalRound := len(reslst)
	totalScore := 0
	for _, e := range reslst {
		totalScore += mathstr.Math2intDefault0(e["score"])
	}
	resmap := map[string]interface{}{}
	resmap["total_score"] = totalScore
	resmap["total_round"] = totalRound
	resmap["uid"] = uid
	resmap["start"] = startDate
	resmap["end"] = endDate
	resmap["vals"] = reslst

	return resmap, -1
}
