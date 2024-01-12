// Initializes a Discord bot that warns users when messages are posted that contain specified
// keywords or regex patterns.
//
// This bot can be configured via a JSON file (config/config.json), it checks each guild message
// for specified keywords or regex atterns, upon a match, it responds with a warning message. It
// can also post a welcome warning message to new members.
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
	JoinWarningMessage string `json:"join_warning_message"`
	AlertRules         []struct {
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

// compileRegex compiles a regex pattern.
func compileRegex(pattern string) (*regexp2.Regexp, error) {
	compiledRegex, err := regexp2.Compile(pattern, regexp2.None)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %v", err)
	}
	return compiledRegex, nil
}

// compileRegexPatterns compiles regex patterns from config keywords and patterns.
func compileRegexPatterns(config *Config) error {
	for i, alertRules := range config.AlertRules {
		for j := range alertRules.RequiredRoles {
			if alertRules.RequiredRoles[j] == "" {
				config.AlertRules[i].RequiredRoles = []string{}
			}
		}

		// If regex patterns are provided, compile them.
		if len(alertRules.RegexPatterns) > 0 {
			for _, pattern := range alertRules.RegexPatterns {
				compiledRegex, err := compileRegex(pattern)
				if err != nil {
					return err
				}
				config.AlertRules[i].CompiledRegexes = append(config.AlertRules[i].CompiledRegexes, compiledRegex)
			}
		} else {
			// If no regex patterns are provided, generate a regex pattern from the keywords.
			pattern := `(?i)\b(` + strings.Join(alertRules.Keywords, "|") + `)\b`
			compiledRegex, err := compileRegex(pattern)
			if err != nil {
				return err
			}
			config.AlertRules[i].CompiledRegexes = append(config.AlertRules[i].CompiledRegexes, compiledRegex)
		}
	}
	return nil
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
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, config Config) {
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
		linkRegex, err := regexp2.Compile(`https?://(?!(www\.)?discord(app)?\.com/channels/`+m.GuildID+`/)[^\s/$.?#].[^\s]*`, regexp2.None)
		if err != nil {
			log.Printf("Error compiling regex: %v", err)
			return
		}
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

// handleMemberJoin handles a new member joining the guild.
func handleMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd, config Config) {
	// Only send a DM if JoinWarningMessage is defined.
	if config.JoinWarningMessage != "" {
		// Send a DM to the new member.
		channel, err := s.UserChannelCreate(m.User.ID)
		if err != nil {
			log.Printf("Error creating DM channel: %v", err)
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, config.JoinWarningMessage)
		if err != nil {
			log.Printf("Error sending DM: %v", err)
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
	err = compileRegexPatterns(&config)
	if err != nil {
		log.Fatalf("Error compiling regex patterns: %v", err)
	}

	// Retrieve bot token.
	_ = godotenv.Load()
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Register the messageCreate callback.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		go handleMessage(s, m, config)
	})

	// Register the guildMemberAdd callback.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		go handleMemberJoin(s, m, config)
	})

	// Add intents and start the bot.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Block the main goroutine until a termination signal is received (CTRL+C).
	select {}
}
