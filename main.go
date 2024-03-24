package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
)

var (
	friendName string = ""
	groupName  string = ""
	userType   string = ""
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if userType == "friend" {
			if msg.IsSendByFriend() {
				sendFriend, err := msg.Sender()
				if err != nil {
					log.Println(err)
					return
				}
				if sendFriend.NickName == friendName {
					if msg.IsText() {
						fmt.Println(sendFriend.NickName + ": " + msg.Content)
					}
				}
			}
		} else if userType == "group" {
			if msg.IsSendByGroup() {
				sendGroup, err := msg.Sender()
				if err != nil {
					log.Println(err)
					return
				}
				gp := openwechat.Group{User: sendGroup}
				if gp.NickName == groupName {
					sender, er := msg.SenderInGroup()
					if er != nil {
						log.Println(er)
						return
					}
					if msg.IsText() {
						fmt.Println(sender.NickName + ": " + msg.Content)
					}
				}
			}
		} else {
			if msg.IsText() && msg.Content == "ping" {
				msg.ReplyText("pong")
			}
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctl := "1"
	var g *openwechat.Group
	var f *openwechat.Friend
	for {
		switch ctl {
		case "1":
			fmt.Println(`请输入命令：
2-群名
3-好友名
q-退出`)
			fmt.Scanln(&ctl)
		case "2":
			fmt.Println("请输入群名:")
			fmt.Scanln(&groupName)
			userType = "group"
			// 获取所有的群组
			groups, err := self.Groups()
			if err != nil {
				panic(err)
			}
			for _, group := range groups {
				if group.NickName == groupName {
					g = group
					ctl = "g"
					break
				}
			}
		case "3":
			fmt.Println("请输入好友名:")
			fmt.Scanln(&friendName)
			userType = "friend"
			// 获取所有的群组
			friends, err := self.Friends()
			if err != nil {
				panic(err)
			}
			for _, friend := range friends {
				if friend.NickName == friendName {
					f = friend
					ctl = "f"
					break
				}
			}
		case "g":
			fmt.Println("开始聊天吧!输入q结束聊天")
			var text string
			for true {
				fmt.Scanln(&text)
				if text != "q" {
					_, er := self.SendTextToGroup(g, text)
					if er != nil {
						panic(er)
					}
				} else {
					ctl = "1"
					userType = ""
					groupName = ""
					break
				}
			}
		case "f":
			fmt.Println("开始聊天吧!输入q结束聊天")
			var text string
			for true {
				fmt.Scanln(&text)
				if text != "q" {
					_, er := self.SendTextToFriend(f, text)
					if er != nil {
						panic(er)
					}
				} else {
					ctl = "1"
					userType = ""
					friendName = ""
					break
				}
			}
		case "q":
			return
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	//bot.Block()
}
