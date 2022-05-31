package dao

import (
	"database/sql"
	"fmt"

	"github.com/go-camp/week02/code"
	"github.com/go-camp/week02/model"
	"github.com/pkg/errors"
)

type UserDao struct{}

func NewUser() *UserDao {
	return &UserDao{}
}

func (u *UserDao) GetByID(id int64) (*model.User, error) {
	//db伪方法getByID()
	user, err := getByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(code.ErrNotFound, fmt.Sprintf("dao：GetByID(%d) failed", id))
	}
	if err != nil {
		// 其他问题，直接抛错
		return nil, errors.Wrap(err, fmt.Sprintf("dao：GetByID(%d) failed", id))

	}
	return user, nil
}

func getByID(id int64) (*model.User, error) {
	return nil, sql.ErrNoRows
}
