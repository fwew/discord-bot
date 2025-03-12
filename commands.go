package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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
	words, err := fwew.Random(amount, restArgs, uint8(1))
	if err != nil {
		sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error getting random words: %s", err), true)
		return
	}

	send1dWordDiscordEmbed(ctx, words)
}

func list(ctx *dgc.Ctx, firstArg int) {
	var args []string

	// get all arguments as array
	arguments := ctx.Arguments
	for i := firstArg; i < arguments.Amount(); i++ {
		argument := arguments.Get(i)
		newArg := argument.Raw()
		if newArg[len(newArg)-1] == ',' {
			newArg = newArg + arguments.Get(i+1).Raw()
			i++
		}
		args = append(args, newArg)
	}

	words, err := fwew.List(args, uint8(1))
	if err != nil {
		sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error executing list command: %s", err), true)
		return
	}

	send1dWordDiscordEmbed(ctx, words)
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
			lenition[1] = "(disappears, except before ll or rr)"
		}
		output += "→ " + lenition[1] + "\n\n"
	}
	output += "leniting prefixes: me+, pxe+, ay+, pe+\n"
	output += "leniting adpositions: fpi, ìlä, lisre, mì, nuä, pxisre, ro, sko, sre, wä\n"
	output += "```"
	sendDiscordMessageEmbed(ctx, output, false)
}

func shortLenition(ctx *dgc.Ctx) {
	lenitionTable := fwew.GetShortLenitionTable()
	const leftSize = 10
	var output string
	output += "```\n"
	for _, lenition := range lenitionTable {
		for i := len(lenition[0]); i < leftSize; i++ {
			output += " "
		}
		output += "" + lenition[0]
		if lenition[1] == "" {
			lenition[1] = "(disappears, except before ll or rr)"
		}
		output += " → " + lenition[1] + "\n\n"
	}
	output += "leniting prefixes: me+, pxe+, ay+, pe+\n"
	output += "leniting adpositions: fpi, ìlä, lisre, mì, nuä, pxisre, ro, sko, sre, wä\n"
	output += "```"
	sendDiscordMessageEmbed(ctx, output, false)
}

func that(ctx *dgc.Ctx) {
	thatTable := fwew.GetThatTable()
	var output string
	output += "```\n"

	//Get column widths
	var lengths = [len(thatTable[2])]int{0, 0, 0, 0, 0}
	for j := 0; j < len(thatTable[2]); j++ {
		lengths[j] = len(thatTable[2][j])
	}

	for _, that := range thatTable {
		for i := 0; i < len(that); i++ {
			var word = that[i]
			if len(word) > 0 {
				output += word
				for j := len(word); j < lengths[i]; j++ {
					output += " "
				}
				output += "|"
			}
		}
		output += "\n"
	}

	output += "\n"

	otherThats := fwew.GetOtherThats()

	//The other ones that don't fit on the chart
	var lineNum = 7
	var lengths2 = [len(otherThats[lineNum])]int{0, 0, 0}
	for j := 0; j < len(otherThats[lineNum]); j++ {
		lengths2[j] = utf8.RuneCountInString(otherThats[lineNum][j])
	}

	for _, that := range otherThats {
		for i := 0; i < len(that); i++ {
			var word = that[i]
			if utf8.RuneCountInString(word) > 0 {
				output += word
				for j := utf8.RuneCountInString(word); j <= lengths2[i]; j++ {
					output += " "
				}
			}
		}
		output += "\n"
	}

	output += "```"
	sendDiscordMessageEmbed(ctx, output, false)
}

