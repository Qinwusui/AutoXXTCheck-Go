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

	"github.com/go-co-op/gocron"
	"gopkg.in/gomail.v2"
)

func main() {
	fmt.Println(`======欢迎使用自动打卡程序======
1.输入个人信息
2.进行打卡
=========使用前须知================
|| 本程序完全开源 请勿二改倒卖      
|| 程序不会发送您的隐私信息到任何地方 
|| 所有个人信息均存放在当前目录下
|| 不放心该程序，请删除并不再使用
=========学生打卡须知==============
|| 打卡时间默认为上午8点(晨检)与下午3点(午检)
|| 可通过修改程序目录下的stuConfig.json文件更改默认打卡时间
|| 该程序需要持续运行，若中途因不可知因素导致程序进程结束，请重新运行该程序
|| 建议放在24小时持续运行的服务器上
==========教师打卡须知=============
|| 打卡一日一次 时间固定为早上九点 亦可自行更改
|| 打卡定位地点为学院内任意地点 不可更改
================================`)
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
		checkList()

	case 0:
		os.Exit(0)

	default:
		fmt.Println("请输入`1`或`2` 输入`0`退出程序")
	}

}
func collectInfo() {
	for {
		fmt.Print("即将输入的用户是学生还是教师？(1:学生 2:教师):")
		var check string
		fmt.Scanln(&check)
		switch check {
		case "1":
			collectStuInfo()
		case "2":
			collectTeaCInfo()
		case "0":
			os.Exit(0)
		default:
			fmt.Println("请输入`1`或`2` 输入`0`退出程序")
		}
		var needContinue string
		fmt.Print("是否继续添加?(default: n/y):")
		fmt.Scanln(&needContinue)
		if needContinue == "y" {
			continue
		} else {
			break
		}
	}
	checkList()
}

func checkList() {
	if len(readStuCheckList()) != 0 {
		startCheckStu()
	} else {
		if len(readTeaCheckList()) != 0 {
			startCheckTea()
		} else {
			fmt.Println("教师和学生待打卡列表均无成员，已跳转至主页面...")
			main()
		}
	}
}

type StuInfo struct {
	UserName string `json:"userName"`
	School   string `json:"school"`
	Major    string `json:"major"`
	Cls      string `json:"class"`
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Pwd      string `json:"pwd"`
	Mail     string `json:"mail"`
}
type TeaCInfo struct {
	UserName   string `json:"userName"`
	DepartMent string `json:"department"`
	JobNumber  string `json:"jobNumber"`
	ID         string `json:"id"`
	Phone      string `json:"phone"`
	Pwd        string `json:"pwd"`
	Mail       string `json:"mail"`
}
type StuTimeConfig struct {
	MorningTime   string `json:"morning"`
	AfternoonTime string `json:"afternoon"`
}
type LoginRes struct {
	Mes    string `json:"mes"`
	Type   int    `json:"type"`
	Url    string `json:"url"`
	Status bool   `json:"status"`
}
type MailConfig struct {
	Address string `json:"address"`
	Pwd     string `json:"pwd"`
	Port    int    `json:"port"`
}
type RespBody struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func init() {
	//判断学生打卡时间配置是否存在

	stuConfigExist, _ := PathExists("./stuConfig.json")
	if !stuConfigExist {
		stuConfig := StuTimeConfig{
			MorningTime:   "08:00:00",
			AfternoonTime: "15:00:00",
		}
		bytes, _ := json.Marshal(&stuConfig)
		os.WriteFile("./stuConfig.json", bytes, 0777)
	}
	configBytes, _ := os.ReadFile("./stuConfig.json")
	_ = json.Unmarshal(configBytes, &stuTimeConfig)

	//判断教师打卡时间配置是否存在
	teaConfigExist, _ := PathExists("./teaConfig.json")
	if !teaConfigExist {
		teaConfig := "09:00:00"
		bytes, _ := json.Marshal(&teaConfig)
		os.WriteFile("./teaConfig.json", bytes, 0777)
	}
	teaConfigBytes, _ := os.ReadFile("./teaConfig.json")
	json.Unmarshal(teaConfigBytes, &teaTimeConfig)
	//检查发件邮箱是否存在
	mailConfigExist, _ := PathExists("./mailConfig.json")
	if !mailConfigExist {
		mailConfig := MailConfig{
			Address: "",
			Pwd:     "",
			Port:    0,
		}
		bytes, _ := json.Marshal(&mailConfig)
		os.WriteFile("./mailConfig.json", bytes, 0777)
	}
}
func collectTeaCInfo() {
	teaCInfo := new(TeaCInfo)
	fmt.Print("请输入姓名：")
	fmt.Scanln(&teaCInfo.UserName)
	fmt.Print("请输入工号：")
	fmt.Scanln(&teaCInfo.JobNumber)
	fmt.Print("请输入身份证号：")
	fmt.Scanln(&teaCInfo.ID)
	fmt.Print("请输入部门（如：人工智能学院）：")
	fmt.Scanln(&teaCInfo.DepartMent)
	fmt.Print("请输入学习通登录手机号：")
	fmt.Scanln(&teaCInfo.Phone)
	fmt.Print("请输入学习通登录密码：")
	fmt.Scanln(&teaCInfo.Pwd)
	fmt.Print("请输入一个邮箱(用于[接收]打卡回执)：")
	fmt.Scanln(&teaCInfo.Mail)
	saveTeaCInfo(*teaCInfo)
}
func collectStuInfo() {

	stuInfo := new(StuInfo)
	fmt.Print("请输入姓名：")
	fmt.Scanln(&stuInfo.UserName)

	fmt.Print("请输入学院名：")
	fmt.Scanln(&stuInfo.School)

	fmt.Print("请输入专业名：")
	fmt.Scanln(&stuInfo.Major)

	fmt.Print("请输入班级名：")
	fmt.Scanln(&stuInfo.Cls)

	fmt.Print("请输入身份证号：")
	fmt.Scanln(&stuInfo.ID)

	fmt.Print("请输入学习通登录手机号：")
	fmt.Scanln(&stuInfo.Phone)

	fmt.Print("请输入学习通密码：")
	fmt.Scanln(&stuInfo.Pwd)

	fmt.Print("请输入一个邮箱(QQ邮箱最佳，用于[接收]打卡成功后的回执消息):")
	fmt.Scanln(&stuInfo.Mail)

	saveStuInfo(*stuInfo)

}

