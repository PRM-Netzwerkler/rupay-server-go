package util

import "github.com/spf13/viper"

type Config struct {
	DbDriver           string `mapstructure:"DB_DRIVER"`
	DbSource           string `mapstructure:"DB_SOURCE"`
	PostgresUser       string `mapstructure:"POSTGRES_USER"`
	PostgresPassword   string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDb         string `mapstructure:"POSTGRES_DB"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	SavaPageAdmin      string `mapstructure:"SAVAPAGE_ADMIN"`
	SavaPagePassword   string `mapstructure:"SAVAPAGE_PASSWORD"`
	SavaPageUrl        string `mapstructure:"SAVAPAGE_API"`
	MqttBroker         string `mapstructure:"MQTT_BROKER"`
	MqttClientId       string `mapstructure:"MQTT_CLIENT_ID"`
	MqttClientName     string `mapstructure:"MQTT_CLIENT_NAME"`
	MqttClientPassword string `mapstructure:"MQTT_CLIENT_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
