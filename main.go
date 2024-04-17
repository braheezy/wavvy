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
	otoCtx *oto.Context
	player *oto.Player
)

func main() {
	c := make(chan struct{})
	js.Global().Set("loadWav", js.FuncOf(loadWav))
	js.Global().Set("playWav", js.FuncOf(playWav))
	js.Global().Set("pauseWav", js.FuncOf(pauseWav))
	<-c // Keep the Go runtime alive
}

func loadWav(this js.Value, args []js.Value) interface{} {
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	var err error
	var ready chan struct{}
	otoCtx, ready, err = oto.NewContext(op)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	go func() {
		<-ready
	}()
	data := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(data, args[0])

	r := bytes.NewReader(data)
	decoder := wav.NewDecoder(r)
	if !decoder.IsValidFile() {
		return map[string]interface{}{"error": "invalid WAV file"}
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return map[string]interface{}{"error": "failed to read WAV buffer: " + err.Error()}
	}

	byteData := make([]byte, len(buf.Data)*2)
	for i, sample := range buf.Data {
		binary.LittleEndian.PutUint16(byteData[i*2:], uint16(sample))
	}

	// Calculate duration and call setAudioDuration in JavaScript
	durationMs := float64(len(buf.Data)) / float64(decoder.Format().SampleRate) * 1000 / float64(decoder.Format().NumChannels)
	js.Global().Call("setAudioDuration", durationMs)

	player = otoCtx.NewPlayer(bytes.NewReader(byteData))

	return nil
}

func pauseWav(this js.Value, args []js.Value) interface{} {
	if player != nil {
		player.Pause()
	}
	return nil
}

func playWav(this js.Value, args []js.Value) interface{} {
	if player != nil {
		player.Play()
	}
	return nil
}
