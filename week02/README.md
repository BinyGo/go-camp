# 极客时间 Go训练营作业

## 作业
我们在做数据库操作的时候，假设在 dao 层中遇到一个 sql.ErrNoRows ，是否应该Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 答：
<p>应该warp这个error，并逐层上抛;</p>
<p>一旦在上抛过程错误被处理了，就不允许再往上传递调用栈，它不能返回错误值。它应该只返回零（比如降级处理中，你返回了降级数据，然后需要 return nil）;</p>
<p>如未处理则在程序顶部或者工作的goroutine顶部(请求入口)进行处理;</p>

### Dao层伪代码
```
func (u *UserDao) GetByID(id int64) (*model.User, error) {
	user := model.User{}
	err := DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "dao:GetByID failed")
	}
	return user, nil
}
```

### service层伪代码--错误上抛
```
func (u *UserService) GetUser(ID int64) (*model.User, error) {
	userDao := dao.NewUser()
	user, err := userDao.GetByID(ID)
	if err != nil {
		return nil, errors.WithMessagef(err, "service:GetUser(%d) failed", ID)
	}
	return user, err
}
```

### service层伪代码--错误降级处理
```
func (u *UserService) GetUser(ID int64) (*model.User, error) {
	userDao := dao.NewUser()
	user, err := userDao.GetByID(ID)
	if errors.Is(err,dao.ErrNoRows) {
		user = dao.GetUserDefault()
		err = nil
	}
	return user, err
}
```