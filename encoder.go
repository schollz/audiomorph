package audiomorph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/braheezy/shine-mp3/pkg/mp3"
	"github.com/go-audio/aiff"
	goaudio "github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// EncodeFile encodes an Audio struct to a file based on the filename extension
func EncodeFile(audio *Audio, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return encodeWAV(audio, filename)
	case ".aif", ".aiff":
		return encodeAIFF(audio, filename)
	case ".mp3":
		return encodeMP3(audio, filename)
	case ".ogg":
		return fmt.Errorf("OGG Vorbis encoding is not yet supported (OGG container library available but Vorbis codec encoder not available)")
	case ".flac":
		return fmt.Errorf("FLAC encoding is not yet supported (library available but encoder not implemented)")
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

// encodeWAV encodes audio data to a WAV file
func encodeWAV(audio *Audio, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %w", err)
	}
	defer f.Close()

	// Create WAV encoder
	encoder := wav.NewEncoder(f, audio.SampleRate, audio.BitDepth, audio.NumChannels, 1)

	// Interlace the audio data from [][]int to []int
	numSamples := len(audio.Data[0])
	interlacedData := make([]int, numSamples*audio.NumChannels)
	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < audio.NumChannels; ch++ {
			interlacedData[i*audio.NumChannels+ch] = audio.Data[ch][i]
		}
	}

	// Create PCM buffer
	buf := &goaudio.IntBuffer{
		Format: &goaudio.Format{
			NumChannels: audio.NumChannels,
			SampleRate:  audio.SampleRate,
		},
		Data:           interlacedData,
		SourceBitDepth: audio.BitDepth,
	}

	// Write the buffer
	if err := encoder.Write(buf); err != nil {
		return fmt.Errorf("failed to write WAV data: %w", err)
	}

	// Close encoder
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close WAV encoder: %w", err)
	}

	return nil
}

// encodeAIFF encodes audio data to an AIFF file
func encodeAIFF(audio *Audio, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create AIFF file: %w", err)
	}
	defer f.Close()

	// Create AIFF encoder
	encoder := aiff.NewEncoder(f, audio.SampleRate, audio.BitDepth, audio.NumChannels)

	// Interlace the audio data from [][]int to []int
	numSamples := len(audio.Data[0])
	interlacedData := make([]int, numSamples*audio.NumChannels)
	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < audio.NumChannels; ch++ {
			interlacedData[i*audio.NumChannels+ch] = audio.Data[ch][i]
		}
	}

	// Create PCM buffer
	buf := &goaudio.IntBuffer{
		Format: &goaudio.Format{
			NumChannels: audio.NumChannels,
			SampleRate:  audio.SampleRate,
		},
		Data:           interlacedData,
		SourceBitDepth: audio.BitDepth,
	}

	// Write the buffer
	if err := encoder.Write(buf); err != nil {
		return fmt.Errorf("failed to write AIFF data: %w", err)
	}

	// Close encoder
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close AIFF encoder: %w", err)
	}

	return nil
}

// encodeMP3 encodes audio data to an MP3 file
func encodeMP3(audio *Audio, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create MP3 file: %w", err)
	}
	defer f.Close()

	// Create MP3 encoder
	encoder := mp3.NewEncoder(audio.SampleRate, audio.NumChannels)

	// Convert audio data to int16 format expected by MP3 encoder
	numSamples := len(audio.Data[0])
	// For MP3 encoding, we need interleaved int16 samples
	int16Data := make([]int16, numSamples*audio.NumChannels)

	// Convert and scale samples to int16 range
	bitDepth := audio.BitDepth
	scale := float64(1<<15) / float64(int64(1)<<uint(bitDepth-1))

	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < audio.NumChannels; ch++ {
			sample := audio.Data[ch][i]
			// Scale to int16 range
			scaledSample := int16(float64(sample) * scale)
			int16Data[i*audio.NumChannels+ch] = scaledSample
		}
	}

	// Write MP3 data
	if err := encoder.Write(f, int16Data); err != nil {
		return fmt.Errorf("failed to write MP3 data: %w", err)
	}

	return nil
}
