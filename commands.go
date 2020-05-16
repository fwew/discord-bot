package main

import (
	"fmt"
	"github.com/Lukaesebrot/dgc"
	fwew "github.com/fwew/fwew_lib"
	"strconv"
)

func registerCommands(router *dgc.Router) {
	// Random command
	router.RegisterCmd(&dgc.Command{
		Name: "random",
		Aliases: []string{
			"rand",
		},
		Description: fwew.Text("/randomDesc"),
		Usage:       fwew.Text("/randomUsage"),
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			var err error

			arguments := ctx.Arguments
			if arguments.Amount() >= 1 {
				firstArg := ctx.CustomObjects["firstArg"].(int)

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
						sendDiscordMessage(ctx, fmt.Sprintf("Argument is not a number: %s", err))
						return
					}
				}

				// rest of the arguments are the filters, if second one is "where"
				var restArgs []string
				whereArg := arguments.Get(firstArg + 1)
				if whereArg.Raw() == "where" {
					for i := 2; i < arguments.Amount(); i++ {
						restArgs = append(restArgs, arguments.Get(i).Raw())
					}
				}

				// Get random words out of dictionary
				words, err := fwew.Random(ctx.CustomObjects["langCode"].(string), amount, restArgs)
				if err != nil {
					sendDiscordMessage(ctx, fmt.Sprintf("Error getting random words: %s", err))
					return
				}

				sendWordDiscordEmbed(ctx, [][]fwew.Word{words})
			}
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
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			// get all arguments as array
			var args []string
			arguments := ctx.Arguments
			for i := 0; i < arguments.Amount(); i++ {
				argument := arguments.Get(i)
				args = append(args, argument.Raw())
			}

			words, err := fwew.List(args, ctx.CustomObjects["langCode"].(string))
			if err != nil {
				sendDiscordMessage(ctx, fmt.Sprintf("Error executing list command: %s", err))
				return
			}

			sendWordDiscordEmbed(ctx, [][]fwew.Word{words})
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
		Usage:       "fwew <word>...\n<word>:\n  - A Na'vi word to translate\n  - With `-l`: A locale word to translate",
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  false,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			arguments := ctx.Arguments

			firstArg := ctx.CustomObjects["firstArg"].(int)
			amount := arguments.Amount() - firstArg
			words := make([][]fwew.Word, amount)

			langCode := ctx.CustomObjects["langCode"].(string)
			// all params are words to search
			for i, j := firstArg, 0; i < arguments.Amount(); i, j = i+1, j+1 {
				var navi []fwew.Word
				if ctx.CustomObjects["reverse"].(bool) {
					navi = fwew.TranslateToNavi(arguments.Get(i).Raw(), langCode)
				} else {
					navi = fwew.TranslateFromNavi(arguments.Get(i).Raw(), langCode)
				}
				words[j] = append(words[j], navi...)
			}

			sendWordDiscordEmbed(ctx, words)
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
		Flags: []string{
			"params",
			"statistic",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			arguments := ctx.Arguments
			argument := arguments.Get(ctx.CustomObjects["firstArg"].(int))
			argInt, err := strconv.ParseInt(argument.Raw(), 8, 16)
			if err == nil {
				// It is an int, try to translate int
				navi, err := fwew.NumberToNavi(int(argInt))
				if err != nil {
					sendDiscordMessage(ctx, fmt.Sprintf("Error writing number: %s", err))
					return
				}

				sendDiscordMessage(ctx, navi)
			} else {
				// Try to translate from navi to number
				number, err := fwew.NaviToNumber(argument.Raw())
				if err != nil {
					sendDiscordMessage(ctx, fmt.Sprintf("Error reading number: %s", err))
					return
				}

				// Write string to print
				output := fmt.Sprintf("Decimal: %d\nOctal: %#o", number, number)

				sendDiscordMessage(ctx, output)
			}
		},
	})

	// just a command to show how to use parameters
	router.RegisterCmd(&dgc.Command{
		Name:        "params",
		Aliases:     nil,
		Description: "Show information about the params, that can be used with \"fwew\", \"list\" and \"random\"",
		Handler: func(ctx *dgc.Ctx) {
			info := "`fwew`, `list` and `random` can have additional optional parameters.\n" +
				"  - `-l=<langCode>`: Set the language\n" +
				"  - `-r`: `fwew` only param, that will mark the translation \"reversed\". If set, translation will be from locale to Na'vi\n" +
				"  - `-i`: Show Infix locations with brackets\n" +
				"  - `-id=false`: Dont show infix dots\n" +
				"  - `-src`: Show Source of this words\n" +
				"  - `-ipa`: Show IPA data\n" +
				"  - `-s=false`: Dont show the dashed syllable stress"

			sendDiscordMessage(ctx, info)
		},
	})

	// version command
	router.RegisterCmd(&dgc.Command{
		Name:        "version",
		Description: "Shows the current version of dict, api and bot.",
		IgnoreCase:  true,
		Handler: func(ctx *dgc.Ctx) {
			sendDiscordMessage(ctx, Version.String())
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
			err := fwew.UpdateDict()
			if err != nil {
				sendDiscordMessage(ctx, fmt.Sprintf("Error updating the dictionary: %s", err))
			} else {
				sendDiscordMessage(ctx, "Updating dictionary successful")
			}
		},
	})
}
