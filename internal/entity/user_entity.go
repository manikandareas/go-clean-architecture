package entity

import "time"

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	Token     string    `gorm:"column:token"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}