func modifyConfig() {
	fmt.Print("是否更改学生默认打卡时间 (default n/y):")
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
			stuTimeConfig = StuTimeConfig{
				MorningTime:   checkM,
				AfternoonTime: checkA,
			}
			bytes, _ := json.Marshal(stuTimeConfig)
			os.WriteFile("./stuConfig.json", bytes, 0777)
			break
		}
		fmt.Println("打卡时间更新完成...")
	}
	mailConfig := new(MailConfig)
	mailConfigBytes, _ := os.ReadFile("./mailConfig.json")
	_ = json.Unmarshal(mailConfigBytes, &mailConfig)
	if mailConfig.Address == "" || mailConfig.Pwd == "" || mailConfig.Port == 0 {
		addMail()
	}
	fmt.Println("当前设定打卡时间为:", stuTimeConfig.MorningTime+"与"+stuTimeConfig.AfternoonTime)
}

func addMail() {
	var needSendMail string
	fmt.Print(`是否需要添加一个邮箱发送邮件？
=>下面是该邮箱的简介，请[认真]阅读
1.程序将使用您提供的邮箱登录并发送邮件，[若您担心隐私泄露，请输入n并回车]
2.该邮箱仅支持QQ邮箱，Outlook邮箱
3.QQ邮箱需要使用[授权码]来登录，而Outlook邮箱可以直接使用邮箱密码登录
3.1.获取授权码前请开启邮箱内的STMP协议，之后即可获取到授权码
4.授权码可以在网页版邮箱->设置中获取
5.由于登陆邮箱还需要提供端口号，一般而言，以下提供邮箱端口样例[仅使用SMTP协议]
		1).QQ邮箱：587
		2).Outlook：587
=========================
是否需要使用邮箱发件功能？(y/默认n):
`)
	fmt.Scanln(&needSendMail)
	if needSendMail == "y" {
		mailConfig := new(MailConfig)
		fmt.Print("请输入邮箱完整地址：")
		fmt.Scanln(&mailConfig.Address)
		fmt.Print("请输入授权码/密码：")
		fmt.Scanln(&mailConfig.Pwd)
		fmt.Print("请输入端口：")
		fmt.Scanln(&mailConfig.Port)
		fmt.Println("发件邮箱已保存...\n[若配置信息填写错误，请在程序文件目录下更改]")
		bytes, _ := json.Marshal(&mailConfig)
		os.WriteFile("./mailConfig.json", bytes, 0777)

	}

}

