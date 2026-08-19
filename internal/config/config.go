package config

import (
	"encoding/json"
	"os"
)

type Config struct {
  DB_URL 			string `json:"db_url"`
  CurrentUserName 	string `json:"current_user_name"`
}

// function that reads the JSON file found at ~/.gatorconfig.json and returns a Config struct. It should read the file from the HOME directory, then decode the JSON string into a new Config struct. I used os.UserHomeDir to get the location of HOME.
func Read () (*Config) {
	dir, error := os.UserHomeDir()
	if error != nil {
		return nil
	}
	data, err := os.ReadFile(dir + "/.gatorconfig.json")
	if err != nil {
		return nil
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil
	}
	return &cfg

}

//  method on the Config struct that writes the config struct to the JSON file after setting the current_user_name field.
func (cfg *Config) SetUser (userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir + "/.gatorconfig.json", nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
	 