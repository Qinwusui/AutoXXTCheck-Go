# 学习通自动定时健康打卡程序（Go）

`该程序仅作学习交流使用，请勿用于商业，请勿二改破解`

## 简介

- 使用Go编写，极大降低内存与CPU占用

- 可进行晨检，午检。打卡时间可自定义（必须先添加成员信息后才能进行自定义）

- 支持多用户同时打卡

## 使用说明

- 程序共两个，无后缀的为Linux（`amd_x64`）版本，可在服务器上运行。有exe后缀的为Windows（`amd_x64`）版本，可在PC上测试。无需安装环境，双击或命令行直接运行即可

- 程序不会泄露你的个人信息，所有信息只会存放在本地，以及发送到学习通服务器

## 使用方法

   1. 在[Release](https://github.com/Qinwusui/AutoXXTCheck-Go/release)中下载对应系统架构的可执行程序
   2. 若是Windows 那么直接双击exe文件运行
   3. 若是Linux 直接在命令窗口中键入以下命令(以Linux x86_64举例)

   ```shell
   sudo chmod -X AutoCheck_linux_x86-64 ; ./AutoCheck_linux_x86-64
   ```

   4. 当程序运行后，可根据程序提示进行操作

# TODO
   
   解决程序在Android Terminal上运行时出现DNS解析失败的问题
