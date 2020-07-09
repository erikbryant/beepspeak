package beepspeak

// $ apt install libasound2-dev

import (
	"bytes"
	"cloud.google.com/go/texttospeech/apiv1"
	"context"
	"flag"
	"fmt"
	"github.com/erikbryant/aes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"log"
	"os"
	"strings"
	"time"
)

var (
	passPhrase = flag.String("passPhrase", "", "Passphrase to unlock API key")

	gcpAuthCrypt = "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"
)

// init decrypts the credentials and writes them to a file
// so that the GCP Speech API can find them.
func init() {
	flag.Parse()

}

// writeGcpAuth saves our GCP Speech API credentials to disk so that the
// GCP speech API can find them.
func writeGcpAuth() error {
	path := os.TempDir() + "/" + "gcp_auth.json"

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	gcpAuth := aes.Decrypt(gcpAuthCrypt, *passPhrase)

	_, err = f.WriteString(gcpAuth + "\n")
	if err != nil {
		return err
	}

	// The Google API looks for the auth using an env var.
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)

	return nil
}

// playStream plays a given audio stream.
func playStream(s beep.StreamSeekCloser, format beep.Format) {
	// Init the Speaker with the SampleRate of the format and a buffer size.
	speaker.Init(format.SampleRate, format.SampleRate.N(3*time.Second))

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))
	<-playing
}

// play plays a given sound file. MP3 and WAV are supported.
func Play(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Could not open audio file", file)
		return
	}

	var (
		s      beep.StreamSeekCloser
		format beep.Format
	)

	if strings.HasSuffix(file, ".wav") {
		s, format, err = wav.Decode(f)
		if err != nil {
			fmt.Println("Could not decode WAV audio file", file, err)
			return
		}
	} else {
		s, format, err = mp3.Decode(f)
		if err != nil {
			fmt.Println("Could not decode MP3 audio file", file, err)
			return
		}
	}

	playStream(s, format)
}

// readable makes a string more human readable by removing some non alphanumeric
// and non-punctuation.
func readable(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, "/", "")
	text = strings.ReplaceAll(text, "^", "")
	text = strings.ReplaceAll(text, "[", " ")
	text = strings.ReplaceAll(text, "]", " ")
	text = strings.ReplaceAll(text, "\"", "")

	return text
}

// Say converts text to speech and then plays it.
func Say(text string) {
	err := writeGcpAuth()
	if err != nil {
		panic(err)
	}

	text = readable(text)

	ctx := context.Background()

	c, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			// https://cloud.google.com/text-to-speech/docs/voices
			LanguageCode: "en-US",
			Name:         "en-US-Standard-C",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
			SpeakingRate:  1.0,
		},
	}
	resp, err := c.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(resp.GetAudioContent())
	s, format, err := wav.Decode(r)
	if err != nil {
		fmt.Println("Could not decode WAV speech stream", text, err)
		return
	}

	playStream(s, format)
}
