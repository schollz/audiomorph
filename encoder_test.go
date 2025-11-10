package audiomorph

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// verifySoxCanReadFile verifies that sox can read the given audio file using "sox --i"
func verifySoxCanReadFile(t *testing.T, filename string) {
	t.Helper()
	cmd := exec.Command("sox", "--i", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Sox failed to read file %s: %v\nOutput: %s", filename, err, string(output))
	}
	t.Logf("Sox successfully verified file: %s", filename)
	t.Logf("Sox output:\n%s", string(output))
}

func TestEncodeWAV(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.wav")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new WAV file
	dstFilename := filepath.Join(os.TempDir(), "test_output.wav")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode WAV file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(dstFilename); os.IsNotExist(err) {
		t.Fatal("Encoded WAV file was not created")
	}

	// Decode the encoded file to verify it's valid
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded WAV file: %v", err)
	}

	// Verify basic properties match
	if decodedAudio.NumChannels != audio.NumChannels {
		t.Errorf("NumChannels mismatch: expected %d, got %d", audio.NumChannels, decodedAudio.NumChannels)
	}
	if decodedAudio.SampleRate != audio.SampleRate {
		t.Errorf("SampleRate mismatch: expected %d, got %d", audio.SampleRate, decodedAudio.SampleRate)
	}
	if decodedAudio.BitDepth != audio.BitDepth {
		t.Errorf("BitDepth mismatch: expected %d, got %d", audio.BitDepth, decodedAudio.BitDepth)
	}

	// Verify that sox can read the encoded file
	verifySoxCanReadFile(t, dstFilename)

	t.Logf("WAV Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeAIFF(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.aiff")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new AIFF file
	dstFilename := filepath.Join(os.TempDir(), "test_output.aiff")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode AIFF file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(dstFilename); os.IsNotExist(err) {
		t.Fatal("Encoded AIFF file was not created")
	}

	// Decode the encoded file to verify it's valid
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded AIFF file: %v", err)
	}

	// Verify basic properties match
	if decodedAudio.NumChannels != audio.NumChannels {
		t.Errorf("NumChannels mismatch: expected %d, got %d", audio.NumChannels, decodedAudio.NumChannels)
	}
	if decodedAudio.SampleRate != audio.SampleRate {
		t.Errorf("SampleRate mismatch: expected %d, got %d", audio.SampleRate, decodedAudio.SampleRate)
	}
	if decodedAudio.BitDepth != audio.BitDepth {
		t.Errorf("BitDepth mismatch: expected %d, got %d", audio.BitDepth, decodedAudio.BitDepth)
	}

	// Verify that sox can read the encoded file
	verifySoxCanReadFile(t, dstFilename)

	t.Logf("AIFF Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeMP3(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.mp3")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new MP3 file
	dstFilename := filepath.Join(os.TempDir(), "test_output.mp3")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode MP3 file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(dstFilename); os.IsNotExist(err) {
		t.Fatal("Encoded MP3 file was not created")
	}

	// Decode the encoded file to verify it's valid
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded MP3 file: %v", err)
	}

	// Verify basic properties match (MP3 may have slight differences due to compression)
	if decodedAudio.NumChannels != audio.NumChannels {
		t.Errorf("NumChannels mismatch: expected %d, got %d", audio.NumChannels, decodedAudio.NumChannels)
	}
	if decodedAudio.SampleRate != audio.SampleRate {
		t.Errorf("SampleRate mismatch: expected %d, got %d", audio.SampleRate, decodedAudio.SampleRate)
	}

	// Verify that sox can read the encoded file
	verifySoxCanReadFile(t, dstFilename)

	t.Logf("MP3 Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeFLAC(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.flac")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new FLAC file
	dstFilename := filepath.Join(os.TempDir(), "test_output.flac")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode FLAC file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(dstFilename); os.IsNotExist(err) {
		t.Fatal("Encoded FLAC file was not created")
	}

	// Decode the encoded file to verify it's valid
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded FLAC file: %v", err)
	}

	// Verify basic properties match
	if decodedAudio.NumChannels != audio.NumChannels {
		t.Errorf("NumChannels mismatch: expected %d, got %d", audio.NumChannels, decodedAudio.NumChannels)
	}
	if decodedAudio.SampleRate != audio.SampleRate {
		t.Errorf("SampleRate mismatch: expected %d, got %d", audio.SampleRate, decodedAudio.SampleRate)
	}
	if decodedAudio.BitDepth != audio.BitDepth {
		t.Errorf("BitDepth mismatch: expected %d, got %d", audio.BitDepth, decodedAudio.BitDepth)
	}

	// Verify that sox can read the encoded file
	verifySoxCanReadFile(t, dstFilename)

	t.Logf("FLAC Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeUnsupportedFormat(t *testing.T) {
	// Create a simple audio structure
	audio := &Audio{
		NumChannels: 2,
		SampleRate:  44100,
		BitDepth:    16,
		Data:        [][]int{{0, 1, 2}, {3, 4, 5}},
		Duration:    0.001,
	}

	// Try to encode to an unsupported format
	dstFilename := filepath.Join(os.TempDir(), "test_output.unknown")
	err := EncodeFile(audio, dstFilename)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}

	t.Logf("Unsupported format correctly returns error: %v", err)
}

