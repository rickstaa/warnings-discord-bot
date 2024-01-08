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
		KeywordRegex         *regexp.Regexp
	} `json:"keyword_lists"`
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
		for j := range keywordList.RequiredRoles {
			if keywordList.RequiredRoles[j] == "" {
				config.KeywordLists[i].RequiredRoles = []string{}
			}
		}
		keywordsPattern := `(?i)\b(` + strings.Join(keywordList.Keywords, "|") + `)\b`
		config.KeywordLists[i].KeywordRegex = regexp.MustCompile(keywordsPattern)
	}

	// Retrieve bot token.
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

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
			linkRegex := regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

			// Check if external link is required and if the message contains one
			if keywordList.ExternalLinkRequired && !linkRegex.MatchString(content) {
				continue // Skip if external link is required but not found
			}

			// Check if RequiredRoles is empty or if the user has one of the required roles
			if len(keywordList.RequiredRoles) == 0 {
				// If RequiredRoles is empty, the bot responds to all messages
				// Continue checking for keywords and sending warnings
			} else {
				// Fetch the roles of the message sender
				member, err := s.GuildMember(m.GuildID, m.Author.ID)
				if err != nil {
					log.Printf("Error fetching member: %v", err)
					continue
				}

				// Check if the message sender has any of the required roles
				hasRequiredRole := false
				for _, roleID := range member.Roles {
					role, err := s.State.Role(m.GuildID, roleID)
					if err != nil {
						log.Printf("Error fetching role: %v", err)
						continue
					}

					// Check if the role name matches any of the required roles
					for _, requiredRole := range keywordList.RequiredRoles {
						if role.Name == requiredRole {
							hasRequiredRole = true
							break
						}
					}
					if hasRequiredRole {
						break // User has a required role, no need to check further
					}
				}

				// If the user doesn't have any of the required roles, skip this message
				if !hasRequiredRole {
					continue
				}
			}

			// Create a regular expression pattern dynamically for all keywords in the list with case-insensitivity
			keywordsPattern := `(?i)\b(` + strings.Join(keywordList.Keywords, "|") + `)\b`
			dmRegex := regexp.MustCompile(keywordsPattern)

			// Check if the message contains any of the specified keywords dynamically
			if dmRegex.MatchString(content) {
				// Reply with the warning message
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
