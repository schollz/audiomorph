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
	"github.com/schollz/goflac"
	interpolators "github.com/schollz/interpolation"
)

// OptionUseChannels specifies which channels to use when encoding audio.
func OptionUseChannels(channels []int) Option {
	return func(a *Audio) {
		a.useChannels = channels
	}
}

// OptionSampleRate specifies the target sample rate for encoding audio.
func OptionSampleRate(sampleRate int) Option {
	return func(a *Audio) {
		a.targetSampleRate = sampleRate
	}
}

// OptionInterpolationMethod specifies the interpolation method to use for sample rate conversion.
// Valid methods are: "linear", "cubic", "hermite", "lanczos2", "lanczos3", "bspline3", "bspline5", "monotonic"
func OptionInterpolationMethod(method string) Option {
	return func(a *Audio) {
		a.interpolationMethod = method
	}
}

// convertSampleRate converts audio data to a different sample rate using interpolation
func convertSampleRate(audio *Audio, targetSampleRate int, method string) error {
	// If no target sample rate is specified or it matches current, no conversion needed
	if targetSampleRate == 0 || targetSampleRate == audio.SampleRate {
		return nil
	}

	// Default to linear interpolation if not specified
	if method == "" {
		method = "linear"
	}

	// Calculate the new number of samples
	numSamples := len(audio.Data[0])
	ratio := float64(targetSampleRate) / float64(audio.SampleRate)
	newNumSamples := int(float64(numSamples) * ratio)

	// Map method string to interpolator type
	var interpType interpolators.InterpolatorType
	switch method {
	case "linear":
		interpType = interpolators.Linear
	case "cubic":
		interpType = interpolators.CubicSpline
	case "hermite":
		interpType = interpolators.Hermite4
	case "lanczos2":
		interpType = interpolators.Lanczos2
	case "lanczos3":
		interpType = interpolators.Lanczos3
	case "bspline3":
		interpType = interpolators.BSpline3
	case "bspline5":
		interpType = interpolators.BSpline5
	case "monotonic":
		interpType = interpolators.MonotonicCubic
	default:
		return fmt.Errorf("unsupported interpolation method: %s", method)
	}

	// Create new data array for resampled audio
	newData := make([][]int, audio.NumChannels)

	// Resample each channel
	for ch := 0; ch < audio.NumChannels; ch++ {
		// Convert int samples to float64 for interpolation
		inData := make([]float64, numSamples)
		for i := 0; i < numSamples; i++ {
			inData[i] = float64(audio.Data[ch][i])
		}

		// Perform interpolation
		outData, err := interpolators.Interpolate(inData, newNumSamples, interpType)
		if err != nil {
			return fmt.Errorf("failed to interpolate channel %d: %w", ch, err)
		}

		// Convert float64 results back to int
		newData[ch] = make([]int, newNumSamples)
		for i := 0; i < newNumSamples; i++ {
			newData[ch][i] = int(outData[i])
		}
	}

	// Update the audio struct with resampled data
	audio.Data = newData
	audio.SampleRate = targetSampleRate
	audio.Duration = float64(newNumSamples) / float64(targetSampleRate)

	return nil
}

