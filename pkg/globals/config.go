package globals

// 和配置相关的结构体

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type App struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// VerifyCodeConfig 验证码配置
type VerifyCodeConfig struct {
	From                        string `yaml:"from"`                        // 发送者邮箱
	Subject                     string `yaml:"subject"`                     // 验证码主题
	VerifyCodeLen               int    `yaml:"verifyCodeLen"`               // 验证码长度
	VerifyCodeEffectiveDuration uint   `yaml:"verifyCodeEffectiveDuration"` // 验证码有效时长，单位秒
	VerifyCodeCoolTime          uint   `yaml:"verifyCodeCoolTime"`          // 发送验证码的冷却时间
	Host                        string `yaml:"host"`                        // SMTP 服务器的主机地址
	Port                        int    `yaml:"port"`                        // 服务器的端口号
	Username                    string `yaml:"username"`                    // SMTP 服务器的用户名（通常是你的邮箱地址）
	AuthorizeCode               string `yaml:"authorizeCode"`               // SMTP 服务器的密码（或者授权码）
}

// Config 总配置
type Config struct {
	Database         DatabaseConfig   `yaml:"database"`
	Redis            RedisConfig      `yaml:"redis"`
	App              App              `yaml:"app"`
	VerifyCodeConfig VerifyCodeConfig `yaml:"verifyCode"`
}
