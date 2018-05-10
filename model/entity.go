package model

type Menber struct {
	Id    int    `sql:"id" json:"id"`
	Phone int    `sql:"phone" json:"phone"`
	Name  string `sql:"name" json:"name"`
	Qq    int    `sql:"qq" json:"qq"`
	Code  int    `sql:"code" json:"code"`
}
type TLimit struct {
	Id              int    `sql:"id" json:"id"`
	Menber_id       int    `sql:"menber_id" json:"menber_id"`
	Mode            int    `sql:"mode" json:"mode"`
	Nickname        string `sql:"nickname" json:"nickname"`
	Big_small_limit int    `sql:"big_small_limit" json:"big_small_limit"`
	Odd_even_limit  int    `sql:"odd_even_limit" json:"odd_even_limit"`
	Star_limit      int    `sql:"star_limit" json:"star_limit"`
	Total_bs_limit  int    `sql:"total_bs_limit" json:"total_bs_limit"`
	Total_oe_limit  int    `sql:"total_oe_limit" json:"total_oe_limit"`
	Pred_limit      int    `sql:"pred_limit" json:"pred_limit"`
	Is_valid        int    `sql:"is_valid" json:"is_valid"`
	Valid_date      string `sql:"valid_date" json:"valid_date"`
	Phone           string `sql:"Phone" json:"phone"`
}

type BJPK struct {
	SSC
	Six_ball   int
	Seven_ball int
	Eight_ball int
	Ninth_ball int
	Ten_ball   int
}
type SSC struct {
	Flowid       int `sql:"flowid"`
	One_ball     int
	Two_ball     int
	Third_ball   int
	Four_ball    int
	Five_ball    int
	Lottery_date string
	Lottery_time string
	Periods      string
	Update_date  string
}
