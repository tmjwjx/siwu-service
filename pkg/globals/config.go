package globals

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
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// App 配置
type App struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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
}

// StaticConfig 静态文件配置
type StaticConfig struct {
	Prefix string `yaml:"prefix"` //  URL 路径前缀
	Path   string `yaml:"path"`   // 本地文件系统中的目录路径
}
