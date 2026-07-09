package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const apiBase = "https://sg-act-public-api.hoyolab.com"
const cardBase = "https://sg-act-public-api.hoyolab.com/game_record/card/wapi/getGameRecordCard"

type apiResp struct {
	Retcode int             `json:"retcode"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type avatarEntry struct {
	ID            int    `json:"id"`
	Level         int    `json:"level"`
	Name          string `json:"name_mi18n"`
	FullName      string `json:"full_name_mi18n"`
	Rarity        string `json:"rarity"`
	Rank          int    `json:"rank"`
	IsChosen      bool   `json:"is_chosen"`
	RoleSquareURL string `json:"role_square_url"`
}

type indexData struct {
	Stats struct {
		AchievementCount   int `json:"achievement_count"`
		ClimbingTowerLayer int `json:"climbing_tower_layer"`
		ActiveDays         int `json:"active_days"`
	} `json:"stats"`
	AvatarList []avatarEntry `json:"avatar_list"`
}

type noteData struct {
	Energy struct {
		Progress struct {
			Current int `json:"current"`
		} `json:"progress"`
	} `json:"energy"`
}

type ledgerData struct {
	MonthData struct {
		List []struct {
			DataType string      `json:"data_type"`
			Count    json.Number `json:"count"`
		} `json:"list"`
	} `json:"month_data"`
}

type cardEntry struct {
	Nickname   string `json:"nickname"`
	Level      int    `json:"level"`
	GameRoleID string `json:"game_role_id"`
	GameID     int    `json:"game_id"`
}

type cardData struct {
	List []cardEntry `json:"list"`
}

func hoyoHeaders(cfg Config) http.Header {
	h := http.Header{}
	h.Set("Cookie", cfg.HoyolabCookie)
	h.Set("Accept", "application/json, text/plain, */*")
	h.Set("Accept-Language", "en-US,en;q=0.9")
	h.Set("x-rpc-client_type", "5")
	h.Set("x-rpc-language", "en-us")
	h.Set("Origin", "https://act.hoyolab.com")
	h.Set("Referer", "https://act.hoyolab.com/")
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return h
}

func getHoyo(client *http.Client, url string, cfg Config, out any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header = hoyoHeaders(cfg)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var wrap apiResp
	if err := json.Unmarshal(body, &wrap); err != nil {
		return fmt.Errorf("bad response from %s: %s", url, string(body))
	}
	if wrap.Retcode != 0 {
		return fmt.Errorf("hoyolab error %d: %s", wrap.Retcode, wrap.Message)
	}
	return json.Unmarshal(wrap.Data, out)
}

func statURLs(cfg Config) (string, string, string, string) {
	month := time.Now().Format("200601")
	index := fmt.Sprintf("%s/event/game_record_zzz/api/zzz/index?server=%s&role_id=%s", apiBase, cfg.Region, cfg.ZenlessUID)
	ledger := fmt.Sprintf("%s/event/nap_ledger/month_info?uid=%s&region=%s&month=%s", apiBase, cfg.ZenlessUID, cfg.Region, month)
	note := fmt.Sprintf("%s/event/game_record_zzz/api/zzz/note?server=%s&role_id=%s", apiBase, cfg.Region, cfg.ZenlessUID)
	card := fmt.Sprintf("%s?uid=%s", cardBase, cfg.HoyolabUID)
	return index, ledger, note, card
}

func polychromes(l ledgerData) string {
	for _, item := range l.MonthData.List {
		if item.DataType == "PolychromesData" {
			return item.Count.String()
		}
	}
	return "0"
}

func zenlessCard(cfg Config, card cardData) (cardEntry, bool) {
	for _, c := range card.List {
		if c.GameRoleID == cfg.ZenlessUID || c.GameID == 8 {
			return c, true
		}
	}
	if len(card.List) > 0 {
		return card.List[0], true
	}

	return cardEntry{}, false
}

func pickCharacter(cfg Config, idx indexData) (avatarEntry, bool) {
	if len(idx.AvatarList) == 0 {
		return avatarEntry{}, false
	}
	if cfg.Character != "" {
		for _, a := range idx.AvatarList {
			if strings.EqualFold(a.Name, cfg.Character) || strings.EqualFold(a.FullName, cfg.Character) || strconv.Itoa(a.ID) == cfg.Character {
				return a, true
			}
		}
	}
	for _, a := range idx.AvatarList {
		if a.IsChosen {
			return a, true
		}
	}
	return idx.AvatarList[0], true
}

func buildPayload(cfg Config, idx indexData, note noteData, led ledgerData, card cardData) map[string]any {
	nickname := "Unknown Proxy"
	level := 0
	if c, ok := zenlessCard(cfg, card); ok {
		if c.Nickname != "" {
			nickname = c.Nickname
		}
		level = c.Level
	}
	poly := polychromes(led)

	char, _ := pickCharacter(cfg, idx)
	charName := char.Name + " main"
	charStats := fmt.Sprintf("LVL %d, S%d", char.Level, char.Rank)

	dynamic := []map[string]any{
		{"type": 3, "name": "avatar", "value": map[string]any{"url": char.RoleSquareURL}},
		{"type": 1, "name": "char", "value": charName},
		{"type": 1, "name": "char_2", "value": charStats},
		{"type": 1, "name": "nickname", "value": nickname},
		{"type": 1, "name": "uid", "value": "UID: " + cfg.ZenlessUID},
		{"type": 1, "name": "polychromes", "value": poly},
		{"type": 1, "name": "IL_str", "value": "Interknot Level"},
		{"type": 1, "name": "IL", "value": strconv.Itoa(level)},
		{"type": 1, "name": "ach_str", "value": "Achievements"},
		{"type": 1, "name": "ach", "value": idx.Stats.AchievementCount},
		{"type": 1, "name": "SBT_str", "value": "Simulated Battle Trial"},
		{"type": 1, "name": "SBT", "value": idx.Stats.ClimbingTowerLayer},
		{"type": 1, "name": "ENER_str", "value": "Energy"},
		{"type": 1, "name": "ENER", "value": note.Energy.Progress.Current},
		{"type": 1, "name": "days_str", "value": "Days Active"},
		{"type": 1, "name": "days", "value": idx.Stats.ActiveDays},
		{"type": 1, "name": "proxylevel", "value": "Legendary Proxy"},
		{"type": 1, "name": "polychromes_str", "value": "Monthly Polychromes"},
		{"type": 1, "name": "mini", "value": fmt.Sprintf("%s: IL %d", nickname, level)},
	}

	return map[string]any{
		"username": nickname,
		"data":     map[string]any{"dynamic": dynamic},
	}
}

func sendDiscord(client *http.Client, cfg Config, payload map[string]any) (int, string, error) {
	url := fmt.Sprintf("https://discord.com/api/v9/applications/%s/users/%s/identities/0/profile", cfg.DiscordClientID, cfg.DiscordUserID)
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Authorization", "Bot "+cfg.DiscordBotToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	msg, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, string(msg), fmt.Errorf("discord error %d: %s", resp.StatusCode, string(msg))
	}
	return resp.StatusCode, string(msg), nil
}

func writeResults(result map[string]any) error {
	if err := os.MkdirAll("json", 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("json/results.json", data, 0644)
}

func fetchCharacters(cfg Config) ([]avatarEntry, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	indexURL, _, _, _ := statURLs(cfg)
	var idx indexData
	if err := getHoyo(client, indexURL, cfg, &idx); err != nil {
		return nil, err
	}
	return idx.AvatarList, nil
}

func listCharacters(cfg Config) error {
	chars, err := fetchCharacters(cfg)
	if err != nil {
		return fmt.Errorf("index: %w", err)
	}
	if len(chars) == 0 {
		fmt.Println("no characters found")
		return nil
	}
	fmt.Printf("%-20s %-6s %-5s %-6s %s\n", "NAME", "RARITY", "LEVEL", "ID", "PORTRAIT")
	for _, a := range chars {
		fmt.Printf("%-20s %-6s %-5d %-6d %s\n", a.Name, a.Rarity, a.Level, a.ID, a.RoleSquareURL)
	}
	fmt.Println("\nput one of the NAME or ID values into config.json as \"character\" to feature it.")
	return nil
}

func syncStats(cfg Config) error {
	client := &http.Client{Timeout: 30 * time.Second}
	indexURL, ledgerURL, noteURL, cardURL := statURLs(cfg)

	var idx indexData
	var led ledgerData
	var note noteData
	var card cardData

	if err := getHoyo(client, indexURL, cfg, &idx); err != nil {
		return fmt.Errorf("index: %w", err)
	}
	if err := getHoyo(client, ledgerURL, cfg, &led); err != nil {
		return fmt.Errorf("ledger: %w", err)
	}
	if err := getHoyo(client, noteURL, cfg, &note); err != nil {
		return fmt.Errorf("note: %w", err)
	}
	if err := getHoyo(client, cardURL, cfg, &card); err != nil {
		return fmt.Errorf("card: %w", err)
	}

	payload := buildPayload(cfg, idx, note, led, card)
	status, response, sendErr := sendDiscord(client, cfg, payload)

	result := map[string]any{
		"updated_at":       time.Now().Format(time.RFC3339),
		"region":           cfg.Region,
		"index":            idx,
		"note":             note,
		"ledger":           led,
		"card":             card,
		"payload":          payload,
		"discord_status":   status,
		"discord_response": response,
	}
	if sendErr != nil {
		result["error"] = sendErr.Error()
	}
	if err := writeResults(result); err != nil {
		log.Printf("results: %v", err)
	}

	if sendErr != nil {
		return fmt.Errorf("discord: %w", sendErr)
	}
	return nil
}
