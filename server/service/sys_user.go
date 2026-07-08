package service

import (
	"errors"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/model"
	"github.com/tanganyu1114/heimdallr-reborn/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.sysUser
//@return: err error, userInter model.sysUser

func Register(u model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	// 否则 附加uuid 密码md5简单加密 注册
	u.Password = utils.MD5V([]byte(u.Password))
	u.UUID = uuid.NewV4()
	err = global.GVA_DB.Create(&u).Error
	return err, u
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Login
//@description: 用户登录
//@param: u *model.sysUser
//@return: err error, userInter *model.sysUser

func Login(u *model.SysUser) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).Preload("Authority").First(&user).Error
	return err, &user
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.sysUser, newPassword string
//@return: err error, userInter *model.sysUser

func ChangePassword(u *model.SysUser, newPassword string) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).Update("password", utils.MD5V([]byte(newPassword))).Error
	return err, u
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.SysUser{})
	var userList []model.SysUser
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
	return err, userList, total
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func SetUserAuthority(uuid uuid.UUID, authorityId string) (err error) {
	err = global.GVA_DB.Where("uuid = ?", uuid).First(&model.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func DeleteUser(id float64) (err error) {
	var user model.SysUser
	err = global.GVA_DB.Where("id = ?", id).Delete(&user).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.sysUser
//@return: err error, user model.sysUser

func SetUserInfo(reqUser model.SysUser) (err error, user model.SysUser) {
	err = global.GVA_DB.Updates(&reqUser).Error
	return err, reqUser
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.sysUser

func FindUserById(id int) (err error, user *model.SysUser) {
	var u model.SysUser
	err = global.GVA_DB.Where("`id` = ?", id).First(&u).Error
	return err, &u
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.sysUser

func FindUserByUuid(uuid string) (err error, user *model.SysUser) {
	var u model.SysUser
	if err = global.GVA_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}

// VerifyAPIKey verifies API Key and Secret
func VerifyAPIKey(apiKey, apiSecret string) (error, *model.SysUser) {
	var user model.SysUser
	if err := global.GVA_DB.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("API Key does not exist"), nil
		}
		return err, nil
	}

	if !utils.BcryptCheck(apiSecret, user.APISecret) {
		return errors.New("invalid API Secret"), nil
	}

	if !user.APIKeyEnabled {
		return errors.New("API Key is not enabled"), nil
	}

	return nil, &user
}

// GenerateAPIKeyForUser generates API Key and Secret for a user
func GenerateAPIKeyForUser(userID uint) (error, *model.SysUser) {
	var user model.SysUser
	if err := global.GVA_DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist"), nil
		}
		return err, nil
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return errors.New("failed to generate API Key"), nil
	}

	apiSecret, err := utils.GenerateAPISecret()
	if err != nil {
		return errors.New("failed to generate API Secret"), nil
	}

	hashedSecret, err := utils.BcryptHash(apiSecret)
	if err != nil {
		return errors.New("failed to hash API Secret"), nil
	}

	if err := global.GVA_DB.Model(&user).Updates(map[string]interface{}{
		"api_key":         apiKey,
		"api_secret":      hashedSecret,
		"api_key_enabled": false,
	}).Error; err != nil {
		return errors.New("failed to save API Key"), nil
	}

	user.APIKey = apiKey
	user.APISecret = apiSecret
	user.APIKeyEnabled = false

	return nil, &user
}

// ToggleAPIKey enables or disables API Key for a user
func ToggleAPIKey(userID uint, enabled bool) (error, *model.SysUser) {
	var user model.SysUser
	if err := global.GVA_DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist"), nil
		}
		return err, nil
	}

	if user.APIKey == "" {
		return errors.New("user has no API Key, please generate one first"), nil
	}

	if err := global.GVA_DB.Model(&user).Update("api_key_enabled", enabled).Error; err != nil {
		return errors.New("failed to update API Key status"), nil
	}

	user.APIKeyEnabled = enabled
	return nil, &user
}

// RegenerateAPISecret regenerates API Secret for a user
func RegenerateAPISecret(userID uint) (error, *model.SysUser) {
	var user model.SysUser
	if err := global.GVA_DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist"), nil
		}
		return err, nil
	}

	if user.APIKey == "" {
		return errors.New("user has no API Key, please generate one first"), nil
	}

	apiSecret, err := utils.GenerateAPISecret()
	if err != nil {
		return errors.New("failed to generate API Secret"), nil
	}

	hashedSecret, err := utils.BcryptHash(apiSecret)
	if err != nil {
		return errors.New("failed to hash API Secret"), nil
	}

	if err := global.GVA_DB.Model(&user).Update("api_secret", hashedSecret).Error; err != nil {
		return errors.New("failed to save API Secret"), nil
	}

	user.APISecret = apiSecret
	return nil, &user
}
