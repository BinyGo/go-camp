package service

import (
	"github.com/go-camp/week02/dao"
	"github.com/go-camp/week02/model"
	"github.com/pkg/errors"
)

type UserService struct{}

func NewUser() *UserService {
	return &UserService{}
}

func (d *UserService) GetUser(ID int64) (*model.User, error) {
	userDao := dao.NewUser()
	user, err := userDao.GetByID(ID)
	if err != nil {
		return nil, errors.WithMessagef(err, "service:GetUser(%d) failed", ID)
	}
	return user, err
}
