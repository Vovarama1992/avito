package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Vovarama1992/avito/internal/avito"
)

func runSourcesCLI(args []string) int {
	if len(args) == 0 {
		printSourcesUsage()
		return 2
	}

	switch args[0] {
	case "list":
		return sourcesList(args[1:])
	case "add":
		return sourcesAdd(args[1:])
	case "enable":
		return sourcesSetEnabled(args[1:], true)
	case "disable":
		return sourcesSetEnabled(args[1:], false)
	case "remove":
		return sourcesRemove(args[1:])
	default:
		printSourcesUsage()
		return 2
	}
}

func printSourcesUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  avito-monitor sources list [-config avito_sources.json]")
	fmt.Fprintln(os.Stderr, "  avito-monitor sources add -name NAME -account-id ID -chat-id ID -client-id ID -client-secret SECRET [-source NAME] [-config FILE]")
	fmt.Fprintln(os.Stderr, "  avito-monitor sources enable -name NAME [-config FILE]")
	fmt.Fprintln(os.Stderr, "  avito-monitor sources disable -name NAME [-config FILE]")
	fmt.Fprintln(os.Stderr, "  avito-monitor sources remove -name NAME [-config FILE]")
}

func sourcesList(args []string) int {
	fs := flag.NewFlagSet("sources list", flag.ContinueOnError)
	configPath := fs.String("config", defaultSourcesConfigPath(), "sources config path")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	cfg, err := loadOrEmptySourcesConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ENABLED\tNAME\tSOURCE\tACCOUNT_ID\tCHAT_ID")
	for _, source := range cfg.Sources {
		fmt.Fprintf(w, "%t\t%s\t%s\t%s\t%s\n", source.Enabled, source.Name, source.Source, source.AccountID, source.ChatID)
	}
	_ = w.Flush()
	return 0
}

func sourcesAdd(args []string) int {
	fs := flag.NewFlagSet("sources add", flag.ContinueOnError)
	configPath := fs.String("config", defaultSourcesConfigPath(), "sources config path")
	name := fs.String("name", "", "source display name")
	sourceName := fs.String("source", "polling: Проверка транспорта", "Matrix source label")
	accountID := fs.String("account-id", "", "Avito account_id")
	chatID := fs.String("chat-id", "", "Avito chat_id")
	accessToken := fs.String("access-token", "", "Avito access token")
	clientID := fs.String("client-id", "", "Avito client_id")
	clientSecret := fs.String("client-secret", "", "Avito client_secret")
	disabled := fs.Bool("disabled", false, "add source disabled")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	newSource := avito.PollSource{
		Name:         *name,
		Source:       *sourceName,
		AccountID:    *accountID,
		ChatID:       *chatID,
		AccessToken:  *accessToken,
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Enabled:      !*disabled,
	}

	if err := validateSource(newSource); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	cfg, err := loadOrEmptySourcesConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	replaced := false
	for i, source := range cfg.Sources {
		if source.Name == newSource.Name {
			cfg.Sources[i] = newSource
			replaced = true
			break
		}
	}
	if !replaced {
		cfg.Sources = append(cfg.Sources, newSource)
	}

	if err := avito.SaveSourcesConfig(*configPath, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if replaced {
		fmt.Println("source updated:", newSource.Name)
	} else {
		fmt.Println("source added:", newSource.Name)
	}
	return 0
}

func sourcesSetEnabled(args []string, enabled bool) int {
	cmdName := "sources disable"
	if enabled {
		cmdName = "sources enable"
	}
	fs := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	configPath := fs.String("config", defaultSourcesConfigPath(), "sources config path")
	name := fs.String("name", "", "source display name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(os.Stderr, "-name is required")
		return 2
	}

	cfg, err := loadOrEmptySourcesConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for i, source := range cfg.Sources {
		if source.Name == *name {
			cfg.Sources[i].Enabled = enabled
			if err := avito.SaveSourcesConfig(*configPath, cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			fmt.Printf("source %s: %s\n", statusWord(enabled), source.Name)
			return 0
		}
	}

	fmt.Fprintln(os.Stderr, "source not found:", *name)
	return 1
}

func sourcesRemove(args []string) int {
	fs := flag.NewFlagSet("sources remove", flag.ContinueOnError)
	configPath := fs.String("config", defaultSourcesConfigPath(), "sources config path")
	name := fs.String("name", "", "source display name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(os.Stderr, "-name is required")
		return 2
	}

	cfg, err := loadOrEmptySourcesConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for i, source := range cfg.Sources {
		if source.Name == *name {
			cfg.Sources = append(cfg.Sources[:i], cfg.Sources[i+1:]...)
			if err := avito.SaveSourcesConfig(*configPath, cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			fmt.Println("source removed:", source.Name)
			return 0
		}
	}

	fmt.Fprintln(os.Stderr, "source not found:", *name)
	return 1
}

func defaultSourcesConfigPath() string {
	if path := os.Getenv("AVITO_SOURCES_CONFIG"); path != "" {
		return path
	}
	return "avito_sources.json"
}

func loadOrEmptySourcesConfig(path string) (avito.SourcesConfig, error) {
	cfg, err := avito.LoadFullSourcesConfig(path)
	if err == nil {
		return cfg, nil
	}
	if os.IsNotExist(err) {
		return avito.SourcesConfig{}, nil
	}
	return avito.SourcesConfig{}, err
}

func validateSource(source avito.PollSource) error {
	if source.Name == "" {
		return fmt.Errorf("-name is required")
	}
	if source.AccountID == "" {
		return fmt.Errorf("-account-id is required")
	}
	if source.ChatID == "" {
		return fmt.Errorf("-chat-id is required")
	}
	if source.AccessToken == "" && (source.ClientID == "" || source.ClientSecret == "") {
		return fmt.Errorf("-client-id and -client-secret are required unless -access-token is set")
	}
	return nil
}

func statusWord(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}
