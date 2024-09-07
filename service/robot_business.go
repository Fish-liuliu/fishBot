package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"main/client"
	"main/constant"
	"main/dao"
	"main/dto"
	"regexp"
	"strconv"
	"time"
)

type RobotBussiness struct {
	Client *client.HttpClient
	Db     *gorm.DB
	cache  *redis.Client
}

func NewRobotBussiness(client *client.HttpClient, db *gorm.DB, cache *redis.Client) *RobotBussiness {
	return &RobotBussiness{
		Client: client,
		Db:     db,
		cache:  cache,
	}
}

func (r *RobotBussiness) AtRoBot(message *dto.Message) error {
	replyMessage, err := r.MsgHandle(message)
	if err != nil {
		log.Fatalf("AtRoBot.mesHandle has err:%v", err)
		return err
	}
	// 日调用不能超过20次，坑！
	r.Client.Post(constant.MessagesURI, "channel_id", message.ChannelID, replyMessage)
	return nil
}

func (r *RobotBussiness) MsgHandle(message *dto.Message) (*dto.ReplyMessage, error) {
	log.Printf("读取的消息内容：%s", message.Content)
	//指令解析
	match := regexp.MustCompile(`^(.+?)\s*/([^ ]+)\s*(\d*)$`).FindStringSubmatch(message.Content)
	var replyContent = "我无法识别你的命令"
	var err error
	if len(match) > 2 {
		replyContent, err = r.cmdHandle(match, message.Author.ID, message.Author.Username)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("replyContent:%s", replyContent)
	return &dto.ReplyMessage{
		Content: replyContent,
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             message.ID,
			IgnoreGetMessageError: true,
		},
	}, nil
}

func (r *RobotBussiness) cmdHandle(match []string, userId, userName string) (string, error) {
	var replyContent = ""
	if match[3] == "" {
		switch match[2] {
		case "猜数字":
			if exists := r.cache.Exists(context.TODO(), constant.NUMKEY).Val(); exists == 1 {
				replyContent = "上一轮游戏还没结束，请继续竞猜"
			} else {
				log.Println("触发游戏，从库中寻找游戏，并且写入缓存")
				replyContent = "开始猜数字游戏，请输入一个1-10的数字，直到猜对为止"
				randomNumbers := &dao.RandomNumbers{}
				randomNum, err := randomNumbers.GetRandomNum(r.Db)
				if err != nil {
					return "", err
				}
				log.Printf("从数据库拿到：%d", randomNum)
				r.cache.Set(context.TODO(), constant.NUMKEY, randomNum, 0)
			}

		case "告诉我猜数字答案":
			log.Println("从缓存中查找答案，并清空缓存")
			if exists := r.cache.Exists(context.TODO(), constant.NUMKEY).Val(); exists == 0 {
				log.Println("缓存中还没数据")
				replyContent = "你还没开始游戏，输入无效"
			} else {
				answer := r.cache.Get(context.TODO(), constant.NUMKEY).Val()
				replyContent = fmt.Sprintf("答案是：%s，游戏结束", answer)
				r.cache.Del(context.TODO(), constant.NUMKEY)
			}

		case "打卡":
			// 查询是否有打卡，没有入库，有则返回打卡
			attendance := &dao.Attendance{}
			checkAttendance, err := attendance.CheckAttendance(r.Db, userId)
			if err != nil {
				return "", err
			}
			if checkAttendance {
				// 如果当天已经打卡了
				replyContent = "当天你已经打卡，无需重复打卡"
			} else {
				attendance.UserId = userId
				attendance.Name = userName
				attendance.AttendanceDate = time.Now()
				attendance.Type = 1
				if err := attendance.Create(r.Db); err != nil {
					return "", err
				}
				// 当天没打卡，落库，打卡成功
				replyContent = "打卡成功"
			}
		case "我的考勤记录":
			// 捞取这个月打卡记录
			attendance := &dao.Attendance{}
			checkAttendance, err := attendance.CheckAttendance(r.Db, userId)
			if err != nil {
				return "", err
			}
			if !checkAttendance {
				// 如果当天已经打卡了
				replyContent = "这个月你还没有打卡"
			} else {
				attendances, err := attendance.GetAttendances(r.Db, userId)
				if err != nil {
					return "", err
				}
				checkedInDays, noCheckInDays := attendance.CountAttendance(attendances)
				replyContent = fmt.Sprintf("你这个月打卡 %d 天，未打卡 %d", checkedInDays, noCheckInDays)
			}
		}
	} else {
		// 先检测键是否过期
		if exists := r.cache.Exists(context.TODO(), constant.NUMKEY).Val(); exists == 0 {
			log.Println("缓存中还没数据")
			replyContent = "你还没开始游戏，输入无效"
		} else {
			// 大了跟他说大了，小了跟他说小了。等于的话游戏结束，清空缓存
			answer := r.cache.Get(context.TODO(), constant.NUMKEY).Val()
			atoiOutput, _ := strconv.Atoi(answer)
			atoiInput, _ := strconv.Atoi(match[3])
			log.Printf("答案：%d", atoiOutput)
			log.Printf("输入：%d", atoiInput)
			if atoiInput == atoiOutput {
				replyContent = "恭喜你答对了，游戏结束"
				r.cache.Del(context.TODO(), constant.NUMKEY)
			} else if atoiInput > atoiOutput {
				replyContent = "你输入的大了"
			} else {
				replyContent = "你输入的小了"
			}
		}
	}
	return replyContent, nil
}
