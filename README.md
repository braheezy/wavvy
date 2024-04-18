# [wavvy](https://wavvy.braheezy.net/)

A toy WAV audio player, [WASM](https://www.wikiwand.com/en/WebAssembly)-style for the browser. You can play a WAV file. You might even be able to play a second one!

I wrote this to learn WASM with Go and specifically audio playback with [`oto`](https://github.com/ebitengine/oto) because I wanted to get ready for porting my [`goqoa`](https://github.com/braheezy/goqoa) project to WASM.

## What's going on
The Go code shown on the site is compiled to a WASM binary and delivered to the browser when you visit. It contains all the logic to read, decode, and perform audio playback of a WAV file. All right in the browser, no weirdo server to talk to. I think this means as long as the file stays cached in the browser, the website will continue to work even if your device is offline.

The HTML frontend provides basic controls to select a WAV file and give it to the WASM executable via explicitly exposed functions. For funs, it also reads the local Go code from [main.go](./main.go) and uses [highlight.js](https://highlightjs.org/) to give it some life.

Other than a few more CDN pulls for fonts, this single-page application is comprised of a single `index.html` file along with the requisite `wasm_exec.js` file that provides the Go glue. And of course, the compiled `main.wasm`.

## Development
To work with this project locally, you need Go and a way to run a web server. Here's the commands I ran:

```bash
# Build it
GOOS=js GOARCH=wasm go build -o main.wasm main.go
# Serve it
python -m http.server
```