// EncodeFile encodes an Audio struct to a file based on the filename extension
func EncodeFile(audio *Audio, filename string, options ...Option) error {
	// Apply options
	for _, option := range options {
		option(audio)
	}

	// Apply sample rate conversion if specified
	if audio.targetSampleRate > 0 && audio.targetSampleRate != audio.SampleRate {
		if err := convertSampleRate(audio, audio.targetSampleRate, audio.interpolationMethod); err != nil {
			return fmt.Errorf("failed to convert sample rate: %w", err)
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".wav":
		return encodeWAV(audio, filename)
	case ".aif", ".aiff":
		return encodeAIFF(audio, filename)
	case ".mp3":
		return encodeMP3(audio, filename)
	case ".flac":
		return encodeFLAC(audio, filename)
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

	// Determine number of channels (mono conversion if requested)
	numChannels := audio.NumChannels
	useChannels := audio.useChannels
	if len(useChannels) > 0 {
		numChannels = len(useChannels)
	}

	// Create WAV encoder
	encoder := wav.NewEncoder(f, audio.SampleRate, audio.BitDepth, numChannels, 1)

	// Interlace the audio data from [][]int to []int
	numSamples := len(audio.Data[0])
	interlacedData := make([]int, numSamples*numChannels)
	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < numChannels; ch++ {
			// Use the specified channel from useChannels, or original channel if not specified
			sourceChannel := ch
			if len(useChannels) > 0 {
				sourceChannel = useChannels[ch]
			}
			interlacedData[i*numChannels+ch] = audio.Data[sourceChannel][i]
		}
	}

	// Create PCM buffer
	buf := &goaudio.IntBuffer{
		Format: &goaudio.Format{
			NumChannels: numChannels,
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

	// Determine number of channels (mono conversion if requested)
	numChannels := audio.NumChannels
	if len(audio.useChannels) > 0 {
		numChannels = len(audio.useChannels)
	}

	// Create AIFF encoder
	encoder := aiff.NewEncoder(f, audio.SampleRate, audio.BitDepth, numChannels)

	// Interlace the audio data from [][]int to []int
	numSamples := len(audio.Data[0])
	interlacedData := make([]int, numSamples*numChannels)
	useChannels := audio.useChannels
	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < numChannels; ch++ {
			// Use the specified channel from useChannels, or original channel if not specified
			sourceChannel := ch
			if len(useChannels) > 0 {
				sourceChannel = useChannels[ch]
			}
			interlacedData[i*numChannels+ch] = audio.Data[sourceChannel][i]
		}
	}

	// Create PCM buffer
	buf := &goaudio.IntBuffer{
		Format: &goaudio.Format{
			NumChannels: numChannels,
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

	// Determine number of channels (mono conversion if requested)
	numChannels := audio.NumChannels
	if len(audio.useChannels) > 0 {
		numChannels = len(audio.useChannels)
	}

	// Create MP3 encoder
	encoder := mp3.NewEncoder(audio.SampleRate, numChannels)

	// Convert audio data to int16 format expected by MP3 encoder
	numSamples := len(audio.Data[0])
	// For MP3 encoding, we need interleaved int16 samples
	int16Data := make([]int16, numSamples*numChannels)

	// Convert and scale samples to int16 range
	bitDepth := audio.BitDepth
	scale := float64(1<<15) / float64(int64(1)<<uint(bitDepth-1))
	useChannels := audio.useChannels
	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < numChannels; ch++ {
			sourceChannel := ch
			if len(useChannels) > 0 {
				sourceChannel = useChannels[ch]
			}
			sample := audio.Data[sourceChannel][i]
			// Scale to int16 range
			scaledSample := int16(float64(sample) * scale)
			int16Data[i*numChannels+ch] = scaledSample
		}
	}

	// Write MP3 data
	if err := encoder.Write(f, int16Data); err != nil {
		return fmt.Errorf("failed to write MP3 data: %w", err)
	}

	return nil
}

// encodeFLAC encodes audio data to a FLAC file
func encodeFLAC(audio *Audio, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create FLAC file: %w", err)
	}
	defer f.Close()

	// Determine number of channels (mono conversion if requested)
	numChannels := audio.NumChannels
	if len(audio.useChannels) > 0 {
		numChannels = len(audio.useChannels)
	}

	// Create FLAC encoder
	encoder, err := goflac.NewEncoder(f, uint32(audio.SampleRate), uint8(numChannels), uint8(audio.BitDepth))
	if err != nil {
		return fmt.Errorf("failed to create FLAC encoder: %w", err)
	}

	// Convert audio data from [][]int to [][]int32
	numSamples := len(audio.Data[0])
	samples := make([][]int32, numChannels)
	useChannels := audio.useChannels
	for ch := 0; ch < numChannels; ch++ {
		samples[ch] = make([]int32, numSamples)
		for i := 0; i < numSamples; i++ {
			sourceChannel := ch
			if len(useChannels) > 0 {
				sourceChannel = useChannels[ch]
			}
			samples[ch][i] = int32(audio.Data[sourceChannel][i])
		}
	}

	// Encode samples
	if err := encoder.Encode(samples); err != nil {
		return fmt.Errorf("failed to encode FLAC data: %w", err)
	}

	return nil
}
