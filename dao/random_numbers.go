package dao

import (
	"gorm.io/gorm"
	"log"
)

// RandomNumbers undefined
type RandomNumbers struct {
	ID  int64 `json:"id" gorm:"id"`
	Num int64 `json:"num" gorm:"num"`
}

// TableName 表名称
func (*RandomNumbers) TableName() string {
	return "random_numbers"
}

func (R *RandomNumbers) GetRandomNum(db *gorm.DB) (int64, error) {
	if err := db.Table("random_numbers").Order("RAND()").Limit(1).First(&R).Error; err != nil {
		log.Fatalf("GetRandomNum has err:%v", err)
		return 0, err
	}
	return R.Num, nil
}
