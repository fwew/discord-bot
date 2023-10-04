package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	fwew "github.com/fwew/fwew-lib/v5"
	"github.com/knoxfighter/dgc"
)

// Config holds configuration options
type Config struct {
	Token     string   `json:"token"`
	Prefixes  []string `json:"prefixes"`
	AdminRole string   `json:"admin_role"`
}

var config Config
var statisticsDir string

func init() {
	// get working dir
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	statisticsDir = filepath.Join(wd, "statistics")

	// assure statistics directory next to the executable
	if _, err := os.Stat(statisticsDir); os.IsNotExist(err) {
		os.Mkdir(statisticsDir, os.ModeDir|0o755)
	}
}

func removeAllReactions(session *discordgo.Session) {
	for _, m := range messages {
		session.MessageReactionsRemoveAll(m.message.ChannelID, m.message.ID)
	}
}

func main() {
	// load json config
	jsonFile, err := ioutil.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		panic(err)
	}

	// Assure Dict, before caching it
	err = fwew.AssureDict()
	if err != nil {
		panic(err)
	}

	// Look up phoneme frequencies once for the phoneme-frequency command
	fwew.PhonemeDistros()

	// cache fwew dictionary
	err = fwew.CacheDictHash()
	if err != nil {
		panic(err)
	}

	// create discord session
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	// open the session and connect to discord
	err = session.Open()
	if err != nil {
		panic(err)
	}

	// set custom status
	err = session.UpdateStatus(0, "DM me to look up Na'vi words")
	if err != nil {
		panic(err)
	}

	// generate a command router
	router := dgc.Create(&dgc.Router{
		Prefixes:         append(config.Prefixes, "<@!"+session.State.User.ID+">"),
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*dgc.Command{},
		//Middlewares:      []dgc.Middleware{},

		// The ping handler will be executed if the message only contains the bot's mention (no arguments)
		PingHandler: func(ctx *dgc.Ctx) {
			sendDiscordMessageEmbed(ctx, "Pong!", false)
		},
	})

	registerCommands(router)

	addMiddleware(router)

	router.RegisterDefaultHelpCommand(session, nil, 0xF1C40E)

	router.Initialize(session)

	// Add Handler for Reactions added
	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		// don't run, when reaction is from myself
		if event.UserID == session.State.User.ID {
			return
		}

		// check message is paginated
		m, ok := messages[event.ChannelID+":"+event.MessageID]
		if !ok {
			return
		}

		// Check which reaction was added
		reactionName := event.Emoji.Name
		switch reactionName {
		case "⬅️":
			// calculate new page
			*m.curPage = (*m.curPage-1+len(m.pages)-1)%len(m.pages) + 1

			// set new stuff to embed
			messageEmbed := m.message.Embeds[0]
			messageEmbed.Title = fmt.Sprintf("%s (Page %d/%d)", m.title, *m.curPage, len(m.pages))
			messageEmbed.Description = m.pages[*m.curPage-1]

			// edit the embed to the new one
			session.ChannelMessageEditEmbed(event.ChannelID, event.MessageID, messageEmbed)

			// Remove the reaction
			session.MessageReactionRemove(event.ChannelID, event.MessageID, reactionName, event.UserID)
		case "➡️":
			// calculate new page
			*m.curPage = (*m.curPage % len(m.pages)) + 1

			// set new stuff to embed
			messageEmbed := m.message.Embeds[0]
			messageEmbed.Title = fmt.Sprintf("%s (Page %d/%d)", m.title, *m.curPage, len(m.pages))
			messageEmbed.Description = m.pages[*m.curPage-1]

			// edit the embed to the new one
			session.ChannelMessageEditEmbed(event.ChannelID, event.MessageID, messageEmbed)

			// Remove the reaction
			session.MessageReactionRemove(event.ChannelID, event.MessageID, reactionName, event.UserID)
		}
	})

	// Add handler for private chat
	session.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		// do not run on my own messages!
		if message.Author.Bot {
			return
		}

		// if in private chat
		channel, err := session.State.Channel(message.ChannelID)
		if err != nil {
			if channel, err = session.Channel(message.ChannelID); err != nil {
				return
			}
		}

		if channel.Type == discordgo.ChannelTypeDM {
			// only run if message not starts with prefix
			for _, prefix := range router.Prefixes {
				if strings.HasPrefix(message.Message.Content, prefix) {
					return
				}
			}

			var fwewCommand *dgc.Command
			for _, command := range router.Commands {
				if command.Name == "fwew" {
					fwewCommand = command
					break
				}
			}

			// use this message as params
			context := &dgc.Ctx{
				Session:       session,
				Event:         message,
				Arguments:     dgc.ParseArguments(message.Content),
				CustomObjects: dgc.NewObjectsMap(),
				Router:        router,
				Command:       fwewCommand,
			}

			fwewCommand.Trigger(context)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Fwew is now running.  Press CTRL-C or send Sigterm/Sigkill to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// cleanup pagination reactions
	removeAllReactions(session)

	// Cleanly close down the Discord session.
	session.Close()
}
