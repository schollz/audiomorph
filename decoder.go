package audiomorph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/go-audio/aiff"
	"github.com/go-audio/wav"
	"github.com/mewkiz/flac"
)

// Audio represents decoded audio data
type Audio struct {
	NumChannels int
	SampleRate  int
	BitDepth    int
	Data        [][]int // Data[channel][sample] - deinterlaced audio data
	Duration    float64 // in seconds
	Mono        bool    // If true, encode only the first channel as mono
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

	// Deinterlace PCM data from []int to [][]int
	numChannels := int(format.NumChannels)
	numSamples := len(buf.Data) / numChannels
	data := make([][]int, numChannels)
	for ch := 0; ch < numChannels; ch++ {
		data[ch] = make([]int, numSamples)
	}

	for i := 0; i < len(buf.Data); i++ {
		ch := i % numChannels
		sample := i / numChannels
		data[ch][sample] = buf.Data[i]
	}

	// Calculate duration
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

	// Deinterlace PCM data from []int to [][]int
	format := buf.Format
	numChannels := int(format.NumChannels)
	numSamples := len(buf.Data) / numChannels
	data := make([][]int, numChannels)
	for ch := 0; ch < numChannels; ch++ {
		data[ch] = make([]int, numSamples)
	}

	for i := 0; i < len(buf.Data); i++ {
		ch := i % numChannels
		sample := i / numChannels
		data[ch][sample] = buf.Data[i]
	}

	// Calculate duration
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
	stream, err := flac.ParseFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FLAC file: %w", err)
	}
	defer stream.Close()

	// Get metadata
	info := stream.Info
	numChannels := int(info.NChannels)
	sampleRate := int(info.SampleRate)
	bitDepth := int(info.BitsPerSample)
	totalSamples := int(info.NSamples)

	// Initialize deinterlaced data structure
	data := make([][]int, numChannels)
	for ch := 0; ch < numChannels; ch++ {
		data[ch] = make([]int, 0, totalSamples)
	}

	// Read all frames
	for {
		frame, err := stream.ParseNext()
		if err != nil {
			break
		}

		// Process each subframe (channel)
		for i := 0; i < len(frame.Subframes[0].Samples); i++ {
			for ch := 0; ch < numChannels; ch++ {
				sample := frame.Subframes[ch].Samples[i]
				data[ch] = append(data[ch], int(sample))
			}
		}
	}

	// Calculate duration
	duration := float64(totalSamples) / float64(sampleRate)

	return &Audio{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
		BitDepth:    bitDepth,
		Data:        data,
		Duration:    duration,
	}, nil
}

// streamToAudio converts a beep.StreamSeekCloser to an Audio struct
func streamToAudio(streamer beep.StreamSeekCloser, format beep.Format) (*Audio, error) {
	// Get the total number of samples
	length := streamer.Len()

	// Initialize deinterlaced data structure
	numChannels := format.NumChannels
	data := make([][]int, numChannels)
	for ch := 0; ch < numChannels; ch++ {
		data[ch] = make([]int, 0, length)
	}

	bufSize := 512
	buf := make([][2]float64, bufSize)
	totalRead := 0

	for totalRead < length {
		n, ok := streamer.Stream(buf)
		if !ok && n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			for ch := 0; ch < numChannels; ch++ {
				// Convert float64 [-1, 1] to int based on precision
				sample := buf[i][ch]
				// Scale by bit depth (precision is in bytes)
				bitDepth := format.Precision * 8
				maxVal := float64(int64(1) << uint(bitDepth-1))
				intVal := int(sample * maxVal)
				data[ch] = append(data[ch], intVal)
			}
		}
		totalRead += n

		if streamer.Err() != nil {
			return nil, fmt.Errorf("error streaming audio: %w", streamer.Err())
		}
	}

	// Calculate duration
	numSamples := len(data[0])
	duration := float64(numSamples) / float64(format.SampleRate)

	return &Audio{
		NumChannels: format.NumChannels,
		SampleRate:  int(format.SampleRate),
		BitDepth:    format.Precision * 8,
		Data:        data,
		Duration:    duration,
	}, nil
}
