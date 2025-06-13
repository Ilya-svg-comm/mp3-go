package service

import (
	"fmt"
	"net/http"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	ctrl        *beep.Ctrl
	initialized bool
)

func PlayAudio(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("decode mp3: %w", err)
	}

	if !initialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(512))
		initialized = true
	}

	StopAudio() // Остановим предыдущий, если был

	ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return nil
}

func StopAudio() {
	speaker.Clear()
	if ctrl != nil {
		if closer, ok := ctrl.Streamer.(beep.StreamSeekCloser); ok {
			_ = closer.Close()
		}
		ctrl = nil
	}
}
