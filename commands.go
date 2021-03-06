package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	fwew "github.com/fwew/fwew-lib/v5"
	"github.com/knoxfighter/dgc"
)

func random(arguments *dgc.Arguments, firstArg int, ctx *dgc.Ctx) {
	var err error

	// only number argument
	argument := arguments.Get(firstArg)

	var amount int

	// if argument is "random" create a random number to run on dict
	argString := argument.Raw()
	if argString == "random" {
		amount = -1
	}

	// If command is not random (amount = 0) get the given number
	if amount == 0 {
		amount, err = argument.AsInt()
		if err != nil {
			sendDiscordMessageEmbed(ctx, fmt.Sprintf("Argument [%s] is not a number: %s", argString, err), true)
			return
		}
	}

	// rest of the arguments are the filters, if second one is "where"
	var restArgs []string
	whereArg := arguments.Get(firstArg + 1)
	if whereArg.Raw() == "where" {
		for i := firstArg + 2; i < arguments.Amount(); i++ {
			restArgs = append(restArgs, arguments.Get(i).Raw())
		}
	}

	// Get random words out of dictionary
	words, err := fwew.Random(amount, restArgs)
	if err != nil {
		sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error getting random words: %s", err), true)
		return
	}

	sendWordDiscordEmbed(ctx, [][]fwew.Word{words})
}

func list(ctx *dgc.Ctx, firstArg int) {
	var args []string

	// get all arguments as array
	arguments := ctx.Arguments
	for i := firstArg; i < arguments.Amount(); i++ {
		argument := arguments.Get(i)
		args = append(args, argument.Raw())
	}

	words, err := fwew.List(args)
	if err != nil {
		sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error executing list command: %s", err), true)
		return
	}

	sendWordDiscordEmbed(ctx, [][]fwew.Word{words})
}

func lenition(ctx *dgc.Ctx) {
	lenitionTable := fwew.GetLenitionTable()
	const leftSize = 3
	var output string
	output += "```\n"
	for _, lenition := range lenitionTable {
		output += "" + lenition[0]
		for i := len(lenition[0]); i < leftSize; i++ {
			output += " "
		}
		if lenition[1] == "" {
			lenition[1] = "(disappears)"
		}
		output += "→ " + lenition[1] + "\n"
	}
	output += "```"
	sendDiscordMessageEmbed(ctx, output, false)
}

