// This Discord bot issues warnings when specific conditions are met.
//
// Configured via a JSON file (config/config.json), it checks each guild message for specified keywords or regex patterns.
// Upon a match, it responds with a warning message.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dlclark/regexp2"
	"github.com/joho/godotenv"
)

// Config represents the configuration for the bot.
type Config struct {
	AlertRules []struct {
		Keywords                 []string `json:"keywords"`
		RegexPatterns            []string `json:"regex_patterns"`
		WarningMessage           string   `json:"warning_message"`
		ExternalLinkRequired     bool     `json:"external_link_required"`
		RequiredRoles            []string `json:"required_roles"`
		ExcludedRoles            []string `json:"excluded_roles"`
		OmitMembersOlderThanDays int      `json:"omit_members_older_than_days"`
		CompiledRegexes          []*regexp2.Regexp
	} `json:"alert_rules"`
}

// loadConfig loads the config file and returns a Config struct.
func loadConfig() (Config, error) {
	// Parse config file.
	configFile, err := os.ReadFile("config/config.json")
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v", err)
	}
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return Config{}, fmt.Errorf("error parsing config JSON: %v", err)
	}
	return config, nil
}

// hasRole checks if a user has any of the specified roles.
func hasRole(memberRoles map[string]bool, roles []string) bool {
	for _, role := range roles {
		if memberRoles[role] {
			return true
		}
	}
	return false
}

// handleMessage handles a message sent in a guild.
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, config Config, linkRegex *regexp2.Regexp) {
	// Ignore messages sent by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Convert the message content to lowercase for case-insensitive comparison.
	content := strings.ToLower(m.Content)

	// Check if the message contains any of the specified keywords or regex patterns.
	for _, alertRules := range config.AlertRules {
		// If both keywords and regex patterns are empty, skip warning.
		if len(alertRules.Keywords) == 0 && len(alertRules.RegexPatterns) == 0 {
			continue
		}

		// If warning message is empty, skip this warning.
		if alertRules.WarningMessage == "" {
			continue
		}

		// Fetch member information.
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			log.Printf("Error fetching member: %v", err)
			continue
		}

		// Check if the user's membership is older than the specified number of days.
		membershipDays := time.Since(member.JoinedAt).Hours() / 24
		if alertRules.OmitMembersOlderThanDays > 0 && membershipDays > float64(alertRules.OmitMembersOlderThanDays) {
			continue // Skip if the user's membership is too old.
		}

		// Check if external link is required and if the message contains one.
		if alertRules.ExternalLinkRequired {
			match, err := linkRegex.MatchString(content)
			if err != nil {
				log.Printf("Error matching regex: %v", err)
				continue
			}
			if !match {
				continue // Skip if external link is required but not found.
			}
		}

		// Convert member roles to a map for efficient lookup
		memberRoles := make(map[string]bool)
		for _, roleID := range member.Roles {
			role, err := s.State.Role(m.GuildID, roleID)
			if err != nil {
				log.Printf("Error fetching role: %v", err)
				continue
			}
			memberRoles[role.Name] = true
		}

		// If excluded roles are specified, check if the user has any of them.
		if len(alertRules.ExcludedRoles) > 0 && hasRole(memberRoles, alertRules.ExcludedRoles) {
			continue // User has an excluded role, skip the warning message.
		}

		// If required roles are specified, ensure the user has any of them.
		for _, regex := range alertRules.CompiledRegexes {
			match, err := regex.MatchString(content)
			if err != nil {
				log.Printf("Error matching regex: %v", err)
				continue
			}
			if match {
				// Create an embed.
				embed := &discordgo.MessageEmbed{
					Description: alertRules.WarningMessage,
					Color:       0xff0000, // Red color.
				}

				// Create a message send struct.
				messageSend := &discordgo.MessageSend{
					Embed: embed,
					Reference: &discordgo.MessageReference{
						MessageID: m.ID,
						ChannelID: m.ChannelID,
						GuildID:   m.GuildID,
					},
				}

				// Reply with the warning embed message.
				_, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend)
				// _, err := s.ChannelMessageSendReply(m.ChannelID, alertRules.WarningMessage, m.Reference())
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				return
			}
		}
	}
}

func main() {
	fmt.Println("Starting bot Warnings Bot...")

	// Parse config file and create regex patterns.
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	for i, alertRules := range config.AlertRules {
		for j := range alertRules.RequiredRoles {
			if alertRules.RequiredRoles[j] == "" {
				config.AlertRules[i].RequiredRoles = []string{}
			}
		}

		// If regex patterns are provided, compile them .
		if len(alertRules.RegexPatterns) > 0 {
			for _, pattern := range alertRules.RegexPatterns {
				compiledRegex, err := regexp2.Compile(pattern, regexp2.None)
				if err != nil {
					log.Fatalf("Error compiling regex: %v", err)
				}
				config.AlertRules[i].CompiledRegexes = append(config.AlertRules[i].CompiledRegexes, compiledRegex)
			}
		} else {
			// If no regex patterns are provided, generate a regex pattern from the keywords.
			pattern := `(?i)\b(` + strings.Join(alertRules.Keywords, "|") + `)\b`
			compiledRegex, err := regexp2.Compile(pattern, regexp2.None)
			if err != nil {
				log.Fatalf("Error compiling regex: %v", err)
			}
			config.AlertRules[i].CompiledRegexes = append(config.AlertRules[i].CompiledRegexes, compiledRegex)
		}
	}

	// Retrieve bot token.
	_ = godotenv.Load()
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Precompile the link regex.
	linkRegex, err := regexp2.Compile(`https?://[^\s/$.?#].[^\s]*`, regexp2.None)
	if err != nil {
		log.Fatalf("Error compiling regex: %v", err)
	}

	// Register the messageCreate callback.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handleMessage(s, m, config, linkRegex)
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
