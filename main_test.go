package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/robfig/cron/v3"
)

func testConfig() Config {
	return Config{
		DiscordUserID:   "111",
		HoyolabUID:      "222",
		ZenlessUID:      "333",
		DiscordClientID: "444",
		DiscordBotToken: "token",
		HoyolabCookie:   "cookie",
		Region:          "prod_gf_eu",
	}
}

func TestValidRegion(t *testing.T) {
	for _, r := range regions {
		if !validRegion(r.Code) {
			t.Errorf("expected %q to be valid", r.Code)
		}
	}
	if validRegion("prod_gf_xx") {
		t.Error("expected prod_gf_xx to be invalid")
	}
}

func TestStatURLs(t *testing.T) {
	cfg := testConfig()
	index, ledger, note, card := statURLs(cfg)
	if !strings.Contains(index, "server=prod_gf_eu") || !strings.Contains(index, "role_id=333") {
		t.Errorf("bad index url: %s", index)
	}
	if !strings.Contains(ledger, "region=prod_gf_eu") || !strings.Contains(ledger, "uid=333") {
		t.Errorf("bad ledger url: %s", ledger)
	}
	if !strings.Contains(note, "server=prod_gf_eu") {
		t.Errorf("bad note url: %s", note)
	}
	if !strings.Contains(card, "uid=222") {
		t.Errorf("bad card url: %s", card)
	}
}

func ledgerWith(count string) ledgerData {
	var led ledgerData
	json.Unmarshal([]byte(`{"month_data":{"list":[{"data_type":"OtherData","count":10},{"data_type":"PolychromesData","count":`+count+`}]}}`), &led)
	return led
}

func TestPolychromes(t *testing.T) {
	if got := polychromes(ledgerWith("1600")); got != "1600" {
		t.Errorf("expected 1600, got %s", got)
	}
	if got := polychromes(ledgerData{}); got != "0" {
		t.Errorf("expected 0, got %s", got)
	}
}

func TestBuildPayload(t *testing.T) {
	cfg := testConfig()
	var idx indexData
	idx.Stats.AchievementCount = 42
	idx.Stats.ClimbingTowerLayer = 7
	idx.Stats.ActiveDays = 100

	var note noteData
	note.Energy.Progress.Current = 120

	led := ledgerWith("800")

	var card cardData
	card.List = append(card.List, cardEntry{Nickname: "Proxy", Level: 55, GameRoleID: cfg.ZenlessUID})

	payload := buildPayload(cfg, idx, note, led, card)
	if payload["username"] != "Proxy" {
		t.Errorf("expected username Proxy, got %v", payload["username"])
	}
	dyn := payload["data"].(map[string]any)["dynamic"].([]map[string]any)
	if len(dyn) != 16 {
		t.Fatalf("expected 16 fields, got %d", len(dyn))
	}
	find := func(name string) any {
		for _, f := range dyn {
			if f["name"] == name {
				return f["value"]
			}
		}
		return nil
	}
	if find("ach") != 42 {
		t.Errorf("expected ach 42, got %v", find("ach"))
	}
	if find("IL") != 55 {
		t.Errorf("expected IL 55, got %v", find("IL"))
	}
	if find("ENER") != 120 {
		t.Errorf("expected ENER 120, got %v", find("ENER"))
	}
	if find("polychromes") != "800" {
		t.Errorf("expected polychromes 800, got %v", find("polychromes"))
	}
	if find("mini") != "Proxy: IL 55" {
		t.Errorf("expected mini, got %v", find("mini"))
	}
}

func TestBuildPayloadDefaults(t *testing.T) {
	payload := buildPayload(testConfig(), indexData{}, noteData{}, ledgerData{}, cardData{})
	if payload["username"] != "Unknown Proxy" {
		t.Errorf("expected Unknown Proxy, got %v", payload["username"])
	}
}

func TestZenlessCard(t *testing.T) {
	cfg := testConfig()
	card := cardData{List: []cardEntry{
		{Nickname: "other", Level: 60, GameRoleID: "999", GameID: 2},
		{Nickname: "me", Level: 56, GameRoleID: cfg.ZenlessUID, GameID: 8},
	}}
	c, ok := zenlessCard(cfg, card)
	if !ok || c.Nickname != "me" {
		t.Errorf("expected to select zenless card, got %+v ok=%v", c, ok)
	}
	if _, ok := zenlessCard(cfg, cardData{}); ok {
		t.Error("expected no card for empty list")
	}
}

func TestDefaultCronValid(t *testing.T) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(defaultCron); err != nil {
		t.Errorf("default cron invalid: %v", err)
	}
}
