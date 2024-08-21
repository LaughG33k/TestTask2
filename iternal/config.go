package iternal

type AppCfg struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	MaxConns        int    `yaml:"maxConns"`
	ReadTimeoutSec  int64  `yaml:"readTimeoutSec"`
	WriteTimeoutSec int64  `yaml:"writeTimeoutSec"`

	NewsDbCfg DbConfig `yaml:"newsDbCfg"`
	UserDbCfg DbConfig `yaml:"userDbCfg"`
}

type DbConfig struct {
	Host        string `yaml:"host"`
	Port        uint16 `yaml:"port"`
	User        string `yaml:"user"`
	DB          string `yaml:"db"`
	Password    string `ymal:"password"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
}
