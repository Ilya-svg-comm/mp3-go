package service

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var (
	ctx        *oto.Context
	player     oto.Player
	initOnce   sync.Once
	playerLock sync.Mutex
)

func initContext(sampleRate int) error {
	var err error
	initOnce.Do(func() {
		ctx, err = oto.NewContext(sampleRate, 2, 2, 8192)
	})
	return err
}

func PlayAudio(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	decoder, err := mp3.NewDecoder(resp.Body)
	if err != nil {
		return fmt.Errorf("decode mp3: %w", err)
	}

	// Инициализируем контекст один раз
	if err := initContext(decoder.SampleRate()); err != nil {
		return fmt.Errorf("init oto context: %w", err)
	}

	playerLock.Lock()
	defer playerLock.Unlock()

	// Остановим предыдущий плеер, если есть

	player = *ctx.NewPlayer()

	buffer := make([]byte, 8192)
	for {
		n, err := decoder.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read mp3: %w", err)
		}
		if n > 0 {
			if _, err := player.Write(buffer[:n]); err != nil {
				return fmt.Errorf("write audio: %w", err)
			}
		}
	}

	return nil
}

func StopAudio() {
	playerLock.Lock()
	defer playerLock.Unlock()

}
