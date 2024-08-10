module fwew-discord-bot

go 1.22

toolchain go1.22.4

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/fwew/fwew-lib/v5 v5.18.2
	github.com/knoxfighter/dgc v0.0.0-20201030020537-397f394c484d
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/karrick/tparse/v2 v2.8.1 // indirect
	github.com/zekroTJA/timedmap v0.0.0-20200518230343-de9b879d109a // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
