package main

import (
	"flag"
	"github.com/erikbryant/beepspeak"
)

var (
	passPhrase = flag.String("passPhrase", "", "Passphrase to unlock API key")
)

func main() {
	flag.Parse()

	beepspeak.InitSay(*passPhrase)
	beepspeak.Say("hello")
}
