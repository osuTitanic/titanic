package charts

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

var (
	backgroundColor = drawing.Color{R: 240, G: 236, B: 250, A: 255}
	osuFillColor    = drawing.Color{R: 255, G: 232, B: 250, A: 255}
	osuLineColor    = drawing.Color{R: 255, G: 205, B: 246, A: 255}
	ircFillColor    = drawing.Color{R: 255, G: 227, B: 249, A: 255}
	ircLineColor    = drawing.Color{R: 144, G: 150, B: 188, A: 50}
	gamesColor      = drawing.Color{R: 239, G: 142, B: 3, A: 255}
	labelColor      = drawing.Color{R: 0, G: 0, B: 0, A: 255}
	gridColor       = drawing.Color{R: 221, G: 221, B: 221, A: 255}
	axisColor       = drawing.Color{R: 0, G: 0, B: 0, A: 255}
	peakDotColor    = drawing.Color{R: 0, G: 0, B: 255, A: 100}
)

// GenerateActivityChart renders the user-activity chart as a PNG.
func GenerateActivityChart(entries []*schemas.BanchoActivity, width int, height int) ([]byte, error) {
	if len(entries) == 0 {
		return nil, errors.New("user activity is empty")
	}

	// Build the series from the entries, ordered from oldest to newest
	x, osu, irc, games := buildSeries(entries)

	// Identify the peak value & its index
	peakValue, peakIndex := activityPeak(osu, irc)

	// Calculate how much headroom to give above the peak
	padding := math.Max(10, math.Floor(peakValue*0.2))
	high := math.Floor((peakValue+padding)/10) * 10
	if high <= 0 {
		high = 10
	}

	graph := chart.Chart{
		Width:  width,
		Height: height,
		// Drop all padding so the plot fills the entire image
		Background: chart.Style{
			FillColor: backgroundColor,
			Padding:   chart.Box{IsSet: true},
		},
		Canvas: chart.Style{FillColor: backgroundColor},
		// Hide the axes, we'll draw them manually later in `drawFrame`
		XAxis: chart.XAxis{Style: chart.Hidden()},
		YAxis: chart.YAxis{
			Style: chart.Hidden(),
			Range: &chart.ContinuousRange{Min: 0, Max: high},
		},
		Series: []chart.Series{
			areaSeries(x, osu, osuFillColor, osuLineColor),
			areaSeries(x, irc, ircFillColor, ircLineColor),
			areaSeries(x, games, gamesColor, gamesColor),
			peakDot(float64(peakIndex), osu[peakIndex]),
			peakLabel(float64(peakIndex), osu[peakIndex], peakValue),
		},
		// Draw the grid, spines & ticks
		Elements: []chart.Renderable{drawFrame},
	}

	buffer := bytes.NewBuffer(nil)
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return nil, fmt.Errorf("failed to render activity chart: %w", err)
	}
	return buffer.Bytes(), nil
}

// buildSeries splits the database activity entries into separated
// series for each "channel", ordered from oldest to newest.
func buildSeries(entries []*schemas.BanchoActivity) (x []float64, osu []float64, irc []float64, games []float64) {
	count := len(entries)
	x = make([]float64, count)
	osu = make([]float64, count)
	irc = make([]float64, count)
	games = make([]float64, count)

	for i := range entries {
		entry := entries[count-1-i]
		x[i] = float64(i)
		osu[i] = float64(entry.OsuCount)
		irc[i] = float64(entry.IrcCount)
		games[i] = float64(entry.MpCount)
	}
	return x, osu, irc, games
}

// areaSeries builds a filled line series for a single activity channel.
func areaSeries(x []float64, y []float64, fill drawing.Color, stroke drawing.Color) chart.ContinuousSeries {
	return chart.ContinuousSeries{
		XValues: x,
		YValues: y,
		Style: chart.Style{
			StrokeColor: stroke,
			StrokeWidth: 1,
			FillColor:   fill,
		},
	}
}

// peakDot marks the peak with a small dot.
func peakDot(x float64, y float64) chart.ContinuousSeries {
	return chart.ContinuousSeries{
		XValues: []float64{x},
		YValues: []float64{y},
		Style: chart.Style{
			DotColor: peakDotColor,
			DotWidth: 3,
		},
	}
}

// peakLabel draws the "Peak: N users" label.
func peakLabel(x float64, y float64, peak float64) chart.AnnotationSeries {
	return chart.AnnotationSeries{
		Annotations: []chart.Value2{{
			XValue: x,
			YValue: y,
			Label:  fmt.Sprintf("Peak: %d users", int(peak)),
			Style: chart.Style{
				FontSize:    8,
				FontColor:   labelColor,
				FillColor:   backgroundColor,
				StrokeColor: backgroundColor,
				Padding:     chart.Box{Top: 1, Left: 2, Right: 2, Bottom: 1, IsSet: true},
			},
		}},
	}
}

// activityPeak returns the highest combined osu! + IRC
// user count & the index at which it occured.
func activityPeak(osu []float64, irc []float64) (value float64, index int) {
	for i := range osu {
		if total := osu[i] + irc[i]; total > value {
			value = total
			index = i
		}
	}
	return value, index
}

// drawFrame renders the background grid, the bottom and
// left lines, and their inward tick marks.
func drawFrame(r chart.Renderer, plot chart.Box, _ chart.Style) {
	left, top := plot.Left, plot.Top
	right, bottom := plot.Right-1, plot.Bottom-1
	width, height := right-left, bottom-top

	// Background grid, 4 columns and 3 rows
	r.SetStrokeColor(gridColor)
	r.SetStrokeWidth(1)
	for i := 1; i < 4; i++ {
		gx := left + width*i/4
		line(r, gx, top, gx, bottom)
	}
	for i := 1; i < 3; i++ {
		gy := top + height*i/3
		line(r, left, gy, right, gy)
	}

	// Bottom and left spines
	r.SetStrokeColor(axisColor)
	r.SetStrokeWidth(1)
	line(r, left, bottom, right, bottom)
	line(r, left, bottom, left, top)

	// Inward ticks
	r.SetStrokeWidth(1)

	// Draw 25 vertical ticks to mark every 4 hours
	for i := 0; i <= 24; i++ {
		tx := left + width*i/24
		length := 3
		if i%6 == 0 {
			// Every 6th tick should be longer
			length = 5
		}
		line(r, tx, bottom, tx, bottom-length)
	}

	// Draw 9 horizontal ticks to mark every 3 hours
	for i := 0; i <= 8; i++ {
		ty := bottom - height*i/8
		length := 3
		if i == 4 {
			// The middle tick should be longer
			length = 5
		}
		line(r, left, ty, left+length, ty)
	}
}

func line(r chart.Renderer, x0 int, y0 int, x1 int, y1 int) {
	r.MoveTo(x0, y0)
	r.LineTo(x1, y1)
	r.Stroke()
}
