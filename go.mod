module fwew-discord-bot

go 1.20

require (
	github.com/bwmarrin/discordgo v0.20.3
	github.com/fwew/fwew-lib/v5 v5.7.1-dev.0.20230704190850-db70d4bd3a5e
	github.com/knoxfighter/dgc v0.0.0-20201030020537-397f394c484d
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/karrick/tparse/v2 v2.8.1 // indirect
	github.com/zekroTJA/timedmap v0.0.0-20200518230343-de9b879d109a // indirect
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
	golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
