package test

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"main/dao"
	"testing"
	"time"
)

var mysqlConnect = "root:123456@tcp(127.0.0.1:3306)/fishbot?charset=utf8mb4&parseTime=True&loc=Local"

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(mysqlConnect), &gorm.Config{})
	if err != nil {
		log.Fatalf("getDb has err:%v", err)
		return nil, err
	}
	return db, nil
}

func TestRandomNumBers(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Error(err)
	}
	t.Run(
		"GetRandomNum", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				randomNumbers := &dao.RandomNumbers{}
				randomNum, err := randomNumbers.GetRandomNum(db)
				if err != nil {
					t.Error(err)
				}
				t.Logf("获取随机数：第%d次，%d", i, randomNum)
			}
		},
	)
}

func TestAttendance(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Error(err)
	}
	t.Run(
		"Create", func(t *testing.T) {
			attendance := &dao.Attendance{
				UserId:         "11667722671262063122",
				Name:           "心温",
				AttendanceDate: time.Now(),
				Type:           1,
			}
			if err := attendance.Create(db); err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"CheckAttendance", func(t *testing.T) {
			attendance := &dao.Attendance{}
			var userIds = []string{"11667722671262063122", "11667722671262063123"}
			for _, id := range userIds {
				checkAttendance, err := attendance.CheckAttendance(db, id)
				if err != nil {
					t.Error(err)
				}
				if checkAttendance == true {
					t.Logf("用户：%s，今日已打卡", id)
				} else {
					t.Logf("用户：%s，今日无打卡", id)
				}
			}
		})
	t.Run(
		"GetAttendances", func(t *testing.T) {
			attendance := &dao.Attendance{}
			getAttendances, err := attendance.GetAttendances(db, "11667722671262063122")
			if err != nil {
				t.Error(err)
			}
			for _, getAttendance := range getAttendances {
				t.Logf("打卡情况：%v", getAttendance)
			}
			checkedInDays, noCheckInDays := attendance.CountAttendance(getAttendances)
			t.Logf("这个月打卡：%d,未打卡：%d", checkedInDays, noCheckInDays)
		})
}
