package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Config represents the configuration for the bot.
type Config struct {
	Token        string `json:"Token"`
	KeywordLists []struct {
		Keywords             []string `json:"Keywords"`
		WarningMessage       string   `json:"WarningMessage"`
		ExternalLinkRequired bool     `json:"ExternalLinkRequired"`
		KeywordRegex         *regexp.Regexp
	} `json:"KeywordLists"`
}

func main() {
	fmt.Println("Starting bot Warnings Bot...")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse config file.
	configFile, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config JSON: %v", err)
	}
	for i, keywordList := range config.KeywordLists {
		keywordsPattern := `(?i)\b(` + strings.Join(keywordList.Keywords, "|") + `)\b`
		config.KeywordLists[i].KeywordRegex = regexp.MustCompile(keywordsPattern)
	}

	// Retrieve bot token.
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		token = config.Token
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Add handler for messages.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		content := strings.ToLower(m.Content)

		for _, keywordList := range config.KeywordLists {
			linkRegex := regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

			if keywordList.ExternalLinkRequired && !linkRegex.MatchString(content) {
				continue
			}

			if keywordList.KeywordRegex.MatchString(content) {
				_, err := s.ChannelMessageSend(m.ChannelID, keywordList.WarningMessage)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				return
			}
		}
	})

	// Start the bot.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Block the main goroutine until a termination signal is received (CTRL+C).
	select {}
}