func registerCommands(router *dgc.Router) {
	// Random command
	router.RegisterCmd(&dgc.Command{
		Name: "random",
		Aliases: []string{
			"rand",
		},
		Description: fwew.Text("/randomDesc"),
		Usage:       fwew.Text("/randomUsage"),
		Example:     fwew.Text("/randomExample"),
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			firstArg, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Wrong usage of `random` command!", true)
				return
			}

			arguments := ctx.Arguments
			random(arguments, firstArg.(int), ctx)
		},
	})

	// list command
	router.RegisterCmd(&dgc.Command{
		Name: "list",
		Aliases: []string{
			"ls",
		},
		Description: fwew.Text("/listDesc"),
		Usage:       fwew.Text("/listUsage"),
		Example:     fwew.Text("/listExample"),
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			firstArg, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Wrong usage of `list` command!", true)
				return
			}
			list(ctx, firstArg.(int))
		},
	})

	// just translation
	router.RegisterCmd(&dgc.Command{
		Name: "fwew",
		Aliases: []string{
			"search",
			"translate",
			"trans",
		},
		Description: "Translate a word",
		Usage:       "fwew <word>...\n<word>:\n  - A Na'vi word to translate\n  - With `-r`: A locale word to translate",
		Example:     "fwew kaltxì",
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  false,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			arguments := ctx.Arguments

			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			// Dont run if firstArg is not set (we have nothing to do in that case)
			firstArgTemp, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Nothing found to translate!", true)
				return
			}

			firstArg := firstArgTemp.(int)
			amount := arguments.Amount() - firstArg
			words := make([][]fwew.Word, amount)

			langCode := ctx.CustomObjects.MustGet("langCode").(string)

			var wordFound bool

			// all params are words to search
			for i, j := firstArg, 0; i < arguments.Amount(); i, j = i+1, j+1 {
				arg := arguments.Get(i).Raw()

				// on first arg, check if this is a known command and fwew-bot is used like the old version
				if j == 0 {
					if strings.HasPrefix(arg, "/") {
						switch arg {
						case "/random":
							random(arguments, firstArg+1, ctx)
						case "/list":
							list(ctx, firstArg+1)
						case "/version":
							sendDiscordMessageEmbed(ctx, Version.String(), false)
						case "/lenition":
							fallthrough
						case "/len":
							lenition(ctx)
						default:
							// unknown command error
							sendEmbed(ctx, ctx.Command.Name, "I dont know this subcommand :(", true)
						}

						break
					}
				}

				// hardcoded stuff override (will send an additional message)
				if arg == "hrh" {
					// KP "HRH" video
					hrh := "https://youtu.be/-AgnLH7Dw3w?t=274\n"
					hrh += "> What would LOL be?\n"
					hrh += "> It would have to do with the word herangham... maybe HRH"
					sendDiscordMessageEmbed(ctx, hrh, false)
					continue
				}
				if arg == "tunayayo" {
					user, err := ctx.Session.User("277818358655877125")
					if err != nil {
						log.Printf("Error getting tunayayo user: %s", err)
						continue
					}
					avatarURL := user.AvatarURL("2048")

					sendEmbedImage(ctx, avatarURL)
					continue
				}

				var navi []fwew.Word
				if ctx.CustomObjects.MustGet("reverse").(bool) {
					navi = fwew.TranslateToNavi(arg, langCode)
				} else {
					var err error
					navi, err = fwew.TranslateFromNavi(arg)
					if err != nil {
						sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error translating: %s", err), true)
					}
				}
				words[j] = navi
				wordFound = true
			}

			if wordFound {
				sendWordDiscordEmbed(ctx, words)
			}
		},
	})

	// number translation
	router.RegisterCmd(&dgc.Command{
		Name: "number",
		Aliases: []string{
			"num",
			"n",
		},
		Description: "Translate a number to Navi and vice-versa",
		Usage: `number <number>
<number>:
  - an octal number to translate to Na'vi
  - the Na'vi word of a number, to read the number
`,
		Example: "number 55",
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			var number int
			var output string
			var err error

			arguments := ctx.Arguments
			firstArg, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Nothing found to translate!", true)
				return
			}
			argument := arguments.Get(firstArg.(int))

			// Parse number
			arg := argument.Raw()

			// check if arg starts with number
			var rune rune
			for _, r := range arg {
				rune = r
				break
			}
			if rune >= '0' && rune <= '9' {
				// try to get number of it
				argInt, err := strconv.ParseInt(arg, 0, 16)
				if err != nil {
					sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error reading number from input: %s", err), true)
					return
				}
				number = int(argInt)

				// It is an int, try to translate int
				output, err = fwew.NumberToNavi(number)
				if err != nil {
					sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error writing number: %s", err), true)
					return
				}
				output = "**" + output + "**\n"
			} else {
				// Try to translate from navi to number
				number, err = fwew.NaviToNumber(arg)
				if err != nil {
					sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error reading number: %s", err), true)
					return
				}
			}
			// Write number to output too
			output = fmt.Sprintf("%sDecimal: %d\nOctal: %#o", output, number, number)

			sendDiscordMessageEmbed(ctx, output, false)
		},
	})

	// just a command to show how to use parameters
	router.RegisterCmd(&dgc.Command{
		Name: "params",
		Aliases: []string{
			"param",
		},
		Description: "Show information about the params, that can be used with \"fwew\", \"list\" and \"random\"",
		Handler: func(ctx *dgc.Ctx) {
			info := "`fwew`, `list` and `random` can have additional optional parameters.\n" +
				"  - `-l=<langCode>`: Set the language (de, en, et, fr, hu, nl, pl, ru, sv). Default: en\n" +
				"  - `-r`: `fwew` only param, that will mark the translation \"reversed\". If set, translation will be from locale to Na'vi\n" +
				"  - `-i`: Show Infix locations with brackets\n" +
				"  - `-id=false`: Dont show infix dots\n" +
				"  - `-src`: Show Source of this words\n" +
				"  - `-ipa`: Show IPA data\n" +
				"  - `-s=false`: Dont show the dashed syllable stress"

			sendDiscordMessageEmbed(ctx, info, false)
		},
	})

	// version command
	router.RegisterCmd(&dgc.Command{
		Name:        "version",
		Description: "Shows the current version of dict, api and bot.",
		IgnoreCase:  true,
		Handler: func(ctx *dgc.Ctx) {
			sendDiscordMessageEmbed(ctx, Version.String(), false)
		},
	})

	// update command
	router.RegisterCmd(&dgc.Command{
		Name:        "update",
		Description: "Update the dictionary (can only be used by admins)",
		IgnoreCase:  true,
		Flags: []string{
			"admin",
		},
		Handler: func(ctx *dgc.Ctx) {
			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			err := fwew.UpdateDict()
			if err != nil {
				sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error updating the dictionary: %s", err), true)
			} else {
				sendDiscordMessageEmbed(ctx, "Updating dictionary successful", false)
			}
		},
	})

	// command to show all possible lenitions
	router.RegisterCmd(&dgc.Command{
		Name: "lenition",
		Aliases: []string{
			"len",
		},
		Description: "Show all possible lenitions",
		IgnoreCase:  true,
		Handler:     lenition,
	})
}
