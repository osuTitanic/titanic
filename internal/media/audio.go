package media

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	ffmpeg "github.com/Lekuruu/ffmpeg-go"
)

// ExtractAudioSnippet extracts a snippet from the audio byte
// slice, compresses it, and saves it into an MP3 file.
func ExtractAudioSnippet(audioData []byte, offset int, duration int, bitrate int) ([]byte, error) {
	beatmapAudio, err := os.CreateTemp("", ".tmp-audio-*")
	if err != nil {
		return nil, err
	}

	inputPath := beatmapAudio.Name()
	defer os.Remove(inputPath)

	_, err = beatmapAudio.Write(audioData)
	if err != nil {
		beatmapAudio.Close()
		return nil, err
	}
	if err := beatmapAudio.Close(); err != nil {
		return nil, err
	}
	return ExtractAudioSnippetFromPath(inputPath, offset, duration, bitrate)
}

// ExtractAudioSnippetFromReader extracts a snippet from the reader
// compresses it, and saves it into an MP3 file.
func ExtractAudioSnippetFromReader(audioReader io.Reader, offset int, duration int, bitrate int) ([]byte, error) {
	beatmapAudio, err := os.CreateTemp("", ".tmp-audio-*")
	if err != nil {
		return nil, err
	}

	inputPath := beatmapAudio.Name()
	defer os.Remove(inputPath)

	_, err = io.Copy(beatmapAudio, audioReader)
	if err != nil {
		beatmapAudio.Close()
		return nil, err
	}
	if err := beatmapAudio.Close(); err != nil {
		return nil, err
	}
	return ExtractAudioSnippetFromPath(inputPath, offset, duration, bitrate)
}

// ExtractAudioSnippetFromReader extracts a snippet from the audio
// filepath compresses it, and saves it into an MP3 file.
func ExtractAudioSnippetFromPath(inputPath string, offset int, duration int, bitrate int) (result []byte, err error) {
	if offset < 0 {
		// Use ~40% of audio duration as offset by default
		offset, err = defaultOffset(inputPath)
		if err != nil {
			return nil, err
		}
	}

	outputPath := inputPath + ".preview.mp3"
	defer os.Remove(outputPath)

	inputArgs := ffmpeg.KwArgs{"ss": offset}
	outputArgs := ffmpeg.KwArgs{
		"t":      duration,
		"ab":     bitrate,
		"vn":     "",
		"acodec": "libmp3lame",
		"format": "mp3",
	}
	ffmpeg.LogCompiledCommand = false

	err = ffmpeg.Input(inputPath, inputArgs).
		Output(outputPath, outputArgs).
		Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(outputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO: Maybe return the handle directly here
	// 		 I think the memory usage should be low enough though
	//		 The biggest preview I have seen is ~250 kb

	return io.ReadAll(file)
}

func defaultOffset(filename string) (int, error) {
	probe, err := ffmpeg.Probe(filename)
	if err != nil {
		return 0, fmt.Errorf("probe audio duration: %w", err)
	}

	var result struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal([]byte(probe), &result); err != nil {
		return 0, fmt.Errorf("decode audio duration: %w", err)
	}

	duration, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("parse audio duration: %w", err)
	}
	return int(duration / 2.5), nil
}
