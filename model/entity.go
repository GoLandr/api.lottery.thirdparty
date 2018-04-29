package model

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
