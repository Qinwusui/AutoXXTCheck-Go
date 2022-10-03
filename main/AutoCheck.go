package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println(`======欢迎使用自动打卡程序======
1.输入个人信息
2.进行打卡
=========使用前须知============
|| 本程序完全开源 请勿二改倒卖      
|| 程序不会发送您的隐私信息到任何地方 
|| 所有个人信息均存放在当前目录下
|| 不放心该程序，请删除并不再使用
==============================
|| 打卡时间默认为上午8点(晨检)与下午3点(午检)
|| 可通过修改程序目录下的config.json文件更改默认打卡时间
|| 该程序需要持续运行，若中途因不可知因素导致程序进程结束，请重新运行该程序
|| 建议放在24小时持续运行的服务器上
==============================
注意:暂不支持教师打卡`)
	input()
}

func input() {
	fmt.Print("请输入序号：")
	var number int
	fmt.Scanln(&number)
	switch number {
	case 1:
		collectInfo()

	case 2:
		startCheck()

	case 0:
		os.Exit(0)

	default:
		fmt.Println("请输入`1`或`2` 输入`0`退出程序")
	}

}

type UserInfo struct {
	UserName string `json:"userName"`
	School   string `json:"school"`
	Major    string `json:"major"`
	Cls      string `json:"class"`
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Pwd      string `json:"pwd"`
}
type Config struct {
	MorningTime   string `json:"morning"`
	AfternoonTime string `json:"afternoon"`
}
type LoginRes struct {
	Mes    string `json:"mes"`
	Type   int    `json:"type"`
	Url    string `json:"url"`
	Status bool   `json:"status"`
}

func init() {
	isExist, _ := PathExists("./config.json")
	if !isExist {
		config := Config{
			MorningTime:   "08:00:00",
			AfternoonTime: "15:00:00",
		}
		bytes, _ := json.Marshal(&config)
		os.WriteFile("./config.json", bytes, 0777)
	}
	bytes, _ := os.ReadFile("./config.json")
	_ = json.Unmarshal(bytes, &config)
}
func collectInfo() {

	for {
		userInfo := new(UserInfo)
		fmt.Print("请输入姓名：")
		fmt.Scanln(&userInfo.UserName)

		fmt.Print("请输入学院名：")
		fmt.Scanln(&userInfo.School)

		fmt.Print("请输入专业名：")
		fmt.Scanln(&userInfo.Major)

		fmt.Print("请输入班级名：")
		fmt.Scanln(&userInfo.Cls)

		fmt.Print("请输入身份证号：")
		fmt.Scanln(&userInfo.ID)

		fmt.Print("请输入学习通登录手机号：")
		fmt.Scanln(&userInfo.Phone)

		fmt.Print("请输入学习通密码：")
		fmt.Scanln(&userInfo.Pwd)
		saveUser(*userInfo)
		var needContinue string
		fmt.Print("是否继续添加?(default: n/y)")
		fmt.Scanln(&needContinue)
		if needContinue == "y" {
			continue
		} else {
			break
		}
	}
	startCheck()

}
func getTime(time string) (string, string, string) {
	var second, min, hour string
	strList := strings.Split(time, ":")
	second = strList[2]
	min = strList[1]
	hour = strList[0]
	return second, min, hour
}
func modifyConfig() {
	fmt.Print("是否更改默认打卡时间 (default n/y):")
	var choose string
	fmt.Scanln(&choose)
	if choose == "y" {
		for {
			fmt.Println("请按样例设置晨检时间(样例: 08:00:00 表示早上八点整打卡  15:00:00 表示下午三点整打卡)")
			fmt.Print("请输入晨检时间:")
			var checkM, checkA string
			fmt.Scanln(&checkM)
			strList := strings.Split(checkM, ":")
			if len(strList) != 3 {
				fmt.Println("输入时间格式错误！")
				continue
			}
			fmt.Print("请输入午检时间:")
			fmt.Scanln(&checkA)
			strList = strings.Split(checkA, ":")
			if len(strList) != 3 {
				fmt.Println("输入时间格式错误！")
				continue
			}
			config = Config{
				MorningTime:   checkM,
				AfternoonTime: checkA,
			}
			bytes, _ := json.Marshal(config)
			os.WriteFile("./config.json", bytes, 0777)
			break
		}
		fmt.Println("打卡时间更新完成...")
	}
	fmt.Println("当前设定打卡时间为:", config.MorningTime+"与"+config.AfternoonTime)
}

// 开始打卡
func startCheck() {
	modifyConfig()
	fmt.Println("即将开始打卡...")
	userList := readCheckList()
	if len(userList) == 0 {
		fmt.Println("列表无待打卡成员,已跳转添加成员")
		collectInfo()
	}
	cronTab := cron.New(cron.WithSeconds())
	seM, minM, hourM := getTime(config.MorningTime)
	seA, minA, hourA := getTime(config.AfternoonTime)

	cronTab.AddFunc(seM+" "+minM+" "+hourM+" * * ?", func() {
		for i := range userList {
			check(userList[i], Morning)
		}
	})
	cronTab.AddFunc(seA+" "+minA+" "+hourA+" * * ?", func() {
		for i := range userList {
			check(userList[i], Afternoon)
		}
	})
	cronTab.Start()
	defer cronTab.Stop()
	select {}
}

var config Config
var client = &http.Client{
	Timeout: 10 * time.Second,
}

const (
	Morning int = iota
	Afternoon
)

