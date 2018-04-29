package spider

import (
	"testing"
)

func Test_main(t *testing.T) {
	//	lottery := new(Lottery)
	//	lottery.SpiderUrl("http://www.j0024.com/lottery/getAllCqsscAutoList")
	//	httpPost("http://www.j0024.com/lottery/getAllCqsscAutoList")
	//	lottery.SpiderUrl("http://www.booktxt.net/2_2096/")
	pj := new(PuJing)
	pj.IsFirst = true
	pj.Cqssc()
}
