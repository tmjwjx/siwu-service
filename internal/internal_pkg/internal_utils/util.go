package internal_utils

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// RandomGenerateStrings 随机生成长度为 l 的数字字母混合的字符串
func RandomGenerateStrings(l int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// HashPassword 加密密码。使用 bcrypt.GenerateFromPassword() 方法来加密密码。
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash 验证密码。使用 bcrypt.CompareHashAndPassword() 方法来验证输入的密码是否正确。
// password: 要比较的密码。
// hashedPassword: 原先的哈希密码。
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// DeleteFile 从文件系统中删除图片
func DeleteFile(home string, homeID uint) error {

	var attachment models.Attachment

	// 开启事务
	tx := globals.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteFile -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询要删除的图片文件路径
	err := tx.Where("home = ? and home_id = ?", home, homeID).First(&attachment).Error
	if err != nil {
		// return fmt.Errorf("deleteFile -> %s", err)
		// 没有查到说明文件系统中没有该图片，直接添加进入文件系统即可
		return nil
	}

	/*path := "./static" + attachment.Path

	// 删除文件系统中的图片
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("deleteFile -> 文件系统中的图片删除失败 -> %s", err)
	}*/

	// 删除 attachments 表中的图片相关信息
	err = tx.Delete(&attachment).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("deleteFile -> 删除 attachments 表中的图片相关信息失败 -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteFile -> 提交事务失败 -> %s", err)
	}

	return nil
}

// CreateFolder 创建存储静态文件的目录路径文件夹
func CreateFolder(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("CreateFolder -> 创建存储静态文件的目录路径文件夹失败 -> %s", err)
	}
	return nil
}

// ChangeType 将string类型的值转换成uint类型
func ChangeType(str string) (uint, error) {

	// 将字符串转换成uint64, 基数为 10, 位大小为 64 位
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("ChangeType -> 将string类型的值转换成uint类型失败 -> %s", err)
	}

	// 将 uint64 转换成 uint 类型
	uintNum := uint(num)

	return uintNum, nil
}

// TimeAgo 函数根据传入的 time.Time 和当前时间计算差值并返回相应的时间描述
func TimeAgo(t time.Time) string {
	duration := time.Since(t) // 计算传入时间和当前时间的差值

	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24
	months := days / 30
	years := days / 365

	if seconds < 60 {
		return fmt.Sprintf("%d秒前", seconds)
	} else if minutes < 60 {
		return fmt.Sprintf("%d分钟前", seconds)
	} else if hours < 24 {
		return fmt.Sprintf("%d小时前", hours)
	} else if days < 30 {
		return fmt.Sprintf("%d天前", days)
	} else if months < 12 {
		return fmt.Sprintf("%d个月前", months)
	} else {
		return fmt.Sprintf("%d年前", years)
	}
}

// TimeFormat 格式化 CreatedAt 为 年-月-日 时:分:秒
func TimeFormat(t time.Time) string {
	formattedTime := t.Format("2006-01-02 15:04:05")
	return formattedTime
}
