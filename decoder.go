package audiomorph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
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

// DecodeFile decodes a WAV, AIF/AIFF, MP3, OGG, or FLAC file and returns an Audio struct
func DecodeFile(filename string) (*Audio, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return decodeWAV(filename)
	case ".aif", ".aiff":
		return decodeAIFF(filename)
	case ".mp3":
		return decodeMP3(filename)
	case ".ogg":
		return decodeOGG(filename)
	case ".flac":
		return decodeFLAC(filename)
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

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MP3 file: %w", err)
	}
	defer streamer.Close()

	return streamToAudio(streamer, format)
}

// decodeOGG decodes an OGG Vorbis file
func decodeOGG(filename string) (*Audio, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open OGG file: %w", err)
	}
	defer f.Close()

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OGG file: %w", err)
	}
	defer streamer.Close()

	return streamToAudio(streamer, format)
}

// decodeFLAC decodes a FLAC file
func decodeFLAC(filename string) (*Audio, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open FLAC file: %w", err)
	}
	defer f.Close()

	streamer, format, err := flac.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode FLAC file: %w", err)
	}
	defer streamer.Close()

	return streamToAudio(streamer, format)
}

// streamToAudio converts a beep.StreamSeekCloser to an Audio struct
func streamToAudio(streamer beep.StreamSeekCloser, format beep.Format) (*Audio, error) {
	// Get the total number of samples
	length := streamer.Len()

	// Read all samples
	data := make([]int, 0, length*format.NumChannels)

	bufSize := 512
	buf := make([][2]float64, bufSize)
	totalRead := 0

	for totalRead < length {
		n, ok := streamer.Stream(buf)
		if !ok && n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			for ch := 0; ch < format.NumChannels; ch++ {
				// Convert float64 [-1, 1] to int based on precision
				sample := buf[i][ch]
				// Scale by bit depth (precision is in bytes)
				bitDepth := format.Precision * 8
				maxVal := float64(int64(1) << uint(bitDepth-1))
				intVal := int(sample * maxVal)
				data = append(data, intVal)
			}
		}
		totalRead += n

		if streamer.Err() != nil {
			return nil, fmt.Errorf("error streaming audio: %w", streamer.Err())
		}
	}

	// Calculate duration
	numSamples := len(data) / format.NumChannels
	duration := float64(numSamples) / float64(format.SampleRate)

	return &Audio{
		NumChannels: format.NumChannels,
		SampleRate:  int(format.SampleRate),
		BitDepth:    format.Precision * 8,
		Data:        data,
		Duration:    duration,
	}, nil
}
