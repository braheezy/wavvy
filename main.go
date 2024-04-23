//go:build js && wasm

package main

import (
	"bytes"
	"encoding/binary"
	"syscall/js"

	"github.com/ebitengine/oto/v3"
	"github.com/go-audio/wav"
)

var (
	// The global player. Should it be global? ¯\_(ツ)_/¯
	player *oto.Player
	// The audio context. There can be only one.
	otoCtx *oto.Context
)

func main() {
	// Declare functions that can be called from js.
	js.Global().Set("loadWav", js.FuncOf(loadWav))
	js.Global().Set("playWav", js.FuncOf(playWav))
	js.Global().Set("pauseWav", js.FuncOf(pauseWav))
	// Block forever, to keep Go runtime alive.
	select {}
}

// loadWav loads a WAV file into memory.
func loadWav(this js.Value, args []js.Value) interface{} {
	if player != nil {
		player.Close()
	}
	// Read WAV data from js.
	data := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(data, args[0])

	// Decode the WAV data.
	r := bytes.NewReader(data)
	decoder := wav.NewDecoder(r)
	if !decoder.IsValidFile() {
		return map[string]interface{}{"error": "invalid WAV file"}
	}
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return map[string]interface{}{"error": "failed to read WAV buffer: " + err.Error()}
	}

	// Load the audio context.
	err = loadContext(decoder)
	if err != nil {
		return map[string]interface{}{"error": "failed to load context: " + err.Error()}
	}

	// Convert to bytes.
	byteData := make([]byte, len(buf.Data)*2)
	for i, sample := range buf.Data {
		binary.LittleEndian.PutUint16(byteData[i*2:], uint16(sample))
	}

	// Once we can figure out the duration of the audio, we can tell the UI about it.
	durationMs := float64(len(buf.Data)) / float64(decoder.Format().SampleRate) * 1000 / float64(decoder.Format().NumChannels)
	js.Global().Call("setAudioDuration", durationMs)

	// Finally, load the data into the player.
	player = otoCtx.NewPlayer(bytes.NewReader(byteData))

	return map[string]interface{}{"msg": "loading new player"}
}

// loadContext loads the audio context.
func loadContext(decoder *wav.Decoder) error {
	if otoCtx != nil {
		return nil
	}
	// Create the Oto options
	op := &oto.NewContextOptions{
		SampleRate:   decoder.Format().SampleRate,
		ChannelCount: decoder.Format().NumChannels,
		Format:       oto.FormatSignedInt16LE,
	}

	// Create the Oto context.
	var err error
	var ready chan struct{}
	otoCtx, ready, err = oto.NewContext(op)
	if err != nil {
		return err
	}

	// Wait for the context to be ready. Must be in a goroutine or things deadlock.
	go func() {
		<-ready
	}()
	return nil
}

// pauseWav pauses the audio player.
func pauseWav(this js.Value, args []js.Value) interface{} {
	if player != nil {
		player.Pause()
	}
	return nil
}

// playWav plays the audio player.
func playWav(this js.Value, args []js.Value) interface{} {
	if player != nil {
		player.Play()
	}
	return nil
}