// 开始打卡 教师版
func startCheckTea() {
	fmt.Println("即将开始教师打卡...")
	teaList := readTeaCheckList()
	if len(teaList) == 0 {
		fmt.Println("教师列表无待打卡成员,是否需要添加？(y/n):")
		var add string
		fmt.Scanln(&add)
		if add == "y" {
			collectTeaCInfo()
		}
	}
	timezone:= time.FixedZone("UTC",+8*60*60)
	cronTab := gocron.NewScheduler(timezone)
	cronTab.Every(1).Day().At(teaTimeConfig).Do(func() {
		for i := range teaList {
			go checkTea(teaList[i])
		}
	})

	cronTab.StartBlocking()
}

// 开始打卡 学生版
func startCheckStu() {
	modifyConfig()
	fmt.Println("即将开始学生打卡...")
	stuList := readStuCheckList()

	if len(stuList) == 0 {
		fmt.Println("列表无待打卡成员,已跳转添加成员")
		collectStuInfo()
	}
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	cronTab := gocron.NewScheduler(timezone)

	cronTab.Every(1).Day().At(stuTimeConfig.MorningTime).Do(func() {
		for i := range stuList {
			go checkStu(stuList[i], Morning)
		}
	})
	cronTab.Every(1).Day().At(stuTimeConfig.AfternoonTime).Do(func() {
		for i := range stuList {
			go checkStu(stuList[i], Afternoon)
		}
	})
	cronTab.StartBlocking()

}

