package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name string   `mapstructure:"name" json:"name"`
	Host string   `mapstructure:"host" json:"host"`
	Tags []string `mapstructure:"tags" json:"tags"`
	Port int      `mapstructure:"port" json:"port"`
	Env  string   `mapstructure:"env" json:"env"`

	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}
