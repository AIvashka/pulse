package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"pulse/structs"
)

// LoadConfig loads the configuration from environment variables and/or config file
func LoadConfig() (*structs.Config, error) {

	viper.SetConfigType("json")
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config structs.SpreadConfig
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Look for .env file
	err = godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Set the config file name and path
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath(".")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Bind environment variables to viper
	viper.BindEnv("DBHOST")
	viper.BindEnv("DBPORT")
	viper.BindEnv("DBUSER")
	viper.BindEnv("DBPASSWORD")
	viper.BindEnv("DBNAME")
	viper.BindEnv("APIKEY")
	viper.BindEnv("APISECRET")

	// Read the values from viper
	cfg := &structs.Config{
		DBHost:       viper.GetString("DBHOST"),
		DBPort:       viper.GetString("DBPORT"),
		DBUser:       viper.GetString("DBUSER"),
		DBPassword:   viper.GetString("DBPASSWORD"),
		DBName:       viper.GetString("DBNAME"),
		SpreadConfig: config,
		APIKey:       viper.GetString("APIKEY"),
		APISecret:    viper.GetString("APISECRET"),
		Interval:     viper.GetInt("interval"),
	}

	return cfg, nil
}