func TestEncodeRightChannel(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.mp3")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Set option to use only the right channel (channel index 1)
	audio.useChannels = []int{1}

	// Encode to a new WAV file
	dstFilename := filepath.Join(os.TempDir(), "test_output_right_channel.wav")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode WAV file with right channel: %v", err)
	}

	// Decode the encoded file to verify it's valid
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded WAV file: %v", err)
	}

	// Verify that the encoded audio has 1 channel
	if decodedAudio.NumChannels != 1 {
		t.Errorf("Expected 1 channel in encoded audio, got %d", decodedAudio.NumChannels)
	}

	// Verify that the samples match the right channel of the original audio
	originalRightChannel := audio.Data[1]
	encodedChannel := decodedAudio.Data[0]

	if len(originalRightChannel) != len(encodedChannel) {
		t.Fatalf("Sample length mismatch: original %d, encoded %d", len(originalRightChannel), len(encodedChannel))
	}

	for i := range originalRightChannel {
		if originalRightChannel[i] != encodedChannel[i] {
			t.Errorf("Sample mismatch at index %d: expected %d, got %d", i, originalRightChannel[i], encodedChannel[i])
		}
	}

	t.Logf("Right channel encoding test passed")
}

func TestSampleRateConversion(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.wav")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	originalSampleRate := audio.SampleRate
	originalSamples := len(audio.Data[0])
	t.Logf("Original sample rate: %d Hz, samples: %d", originalSampleRate, originalSamples)

	// Test different target sample rates
	testCases := []struct {
		name                string
		targetSampleRate    int
		interpolationMethod string
	}{
		{"Upsample to 48kHz with linear", 48000, "linear"},
		{"Downsample to 22.05kHz with linear", 22050, "linear"},
		{"Upsample to 48kHz with cubic", 48000, "cubic"},
		{"Downsample to 22.05kHz with hermite", 22050, "hermite"},
		{"Upsample to 48kHz with lanczos3", 48000, "lanczos3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Decode fresh audio for each test
			audio, err := DecodeFile(srcFilename)
			if err != nil {
				t.Fatalf("Failed to decode source file: %v", err)
			}

			// Encode to a new WAV file with sample rate conversion
			dstFilename := filepath.Join(os.TempDir(), "test_sample_rate_"+tc.interpolationMethod+".wav")
			defer os.Remove(dstFilename)

			err = EncodeFile(audio, dstFilename,
				OptionSampleRate(tc.targetSampleRate),
				OptionInterpolationMethod(tc.interpolationMethod))
			if err != nil {
				t.Fatalf("Failed to encode with sample rate conversion: %v", err)
			}

			// Decode the encoded file to verify the sample rate was converted
			decodedAudio, err := DecodeFile(dstFilename)
			if err != nil {
				t.Fatalf("Failed to decode encoded file: %v", err)
			}

			// Verify the sample rate matches the target
			if decodedAudio.SampleRate != tc.targetSampleRate {
				t.Errorf("Sample rate mismatch: expected %d, got %d", tc.targetSampleRate, decodedAudio.SampleRate)
			}

			// Verify the number of samples is approximately correct (accounting for ratio)
			expectedSamples := int(float64(originalSamples) * float64(tc.targetSampleRate) / float64(originalSampleRate))
			actualSamples := len(decodedAudio.Data[0])
			tolerance := int(float64(expectedSamples) * 0.01) // 1% tolerance
			if actualSamples < expectedSamples-tolerance || actualSamples > expectedSamples+tolerance {
				t.Errorf("Sample count mismatch: expected ~%d, got %d", expectedSamples, actualSamples)
			}

			t.Logf("Successfully converted sample rate from %d to %d Hz using %s interpolation",
				originalSampleRate, tc.targetSampleRate, tc.interpolationMethod)
			t.Logf("  Original samples: %d, New samples: %d", originalSamples, actualSamples)

			// check that sox can read the encoded file
			verifySoxCanReadFile(t, dstFilename)
		})
	}
}

func TestSampleRateConversionPreservesChannels(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.mp3")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	originalChannels := audio.NumChannels
	originalSampleRate := audio.SampleRate

	// Encode with sample rate conversion
	dstFilename := filepath.Join(os.TempDir(), "test_sample_rate_channels.wav")
	defer os.Remove(dstFilename)

	targetSampleRate := 48000
	err = EncodeFile(audio, dstFilename, OptionSampleRate(targetSampleRate))
	if err != nil {
		t.Fatalf("Failed to encode with sample rate conversion: %v", err)
	}

	// Decode and verify
	decodedAudio, err := DecodeFile(dstFilename)
	if err != nil {
		t.Fatalf("Failed to decode encoded file: %v", err)
	}

	// Verify channels are preserved
	if decodedAudio.NumChannels != originalChannels {
		t.Errorf("Channel count mismatch: expected %d, got %d", originalChannels, decodedAudio.NumChannels)
	}

	// Verify sample rate was converted
	if decodedAudio.SampleRate != targetSampleRate {
		t.Errorf("Sample rate mismatch: expected %d, got %d", targetSampleRate, decodedAudio.SampleRate)
	}

	t.Logf("Sample rate conversion preserved %d channels", originalChannels)
	t.Logf("Converted from %d Hz to %d Hz", originalSampleRate, targetSampleRate)
}

func TestInvalidInterpolationMethod(t *testing.T) {
	// Create a simple audio structure
	audio := &Audio{
		NumChannels: 1,
		SampleRate:  44100,
		BitDepth:    16,
		Data:        [][]int{{0, 100, 200, 300, 400}},
		Duration:    0.0001,
	}

	// Try to convert with an invalid interpolation method
	err := convertSampleRate(audio, 48000, "invalid_method")
	if err == nil {
		t.Fatal("Expected error for invalid interpolation method, got nil")
	}

	t.Logf("Invalid interpolation method correctly returns error: %v", err)
}
