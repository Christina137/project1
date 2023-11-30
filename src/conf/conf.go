package conf

import "github.com/spf13/viper"

type MysqlConf struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
	DbName   string `yaml:"dbname"`
	Port     string `yaml:"port"`
}

type JwtConf struct {
	Secret string
}

type Configs struct {
	Mysql MysqlConf
	Jwt   JwtConf
}

var Config Configs

func InitConfig() {
	viper.SetConfigFile("./resources/application.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	mysql := MysqlConf{
		Url:      viper.GetString("mysql.url"),
		UserName: viper.GetString("mysql.userName"),
		PassWord: viper.GetString("mysql.passWord"),
		DbName:   viper.GetString("mysql.dbname"),
		Port:     viper.GetString("mysql.port"),
	}

	jwt := JwtConf{
		Secret: viper.GetString("jwt.secret"),
	}

	Config = Configs{
		Mysql: mysql,
		Jwt:   jwt,
	}
}

func GetConfig() Configs {
	return Config
}
