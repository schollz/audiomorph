package audiomorph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/aiff"
	"github.com/go-audio/wav"
)

// Audio represents decoded audio data
type Audio struct {
	NumChannels int
	SampleRate  int
	BitDepth    int
	Data        []int
	Duration    float64 // in seconds
}

// DecodeFile decodes a WAV or AIF/AIFF file and returns an Audio struct
func DecodeFile(filename string) (*Audio, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return decodeWAV(filename)
	case ".aif", ".aiff":
		return decodeAIFF(filename)
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
