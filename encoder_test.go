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

func TestEncodeOGG(t *testing.T) {
	// Decode an existing audio file
	srcFilename := filepath.Join("data", "wilhelm.ogg")
	audio, err := DecodeFile(srcFilename)
	if err != nil {
		t.Fatalf("Failed to decode source file: %v", err)
	}

	// Encode to a new OGG file
	dstFilename := filepath.Join(os.TempDir(), "test_output.ogg")
	defer os.Remove(dstFilename)

	err = EncodeFile(audio, dstFilename)
	if err != nil {
		t.Fatalf("Failed to encode OGG file: %v", err)
	}

	// Verify the file was created
	stat, err := os.Stat(dstFilename)
	if err != nil {
		t.Fatal("Encoded OGG file was not created")
	}

	// Verify the file has reasonable size (should be non-zero)
	if stat.Size() == 0 {
		t.Fatal("Encoded OGG file is empty")
	}

	t.Logf("OGG Encode test passed")
	t.Logf("  NumChannels: %d", audio.NumChannels)
	t.Logf("  SampleRate: %d", audio.SampleRate)
	t.Logf("  BitDepth: %d", audio.BitDepth)
	t.Logf("  Samples: %d", len(audio.Data[0]))
	t.Logf("  Output file size: %.2f KB", float64(stat.Size())/1024)
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
