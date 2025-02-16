package internalUtils

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// RandomGenerateNickname 随机生成唯一名字
func RandomGenerateNickname() (string, error) {
	nickname, err := gonanoid.Generate(CustomAlphabetNickname, UserNameLen)
	if err != nil {
		fmt.Println("Error generating ID:", err)
		return "", err
	}
	return nickname, nil
}

// RandomGenerateVerifyCode 随机生成验证码
func RandomGenerateVerifyCode() (string, error) {
	nickname, err := gonanoid.Generate(CustomAlphabetVerifyCode, VerifyCodeLen)
	if err != nil {
		fmt.Println("Error generating ID:", err)
		return "", err
	}
	return nickname, nil
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

// EmptyFilled
// @Description: 空接口填充
// @param        data interface{}
func EmptyFilled(data interface{}) interface{} {
	if data == nil {
		data = gin.H{}
	}
	return data
}

// ArticlesOrder
// @Description: 选择排序方式 0热度 1时间
// @param        kind int
// @return       string
func ArticlesOrder(kind int) string {
	var condition string
	if kind == 0 { // 0 代表按照热度排序
		condition = "heat DESC"
	} else if kind == 1 { // 1 代表按照发布时间排序
		condition = "published_at DESC"
	}
	return condition
}

// TimeFormatDaily
// @Description: 格式化标准时间为日常时间
// @param        time.Time
// @return       date
// @return       err
// @Author tianjiajie 2025-01-21 11:35:28
func TimeFormatDaily(t *time.Time) (date string) {
	if t == nil {
		return ""
	}
	duration := time.Since(*t)
	// 根据时间差判断返回值
	if duration < time.Minute {
		// 小于1分钟
		return "刚刚"
	} else if duration < time.Hour {
		// 小于1小时
		return fmt.Sprintf("%d分钟前", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		// 小于1天
		return fmt.Sprintf("%d小时前", int(duration.Hours()))
	} else if duration < 7*24*time.Hour {
		// 小于1周
		return fmt.Sprintf("%d天前", int(duration.Hours()/24))
	} else {
		// 大于1周
		return t.Format("2006-01-02")
	}
}

// TimeFormat
// @Description: 格式化时间
// @param        t *time.Time
// @return       format
// @Author tianjiajie 2025-01-22 12:07:57
func TimeFormat(t *time.Time) (format string) {
	if t == nil {
		return ""
	}
	format = t.Format("2006-01-02")
	return format
}

// MessagePush
// @Description: 向用户实时发送更新数据
// @param        data string
// @param        userId string
// @Author tianjiajie 2024-10-10 21:45:13
func MessagePush(data string, userId string) {
	notifyChan, exist := globals.SubscriberChannels[userId]
	if exist {
		notifyChan <- data
	}
}

// Highlight
// @Description: 高亮处理
// @param        content string
// @param        keyword string
// @return       string
// @Author tianjiajie 2025-02-08 20:28:05
func Highlight(content string, keyword string) string {
	// 高亮处理
	if keyword == "" {
		return content
	}

	// keywordLower := strings.ToLower(keyword)
	// contentLower := strings.ToLower(content)
	//
	// if strings.Contains(contentLower, keywordLower) {
	//	// 高亮显示匹配项
	// }

	// content = strings.ReplaceAll(content, keyword, fmt.Sprintf("<mark>%s</mark>", keyword))

	// 转换为小写，用于大小写不敏感的匹配
	keywordLower := strings.ToLower(keyword)
	contentLower := strings.ToLower(content)

	if strings.Contains(contentLower, keywordLower) {
		// 使用正则表达式忽略大小写替换
		regex := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
		// content = regex.ReplaceAllString(content, fmt.Sprintf("<mark>%s</mark>", keyword))
		return regex.ReplaceAllStringFunc(content, func(match string) string {
			return fmt.Sprintf("<mark>%s</mark>", match)
		})
	}

	return content
}

/*// ChangeStringToUint
// @Description: 将string类型的值转换成uint类型
// @Author wangyulong 2024-10-09 15:26:15
// @param        str string
// @return       uint
// @return       error
func ChangeStringToUint(str string) (uint, error) {

	// 将字符串转换成uint64, 基数为 10, 位大小为 64 位
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("ChangeType -> 将string类型的值转换成uint类型失败 -> %s", err)
	}

	// 将 uint64 转换成 uint 类型
	uintNum := uint(num)

	return uintNum, nil
}*/

// ChangeStringToInt
// @Description: 将string类型的值转换成int类型
// @Author wangyulong 2024-10-10 12:32:45
// @param        str string
// @return       int
// @return       error
func ChangeStringToInt(str string) (int, error) {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("ChangeStringToInt -> 将string类型的值转换成int类型失败 -> %s", err)
	}
	return intValue, nil
}

// TimeAgo 函数根据传入的 time.Time 和当前时间计算差值并返回相应的时间描述
func TimeAgo(t time.Time) string {
	duration := time.Since(t) // 计算传入时间和当前时间的差值
	// y, m, d := t.Date()
	// fmt.Println(y, m, d)

	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24
	months := days / 30
	years := days / 365

	if seconds < 60 {
		return fmt.Sprintf("%d秒前", seconds)
	} else if minutes < 60 {
		return fmt.Sprintf("%d分钟前", minutes)
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
// func TimeFormat(t time.Time) string {
//	formattedTime := t.Format("2006-01-02 15:04:05")
//	return formattedTime
// }

// ChangeAnyToUint
// @Description: 将Any类型转换成Uint类型
// @Author wangyulong 2025-01-16 20:50:05
// @param        value any
// @return       uint
// @return       error
func ChangeAnyToUint(v any) (uint, error) {
	var uintValue uint
	// 使用类型断言进行转换
	if num, ok := v.(uint); ok {
		uintValue = num
	} else {
		fmt.Println("ChangeAnyToUint -> 将Any类型转换成Uint类型失败")
	}
	return uintValue, nil
}

// RemoveDuplicates
// @Description: uint类型的切片去重
// @Author wangyulong 2025-02-07 19:45:04
// @param        seen map[int]struct{}
// @param        slice []int
// @return       []int
func RemoveDuplicates(seen *map[uint]struct{}, slice []uint) {
	for _, v := range slice {
		_, ok := (*seen)[v]
		if !ok {
			(*seen)[v] = struct{}{}
		}
	}
}

// RemoveDuplicates2
// @Description: 获取一个切片相对于另一个切片中没有的元素
// @Author wangyulong 2025-02-08 20:48:25
// @param        seen map[uint]struct{}
// @param        slice []uint
// @return       []uint
func RemoveDuplicates2(seen *map[uint]struct{}, slice []uint) []uint {
	var result []uint

	for _, v := range slice {
		_, ok := (*seen)[v]
		if !ok {
			// seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// ProcessImagePath
// @Description: 将上传图片的url路径处理成其在服务器中的路径
// @Author wangyulong 2025-02-13 16:02:59
func ProcessImagePath(url string) (string, error) {
	var outFile string // 图片在服务器的存储路径
	if url == "" {
		return "", nil
	}

	originalURL := url

	// 获取最后一个 '/' 的位置
	lastSlashIndex := strings.LastIndex(originalURL, "/")
	// 如果找到了 '/', 进行截取
	if lastSlashIndex != -1 {

		// 获取最后一个 '/' 前面的部分（不包含最后一个 '/'）
		prefix := originalURL[:lastSlashIndex]
		// 替换前缀为 './static/images'
		outFile = strings.Replace(originalURL, prefix, "./static/images", 1)

	} else {
		return "", fmt.Errorf("ProcessImagePath -> 获取最后一个 '/' 的位置失败")
	}

	return outFile, nil
}
