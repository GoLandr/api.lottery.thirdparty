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
func GetStars() {

}
