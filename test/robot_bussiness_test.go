package test

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"main/dto"
	"main/service"
	"testing"
)

var redisConnect = "localhost:6379"

func InitCache() *redis.Client {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     redisConnect, // Redis 地址，默认为 localhost:6379
		Password: "",           // Redis 密码，如果没有则留空
		DB:       0,            // 默认数据库索引
	})
	return client
}

func TestRobtBusiness(t *testing.T) {
	t.Run(
		/*
			修改json文件  content 字段，进行不同场景的测试
			/猜数字
			/猜数字 x
			/告诉我猜数字的答案
		*/
		"Guess numbers", func(t *testing.T) {
			// json 文件 反序列
			message, err := GetMessage("guess_number_msg.json")
			if err != nil {
				t.Error(err)
			}
			db, err := InitDB()
			if err != nil {
				t.Error(err)
			}
			cache := InitCache()
			// 测验猜数游戏
			robotBussiness := service.NewRobotBussiness(nil, db, cache)
			handle, err := robotBussiness.MsgHandle(message)
			if err != nil {
				t.Error()
			}
			t.Logf("猜数字回复内容：%s", handle.Content)

		})
	t.Run("User Attendance", func(t *testing.T) {
		/*
			修改json文件  content 字段，进行不同场景的测试
			/打卡
			/我的考勤记录
		*/
		message, err := GetMessage("user_attendance_msg.json")
		if err != nil {
			t.Error(err)
		}
		db, err := InitDB()
		if err != nil {
			t.Error(err)
		}
		cache := InitCache()
		// 测试打卡功能
		robotBussiness := service.NewRobotBussiness(nil, db, cache)
		handle, err := robotBussiness.MsgHandle(message)
		if err != nil {
			t.Error()
		}
		t.Logf("打卡回复内容：%s", handle.Content)
	})
}

func GetMessage(fileName string) (*dto.Message, error) {
	// json 文件 反序列
	data, err := ioutil.ReadFile("user_attendance_msg.json")
	if err != nil {
		return nil, err
	}
	messageJson := gjson.Get(string(data), "d")
	message := &dto.Message{}
	if err = json.Unmarshal([]byte(messageJson.String()), message); err != nil {
		return nil, err
	}
	return message, err
}
