package service

import (
	"zeus/pkg/api/dao"
	"zeus/pkg/api/domain/account"
	"zeus/pkg/api/domain/account/login"
	"zeus/pkg/api/domain/user"
	"zeus/pkg/api/dto"
	"zeus/pkg/api/log"
	"zeus/pkg/api/model"
)

const pwHashBytes = 64

var userDao = dao.User{}

type UserService struct {
}

func (us UserService) InfoOfId(dto dto.GeneralGetDto) model.User {
	return userDao.Get(dto.Id,true)
}

// List - users list with pagination
func (us UserService) List(dto dto.GeneralListDto) ([]model.User, int64) {
	return userDao.List(dto)
}

// Create - create a new account
func (us UserService) Create(dto dto.UserCreateDto) model.User {
	salt, _ := account.MakeSalt()
	pwd, _ := account.HashPassword(dto.Password, salt)
	userModel := model.User{
		Username:     dto.Username,
		Mobile:       dto.Mobile,
		Password:     pwd,
		DepartmentId: dto.DepartmentId,
		Salt:         salt,
	}
	c := userDao.Create(&userModel)
	if c.Error != nil {
		log.Error(c.Error.Error())
	}
	return userModel
}

// Update - update user's information
func (us UserService) Update(dto dto.UserEditDto) int64 {
	userModel := model.User{
		Id:           dto.Id,
		Username:     dto.Username,
		Mobile:       dto.Mobile,
		DepartmentId: dto.DepartmentId,
		Status: dto.Status,
		Title: dto.Title,
		Realname: dto.Realname,
		Email: dto.Email,
	}

	c := userDao.Update(&userModel)
	return c.RowsAffected
}

// UpdateStatus - update user's status only
func (UserService) UpdateStatus(dto dto.UserEditStatusDto) int64  {
	user := userDao.Get(dto.Id,false)
	user.Status = dto.Status
	c := userDao.Update(&user)
	return c.RowsAffected
}

// UpdatePassword - update password only
func (UserService) UpdatePassword(dto dto.UserEditPasswordDto) int64 {
	salt,_ := account.MakeSalt()
	pwd,_ := account.HashPassword(dto.Password,salt)
	user := userDao.Get(dto.Id,false)
	user.Password = pwd
	user.Salt = salt
	c := userDao.Update(&user)
	return c.RowsAffected
}

// Delete - delete user
func (us UserService) Delete(dto dto.GeneralDelDto) int64 {
	userModel := model.User{
		Id: dto.Id,
	}
	c := userDao.Delete(&userModel)
	return c.RowsAffected
}

//VerifyAndReturnUserInfo - login and return user info
func (UserService) VerifyAndReturnUserInfo(dto dto.LoginDto) (bool, model.User) {
	userModel := userDao.GetByUserName(dto.Username)
	if login.VerifyPassword(dto.Password, userModel) {
		return true, userModel
	}
	return false, model.User{}
}

//AssignRole - assign roles to specific user
func (UserService) AssignRole(userId string, roleNames []string) {
	var groups [][]string
	for _, role := range roleNames {
		groups = append(groups, []string{userId, role})
	}
	user.OverwriteRoles(userId, groups)
}