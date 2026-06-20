package media

import (
	"io"
	"os"

	ffmpeg "github.com/Lekuruu/ffmpeg-go"
)

// ExtractAudioSnippet extracts a snippet from the audio byte slice, compresses it, and saves it into an MP3 file
func ExtractAudioSnippet(audioData []byte, offset int, duration int, bitrate int) ([]byte, error) {
	beatmapAudio, err := os.CreateTemp("", ".tmp_audio_*")
	if err != nil {
		return nil, err
	}

	_, err = beatmapAudio.Write(audioData)
	if err != nil {
		return nil, err
	}

	inputArgs := ffmpeg.KwArgs{"ss": offset}
	outputArgs := ffmpeg.KwArgs{"t": duration, "ab": bitrate}
	ffmpeg.LogCompiledCommand = false

	err = ffmpeg.Input(beatmapAudio.Name(), inputArgs).
		Output(beatmapAudio.Name()+".mp3", outputArgs).
		Run()

	if err != nil {
		return nil, err
	}

	file, err := os.Open(beatmapAudio.Name() + ".mp3")
	if err != nil {
		return nil, err
	}

	defer func() {
		file.Close()
		os.Remove(beatmapAudio.Name())
		os.Remove(beatmapAudio.Name() + ".mp3")
	}()

	return io.ReadAll(file)
}
