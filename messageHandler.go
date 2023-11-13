package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"

	"github.com/bwmarrin/discordgo"
	fwew "github.com/fwew/fwew-lib/v5"
	"github.com/knoxfighter/dgc"
)

func sendEmbed(ctx *dgc.Ctx, title string, message string, isErr bool) *discordgo.Message {
	// set color to use
	var color int
	if isErr {
		color = 0xFF0000
	} else {
		color = 0x607CA3
	}

	// create embed to send
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.Event.Author.Username,
			IconURL: ctx.Event.Author.AvatarURL("1024"),
		},
		Title:       title,
		Color:       color,
		Description: message,
	}

	// send the Embed
	dcMessage, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embed)
	if err != nil {
		log.Printf("Something went wrong sending message to discord: %s", err)
	}

	return dcMessage
}

// func sendEmbedImage(ctx *dgc.Ctx, imageURL string) {
// 	// create embed to send
// 	embed := &discordgo.MessageEmbed{
// 		Author: &discordgo.MessageEmbedAuthor{
// 			Name:    ctx.Event.Author.Username,
// 			IconURL: ctx.Event.Author.AvatarURL("1024"),
// 		},
// 		Image: &discordgo.MessageEmbedImage{
// 			URL: imageURL,
// 		},
// 		Color: 0x607CA3,
// 	}

// 	// send the Embed
// 	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, embed)
// 	if err != nil {
// 		log.Printf("Something went wrong sending message to discord: %s", err)
// 	}
// }

// Send the message to discord within the `fwew` layout of an embed.
func sendDiscordMessageEmbed(ctx *dgc.Ctx, message string, isErr bool) {
	// create title from executed command
	title := ctx.Command.Name
	arguments := ctx.Arguments.Raw()
	if arguments != "" {
		title += " " + arguments
	}

	sendEmbed(ctx, title, message, isErr)
}

type message struct {
	message *discordgo.Message
	title   string
	curPage *int
	pages   []string
}

var messages = map[string]message{}

func sendDiscordMessagePaginated(ctx *dgc.Ctx, pages []string) {
	// create title from executed command with pages count
	titleSimple := ctx.Command.Name
	arguments := ctx.Arguments.Raw()
	if arguments != "" {
		titleSimple += " " + arguments
	}

	var title = titleSimple
	// add pages to
	if len(pages) > 1 {
		title = fmt.Sprintf(" (Page %d/%d)", 1, len(pages))
	}

	// post first page
	dcMessage := sendEmbed(ctx, title, pages[0], false)
	session := ctx.Session

	if len(pages) > 1 {
		// add arrows as reaction to pagination
		session.MessageReactionAdd(dcMessage.ChannelID, dcMessage.ID, "⬅️")
		session.MessageReactionAdd(dcMessage.ChannelID, dcMessage.ID, "➡️")

		// save message so pagination can work
		p := 1
		messages[dcMessage.ChannelID+":"+dcMessage.ID] = message{
			message: dcMessage,
			title:   titleSimple,
			pages:   pages,
			curPage: &p,
		}
	}
}

func send1dWordDiscordEmbed(ctx *dgc.Ctx, words []fwew.Word) {
	var output []string
	var outTemp string

	for j, word := range words {
		iString := strconv.Itoa(j + 1)
		line, err := word.ToOutputLine(
			iString,
			true, // were discord-bot, always with markdown
			ctx.CustomObjects.MustGet("showIPA").(bool),
			ctx.CustomObjects.MustGet("showInfix").(bool),
			ctx.CustomObjects.MustGet("showDashed").(bool),
			ctx.CustomObjects.MustGet("showInfixDots").(bool),
			ctx.CustomObjects.MustGet("showSource").(bool),
			ctx.CustomObjects.MustGet("langCode").(string),
		)
		if err != nil {
			sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error creating output line: %s", err), true)
			return
		}

		if (len(outTemp) + len(line)) > 2000 {
			// add to output
			output = append(output, outTemp)
			outTemp = ""
		}
		outTemp += line
	}
	// add last outTemp also
	output = append(output, outTemp)
	sendDiscordMessagePaginated(ctx, output)
}

func sendWordDiscordEmbed(ctx *dgc.Ctx, words [][]fwew.Word) {
	var output []string
	var outTemp string
	for i, words := range words {
		if len(words) == 0 {
			outTemp += fwew.Text("none")
		}
		for j, word := range words {
			iString := strconv.Itoa(i + 1)
			if len(words) > 1 {
				iString += "-" + strconv.Itoa(j+1)
			}
			line, err := word.ToOutputLine(
				iString,
				true, // were discord-bot, always with markdown
				ctx.CustomObjects.MustGet("showIPA").(bool),
				ctx.CustomObjects.MustGet("showInfix").(bool),
				ctx.CustomObjects.MustGet("showDashed").(bool),
				ctx.CustomObjects.MustGet("showInfixDots").(bool),
				ctx.CustomObjects.MustGet("showSource").(bool),
				ctx.CustomObjects.MustGet("langCode").(string),
			)
			if err != nil {
				sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error creating output line: %s", err), true)
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

func sendErrorWhenRecovered(ctx *dgc.Ctx) {
	sendDiscordMessageEmbed(
		ctx,
		fmt.Sprintf("PANIC!! Please report this error.\ncommand: %s\nargs: %s\nstacktrace: %s", ctx.Command.Name, ctx.Arguments.Raw(), debug.Stack()),
		true,
	)
}
