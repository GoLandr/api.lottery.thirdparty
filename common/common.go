package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"code.google.com/p/mahonia"
)

//post 請求
func Httppost(urlstr string, data url.Values) (string, error) {
	resp, err := http.PostForm(urlstr,
		data)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	//	fmt.Println(string(body))
	return string(body), err
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
func GbkToUtf8(src string) string {
	return ConvertToString(src, "gbk", "utf-8")
}
