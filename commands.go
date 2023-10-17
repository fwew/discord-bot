package main

import (
	"fmt"
	"sort"
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
	words, err := fwew.Random(amount, restArgs)
	if err != nil {
		sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error getting random words: %s", err), true)
		return
	}

	send1dWordDiscordEmbed(ctx, [][]fwew.Word{words}[0])
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

	send1dWordDiscordEmbed(ctx, [][]fwew.Word{words}[0])
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
		output += "→ " + lenition[1] + "\n"
	}
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
		output += " → " + lenition[1] + "\n"
	}
	output += "```"
	sendDiscordMessageEmbed(ctx, output, false)
}

func that(ctx *dgc.Ctx) {
	thatTable := fwew.GetThatTable()
	const leftSize = 3
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
	var output string = "- **A1 Names:** Akwey, Ateyo, Eytukan, Eywa," +
		" Mo'at, Na'vi, Newey, Neytiri, Ninat, Omatikaya," +
		" Otranyu, Rongloa, Silwanin, Tskaha, Tsu'tey, Tsumongwi\n" +
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
func chart_entry(entry string, amount string, length int) (output string) {
	output = entry
	for i := utf8.RuneCountInString(entry); i < length-utf8.RuneCountInString(amount); i++ {
		output += " "
	}
	output += amount + "|"
	return output
}

// helper classes for phoneme_frequency
type phoneme struct {
	Freq int
	Name string
}

type phonemes []phoneme

func (e phonemes) Len() int {
	return len(e)
}

func (e phonemes) Less(i, j int) bool {
	if e[i].Freq == e[j].Freq {
		return e[i].Name < e[j].Name
	}
	return e[i].Freq > e[j].Freq
}

func (e phonemes) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
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
			if arg == "hrh" {
				// KP "HRH" video
				hrh := "https://youtu.be/-AgnLH7Dw3w?t=274\n"
				hrh += "> What would LOL be?\n"
				hrh += "> It would have to do with the word herangham... maybe HRH"
				sendDiscordMessageEmbed(ctx, hrh, false)
				//continue
			}

			argString := ""
			for i := 0; i < arguments.Amount(); i++ {
				argString += arguments.Get(i).Raw() + " "
			}
			argString = argString[:len(argString)-1]

			var navi [][]fwew.Word

			var err error
			navi, err = fwew.BidirectionalSearch(argString, true, langCode)
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

			fmt.Println(arguments)

			// all params are words to search
			argString := ""
			for i := 0; i < arguments.Amount(); i++ {
				argString += arguments.Get(i).Raw() + " "
			}
			argString = argString[:len(argString)-1]

			var navi [][]fwew.Word
			if ctx.CustomObjects.MustGet("reverse").(bool) {
				navi = fwew.TranslateToNaviHash(argString, langCode)
			} else {
				var err error
				navi, err = fwew.TranslateFromNaviHash(argString, false)
				if err != nil {
					sendDiscordMessageEmbed(ctx, fmt.Sprintf("Error translating: %s", err), true)
				}
			}
			words = navi
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
				"  - `-l=<langCode>`: Set the language (de, en, et, fr, hu, nl, pl, ru, sv, tr). Default: en\n" +
				"  - `-r`: `fwew` only param, that will mark the translation \"reversed\". If set, translation will be from locale to Na'vi\n" +
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

	// command to show all possible thats
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
		Handler: func(ctx *dgc.Ctx) {
			all_frequencies := fwew.GetPhonemeDistrosMap()
			entries := []string{"| Onset:|Nuclei:|Ending:|", "|=======|=======|=======|"}

			onset_letters := [21]string{"t", "", "n", "k", "l", "s", "'", "p", "r", "y",
				"ts", "m", "tx", "v", "w", "h", "ng", "z", "kx", "px", "f"}
			nucleus_letters := [14]string{"a", "e", "ì", "o", "u", "i", "ä", "aw", "ey", "ù", "rr", "ay", "ew", "ll"}
			coda_letters := [13]string{"", "n", "m", "ng", "l", "k", "p", "'", "r", "t", "kx", "px", "tx"}

			// Onsets
			onset_tuples := []phoneme{}
			for i := 0; i < len(onset_letters); i++ {
				var a phoneme
				a.Name = onset_letters[i]
				a.Freq = all_frequencies["Others"]["Onsets"][onset_letters[i]]
				onset_tuples = append(onset_tuples, a)
			}

			sort.Sort(phonemes(onset_tuples))

			for i := 0; i < len(onset_tuples); i++ {
				entries = append(entries, "|"+chart_entry(onset_tuples[i].Name, strconv.Itoa(onset_tuples[i].Freq), 7))
			}

			// Nuclei
			nuclei_tuples := []phoneme{}
			for i := 0; i < len(nucleus_letters); i++ {
				var a phoneme
				a.Name = nucleus_letters[i]
				a.Freq = all_frequencies["Others"]["Nuclei"][nucleus_letters[i]]
				nuclei_tuples = append(nuclei_tuples, a)
			}

			sort.Sort(phonemes(nuclei_tuples))

			i := 2
			for ; i < len(nuclei_tuples)+2; i++ {
				entries[i] += chart_entry(nuclei_tuples[i-2].Name, strconv.Itoa(nuclei_tuples[i-2].Freq), 7)
			}
			for ; i < len(entries); i++ {
				entries[i] += "       |"
			}

			// Ends
			codaTuples := []phoneme{}
			for i := 0; i < len(coda_letters); i++ {
				var a phoneme
				a.Name = coda_letters[i]
				a.Freq = all_frequencies["Others"]["Codas"][coda_letters[i]]
				codaTuples = append(codaTuples, a)
			}

			sort.Sort(phonemes(codaTuples))

			i = 2
			for ; i < len(codaTuples)+2; i++ {
				entries[i] += chart_entry(codaTuples[i-2].Name, strconv.Itoa(codaTuples[i-2].Freq), 7)
			}
			for ; i < len(entries); i++ {
				entries[i] += "       |"
			}

			// Top
			//entries_2 = "## Phoneme distributions:\n```\n"
			entries = append(entries, "")
			entries = append(entries, "Clusters:")
			entries = append(entries, "  | f:| s:|ts:|")
			entries = append(entries, "==|===|===|===|")

			// Clusters
			cluster_starts := []string{"f", "s", "ts"}
			cluster_ends := []string{"k", "kx", "l", "m", "n", "ng", "p", "px", "r", "t", "tx", "w", "y"}

			for i := 0; i < len(cluster_ends); i++ {
				entries = append(entries, chart_entry("", cluster_ends[i], 2))
			}

			// clusters
			for i := 0; i < len(cluster_starts); i++ {
				for j := 0; j < len(cluster_ends); j++ {
					entries[j+len(entries)-len(cluster_ends)] += chart_entry("", strconv.Itoa(all_frequencies["Clusters"][cluster_starts[i]][cluster_ends[j]]), 3)
				}
			}

			results := "```\n"

			for i := 0; i < len(entries); i++ {
				results += entries[i] + "\n"
			}

			results += "```"

			sendDiscordMessageEmbed(ctx, results, false)
		},
	})
}
