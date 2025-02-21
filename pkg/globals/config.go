package globals

import (
	"time"
)

// 和配置相关的结构体

// DatabaseConfig mysql配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`      // Redis 连接池中的最大连接数
	MinIdleConns int           `yaml:"min_idle_conns"` // Redis 连接池中的最小空闲连接数
	IdleTimeout  time.Duration `yaml:"idle_timeout"`   // Redis 连接池中空闲连接的最大超时时间
	DialTimeout  time.Duration `yaml:"dial_timeout"`   // Redis 连接超时
	ReadTimeout  time.Duration `yaml:"read_timeout"`   // Redis 读取数据超时
	WriteTimeout time.Duration `yaml:"write_timeout"`  // Redis 写入数据超时
	MaxRetries   int           `yaml:"max_retries"`    // Redis 最大重试次数
}

// App 配置
type App struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Domain string `yaml:"domain"` // 域名
}

// SendEmailConfig 验证码配置
type SendEmailConfig struct {
	From          string `yaml:"from"`          // 发送者邮箱
	Host          string `yaml:"host"`          // SMTP 服务器的主机地址
	Port          int    `yaml:"port"`          // 服务器的端口号
	Username      string `yaml:"username"`      // SMTP 服务器的用户名（通常是你的邮箱地址）
	AuthorizeCode string `yaml:"authorizeCode"` // SMTP 服务器的密码（或者授权码）
}

// Config 总配置
type Config struct {
	Database     DatabaseConfig  `yaml:"database"`
	Redis        RedisConfig     `yaml:"redis"`
	App          App             `yaml:"app"`
	SendEmailCfg SendEmailConfig `yaml:"verifyCode"`
	Log          LogConfig       `yaml:"log"`
}

// LogConfig
// @Description: 日志配置
// @Author tianjiajie 2024-10-11 21:59:41
type LogConfig struct {
	Level   string `yaml:"level"`
	LogPath string `yaml:"logPath"`
	AppName string `yaml:"appName"`
}

// StaticConfig 静态文件配置
type StaticConfig struct {
	Prefix                    string  `yaml:"prefix"`                                                                 //  URL 路径前缀
	Path                      string  `yaml:"path"`                                                                   // 本地文件系统中的目录路径
	Limit                     int64   `yaml:"limit"`                                                                  // 上传图片文件的大小限制，单位是MB
	Watermark                 string  `yaml:"watermark"`                                                              // 图片水印的内容
	Size                      float64 `yaml:"size"`                                                                   // 水印的大小
	Scale                     float64 `yaml:"scale"`                                                                  // 图片压缩时的缩放比例
	Margin                    float64 `yaml:"margin"`                                                                 // 水印距离左下角的距离 ，规则: 从右向左移动margin个位置的距离，从下往上移动3*margin个位置的距离
	UserDefaultImage          string  `yaml:"user_default_image" mapstructure:"user_default_image"`                   //  默认用户头像路径
	ArticleDefaultImage       string  `yaml:"article_default_image" mapstructure:"article_default_image"`             // 默认文章图片路径
	TagDefaultImage           string  `yaml:"tag_default_image" mapstructure:"tag_default_image"`                     //  默认标签图片
	AdvertisementDefaultImage string  `yaml:"advertisement_default_image" mapstructure:"advertisement_default_image"` // 默认广告图片
	CommentDefaultImage       string  `yaml:"comment_default_image" mapstructure:"comment_default_image"`             // 默认评论图片
	CategoryDefaultImage      string  `yaml:"category_default_image" mapstructure:"category_default_image"`           // 默认类目图片
}
