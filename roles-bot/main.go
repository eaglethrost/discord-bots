package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot Token
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create Discord Session, API Client
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating bot: ", err.Error())
		return
	}

	// Register function to print member roles
	dg.AddHandler(printRoles)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Connect to DiscordBot Socket and start listeing
	err = dg.Open()
	if err != nil {
		fmt.Println("Error listening to bot: ", err.Error())
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// Function is called everytime message is sent
func printRoles(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Get all the roles in this server
	if m.Content == "!roles" {

		// Get Guild ID from message
		guildID := m.GuildID

		// Get members of Guild
		members, err := s.GuildMembers(guildID, "0", 1000)
		if err != nil {
			fmt.Println("Failed to get members")
			return
		}

		// Map Role ID to its names
		rolesDefinition := make(map[string]string)
		roles, err := s.GuildRoles(guildID)
		if err != nil {
			fmt.Println("Failed to get roles")
			return
		}
		for _, role := range roles {
			rolesDefinition[role.ID] = role.Name
		}

		// Map members to each role they have
		memberRoles := make(map[string][]string)
		for _, member := range members {
			for _, role := range member.Roles {
				// Since the role of members is in ID form, we need to translate it to its name
				roleName := rolesDefinition[role]
				memberRoles[roleName] = append(memberRoles[roleName], member.User.Username)
			}
		}

		// Print Roles and send it back to discord chat
		channelID := m.ChannelID
		var message = ""
		for role, members := range memberRoles {
			message += fmt.Sprintf("%v:\n", role)
			for _, member := range members {
				message += fmt.Sprintf("- %v\n", member)
			}
			message += fmt.Sprintf("\n")
		}
		_, err = s.ChannelMessageSend(channelID, message)
		if err != nil {
			fmt.Println("Failed to send message")
			return
		}

		// TODO: Add custom HTML
	}
}
