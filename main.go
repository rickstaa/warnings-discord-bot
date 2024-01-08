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
	KeywordLists []struct {
		Keywords             []string `json:"keywords"`
		WarningMessage       string   `json:"warning_message"`
		ExternalLinkRequired bool     `json:"external_link_required"`
		RequiredRoles        []string `json:"required_roles"`
		ExcludedRoles        []string `json:"excluded_roles"`
		KeywordRegex         *regexp.Regexp
	} `json:"keyword_lists"`
}

const (
	LINK_REGEX = `https?://[^\s/$.?#].[^\s]*`
)

// fetchRoles fetches the roles of a member and returns them as a map for easy lookup.
func fetchRoles(s *discordgo.Session, guildID, memberID string) (map[string]bool, error) {
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return nil, fmt.Errorf("error fetching member: %w", err)
	}

	roles := make(map[string]bool)
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return nil, fmt.Errorf("error fetching role: %w", err)
		}
		roles[role.Name] = true
	}

	return roles, nil
}

func main() {
	fmt.Println("Starting bot Warnings Bot...")

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
		for j := range keywordList.RequiredRoles {
			if keywordList.RequiredRoles[j] == "" {
				config.KeywordLists[i].RequiredRoles = []string{}
			}
		}
		keywordsPattern := `(?i)\b(` + strings.Join(keywordList.Keywords, "|") + `)\b`
		config.KeywordLists[i].KeywordRegex = regexp.MustCompile(keywordsPattern)
	}

	// Retrieve bot token.
	_ = godotenv.Load()
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Precompile the link regex.
	linkRegex := regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

	// Register the messageCreate callback
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages sent by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Convert the message content to lowercase for case-insensitive comparison
		content := strings.ToLower(m.Content)

		// Check if the message contains any of the keywords from any keyword list
		for _, keywordList := range config.KeywordLists {
			// Check if external link is required and if the message contains one
			if keywordList.ExternalLinkRequired && !linkRegex.MatchString(content) {
				continue // Skip if external link is required but not found
			}

			// Fetch the roles of the message sender
			roles, err := fetchRoles(s, m.GuildID, m.Author.ID)
			if err != nil {
				log.Printf("%v", err)
				return
			}

			// Check if the user has any of the required roles
			hasRequiredRole := len(keywordList.RequiredRoles) == 0
			for _, requiredRole := range keywordList.RequiredRoles {
				if requiredRole != "" && roles[requiredRole] {
					hasRequiredRole = true
					break
				}
			}

			// Check if the user has any of the excluded roles
			hasExcludedRole := false
			for _, excludedRole := range keywordList.ExcludedRoles {
				if excludedRole != "" && roles[excludedRole] {
					hasExcludedRole = true
					break
				}
			}

			// If the user has an excluded role or doesn't have any of the required roles, skip this message
			if hasExcludedRole || !hasRequiredRole {
				continue
			}

			// Create a regular expression pattern dynamically for all keywords in the list with case-insensitivity.
			keywordsPattern := `(?i)\b(` + strings.Join(keywordList.Keywords, "|") + `)\b`
			dmRegex := regexp.MustCompile(keywordsPattern)

			// Check if the message contains any of the specified keywords dynamically.
			if dmRegex.MatchString(content) {
				// Reply with the warning message.
				_, err := s.ChannelMessageSendReply(m.ChannelID, keywordList.WarningMessage, m.Reference())
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
