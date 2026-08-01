package avito

import (
	"encoding/json"
	"fmt"
	"os"
)

type PollSource struct {
	Name         string `json:"name"`
	Source       string `json:"source"`
	AccountID    string `json:"account_id"`
	ChatID       string `json:"chat_id"`
	AccessToken  string `json:"access_token,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Enabled      bool   `json:"enabled"`
}

type SourcesConfig struct {
	Sources []PollSource `json:"sources"`
}

func LoadSourcesConfig(path string) ([]PollSource, error) {
	cfg, err := LoadFullSourcesConfig(path)
	if err != nil {
		return nil, err
	}

	return EnabledSources(cfg)
}

func LoadFullSourcesConfig(path string) (SourcesConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return SourcesConfig{}, err
	}

	var cfg SourcesConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return SourcesConfig{}, fmt.Errorf("decode avito sources config: %w", err)
	}

	return cfg, nil
}

func SaveSourcesConfig(path string, cfg SourcesConfig) error {
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode avito sources config: %w", err)
	}
	raw = append(raw, '\n')

	return os.WriteFile(path, raw, 0600)
}

func EnabledSources(cfg SourcesConfig) ([]PollSource, error) {
	var sources []PollSource
	for i, source := range cfg.Sources {
		if !source.Enabled {
			continue
		}
		if source.Name == "" {
			return nil, fmt.Errorf("source #%d: name is required", i+1)
		}
		if source.AccountID == "" {
			return nil, fmt.Errorf("source %q: account_id is required", source.Name)
		}
		if source.ChatID == "" {
			return nil, fmt.Errorf("source %q: chat_id is required", source.Name)
		}
		if source.AccessToken == "" && (source.ClientID == "" || source.ClientSecret == "") {
			return nil, fmt.Errorf("source %q: access_token or client_id/client_secret is required", source.Name)
		}
		if source.Source == "" {
			source.Source = "polling"
		}
		sources = append(sources, source)
	}

	return sources, nil
}
