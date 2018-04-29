package webhandler

import (
	"mycommon/encode"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
)

type RankPrize struct{}

func (this *RankPrize) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "insert":
		return this.Insert(param)
	case "update":
		return this.update(param)
	case "del":
		return this.del(param)
	case "select":
		return this.selectWhere(param)
	case "selectOne":
		return this.selectOne(param)

	case "setPrize":
		return this.setPrize(param)
	case "selectSet":
		return this.selectSet(param)

	}

	return "no func name:" + do, 999
}

func (this *RankPrize) Insert(vals map[string]interface{}) (interface{}, int) {

	utils.S_CRPM(vals, "prize_code", "prize_name", "kind", "count")

	orm := global.GetOrm()

	nmp := map[string]interface{}{}
	nmp["id"] = encode.UUID()
	mathstr.S_CME(vals, nmp, "prize_code", "prize_name", "kind", "count")

	_, err := orm.SetTable("z_rank_prize").Insert(nmp)
	utils.ThrowError(err)

	return nil, -1
}

func (this *RankPrize) update(vals map[string]interface{}) (interface{}, int) {

	utils.S_CRPM(vals, "id")

	orm := global.GetOrm()

	nmp := map[string]interface{}{}
	mathstr.S_CME(vals, nmp, "id", "prize_code", "prize_name", "kind", "count")

	_, err := orm.SetTable("z_rank_prize").SetPK("id").UpdateOnly(nmp)
	utils.ThrowError(err)

	return nil, -1
}

func (this *RankPrize) del(vals map[string]interface{}) (interface{}, int) {
	id := fmt.Sprint(vals["id"])
	utils.S_CRPM(vals, "id")

	orm := global.GetOrm()
	_, err := orm.SetTable("z_rank_prize").Where("id=?", id).DeleteRow()
	utils.ThrowError(err)

	return nil, -1
}

func (this *RankPrize) selectWhere(vals map[string]interface{}) (interface{}, int) {
	name := fmt.Sprint(vals["name"]) // 模糊
	kind := fmt.Sprint(vals["kind"]) // 精确

	wherestr := mathstr.S_SFT(`1=1`)
	if "" != name {
		wherestr += mathstr.S_SFT(` and (prize_name like '%{0}%' or prize_code like '%{0}%') `, name)
	}
	if "" != kind {
		wherestr += mathstr.S_SFT(` and kind={0} `, kind)
	}

	orm := global.GetOrm()

	lst, err := orm.SetTable("z_rank_prize").Where(wherestr).FindStringInterface()
	utils.ThrowError(err)

	return lst, -1
}

func (this *RankPrize) selectOne(vals map[string]interface{}) (interface{}, int) {
	id := fmt.Sprint(vals["id"])

	orm := global.GetOrm()

	one, err := orm.SetTable("z_rank_prize").Where("id=?", id).FindOneMap()
	utils.ThrowError(err)

	return one, -1
}

func (this *RankPrize) setPrize(vals map[string]interface{}) (interface{}, int) {
	rankSort := mathstr.Math2intDefault0(vals["rank_sort"])
	prizeCode := fmt.Sprint(vals["prize_code"])

	sqlstr := mathstr.S_SFT(`
		insert into z_rank_set(id,rank_sort,prize_code,c_date)
		values('{0}','{1}','{2}',now())
		ON DUPLICATE KEY UPDATE prize_code='{2}';
		`, encode.UUID(), rankSort, prizeCode)

	orm := global.GetOrm()
	_, err := orm.Exec(sqlstr)
	utils.ThrowError(err)

	sqlUpdate := `
		update z_rank_set 
		left join z_rank_prize p on z_rank_set.prize_code=p.prize_code
		set z_rank_set.prize_name=p.prize_name
		where ifnull(z_rank_set.prize_name,'') != p.prize_name
		`
	_, err = orm.Exec(sqlUpdate)
	utils.ThrowError(err)

	return nil, -1
}

func (this *RankPrize) selectSet(vals map[string]interface{}) (interface{}, int) {

	orm := global.GetOrm()
	lst, err := orm.SetTable("z_rank_set").OrderBy("rank_sort asc").FindStringInterface()
	utils.ThrowError(err)

	return lst, -1
}
