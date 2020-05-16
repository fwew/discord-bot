package main

import (
	"encoding/json"
	"fmt"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	fwew "github.com/fwew/fwew_lib"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func sendWordDiscordEmbed(ctx *dgc.Ctx, words [][]fwew.Word) {
	var output []string
	var outTemp string
	for _, worda := range words {
		for i, word := range worda {
			line, err := word.ToOutputLine(
				i,
				true, // were discord-bot, always with markdown
				ctx.CustomObjects["showIPA"].(bool),
				ctx.CustomObjects["showInfix"].(bool),
				ctx.CustomObjects["showDashed"].(bool),
				ctx.CustomObjects["showInfixDots"].(bool),
				ctx.CustomObjects["showSource"].(bool),
			)
			if err != nil {
				sendDiscordMessage(ctx, fmt.Sprintf("Error creating output line: %s", err))
				return
			}

			if (len(outTemp) + len(line)) > 2000 {
				// add to output
				output = append(output, outTemp)
				outTemp = ""
			}
			outTemp += line
		}
		outTemp += "\n"
	}
	// add last outTemp also
	output = append(output, outTemp)
	sendDiscordMessagePaginated(ctx, output)
}

type Config struct {
	Token string `json:"token"`
}

var config Config

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

	// cache fwew dictionary
	err = fwew.CacheDict()
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

	// generate a command router
	router := dgc.Create(&dgc.Router{
		Prefixes: []string{
			"$",
			"<@!" + session.State.User.ID + ">",
		},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*dgc.Command{},
		Middlewares:      map[string][]dgc.Middleware{},

		// The ping handler will be executed if the message only contains the bot's mention (no arguments)
		PingHandler: func(ctx *dgc.Ctx) {
			sendDiscordMessage(ctx, "Pong!")
		},
	})

	registerCommands(router)

	addMiddleware(router)

	router.RegisterDefaultHelpCommand(session)

	router.Initialize(session)

	// Add Handler for Reactions added
	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		// dont run, when reaction is from myself
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
			*m.curPage = (*m.curPage-1+1)%(len(m.pages)) + 1
			//m.curPage = (m.curPage-1)%len(m.pages) + 1

			// set new stuff to embed
			messageEmbed := m.message.Embeds[0]
			messageEmbed.Title = fmt.Sprintf("%s (Page%d/%d)", m.title, *m.curPage, len(m.pages))
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
			messageEmbed.Title = fmt.Sprintf("%s (%d/%d)", m.title, *m.curPage, len(m.pages))
			messageEmbed.Description = m.pages[*m.curPage-1]

			// edit the embed to the new one
			session.ChannelMessageEditEmbed(event.ChannelID, event.MessageID, messageEmbed)

			// Remove the reaction
			session.MessageReactionRemove(event.ChannelID, event.MessageID, reactionName, event.UserID)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Fwew is now running.  Press CTRL-C or send Sigterm/Sigkill to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}
