# [wavvy](https://wavvy.braheezy.net/)

A toy WAV audio player, WASM-style for the browser. You can play a WAV file. You might even be able to play a second one!

I wrote this to learn WASM with Go and it's specifically audio playback because I wanted to get ready for porting
my [`goqoa`](https://github.com/braheezy/goqoa) project to WASM.

## Development
To work with this project locally, you need Go and a way to run a web server. Here's the commands I ran:

    # Build it
    GOOS=js GOARCH=wasm go build -o main.wasm main.go
    # Serve it
    python -m http.server
