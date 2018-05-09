package utils

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
		tiger += 1
		dragon = 0
		draw = 0
	} else if first_ball < second_ball {
		dragon += 1
		tiger = 0
		draw = 0
	} else {
		draw += 1
		tiger = 0
		dragon = 0
	}
	return dragon, tiger, draw
}
