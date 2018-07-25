package idiom

import (
	"errors"
	"fmt"
	"strings"
)

// func main() {
// 	var jl idiom.JL = idiom.JL{UserName: "t2"}
// 	var input string
// 	fmt.Println("输入成语")
// 	for {
// 		fmt.Scanln(&input)
// 		if input == "不玩了" {
// 			break
// 		}
// 		msg, err := idiom.Jielong(&jl, input)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Println(msg)
// 	}
// 	fmt.Println("done")
// }

func Jielong(jl *JL, str string) (msg string, err error) {
	var code int
	var count int
	for {
		code, msg, err = jl.Jielong(str)
		if err != nil {
			return
		}
		if code == 2 {
			continue
		}
		if strings.Contains(jl.Used, msg) {
			count += 1
			if count > 4 {
				msg = "不玩了"
				jl.Started = false
				jl.Used = ""
				return
			}
		}
		jl.Used += str
		break
	}
	return
}

type JL struct {
	Started  bool
	Used     string
	UserName string
}

func (j *JL) Jielong(str string) (code int, msg string, err error) {
	if !j.Started {
		str = "成语接龙" + str
	}
	resp, err := Chengyujielong(str, j.UserName)
	if err != nil {
		return
	}
	fmt.Println(resp)
	if resp.Code != 100000 {
		err = errors.New(resp.Text)
		return
	}
	if strings.Contains(resp.Text, "进入成语接龙模式") {
		j.Started = true
		msg = resp.Text[len(resp.Text)-12:]
		return
	} else if strings.Contains(resp.Text, "退出成语接龙模式") {
		j.Started = false
		code = 2
		return
	}
	if len(resp.Text) != 12 {
		j.Started = false
		code = 2
		return
	}
	msg = resp.Text
	return
}
