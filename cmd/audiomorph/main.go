package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/schollz/audiomorph"
	"github.com/spf13/cobra"
)

// Version is the version of the audiomorph utility
var Version = "dev"

var (
	flagMono bool
)

var rootCmd = &cobra.Command{
	Use:   "audiomorph [input-file] [output-file]",
	Short: "A utility for audio file transformation and analysis",
	Long: `audiomorph is a command-line utility for transforming audio files between different formats
and analyzing audio file properties.

When provided with only an input file, it displays statistics about the audio file.
When provided with both input and output files, it transforms the audio from one format to another.

Supported formats: WAV, AIFF, MP3, OGG, FLAC (for input)
                  WAV, AIFF, MP3, OGG, FLAC (for output)`,
	Version: Version,
	Args:    cobra.RangeArgs(1, 2),
	RunE:    run,
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "audiomorph version %s\n" .Version}}`)
	rootCmd.Flags().BoolVar(&flagMono, "mono", false, "Convert to mono by using only the first channel")
}

func run(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Decode the input file
	audio, err := audiomorph.DecodeFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode input file: %w", err)
	}

	// If no output file is specified, display statistics
	if len(args) == 1 {
		displayStatistics(inputFile, audio)
		return nil
	}

	// Apply mono conversion if requested
	if flagMono {
		audio.Mono = true
	}

	// Transform audio to output file
	outputFile := args[1]
	if err := audiomorph.EncodeFile(audio, outputFile); err != nil {
		return fmt.Errorf("failed to encode output file: %w", err)
	}

	fmt.Printf("Successfully transformed %s to %s\n", inputFile, outputFile)
	return nil
}

func displayStatistics(filename string, audio *audiomorph.Audio) {
	fmt.Printf("Audio File Statistics\n")
	fmt.Printf("=====================\n")
	fmt.Printf("File:         %s\n", filepath.Base(filename))
	fmt.Printf("Format:       %s\n", filepath.Ext(filename))
	fmt.Printf("Channels:     %d\n", audio.NumChannels)
	fmt.Printf("Sample Rate:  %d Hz\n", audio.SampleRate)
	fmt.Printf("Bit Depth:    %d bits\n", audio.BitDepth)
	fmt.Printf("Duration:     %.2f seconds\n", audio.Duration)
	fmt.Printf("Samples:      %d per channel\n", len(audio.Data[0]))
	
	// Calculate file size
	fileInfo, err := os.Stat(filename)
	if err == nil {
		fmt.Printf("File Size:    %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
