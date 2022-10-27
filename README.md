# 学习通自动定时健康打卡程序（Go）

`该程序仅作学习交流使用，请勿用于商业，请勿二改破解`

## 简介

- 使用Go编写，极大降低内存与CPU占用

- 可进行晨检，午检。打卡时间可自定义（必须先添加成员信息后才能进行自定义）

- 支持多用户同时打卡

## 使用说明

- 程序共两个，无后缀的为Linux版本，可在服务器上运行。有exe后缀的为Windows版本，可在PC上测试。无需安装环境，双击或命令行直接运行即可

- 程序不会泄露你的个人信息，所有信息只会存放在本地，以及发送到学习通服务器

## 使用方法

   1. 在[Release](https://github.com/Qinwusui/AutoXXTCheck-Go/release)中下载对应系统架构的可执行程序
   2. 若是Windows 那么直接双击exe文件运行
   3. 若是Linux 直接在命令窗口中键入以下命令(以Linux x86_64举例)

   ```shell
   sudo chmod +x AutoCheck_linux_x86-64 ; ./AutoCheck_linux_x86-64
   ```

   4. 当程序运行后，可根据程序提示进行操作

### 发件箱相关
   
   1. 用户可以通过自定义发件箱达到打卡回执自动发送到指定邮箱
   2. 程序会将发件箱地址、密码（授权码）存储到可执行文件同一目录下，文件名为mailConfig.json
   3. 程序目前只支持SMTP协议发送邮件，且仅支持QQMail，OutlookMail，个人建议使用OutlookMail
   4. 若需使用QQMail发送邮件，则需要在QQMail Web中打开设置，并生成授权码。[参看](https://service.mail.qq.com/cgi-bin/help?subtype=1&&id=28&&no=369)
   5. 若使用`OutlookMail发送邮件，不需要设置授权码`直接使用邮件地址和密码即可。
   6. QQMail和OutlookMail都支持`587`端口
   
## 在Android Termux(以下简称终端)上运行

   - [参见](https://qa.1r1g.com/sf/ask/2727134721/#)
   - 在终端运行Golang程序时，Go会判断/etc/resolv.conf文件是否存在。当该文件不存在时，会默认使用localhost:53进行dns解析，所以才会导致域名解析失败。
   - 解决方案
     - termux安装zsh(网上有教程)
     - 在termux中输入`tsu`后回车
     - 输入`echo -e "nameserver 8.8.8.8" > /etc/resolv.conf`后回车
     - 重新启动打卡程序 
   - 由于tsu并不能作持久性更改，所以每次启动时都需要执行如上语句
