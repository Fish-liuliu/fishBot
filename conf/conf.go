package conf

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"path"
	"runtime"
)

var configStr = "config.yaml"

type Token struct {
	AppID       uint64
	AccessToken string
}

func GetToken() (*Token, error) {
	token := &Token{}
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		file := fmt.Sprintf("%s/%s", path.Dir(filename), configStr)
		var conf struct {
			AppID uint64 `yaml:"appid"`
			Token string `yaml:"token"`
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err = yaml.Unmarshal(content, &conf); err != nil {
			return nil, err
		}
		token.AppID = conf.AppID
		token.AccessToken = conf.Token
	}
	return token, nil
}

func (t *Token) GetString() string {
	return fmt.Sprintf("%v.%s", t.AppID, t.AccessToken)
}

func GetGormConnect() (*gorm.DB, error) {
	connectUrl, err := getDbConnectUrl(configStr)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	db, err := getDb(connectUrl)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}

func getDb(connectUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("getDb has err:%v", err)
		return nil, err
	}
	return db, nil
}

func getDbConnectUrl(name string) (string, error) {
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		file := fmt.Sprintf("%s/%s", path.Dir(filename), name)
		var conf struct {
			MysqlConnect string `yaml:"mysqlConnect"`
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		if err = yaml.Unmarshal(content, &conf); err != nil {
			return "", err
		}
		return conf.MysqlConnect, err
	}
	return "", nil
}

func GetRedisConnect() (*redis.Client, error) {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		file := fmt.Sprintf("%s/%s", path.Dir(filename), configStr)
		var RedisConnect struct {
			Addr string `yaml:"Addr"`
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err = yaml.Unmarshal(content, &RedisConnect); err != nil {
			return nil, err
		}
		// 创建 Redis 客户端
		client := redis.NewClient(&redis.Options{
			Addr:     RedisConnect.Addr, // Redis 地址，默认为 localhost:6379
			Password: "",                // Redis 密码，如果没有则留空
			DB:       0,                 // 默认数据库索引
		})
		return client, nil
	}
	return nil, errors.New("无法连接redis")
}
