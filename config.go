package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const configPath = "config.json"

type Config struct {
	DiscordUserID   string `json:"discord_user_id"`
	HoyolabUID      string `json:"hoyolab_uid"`
	ZenlessUID      string `json:"zenless_uid"`
	DiscordClientID string `json:"discord_client_id"`
	DiscordBotToken string `json:"discord_bot_token"`
	HoyolabCookie   string `json:"hoyolab_cookie"`
	Region          string `json:"region"`
}

var regions = []struct {
	Label string
	Code  string
}{
	{"Japan", "prod_gf_jp"},
	{"EU", "prod_gf_eu"},
	{"NA", "prod_gf_us"},
	{"TW/HK/MO", "prod_gf_sg"},
}

func validRegion(code string) bool {
	for _, r := range regions {
		if r.Code == code {
			return true
		}
	}
	return false
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func saveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}

func promptConfig() (Config, error) {
	reader := bufio.NewReader(os.Stdin)
	ask := func(label string) string {
		fmt.Printf("%s: ", label)
		line, _ := reader.ReadString('\n')
		return strings.TrimSpace(line)
	}

	var cfg Config
	cfg.ZenlessUID = ask("Zenless UID")

	fmt.Println("Select region:")
	for i, r := range regions {
		fmt.Printf("  %d) %s (%s)\n", i+1, r.Label, r.Code)
	}
	choice := ask("Region number")
	idx := 0
	fmt.Sscanf(choice, "%d", &idx)
	if idx < 1 || idx > len(regions) {
		return cfg, fmt.Errorf("invalid region selection: %q", choice)
	}
	cfg.Region = regions[idx-1].Code

	cfg.DiscordUserID = ask("Discord user ID")
	cfg.DiscordBotToken = ask("Widget bot token")
	cfg.DiscordClientID = ask("Widget bot ID")
	cfg.HoyolabUID = ask("Hoyolab ID")
	cfg.HoyolabCookie = ask("Hoyolab cookie")

	if err := saveConfig(cfg); err != nil {
		return cfg, err
	}
	fmt.Printf("Saved %s\n", configPath)
	return cfg, nil
}

func getConfig() (Config, error) {
	cfg, err := loadConfig()
	if err == nil {
		if !validRegion(cfg.Region) {
			return cfg, fmt.Errorf("invalid region in config: %q", cfg.Region)
		}
		return cfg, nil
	}
	if !os.IsNotExist(err) {
		return cfg, err
	}
	return promptConfig()
}