var teaTimeConfig string
var stuTimeConfig StuTimeConfig
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
func checkTea(teaInfo TeaCInfo) {
	cookies := teaLogin(teaInfo)
	checkCode, cookies := getCheckCode(88444, cookies) // 表单固定为88444
	data := "formId=88444&formAppId=&version=10&formData=" +
		getTeaFormData(teaInfo, cookies) + "&ext=&t=1&enc=59f8da2a3a871719372dfec0e33fd33f&checkCode=" + checkCode +
		"&gatherId=0&anonymous=0&uuid=" + getUid(cookies) +
		"&uniqueCondition=%5B%5D&gverify="
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
	respBody := new(RespBody)
	_ = json.Unmarshal(bytes, &respBody)

	sendTeaMail(teaInfo, respBody.Msg)

	defer resp.Body.Close()
}
func checkStu(userInfo StuInfo, num int) {
	cookies := stuLogin(userInfo)
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
	respBody := new(RespBody)
	_ = json.Unmarshal(bytes, &respBody)

	sendStuMail(userInfo, respBody.Msg)

	defer resp.Body.Close()

}
func sendTeaMail(teaInfo TeaCInfo, msg string) {
	senderConfig, _ := os.ReadFile("./mailConfig.json") //TODO 检查mailConfig完整性
	sender := new(MailConfig)
	_ = json.Unmarshal(senderConfig, &sender)
	m := gomail.NewMessage()
	m.SetHeader("From", sender.Address)
	m.SetHeader("To", teaInfo.Mail)
	m.SetHeader("Subject", "健康打卡回执")
	m.SetBody("text/html", `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>
<body>
<h2>`+msg+`<h2>
</body>
</html>
`)
	if strings.Contains(sender.Address, "qq.com") {
		d := gomail.NewDialer("smtp.qq.com", sender.Port, sender.Address, sender.Pwd)
		if e := d.DialAndSend(m); e != nil {
			log.Fatalln(e.Error())
		}
		return
	}

	if strings.Contains(sender.Address, "outlook.com") {
		d := gomail.NewDialer("smtp.office365.com", sender.Port, sender.Address, sender.Pwd)
		if e := d.DialAndSend(m); e != nil {
			log.Fatalln(e.Error())
		}
		return
	}
}
func sendStuMail(userInfo StuInfo, msg string) {
	senderConfig, _ := os.ReadFile("./mailConfig.json")
	sender := new(MailConfig)
	_ = json.Unmarshal(senderConfig, &sender)
	m := gomail.NewMessage()
	m.SetHeader("From", sender.Address)
	m.SetHeader("To", userInfo.Mail)
	m.SetHeader("Subject", "健康打卡回执")
	m.SetBody("text/html", `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>
<body>
<h2>`+msg+`<h2>
</body>
</html>
`)
	if strings.Contains(sender.Address, "qq.com") {
		d := gomail.NewDialer("smtp.qq.com", sender.Port, sender.Address, sender.Pwd)
		if e := d.DialAndSend(m); e != nil {
			log.Fatalln(e.Error())
		}
		return
	}

	if strings.Contains(sender.Address, "outlook.com") {
		d := gomail.NewDialer("smtp.office365.com", sender.Port, sender.Address, sender.Pwd)
		if e := d.DialAndSend(m); e != nil {
			log.Fatalln(e.Error())
		}
		return
	}
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
		return "", nil
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

func getTeaFormData(teaInfo TeaCInfo, cookies []*http.Cookie) string {
	originText, _ := url.QueryUnescape(teaForm)
	teaFormData := TeaFormData{}
	json.Unmarshal([]byte(originText), &teaFormData)
	Puid, _ := strconv.Atoi(getUid(cookies))
	teaFormData[2].Fields[0].Values[0].Puid = Puid              //PUID
	teaFormData[2].Fields[0].Values[0].Uname = teaInfo.UserName //姓名
	teaFormData[4].Fields[0].Values[0].Val = teaInfo.JobNumber  // 工号
	teaFormData[7].Fields[0].Values[0].Val = teaInfo.ID         // 身份证号
	teaFormData[6].Fields[0].Values[0].Val = teaInfo.Phone      // 手机号
	// teaFormData[10].Fields[0].Values[0].Address = "武昌理工学院" // 位置
	// teaFormData[10].Fields[0].Values[0].Lng = 0.0          // 位置
	// teaFormData[10].Fields[0].Values[0].Lat = 0.0          // 位置

	teaFormData[17].Fields[0].Values[0].Val = fmt.Sprintf("%d-%d-%d %d:%d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute()) // 手机号
	maStr, _ := json.Marshal(&teaFormData)
	return url.PathEscape(string(maStr))
}
func getStuFormData(userInfo StuInfo, cookies []*http.Cookie) string {

	originText, _ := url.QueryUnescape(stuForm)
	stuFormData := StuFormData{}
	_ = json.Unmarshal([]byte(originText), &stuFormData)
	Puid, _ := strconv.Atoi(getUid(cookies))
	stuFormData[0].Fields[0].Values[0].Puid = Puid               //UID
	stuFormData[0].Fields[0].Values[0].Uname = userInfo.UserName //姓名
	stuFormData[1].Fields[0].Values[0].Val = userInfo.School     //学院
	stuFormData[2].Fields[0].Values[0].Val = userInfo.Major      //专业
	stuFormData[3].Fields[0].Values[0].Val = userInfo.Cls        //班级
	stuFormData[4].Fields[0].Values[0].Val = userInfo.ID         //身份证号码
	stuFormData[5].Fields[0].Values[0].Val = userInfo.Phone      //电话号码
	stuFormData[6].Fields[0].Values[0].Val =
		fmt.Sprintf("%d-%d-%d",
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day()) //日期
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
func teaLogin(teaInfo TeaCInfo) []*http.Cookie {
	data := &url.Values{}
	data.Set("uname", teaInfo.Phone)
	data.Set("code", teaInfo.Pwd)
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

// 获取登录Cookie
func stuLogin(userInfo StuInfo) []*http.Cookie {
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
func readStuCheckList() []StuInfo {
	var userInfo []StuInfo
	isExist, _ := PathExists("./stuCheckList.json")
	if !isExist {
		//生成文件
		bytes, _ := json.Marshal([]StuInfo{})
		os.WriteFile("./stuCheckList.json", bytes, 0777)
	}
	bytes, _ := os.ReadFile("./stuCheckList.json")

	e := json.Unmarshal(bytes, &userInfo)
	if e != nil {
		log.Fatalln(e.Error())
	}
	return userInfo
}
func saveStuInfo(userInfo StuInfo) {
	userList := readStuCheckList()
	for i := range userList {
		if userList[i].Phone == userInfo.Phone {
			fmt.Println("该手机号已存在！")
			return
		}
	}
	userList = append(userList, userInfo)
	bytes, _ := json.Marshal(&userList)
	os.WriteFile("./stuCheckList.json", bytes, 0777)
}
func readTeaCheckList() []TeaCInfo {
	var teacherList []TeaCInfo
	isExist, _ := PathExists("./teacherCheckList.json")
	if !isExist {
		bytes, _ := json.Marshal([]TeaCInfo{})
		os.WriteFile("./teacherCheckList.json", bytes, 0777)

	}
	bytes, _ := os.ReadFile("./teacherCheckList.json")
	e := json.Unmarshal(bytes, &teacherList)
	if e != nil {
		log.Fatalln(e.Error())
	}
	return teacherList
}
func saveTeaCInfo(teaCInfo TeaCInfo) {
	teacherList := readTeaCheckList()
	for i := range teacherList {
		if teacherList[i].ID == teaCInfo.ID {
			fmt.Println("该工号已经存在！")
			return
		}
	}
	teacherList = append(teacherList, teaCInfo)
	bytes, _ := json.Marshal(&teacherList)
	os.WriteFile("./teacherCheckList.json", bytes, 0777)
}
