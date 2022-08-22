![go fmt](https://github.com/erikbryant/beepspeak/actions/workflows/fmt.yml/badge.svg)
![go vet](https://github.com/erikbryant/beepspeak/actions/workflows/vet.yml/badge.svg)
![go test](https://github.com/erikbryant/beepspeak/actions/workflows/test.yml/badge.svg)

# Beep Speak

Takes audio files as input and plays them. Also takes text as input, converts it to speech, and plays it.

If you use the text-to-speech functions then you need to provide a GCP Speech API Token.

## Usage (playing sound files)

```golang
import (
  "github.com/erikbryant/beepspeak"
)

err := beepspeak.Play("mysong.wav")
if err != nil {
  return err
}

beepspeak.Play("mymusic.mp3")
if err != nil {
  return err
}
```

# Usage (text to speech)

```golang
import (
  "github.com/erikbryant/aes"
  "github.com/erikbryant/beepspeak"
)

// Put your GCP Speech API credentials in plainText.
plainText := "<redacted>"
passphrase := "MySuperSecurePassword"

cipherText, err := aes.Encrypt(plainText, passphrase)
if err != nil {
  return err
}

err = beepspeak.InitSay(cipherText, passphrase)
if err != nil {
  return err
}

err = beepspeak.Say("Hello, world!")
if err != nil {
  return err
}
```
