package dao

import (
	"database/sql"

	"github.com/go-camp/week02/model"
	"github.com/pkg/errors"
)

type UserDao struct{}

func NewUser() *UserDao {
	return &UserDao{}
}

func (d *UserDao) GetByID(id int64) (*model.User, error) {
	//db伪方法getByID()
	user, err := getByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "dao:GetByID failed")
	}
	return user, nil
}

func getByID(id int64) (*model.User, error) {
	return nil, sql.ErrNoRows
}
