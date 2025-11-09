# audiomorph

[![CI](https://github.com/schollz/audiomorph/actions/workflows/ci.yml/badge.svg)](https://github.com/schollz/audiomorph/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/schollz/audiomorph/branch/main/graph/badge.svg)](https://codecov.io/gh/schollz/audiomorph)
[![Release](https://img.shields.io/github/v/release/schollz/audiomorph)](https://github.com/schollz/audiomorph/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/schollz/audiomorph.svg)](https://pkg.go.dev/github.com/schollz/audiomorph)

A Go library and CLI tool for decoding and encoding audio files across multiple formats. audiomorph provides a unified interface for reading audio data from WAV, AIFF, MP3, OGG, and FLAC files, and encoding to WAV, AIFF, MP3, and FLAC formats.

## How It Works

audiomorph decodes audio files into a common in-memory representation (`Audio` struct) containing deinterlaced PCM data, then encodes that data to your desired output format. The library handles format-specific quirks and provides a consistent API regardless of the underlying codec.

### Dependencies

This project is grateful for and depends on these excellent Go libraries:

- [faiface/beep](https://github.com/faiface/beep) - Audio decoding for MP3 and OGG Vorbis
- [go-audio/wav](https://github.com/go-audio/wav) - WAV file encoding/decoding
- [go-audio/aiff](https://github.com/go-audio/aiff) - AIFF file encoding/decoding
- [mewkiz/flac](https://github.com/mewkiz/flac) - FLAC file decoding
- [schollz/goflac](https://github.com/schollz/goflac) - FLAC file encoding
- [braheezy/shine-mp3](https://github.com/braheezy/shine-mp3) - MP3 file encoding

## API

The library exposes two primary functions:

```go
// Decode audio from file (supports WAV, AIFF, MP3, OGG, FLAC)
audio, err := audiomorph.DecodeFile("input.mp3")

// Encode audio to file (supports WAV, AIFF, MP3, FLAC)
err = audiomorph.EncodeFile(audio, "output.wav")
```

The `Audio` struct provides access to all audio properties:

```go
type Audio struct {
    NumChannels int      // Number of audio channels
    SampleRate  int      // Sample rate in Hz
    BitDepth    int      // Bit depth (bits per sample)
    Data        [][]int  // Deinterlaced PCM data [channel][sample]
    Duration    float64  // Duration in seconds
}
```

## Usage

### Installation

```bash
go install github.com/schollz/audiomorph/cmd/audiomorph@latest
```

### Command Line

Display audio file statistics:

```bash
audiomorph input.mp3
```

Transform audio between formats:

```bash
audiomorph input.mp3 output.wav
audiomorph input.flac output.mp3
audiomorph input.wav output.aiff
```

### Library Usage

```go
import "github.com/schollz/audiomorph"

// Decode audio file
audio, err := audiomorph.DecodeFile("input.mp3")
if err != nil {
    log.Fatal(err)
}

// Process audio data...
fmt.Printf("Duration: %.2f seconds\n", audio.Duration)
fmt.Printf("Sample rate: %d Hz\n", audio.SampleRate)

// Encode to different format
err = audiomorph.EncodeFile(audio, "output.wav")
if err != nil {
    log.Fatal(err)
}
```

## License

MIT