func getFormId(num int) string {
	switch num {
	case Morning:
		return "89398"
	case Afternoon:
		return "305014"
	default:
		return ""
	}
}
func getEnc(num int) string {
	switch num {
	case Morning:
		return "75244f75384287c902e57b080c4d1c6d"
	case Afternoon:
		return "e7fc8869547dc0b3922c6453e65b509f"
	default:
		return ""
	}
}

func check(userInfo UserInfo, num int) {
	cookies := login(userInfo)
	checkCode, cookies := getCheckCode(num, cookies)
	data := "formId=" + getFormId(num) + "&formAppId=&version=3&formData=" +
		getStuFormData(userInfo, cookies) +
		"&ext=&t=1&enc=" + getEnc(num) + "&checkCode=" + checkCode + "&gatherId=0&anonymous=0&uuid=&uniqueCondition=%5B%5D&gverify="

	req, e := http.NewRequest("POST", "https://office.chaoxing.com/data/apps/forms/fore/user/save", bytes.NewBufferString(data))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 12;) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.41 Mobile Safari/537.36 Language/zh_CN com.chaoxing.mobile/ChaoXingStudy_3_5.1.4_android_phone_614_74 (@Kalimdor)_482bfb22af77461b96e77e64aa40abc2")
	req.Header.Add("Host", "office.chaoxing.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Origin", "https://office.chaoxing.com")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
	if e != nil {
		return
	}

	resp, e := client.Do(req)
	if e != nil {
		log.Fatalln(e.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	defer resp.Body.Close()

}
func getCheckCode(num int, cookies []*http.Cookie) (string, []*http.Cookie) {

	url := "https://office.chaoxing.com/front/web/apps/forms/fore/apply?"
	req, _ := http.NewRequest("GET", url+`id=`+getFormId(num)+`&enc=`+getEnc(num), nil)
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
	resp, e := client.Do(req)
	if e != nil {
		log.Fatalln(e.Error())
		return "",nil
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	respCookies := resp.Cookies()
	for range respCookies {
		cookies = append(cookies, respCookies...)
	}
	resStr := string(bytes)
	indexStr := "checkCode = '"
	index := strings.Index(resStr, indexStr)
	lastIndex := strings.Index(resStr[index:], "',")

	return resStr[index+len(indexStr) : index+lastIndex], cookies

}
func getUid(cookies []*http.Cookie) string {
	for i := range cookies {
		if cookies[i].Name == "_uid" {
			return cookies[i].Value
		}
	}
	fmt.Println("未获取到UID!")
	return ""
}
func getStuFormData(userInfo UserInfo, cookies []*http.Cookie) string {

	originText, _ := url.QueryUnescape(data)
	stuFormData := FormData{}
	_ = json.Unmarshal([]byte(originText), &stuFormData)
	Puid, _ := strconv.Atoi(getUid(cookies))
	stuFormData[0].Fields[0].Values[0].Puid = Puid                                                                            //UID
	stuFormData[0].Fields[0].Values[0].Uname = userInfo.UserName                                                              //姓名
	stuFormData[1].Fields[0].Values[0].Val = userInfo.School                                                                  //学院
	stuFormData[2].Fields[0].Values[0].Val = userInfo.Major                                                                   //专业
	stuFormData[3].Fields[0].Values[0].Val = userInfo.Cls                                                                     //班级
	stuFormData[4].Fields[0].Values[0].Val = userInfo.ID                                                                      //身份证号码
	stuFormData[5].Fields[0].Values[0].Val = userInfo.Phone                                                                   //电话号码
	stuFormData[6].Fields[0].Values[0].Val = fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()) //日期
	// stuFormData[7].Fields[0].Values[0].Val                                                                                    //省
	// stuFormData[7].Fields[0].Values[1].Val                                                                                    //市
	// stuFormData[8].Fields[0].Values[0].Val                                                                                    //湖北省武汉市江夏区江夏区经济开发区庙山街道武昌理工学院人工智能学院
	stuFormData[14].Fields[0].Values[0].Val =
		fmt.Sprintf("%d-%d-%d+%d:%d",
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			time.Now().Hour(),
			time.Now().Minute())
	marshalBytes, _ := json.Marshal(&stuFormData)
	return url.PathEscape(string(marshalBytes))
}

// 获取登录Cookie
func login(userInfo UserInfo) []*http.Cookie {
	data := &url.Values{}
	data.Set("uname", userInfo.Phone)
	data.Set("code", userInfo.Pwd)
	data.Set("loginType", "1")
	data.Set("roleSelect", "true")
	req, _ := http.NewRequest("POST",
		"https://passport2-api.chaoxing.com/v11/loginregister?cx_xxt_passport=json",
		bytes.NewBufferString(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	resp, e := client.Do(req)
	if e != nil {
		log.Fatalln(e.Error())
	}
	defer resp.Body.Close()

	return resp.Cookies()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func readCheckList() []UserInfo {
	var userInfo []UserInfo
	isExist, _ := PathExists("./checkList.json")
	if !isExist {
		//生成文件
		bytes, _ := json.Marshal([]UserInfo{})
		os.WriteFile("./checkList.json", bytes, 0777)
	}
	bytes, _ := os.ReadFile("./checkList.json")

	e := json.Unmarshal(bytes, &userInfo)
	if e != nil {
		log.Fatalln(e.Error())
	}
	return userInfo
}
func saveUser(userInfo UserInfo) {
	userList := readCheckList()
	for i := range userList {
		if userList[i].Phone == userInfo.Phone {
			fmt.Println("该手机号已存在！")
			return
		}
	}
	userList = append(userList, userInfo)
	bytes, _ := json.Marshal(&userList)
	os.WriteFile("./checkList.json", bytes, 0777)
}
