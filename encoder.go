package audiomorph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// EncodeFile encodes an Audio struct to a file based on the file extension
func EncodeFile(a *Audio, filename string) error {
	if a == nil {
		return fmt.Errorf("audio is nil")
	}

	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return encodeWAV(a, filename)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

// encodeWAV encodes an Audio struct to a WAV file
func encodeWAV(a *Audio, filename string) error {
	// Create the output file
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %w", err)
	}
	defer f.Close()

	// Create a WAV encoder
	encoder := wav.NewEncoder(f, a.SampleRate, a.BitDepth, a.NumChannels, 1)
	defer encoder.Close()

	// Convert []int to audio.IntBuffer
	buf := &audio.IntBuffer{
		Data: a.Data,
		Format: &audio.Format{
			NumChannels: a.NumChannels,
			SampleRate:  a.SampleRate,
		},
	}

	// Write the audio data
	if err := encoder.Write(buf); err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	return nil
}
