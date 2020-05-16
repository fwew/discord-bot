package main

import (
	"fmt"
	"github.com/Lukaesebrot/dgc"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func addMiddleware(router *dgc.Router) {
	// add middleware to parse additional params (-r -l=de and more)
	router.AddMiddleware("params", func(ctx *dgc.Ctx) bool {
		amount := ctx.Arguments.Amount()

		// set up default values of params
		ctx.CustomObjects["langCode"] = "en"
		ctx.CustomObjects["reverse"] = false      // translate from navi to locale
		ctx.CustomObjects["showInfix"] = false    // dont show Infix data
		ctx.CustomObjects["showInfixDots"] = true // dont show infix data dotted
		ctx.CustomObjects["showSource"] = false   // dont show source
		ctx.CustomObjects["showDashed"] = true    // dont show syllable stress
		ctx.CustomObjects["showIPA"] = false      // dont show IPA data

		// read the real values from the user input
		for i := 0; i < amount; i++ {
			argument := ctx.Arguments.Get(i)
			arg := argument.Raw()
			if arg == "-r" {
				// mark as reverse (local to navi)
				ctx.CustomObjects["reverse"] = true
			} else if strings.HasPrefix(arg, "-l=") {
				ctx.CustomObjects["langCode"] = strings.TrimPrefix(arg, "-l=")
			} else if arg == "-i" {
				ctx.CustomObjects["showInfix"] = true
			} else if arg == "-id=false" {
				ctx.CustomObjects["showInfixDots"] = false
			} else if arg == "-src" {
				ctx.CustomObjects["showSource"] = true
			} else if arg == "-ipa" {
				ctx.CustomObjects["showIPA"] = true
			} else if arg == "-s=false" {
				ctx.CustomObjects["showDashed"] = false
			} else if strings.HasPrefix(arg, "-") {
				// ignore every other parameter
			} else {
				ctx.CustomObjects["firstArg"] = i
				break
			}
		}

		return true
	})

	// check if user is allowed to use this command (developer@Fwew Bot Testing discord)
	router.AddMiddleware("admin", func(ctx *dgc.Ctx) bool {
		author := ctx.Event.Author.ID
		guild := ctx.Event.GuildID

		member, err := ctx.Session.State.Member(guild, author)
		if err != nil {
			sendDiscordMessage(ctx, fmt.Sprintf("Couldnt get member from guild: %s", err))
			return false
		}

		for _, role := range member.Roles {
			if role == "396942792892481536" {
				// user is allowed to do it!
				log.Printf("User is allowed to do that :)")
				sendDiscordMessage(ctx, "You are not allowed to use this command!")
				return true
			}
		}

		return false
	})

	// write command and params in statistics file
	router.AddMiddleware("statistic", func(ctx *dgc.Ctx) bool {
		go func() {
			// one file for every command
			filename := filepath.Join("statistics", ctx.Command.Name+".log")

			// open statistics file to append call
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("error opening statistics.log: %s\n", err)
				return
			}
			defer file.Close()

			// only save Arguments to statistics file
			output := ctx.Arguments.Raw() + "\n"

			if _, err = file.WriteString(output); err != nil {
				log.Printf("Error writing string to statistics.log: %s\n", err)
				return
			}
		}()
		return true
	})
}
