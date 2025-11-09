package audiomorph

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/aiff"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
)

// Audio represents decoded audio data
type Audio struct {
	NumChannels int
	SampleRate  int
	BitDepth    int
	Data        []int
	Duration    float64 // in seconds
}

// DecodeFile decodes a WAV, AIF/AIFF, or MP3 file and returns an Audio struct
func DecodeFile(filename string) (*Audio, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return decodeWAV(filename)
	case ".aif", ".aiff":
		return decodeAIFF(filename)
	case ".mp3":
		return decodeMP3(filename)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// decodeWAV decodes a WAV file
func decodeWAV(filename string) (*Audio, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAV file: %w", err)
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	// Read the format information
	if err := decoder.FwdToPCM(); err != nil {
		return nil, fmt.Errorf("failed to forward to PCM data: %w", err)
	}

	// Get audio format
	format := decoder.Format()

	// Read all PCM data
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to read PCM buffer: %w", err)
	}

	// Convert PCM data to []int
	data := make([]int, len(buf.Data))
	for i, v := range buf.Data {
		data[i] = v
	}

	// Calculate duration
	numSamples := len(data) / int(format.NumChannels)
	duration := float64(numSamples) / float64(format.SampleRate)

	return &Audio{
		NumChannels: int(format.NumChannels),
		SampleRate:  int(format.SampleRate),
		BitDepth:    int(decoder.BitDepth),
		Data:        data,
		Duration:    duration,
	}, nil
}

// decodeAIFF decodes an AIFF/AIF file
func decodeAIFF(filename string) (*Audio, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open AIFF file: %w", err)
	}
	defer f.Close()

	decoder := aiff.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid AIFF file")
	}

	// Read all PCM data
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to read PCM buffer: %w", err)
	}

	// Convert PCM data to []int
	data := make([]int, len(buf.Data))
	for i, v := range buf.Data {
		data[i] = v
	}

	// Get format from the buffer
	format := buf.Format

	// Calculate duration
	numSamples := len(data) / int(format.NumChannels)
	duration := float64(numSamples) / float64(format.SampleRate)

	return &Audio{
		NumChannels: int(format.NumChannels),
		SampleRate:  int(format.SampleRate),
		BitDepth:    int(decoder.BitDepth),
		Data:        data,
		Duration:    duration,
	}, nil
}

// decodeMP3 decodes an MP3 file
func decodeMP3(filename string) (*Audio, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open MP3 file: %w", err)
	}
	defer f.Close()

	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create MP3 decoder: %w", err)
	}

	// Get audio format
	sampleRate := decoder.SampleRate()

	// MP3 is typically 16-bit (but go-mp3 doesn't expose this directly)
	bitDepth := 16

	// Read all audio data
	// The decoder returns interleaved PCM data as bytes
	data := make([]byte, 0)
	buf := make([]byte, 8192)
	for {
		n, err := decoder.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read MP3 data: %w", err)
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	// Convert byte data to []int (16-bit samples, little-endian)
	samples := make([]int, len(data)/2)
	for i := 0; i < len(samples); i++ {
		// Read 16-bit little-endian signed integer
		sample := int(int16(data[i*2]) | int16(data[i*2+1])<<8)
		samples[i] = sample
	}

	// Calculate duration and number of channels
	// go-mp3 always outputs stereo (2 channels)
	numChannels := 2
	numSamples := len(samples) / numChannels
	duration := float64(numSamples) / float64(sampleRate)

	return &Audio{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
		BitDepth:    bitDepth,
		Data:        samples,
		Duration:    duration,
	}, nil
}
