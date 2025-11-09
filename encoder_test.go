package audiomorph

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeWAV(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "windchimes.wav")
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

	t.Logf("WAV Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeAIFF(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "windchimes.aiff")
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

	t.Logf("AIFF Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeMP3(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "windchimes.mp3")
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

	t.Logf("MP3 Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeFLAC(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "windchimes.flac")
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

	t.Logf("FLAC Encode/Decode test passed")
	t.Logf("  NumChannels: %d", decodedAudio.NumChannels)
	t.Logf("  SampleRate: %d", decodedAudio.SampleRate)
	t.Logf("  BitDepth: %d", decodedAudio.BitDepth)
	t.Logf("  Samples: %d", len(decodedAudio.Data[0]))
}

func TestEncodeOGG(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "windchimes.ogg")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new OGG file (should fail with unsupported error)
	dstFilename := filepath.Join(os.TempDir(), "test_output.ogg")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err == nil {
		t.Fatal("Expected error for OGG encoding (not yet supported), got nil")
	}

	t.Logf("OGG encoding correctly returns unsupported error: %v", err)
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
