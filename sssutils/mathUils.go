package sssutils

import (
	"regexp"
	"strings"

	"api.lottery.thirdparty/config"

	"mycommon/encode"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"

	"mycommon/utils"
)

// 第一个返回值： 前面的内容，第二个：后面的内容，第三个 ：中间的内容
func GetUrlFromContentByRegix(content string) []string {
	reg := regexp.MustCompile(`(h|H)(r|R)(e|E)(f|F) *= *('|")?((\w|\\|\/|\.|:|-|_|\s|\?|\&|=)+)('|"| *|>)?`)
	return reg.FindAllString(content, -1)
}
func GetUrlFromContent(content string) (fStr string, eStr string, obj string) {
	var furl string
	if strings.Contains(content, "<a ") {
		sdx := strings.Index(content, "<a ")
		edex := strings.Index(content, "</a>")
		fmt.Println(sdx, ",", edex)
		furl = content[sdx:edex]
		if strings.Contains(content, "href") {
			sdx := strings.Index(content, "href")
			furl = content[sdx:]
		}

		// 获取第一个“与第二个”之间的内容
		fdex := strings.Index(furl, "\"")
		bIsS := false
		if fdex < 0 {
			fdex = strings.Index(furl, "'")
			bIsS = true
		}
		furl = furl[fdex+1:]
		var sdex int
		if !bIsS {
			sdex = strings.Index(furl, "\"")
		} else {
			sdex = strings.Index(furl, "'")
		}

		furl = furl[0:sdex]

		ffdex := strings.Index(content, furl)
		eedex := ffdex + len(furl)

		return content[:ffdex], content[eedex:], furl
	}
	return "", "", ""
}

func GenerateToken() string {
	randStr := mathstr.GetRandNum(6) + "_" + utils.NowStr()
	return encode.Base64Encode([]byte(randStr))
}
func GenerageWXZFSign(parammap map[string]interface{}, keys []string, authkey string) string {

	sign := ``
	for _, key := range keys {
		v := fmt.Sprint(parammap[key])
		if `` == v {
			continue
		}
		if `` != sign {
			sign = sign + `&`
		}

		sign = fmt.Sprint(sign, key, "=", v)
	}

	sign = fmt.Sprint(sign, `&key=`, config.GetWxzfSignkey())
	finalSign := encode.GetMd5String(sign)
	finalSign = strings.ToUpper(finalSign)

	return finalSign

}
