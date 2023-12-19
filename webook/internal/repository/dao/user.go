package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var ErrUserDuplicateEmail = errors.New("邮箱冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail

		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Edit(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	return dao.db.WithContext(ctx).Model(&u).Updates(User{
		Utime:    u.Utime,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	}).Error
}

func (dao *UserDAO) Profile(ctx context.Context, u User) (User, error) {
	key := strconv.FormatInt(u.Id, 10)
	err := dao.db.WithContext(ctx).Where("id=?", key).First(&u).Error
	resUser := User{
		Id:       u.Id,
		Email:    u.Email,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	}
	return resUser, err
}

// 对应数据库表结构
type User struct {
	Id int64 `gorm:"primaryKer,autoIncrement"`
	// 用户唯一标识
	Email    string `gorm:"unique"`
	Password string

	// 创建时间(毫秒)
	Ctime int64
	// 更新时间(毫秒)
	Utime int64

	NickName string
	Birthday string
	Info     string
}
