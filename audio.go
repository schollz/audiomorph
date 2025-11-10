package audiomorph

// Audio represents decoded audio data
type Audio struct {
	NumChannels int
	SampleRate  int
	BitDepth    int
	Data        [][]int // Data[channel][sample] - deinterlaced audio data
	Duration    float64 // in seconds
	useChannels []int
}

// Option is the type all options need to adhere to
type Option func(a *Audio)
