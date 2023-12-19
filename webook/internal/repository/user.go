package repository

import (
	"context"
	"geekcamp/webook/internal/domain"
	"geekcamp/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (repo *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return repo.dao.Edit(ctx, dao.User{
		Id:       u.Id,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Info:     u.Info,
	})
}

func (repo *UserRepository) Profile(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := repo.dao.Profile(ctx, dao.User{
		Id: u.Id,
	})
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
		NickName: user.NickName,
		Birthday: user.Birthday,
		Info:     user.Info,
	}, err
}
