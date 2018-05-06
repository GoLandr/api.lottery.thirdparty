package base

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func Test_main(t *testing.T) {
	//	s := "tencent://message/?uin=97106428&Site=在线客服&Menu=yes"
	//	//解析这个 URL 并确保解析没有出错。
	//	_, err := url.Parse(s)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("11")
	//	c := exec.Command("ping", "127.0.0.1")
	c := exec.Command("cmd", "rundll32 url.dll,FileProtocolHandler  tencent://message/?uin=97106428&Site=在线客服&Menu=yes")
	//	c.Start()
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
