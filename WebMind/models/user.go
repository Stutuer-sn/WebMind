// 负责处理用户相关的数据和逻辑
package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// 密码哈希存储
func (u *User) SetPassword(password string) error {
	cost := 10                                                                 //成本参数
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost) //使用bcrypt.GenerateFromPassword函数
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码-----待办
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) //使用bcrypt.CompareHashAndPassword函数
	if err != nil {                                                            //密码匹配不成功
		return err
	}
	return nil //密码匹配成功
}
