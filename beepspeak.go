package main

// $ apt install libasound2-dev

import (
	"bytes"
	"cloud.google.com/go/texttospeech/apiv1"
	"context"
	"github.com/erikbryant/aes"
	"github.com/erikbryant/web"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var (
	passPhrase     = flag.String("passPhrase", "", "Passphrase to unlock API key")
)

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
func play(file string) {
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

// readable makes a string more human readable by removing all non alphanumeric and non-punctuation.
func readable(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, "/", "")
	text = strings.ReplaceAll(text, "^", "")
	text = strings.ReplaceAll(text, "[", " ")
	text = strings.ReplaceAll(text, "]", " ")
	text = strings.ReplaceAll(text, "\"", "")

	// Specific ships we have seen that read poorly.
	text = strings.ReplaceAll(text, "SEASPAN HAMBURG", "SEA SPAN HAMBURG")

	return text
}

// say converts text to speech and then plays it.
func say(text string) {
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

// alert prints a message and plays an alert tone.
func alert(details database.Ship) {
	fmt.Printf(
		"\nShip Ahoy!  %s  %s\n%+v\n\n",
		time.Now().Format("Mon Jan 2 15:04:05"),
		decodeMmsi(details.MMSI),
		prettify(details),
	)

	if strings.Contains(strings.ToLower(details.Type), "vehicle") {
		play("meep.wav")
	} else if strings.Contains(strings.ToLower(details.Type), "pilot") {
		play("pilot.mp3")
	} else {
		play("ship_horn.mp3")
	}

	summary := fmt.Sprintf("Ship ahoy! %s. %s. Course %3.f degrees.", details.Name, details.Type, details.ShipCourse)

	// Hearing, "eleven point zero knots" sounds awkward. Remove the "point zero".
	if math.Trunc(details.Speed) == details.Speed {
		summary = fmt.Sprintf("%s Speed %3.0f knots.", summary, math.Trunc(details.Speed))
	} else {
		summary = fmt.Sprintf("%s Speed %3.1f knots.", summary, details.Speed)
	}

	switch details.Sightings {
	case 0:
		summary = fmt.Sprintf("%s This is the first sighting.", summary)
	case 1:
		summary = fmt.Sprintf("%s One previous sighting.", summary)
	default:
		summary = fmt.Sprintf("%s %d previous sightings.", summary, details.Sightings)
	}

	say(summary)
}
