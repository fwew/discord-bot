# fwew discord bot
[![License: GPL v2](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)

The Best Na'vi Dictionary as a discord Bot.

To get the usage, write `$help` to the bot.

## Development
This option is mostly for Contributors and Developers. Or people who like to compile stuff themselves.
You will need the [GO Programming Language](https://golang.org/) and [Git](https://git-scm.com/) installed. 

### Setup
We are using go modules so no GOPATH setup is needed.
To compile the bot simply run:
```shell script
cd ~/wherever/you/want
git clone https://github.com/fwew/discord-bot
cd discord-bot
go build ./...
```

### Config
The discord bot token has to be placed in conf.json.
Just copy the conf.json.example as conf.json and place your token as value of the token field.

After these steps you have an executable, that can be run directly, for the operating system your currently on.

To run the Bot correctly, you have to put the `dictionary.txt` file in one of the following directories:
- `.` (next to the executable)
- `./.fwew/` (into a .fwew directory, next to the executable)
- `~/.fwew/` (into a .fwew directory in the home dir.)

Dictionary can be downloaded from the [main repository](https://github.com/fwew/fwew_lib/tree/master/.fwew/dictionary.txt) or from [tireas Learnnavi page](https://tirea.learnnavi.org/dictionarydata/dictionary.txt)

### Misc
To cross compile:
```shell script
GOOS=darwin go build -o mac_fwew_ ./...
GOOS=linux go build -o bin/linux/fwew ./...
GOOS=windows go build -o bin/windows/fwew.exe ./...
```

## Statistics
The bot by itself will create anonymized statistics, about the calls to each call inside the statistics directory.
These statistics only save, what Parameters were used per command. Each command has its own file where it is saved.
