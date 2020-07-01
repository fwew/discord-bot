package main

import (
	"fmt"
	"github.com/knoxfighter/dgc"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func addMiddleware(router *dgc.Router) {
	// add middleware to parse additional params (-r -l=de and more)
	router.RegisterMiddleware(dgc.Middleware(func(following dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			execute := false
			for _, flag := range ctx.Command.Flags {
				if flag == "params" {
					execute = true
				}
			}

			if !execute {
				following(ctx)
				return
			}

			amount := ctx.Arguments.Amount()

			// set up default values of params
			ctx.CustomObjects.Set("langCode", "en")
			ctx.CustomObjects.Set("reverse", false)      // translate from navi to locale
			ctx.CustomObjects.Set("showInfix", false)    // dont show Infix data
			ctx.CustomObjects.Set("showInfixDots", true) // dont show infix data dotted
			ctx.CustomObjects.Set("showSource", false)   // dont show source
			ctx.CustomObjects.Set("showDashed", true)    // dont show syllable stress
			ctx.CustomObjects.Set("showIPA", false)      // dont show IPA data

			var nextLanguage, nextInfixDots, nextDashed bool
			// read the real values from the user input
			for i := 0; i < amount; i++ {
				argument := ctx.Arguments.Get(i)
				arg := argument.Raw()
				if arg == "-r" {
					// mark as reverse (local to navi)
					ctx.CustomObjects.Set("reverse", true)
				} else if strings.HasPrefix(arg, "-l=") {
					ctx.CustomObjects.Set("langCode", strings.TrimPrefix(arg, "-l="))
				} else if arg == "-l" {
					// next arg is language code
					nextLanguage = true
				} else if arg == "-i" {
					ctx.CustomObjects.Set("showInfix", true)
				} else if arg == "-id=false" {
					ctx.CustomObjects.Set("showInfixDots", false)
				} else if arg == "-id" {
					// next is infix dots
					nextInfixDots = true
				} else if arg == "-src" {
					ctx.CustomObjects.Set("showSource", true)
				} else if arg == "-ipa" {
					ctx.CustomObjects.Set("showIPA", true)
				} else if arg == "-s=false" {
					ctx.CustomObjects.Set("showDashed", false)
				} else if arg == "-s" {
					// next is dashed
					nextDashed = true
				} else if strings.HasPrefix(arg, "-") {
					// ignore every other parameter
				} else if nextLanguage {
					ctx.CustomObjects.Set("langCode", arg)
					nextLanguage = false
				} else if nextInfixDots {
					if arg == "true" {
						ctx.CustomObjects.Set("showInfixDots", true)
					} else if arg == "false" {
						ctx.CustomObjects.Set("showInfixDots", false)
					}
					nextInfixDots = false
				} else if nextDashed {
					if arg == "true" {
						ctx.CustomObjects.Set("showDashed", true)
					} else if arg == "false" {
						ctx.CustomObjects.Set("showDashed", false)
					}
					nextDashed = false
				} else {
					ctx.CustomObjects.Set("firstArg", i)
					break
				}
			}

			following(ctx)
		}
	}))

	// check if user is allowed to use this command (developer@Fwew Bot Testing discord)
	router.RegisterMiddleware(func(following dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			execute := false
			for _, flag := range ctx.Command.Flags {
				if flag == "admin" {
					execute = true
				}
			}

			if !execute {
				following(ctx)
				return
			}

			author := ctx.Event.Author.ID
			guild := ctx.Event.GuildID

			member, err := ctx.Session.State.Member(guild, author)
			if err != nil {
				sendDiscordMessageEmbed(ctx, fmt.Sprintf("Couldnt get member from guild: %s", err), true)
				return
			}

			for _, role := range member.Roles {
				if role == config.AdminRole {
					// user is allowed to do it!
					log.Printf("User is allowed to do that :)")
					following(ctx)
					return
				}
			}

			sendDiscordMessageEmbed(ctx, "You are not allowed to use this command!", true)
			return
		}
	})

	// write command and params in statistics file
	router.RegisterMiddleware(func(following dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			execute := false
			for _, flag := range ctx.Command.Flags {
				if flag == "statistic" {
					execute = true
				}
			}

			if !execute {
				following(ctx)
				return
			}

			go func() {
				// one file for every command
				filename := filepath.Join(statisticsDir, ctx.Command.Name+".log")

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
					log.Printf("Error writing string to statistics log: %s\n", err)
					return
				}
			}()

			following(ctx)
		}
	})
}
