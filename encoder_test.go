package audiomorph

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeWAV(t *testing.T) {
	// First decode an existing WAV file
	inputFile := filepath.Join("data", "windchimes.wav")
	audio, err := DecodeFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to decode WAV file: %v", err)
	}

	// Encode to a new WAV file
	outputFile := filepath.Join(t.TempDir(), "output.wav")
	err = EncodeFile(audio, outputFile)
	if err != nil {
		t.Fatalf("Failed to encode WAV file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Decode the output file to verify it's valid
	decodedAudio, err := DecodeFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to decode output WAV file: %v", err)
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
	if len(decodedAudio.Data) != len(audio.Data) {
		t.Errorf("Data length mismatch: expected %d, got %d", len(audio.Data), len(decodedAudio.Data))
	}

	t.Logf("Successfully encoded and verified WAV file")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Data length: %d samples", len(decodedAudio.Data))
	t.Logf("  Duration: %.2f seconds", decodedAudio.Duration)
}

func TestEncodeWAVRoundTrip(t *testing.T) {
	// First decode an existing WAV file
	inputFile := filepath.Join("data", "windchimes.wav")
	original, err := DecodeFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to decode WAV file: %v", err)
	}

	// Encode to a new WAV file
	outputFile := filepath.Join(t.TempDir(), "roundtrip.wav")
	err = EncodeFile(original, outputFile)
	if err != nil {
		t.Fatalf("Failed to encode WAV file: %v", err)
	}

	// Decode the output file
	decoded, err := DecodeFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to decode output WAV file: %v", err)
	}

	// Compare the audio data sample by sample
	if len(decoded.Data) != len(original.Data) {
		t.Fatalf("Data length mismatch: expected %d, got %d", len(original.Data), len(decoded.Data))
	}

	// Check first few samples to ensure data integrity
	samplesToCheck := 100
	if len(decoded.Data) < samplesToCheck {
		samplesToCheck = len(decoded.Data)
	}

	for i := 0; i < samplesToCheck; i++ {
		if decoded.Data[i] != original.Data[i] {
			t.Errorf("Sample %d mismatch: expected %d, got %d", i, original.Data[i], decoded.Data[i])
		}
	}

	t.Logf("Round-trip test successful for %d samples", len(decoded.Data))
}

func TestEncodeNilAudio(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "nil.wav")
	err := EncodeFile(nil, outputFile)
	if err == nil {
		t.Fatal("Expected error when encoding nil audio, got nil")
	}
	t.Logf("Got expected error: %v", err)
}

func TestEncodeUnsupportedFormat(t *testing.T) {
	// Create a simple audio struct
	audio := &Audio{
		NumChannels: 2,
		SampleRate:  44100,
		BitDepth:    16,
		Data:        []int{0, 0, 100, 100, 0, 0},
		Duration:    0.001,
	}

	outputFile := filepath.Join(t.TempDir(), "output.mp3")
	err := EncodeFile(audio, outputFile)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}
	t.Logf("Got expected error: %v", err)
}

func TestEncodeSimpleAudio(t *testing.T) {
	// Create a simple audio struct with a short sine-wave-like pattern
	audio := &Audio{
		NumChannels: 1,
		SampleRate:  44100,
		BitDepth:    16,
		Data:        make([]int, 44100), // 1 second of audio
		Duration:    1.0,
	}

	// Fill with a simple pattern
	for i := range audio.Data {
		audio.Data[i] = int(float64(i) * 100.0 / float64(len(audio.Data)))
	}

	// Encode to WAV
	outputFile := filepath.Join(t.TempDir(), "simple.wav")
	err := EncodeFile(audio, outputFile)
	if err != nil {
		t.Fatalf("Failed to encode simple audio: %v", err)
	}

	// Verify the file was created and can be decoded
	decoded, err := DecodeFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to decode simple audio: %v", err)
	}

	// Verify properties
	if decoded.NumChannels != audio.NumChannels {
		t.Errorf("NumChannels mismatch: expected %d, got %d", audio.NumChannels, decoded.NumChannels)
	}
	if decoded.SampleRate != audio.SampleRate {
		t.Errorf("SampleRate mismatch: expected %d, got %d", audio.SampleRate, decoded.SampleRate)
	}

	t.Logf("Successfully encoded and decoded simple audio")
}
