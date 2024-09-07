package dao

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

// Attendance undefined
type Attendance struct {
	ID             int64     `json:"id" gorm:"id"`                           // 主键
	UserId         string    `json:"user_id" gorm:"user_id"`                 // 成员id
	Name           string    `json:"name" gorm:"name"`                       // 名称
	AttendanceDate time.Time `json:"attendance_date" gorm:"attendance_date"` // 打卡日期
	Type           int       `json:"type" gorm:"type"`                       // 0-未打卡，1-已打卡
}

// TableName 表名称
func (*Attendance) TableName() string {
	return "attendance"
}

func (a *Attendance) Create(db *gorm.DB) error {
	return db.Table(a.TableName()).Create(&a).Error
}

func (a *Attendance) CheckAttendance(db *gorm.DB, userId string) (bool, error) {
	// 检查的日期为今天
	dateToCheck := time.Now().Format("2006-01-02")
	db = db.Table(a.TableName()).
		Where("user_id", userId).
		Where("date(attendance_date)=?", dateToCheck).
		Where("type=?", 1)
	if err := db.First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (a *Attendance) GetAttendances(db *gorm.DB, userId string) ([]*Attendance, error) {
	// 获取当前月份的第一天和最后一天
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)
	log.Printf("第一天：firstDayOfMonth：%s,最后一天：lastDayOfMonth：%s", firstDayOfMonth, lastDayOfMonth)

	attendances := make([]*Attendance, 0)

	// 获取这个月的打卡记录列表
	db = db.Table(a.TableName()).
		Where("user_id", userId).
		Where("attendance_date BETWEEN ? AND ?", firstDayOfMonth, lastDayOfMonth)

	if err := db.Find(&attendances).Error; err != nil {
		log.Fatalf("GetAttendance has err:%v", err)
		return nil, err
	}

	return attendances, nil
}

func (a *Attendance) CountAttendance(attendances []*Attendance) (int64, int64) {
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	// 使用time.Date创建一个该月第一天之后一个月的日期，但日期设为0，让它回退到上一个月的最后一天
	lastDay := time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.UTC)
	log.Printf("这个月总共有%d天", lastDay.Day())
	// 获取当前月份所有日期
	dates := make([]time.Time, 0)
	for day := 0; day < lastDay.Day(); day++ {
		date := firstDayOfMonth.AddDate(0, 0, day)
		dates = append(dates, date)
	}

	// 统计打卡天数和未打卡天数
	checkedInDates := make(map[time.Time]bool)
	for _, a := range attendances {
		checkedInDates[a.AttendanceDate] = true
	}

	var checkedInDays, noCheckInDays int64

	for _, date := range dates {
		if _, exists := checkedInDates[date]; exists {
			checkedInDays++
		} else {
			noCheckInDays++
		}
	}

	return checkedInDays, noCheckInDays
}
