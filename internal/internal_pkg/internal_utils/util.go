package internal_utils

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// RandomGenerateStrings 随机生成长度为 l 的数字字母混合的字符串
func RandomGenerateStrings(l int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
		//return fmt.Errorf("deleteFile -> %s", err)
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
