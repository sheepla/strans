package audio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/sheepla/strans/api"
)

const (
	samplingRate     = 24000
	numberOfChannels = 2
	bitDepth         = 2
)

var (
	ErrEncode    = errors.New("an error occurred on encoding audio")
	ErrPlayAudio = errors.New("an error occurred on playing audio")
)

func encodeWithMP3(r io.Reader) (*mp3.Decoder, error) {
	dec, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrEncode, err)
	}

	return dec, nil
}

//nolint:varnamelen,nolintireturn
func newPlayer(r io.Reader) (oto.Player, error) {
	ctx, ready, err := oto.NewContext(samplingRate, numberOfChannels, bitDepth)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPlayAudio, err)
	}

	// await audio device ready
	<-ready

	player := ctx.NewPlayer(r)

	return player, nil
}

func Play(r io.Reader) error {
	dec, err := encodeWithMP3(r)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrPlayAudio, err)
	}

	player, err := newPlayer(dec)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrPlayAudio, err)
	}

	player.Play()

	// start playing
	for player.IsPlaying() {
		time.Sleep(1 * time.Millisecond)
	}

	// close audio device
	if err := player.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrPlayAudio, err)
	}

	return nil
}

//nolint:wrapcheck
func FetchAndPlay(lang, text, instance string) error {
	param, err := api.NewVoiceParam(lang, text, instance)
	if err != nil {
		return err
	}

	data, err := api.FetchVoice(param)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data.Audio)

	if err := Play(buf); err != nil {
		return err
	}

	return nil
}