func cameronWords(ctx *dgc.Ctx) {
	var output = "- **A1 Names:** Akwey, Ateyo, Eytukan, Eywa," +
		" Mo'at, Na'vi, Newey, Neytiri, Ninat, Omatikaya," +
		" Otranyu, Rongloa, Silwanin, Tskaha, Tsu'tey\n" +
		"- **A2 Names:** Aonung, Kiri, Lo'ak, Neteyam," +
		" Ronal, Rotxo, Tonowari, Tuktirey, Tsireya\n" +
		"- **Nouns:** 'itan, 'ite, atan, au *(drum)*, eyktan, i'en," +
		" Iknimaya, mikyun, ontu, seyri, tsaheylu, tsahìk, unil\n" +
		"- **Life:** Atokirina', Ikran, Palulukan," +
		" Riti, talioang, teylu, Toruk\n" +
		"- **Other:** eyk, irayo, makto, taron, te"
	sendDiscordMessageEmbed(ctx, output, false)
}

// Helper function for phoneme_frequency
func chartEntry(entry string, amount string, length int) (output string) {
	output = entry
	for i := utf8.RuneCountInString(entry); i < length-utf8.RuneCountInString(amount); i++ {
		output += " "
	}
	output += amount + "|"
	return output
}

func phonemeFrequency(ctx *dgc.Ctx) {
	all_frequencies := fwew.GetPhonemeDistrosMap("en") // English only

	results := "```\n"

	for _, a := range all_frequencies[0] {
		results += "|"
		for _, b := range a {
			entries := strings.Split(b, " ")
			if len(entries) == 2 {
				results += chartEntry(entries[0], entries[1], 8)
			} else {
				results += chartEntry("", b, 8)
			}
		}
		results += "\n"
	}

	results += "\n" + all_frequencies[1][0][0] + ":\n"
	all_frequencies[1][0][0] = ""

	for _, a := range all_frequencies[1] {
		newLine := ""
		for _, b := range a {
			newLine += chartEntry("", b, 3)
		}
		newLine = strings.TrimPrefix(newLine, " ")
		results += newLine + "\n"
	}

	results += "```"

	sendDiscordMessageEmbed(ctx, results, false)
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
		Description: "Translate a word or several",
		Usage:       "fwew <word>...\n<word>:\n  - A word to translate\n  - With `-r`: A locale word to translate",
		Example:     "fwew kaltxì run",
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

			// Don't run if firstArg is not set (we have nothing to do in that case)
			firstArgTemp, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Nothing found to translate!", true)
				return
			}

			firstArg := firstArgTemp.(int)

			langCode := ctx.CustomObjects.MustGet("langCode").(string)

			// all params are words to search
			arg := arguments.Get(0).Raw()

			// on first arg, check if this is a known command and fwew-bot is used like the old version
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
				case "/that":
					that(ctx)
				case "/cameronWords":
					cameronWords(ctx)
				default:
					// unknown command error
					sendEmbed(ctx, ctx.Command.Name, "I don't know this subcommand :(", true)
				}
			}

			// hardcoded stuff override (will send an additional message)
			if strings.ToLower(arg) == "hrh" {
				// KP "HRH" video
				hrh := "https://youtu.be/-AgnLH7Dw3w?t=274\n"
				hrh += "> What would LOL be?\n"
				hrh += "> It would have to do with the word herangham... maybe HRH"
				sendDiscordMessageEmbed(ctx, hrh, false)
				//continue
			}

			argString := ""
			collect := false
			for i := 0; i < arguments.Amount(); i++ {
				if collect {
					argString += " "
				} else if arguments.Get(i).Raw()[0] != '-' {
					collect = true
				} else {
					continue
				}
				argString += arguments.Get(i).Raw()
			}

			var navi [][]fwew.Word

			var err error
			navi, err = fwew.BidirectionalSearch(argString, true, langCode, false)
			if err != nil {
				sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error translating: %s", err), true)
			}

			sendWordDiscordEmbed(ctx, navi)
		},
	})

	// translation and skipping any affix checks
	router.RegisterCmd(&dgc.Command{
		Name:        "fwew-simple",
		Description: "Translate a word (no checking for affixes or natural language words)",
		Usage:       "fwew-simple <word>...\n<word>:\n  - A Na'vi word to translate\n  - With `-r`: A locale word to translate",
		Example:     "fwew-simple uturu",
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

			// Don't run if firstArg is not set (we have nothing to do in that case)
			_, b := ctx.CustomObjects.Get("firstArg")
			if !b {
				sendDiscordMessageEmbed(ctx, "Nothing found to translate!", true)
				return
			}

			langCode := ctx.CustomObjects.MustGet("langCode").(string)

			var wordFound bool

			// all params are words to search
			argString := ""
			collect := false
			for i := 0; i < arguments.Amount(); i++ {
				if collect {
					argString += " "
				} else if arguments.Get(i).Raw()[0] != '-' {
					collect = true
				} else {
					continue
				}
				argString += arguments.Get(i).Raw()
			}

			var navi [][]fwew.Word
			if ctx.CustomObjects.MustGet("reverse").(bool) {
				navi = fwew.TranslateToNaviHash(argString, langCode)
			} else {
				var err error
				navi, err = fwew.TranslateFromNaviHash(argString, false, false, false)
				if err != nil {
					sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error translating: %s", err), true)
				}
			}
			words := navi
			wordFound = true

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
			var argRune rune
			for _, r := range arg {
				argRune = r
				break
			}
			if argRune >= '0' && argRune <= '9' {
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
				"  - `-l=<langCode>`: Set the language (de, en, es, et, fr, hu, ko, nl, pl, pt, ru, sv, tr, uk). Default: en\n" +
				"  - `-r`: `fwew` only param, that will mark the translation \"reversed\". If set, translation will be from locale to Na'vi\n" +
				"  - `-reef`: Show Reef dialect information\n" +
				"  - `-i`: Show Infix locations with brackets\n" +
				"  - `-id=false`: Don't show infix dots\n" +
				"  - `-src`: Show Source of this words\n" +
				"  - `-ipa`: Show IPA data\n" +
				"  - `-s=false`: Don't show the dashed syllable stress"

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

	// command to show all possible lenition
	router.RegisterCmd(&dgc.Command{
		Name:        "lenition",
		Description: "Show all possible lenition",
		IgnoreCase:  true,
		Handler:     lenition,
	})

	// command to show all possible len
	router.RegisterCmd(&dgc.Command{
		Name:        "len",
		Description: "Show all possible lenition (short)",
		IgnoreCase:  true,
		Handler:     shortLenition,
	})

	// command to show all possible "that"s
	router.RegisterCmd(&dgc.Command{
		Name:        "that",
		Description: "Show all possible thats",
		IgnoreCase:  true,
		Handler:     that,
	})

	// command to show words James Cameron invented
	router.RegisterCmd(&dgc.Command{
		Name:        "Cameron Words",
		Description: "Show words James Cameron invented",
		IgnoreCase:  true,
		Handler:     cameronWords,
	})

	// command to show how often each phoneme appears
	router.RegisterCmd(&dgc.Command{
		Name:        "phoneme-frequency",
		Description: "Show how often a phoneme appears",
		IgnoreCase:  true,
		Handler:     phonemeFrequency,
	})

	// Tell the user if given word(s) are valid in Na'vi
	router.RegisterCmd(&dgc.Command{
		Name:        "valid",
		Description: "See if a word would be valid in Na'vi",
		Usage:       "valid <word>...\n<word>:\n  - A word to validate",
		Example:     "valid omati s'ampta",
		Flags: []string{
			"params",
		},
		IgnoreCase:  true,
		SubCommands: nil,
		Handler: func(ctx *dgc.Ctx) {
			arguments := ctx.Arguments

			defer func() {
				if err := recover(); err != nil {
					sendErrorWhenRecovered(ctx)
				}
			}()

			argString := ""
			for i := 0; i < arguments.Amount(); i++ {
				argString += arguments.Get(i).Raw() + " "
			}
			argString = argString[:len(argString)-1]

			navi := fwew.IsValidNavi(argString, "en", true) // Not sure how to enable language support easily
			sendDiscordMessageEmbed(ctx, navi, false)
		},
	})
}
