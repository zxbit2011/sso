package main

import (
	"github.com/beewit/sso/global"
	//"github.com/beewit/sso/router"

	"time"
)

func main() {

	global.RD.SetAndExpire("times", "1", 1)
	k, _ := global.RD.GetString("times")
	println(k)
	time.Sleep(3 * time.Second)
	k, _ = global.RD.GetString("times")
	println(k)

	//router.Start()

	//e := make(map[string]interface{})
	//e["pwd"] = 123456
	//switch v := e["pwd"].(type) {
	//case int:
	//	fmt.Println("整型", v)
	//	var s int
	//	s = v
	//	fmt.Println(s)
	//case string:
	//	fmt.Println("字符串", v)
	//}
}
