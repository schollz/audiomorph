package audiomorph

import (
	"path/filepath"
	"testing"
)

func TestDecodeWAV(t *testing.T) {
	filename := filepath.Join("data", "windchimes.wav")

	audio, err := DecodeFile(filename)
	if err != nil {
		t.Fatalf("Failed to decode WAV file: %v", err)
	}

	// Verify that we got some data
	if audio == nil {
		t.Fatal("Audio is nil")
	}

	// Check that basic fields are populated
	if audio.NumChannels <= 0 {
		t.Errorf("Expected NumChannels > 0, got %d", audio.NumChannels)
	}
	if audio.SampleRate <= 0 {
		t.Errorf("Expected SampleRate > 0, got %d", audio.SampleRate)
	}
	if audio.BitDepth <= 0 {
		t.Errorf("Expected BitDepth > 0, got %d", audio.BitDepth)
	}
	if len(audio.Data) == 0 {
		t.Error("Expected audio Data to be non-empty")
	}
	if audio.Duration <= 0 {
		t.Errorf("Expected Duration > 0, got %f", audio.Duration)
	}

	t.Logf("WAV Audio Info:")
	t.Logf("  NumChannels: %d", audio.NumChannels)
	t.Logf("  SampleRate: %d", audio.SampleRate)
	t.Logf("  BitDepth: %d", audio.BitDepth)
	t.Logf("  Data length: %d samples", len(audio.Data))
	t.Logf("  Duration: %.2f seconds", audio.Duration)
}

func TestDecodeAIFF(t *testing.T) {
	filename := filepath.Join("data", "windchimes.aiff")

	audio, err := DecodeFile(filename)
	if err != nil {
		t.Fatalf("Failed to decode AIFF file: %v", err)
	}

	// Verify that we got some data
	if audio == nil {
		t.Fatal("Audio is nil")
	}

	// Check that basic fields are populated
	if audio.NumChannels <= 0 {
		t.Errorf("Expected NumChannels > 0, got %d", audio.NumChannels)
	}
	if audio.SampleRate <= 0 {
		t.Errorf("Expected SampleRate > 0, got %d", audio.SampleRate)
	}
	if audio.BitDepth <= 0 {
		t.Errorf("Expected BitDepth > 0, got %d", audio.BitDepth)
	}
	if len(audio.Data) == 0 {
		t.Error("Expected audio Data to be non-empty")
	}
	if audio.Duration <= 0 {
		t.Errorf("Expected Duration > 0, got %f", audio.Duration)
	}

	t.Logf("AIFF Audio Info:")
	t.Logf("  NumChannels: %d", audio.NumChannels)
	t.Logf("  SampleRate: %d", audio.SampleRate)
	t.Logf("  BitDepth: %d", audio.BitDepth)
	t.Logf("  Data length: %d samples", len(audio.Data))
	t.Logf("  Duration: %.2f seconds", audio.Duration)
}

func TestDecodeUnsupportedFormat(t *testing.T) {
	filename := filepath.Join("data", "windchimes.mp3")

	_, err := DecodeFile(filename)
	if err == nil {
		t.Fatal("Expected error for unsupported format, got nil")
	}
}
