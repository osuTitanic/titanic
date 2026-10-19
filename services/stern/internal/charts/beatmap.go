package charts

import (
	"bytes"
	"fmt"
	"math"

	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

const (
	BeatmapRatingChartWidth  = 200
	BeatmapRatingChartHeight = 40
	BeatmapFailChartWidth    = 200
	BeatmapFailChartHeight   = 60
)

var (
	beatmapChartBackground = drawing.Color{R: 228, G: 219, B: 251, A: 255}
	failColor              = drawing.Color{R: 255, A: 255}
	retryColor             = drawing.Color{R: 22, G: 125, B: 212, A: 255}
	failGridColor          = drawing.Color{R: 221, G: 208, B: 249, A: 255}
)

var (
	ratingGradientStart = drawing.Color{R: 104, G: 112, B: 228, A: 255}
	ratingGradientEnd   = drawing.Color{R: 156, G: 162, B: 245, A: 255}
	ratingOutlineColor  = drawing.ColorBlack
)

// GenerateBeatmapRatingChart renders the distribution of ratings 1 through 10
func GenerateBeatmapRatingChart(ratings [10]int64) ([]byte, error) {
	const width, height = BeatmapRatingChartWidth, BeatmapRatingChartHeight

	renderer, err := newRenderer(width, height)
	if err != nil {
		return nil, fmt.Errorf("create beatmap rating chart renderer: %w", err)
	}

	maxCount := int64(0)
	for _, count := range ratings {
		maxCount = max(maxCount, count)
	}

	barWidth := width / len(ratings)
	for index, count := range ratings {
		barHeight := 0
		if maxCount > 0 {
			barHeight = min(scale(count, maxCount, height), height)
		}

		left := index * barWidth
		top := height - barHeight
		outlineTop := min(top, height-1)

		// Draw the vertical gradient
		if barHeight > 0 {
			for offset := range barHeight {
				drawBox(
					renderer,
					left, top+offset, left+barWidth, top+offset+1,
					interpolateColor(ratingGradientStart, ratingGradientEnd, offset, barHeight-1),
				)
			}
		}

		// Add outline to each bar (empty ones too)
		drawOutline(
			renderer,
			left, outlineTop, left+barWidth, height,
			ratingOutlineColor,
		)
	}

	return saveBeatmapChart(renderer, "beatmap rating")
}

// GenerateBeatmapFailChart renders accumulated fail and exit (retry) counts over the percentage of the beatmap completed.
func GenerateBeatmapFailChart(fails [100]int64, exits [100]int64) ([]byte, error) {
	const width, height = BeatmapFailChartWidth, BeatmapFailChartHeight

	renderer, err := newRenderer(width, height)
	if err != nil {
		return nil, err
	}

	// Draw the grey grid in the background at around 25%-ish intervals
	for x := width / 4; x < width; x += width / 4 {
		drawBox(renderer, x, 0, x+1, height, failGridColor)
	}

	maxTotal := int64(0)
	for index := range fails {
		maxTotal = max(maxTotal, fails[index]+exits[index])
	}

	if maxTotal > 0 {
		barWidth := width / len(fails)

		for i, count := range fails {
			failHeight := scale(count, maxTotal, height)
			totalHeight := scale(count+exits[i], maxTotal, height)
			x := i * barWidth

			// Draw both fails & retry bars
			// Fails should be on top of the retry bars
			drawBox(renderer, x, height-failHeight, x+barWidth, height, failColor)
			drawBox(renderer, x, height-totalHeight, x+barWidth, height-failHeight, retryColor)
		}
	}

	if err := drawFailLegend(renderer); err != nil {
		return nil, err
	}
	return saveBeatmapChart(renderer, "beatmap fail")
}

func newRenderer(width int, height int) (chart.Renderer, error) {
	renderer, err := chart.PNG(width, height)
	if err != nil {
		return nil, err
	}
	renderer.SetDPI(chart.DefaultDPI)
	drawBox(renderer, 0, 0, width, height, beatmapChartBackground)
	return renderer, nil
}

func saveBeatmapChart(renderer chart.Renderer, name string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := renderer.Save(buffer); err != nil {
		return nil, fmt.Errorf("render %s chart: %w", name, err)
	}
	return buffer.Bytes(), nil
}

func drawFailLegend(renderer chart.Renderer) error {
	font, err := chart.GetDefaultFont()
	if err != nil {
		return fmt.Errorf("load beatmap fail chart font: %w", err)
	}

	drawBox(renderer, 0, 0, 102, 12, beatmapChartBackground)
	drawBox(renderer, 0, 1, 10, 11, retryColor)
	drawBox(renderer, 64, 1, 74, 11, failColor)

	drawOutline(renderer, 0, 1, 10, 11, drawing.ColorBlack)
	drawOutline(renderer, 64, 1, 74, 11, drawing.ColorBlack)

	renderer.SetFont(font)
	renderer.SetFontColor(drawing.ColorBlack)
	renderer.SetFontSize(7)
	renderer.Text("Retry", 14, 10)
	renderer.Text("Fail", 78, 10)
	return nil
}

func drawBox(renderer chart.Renderer, left, top, right, bottom int, color drawing.Color) {
	if left >= right || top >= bottom {
		return
	}

	chart.Draw.Box(
		renderer,
		chart.Box{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		},
		chart.Style{
			FillColor: color,
		},
	)
}

func drawOutline(renderer chart.Renderer, left int, top int, right int, bottom int, color drawing.Color) {
	drawBox(renderer, left, top, right, top+1, color)
	drawBox(renderer, left, bottom-1, right, bottom, color)
	drawBox(renderer, left, top, left+1, bottom, color)
	drawBox(renderer, right-1, top, right, bottom, color)
}

func interpolateColor(start, end drawing.Color, position, maximum int) drawing.Color {
	if maximum <= 0 {
		return start
	}

	mix := func(a, b uint8) uint8 {
		return uint8(int(a) + (int(b)-int(a))*position/maximum)
	}

	return drawing.Color{
		R: mix(start.R, end.R),
		G: mix(start.G, end.G),
		B: mix(start.B, end.B),
		A: mix(start.A, end.A),
	}
}

func scale(value, maximum int64, height int) int {
	return int(math.Round(float64(height) * float64(value) / float64(maximum)))
}
