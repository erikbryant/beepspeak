package beepspeak

// $ apt install libasound2-dev

import (
	"bytes"
	"cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/erikbryant/aes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"os"
	"strings"
	"time"
)

// InitSay saves our GCP Speech API credentials to disk so that the
// GCP speech API can find them.
func InitSay(gcpAuthCrypt, passPhrase string) error {
	path := os.TempDir() + "/" + "gcp_auth.json"

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	gcpAuth, err := aes.Decrypt(gcpAuthCrypt, passPhrase)
	if err != nil {
		return err
	}

	_, err = f.WriteString(gcpAuth + "\n")
	if err != nil {
		return err
	}

	// The Google API looks for the auth using an env var.
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)
	if err != nil {
		return err
	}

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

// Play plays a given sound file. MP3 and WAV are supported.
func Play(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	var (
		s      beep.StreamSeekCloser
		format beep.Format
	)

	if strings.HasSuffix(file, ".wav") {
		s, format, err = wav.Decode(f)
		if err != nil {
			return err
		}
	} else {
		s, format, err = mp3.Decode(f)
		if err != nil {
			return err
		}
	}

	playStream(s, format)

	return nil
}

// readable makes a string more human readable by removing some non alphanumeric
// and non-punctuation.
func readable(text string) string {
	text = strings.TrimSpace(text)

	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, "/", " ")
	text = strings.ReplaceAll(text, "[", " ")
	text = strings.ReplaceAll(text, "]", " ")

	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "^", "")

	return text
}

// Say converts text to speech and then plays it.
func Say(text string) error {
	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALSxx")
	if !ok {
		return fmt.Errorf("ERROR: env var is not set; did you call InitSay()")
	}

	text = readable(text)

	ctx := context.Background()

	c, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
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
			SpeakingRate:  1.1,
		},
	}
	resp, err := c.SynthesizeSpeech(ctx, req)
	if err != nil {
		return err
	}

	r := bytes.NewReader(resp.GetAudioContent())
	s, format, err := wav.Decode(r)
	if err != nil {
		return err
	}

	playStream(s, format)

	return nil
}
