package main

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	maxFunctionValue = 8.0
	plotGrid         = 120
	plotSize         = 560
	canvasWidth      = 860
	canvasHeight     = 720
)

func generateResultImages(best individual) (string, string, error) {
	svgPath := outputPath("wykres.svg")
	pngPath := outputPath("wykres.png")

	if err := drawSVG(svgPath, best); err != nil {
		return "", "", fmt.Errorf("nie mozna zapisac wykresu SVG: %w", err)
	}
	if err := convertSVGToPNG(svgPath, pngPath); err != nil {
		return "", "", fmt.Errorf("nie mozna przekonwertowac SVG do PNG: %w", err)
	}

	return svgPath, pngPath, nil
}

func outputPath(fileName string) string {
	taskDir := "zadanie6_algorytm_ewolucyjny"
	if _, err := os.Stat(taskDir); err == nil {
		return filepath.Join(taskDir, fileName)
	}
	return fileName
}

func drawSVG(path string, best individual) error {
	const (
		plotLeft   = 90
		plotTop    = 80
		legendLeft = plotLeft + plotSize + 50
	)

	var svg strings.Builder
	svg.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	svg.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+"\n", canvasWidth, canvasHeight, canvasWidth, canvasHeight))
	svg.WriteString(`<rect width="100%" height="100%" fill="#f8fafc"/>` + "\n")
	svg.WriteString(`<text x="90" y="36" font-family="Verdana, sans-serif" font-size="22" font-weight="700" fill="#111827">Optimum funkcji znalezione algorytmem ewolucyjnym</text>` + "\n")
	svg.WriteString(`<text x="90" y="60" font-family="Verdana, sans-serif" font-size="13" fill="#475569">Kolor pokazuje wartosc f(x,y), punkt oznacza najlepsze znalezione rozwiazanie.</text>` + "\n")

	cellSize := float64(plotSize) / float64(plotGrid)
	svg.WriteString(`<g shape-rendering="crispEdges">` + "\n")
	for row := 0; row < plotGrid; row++ {
		for col := 0; col < plotGrid; col++ {
			x := (float64(col) + 0.5) / float64(plotGrid) * period
			y := (float64(plotGrid-row-1) + 0.5) / float64(plotGrid) * period
			value := objective(x, y)
			svg.WriteString(fmt.Sprintf(
				`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`+"\n",
				float64(plotLeft)+float64(col)*cellSize,
				float64(plotTop)+float64(row)*cellSize,
				cellSize+0.02,
				cellSize+0.02,
				heatColor(value),
			))
		}
	}
	svg.WriteString(`</g>` + "\n")

	svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#111827" stroke-width="2"/>`+"\n", plotLeft, plotTop, plotSize, plotSize))
	drawAxes(&svg, plotLeft, plotTop)
	drawLegend(&svg, legendLeft, plotTop)
	drawBestPoint(&svg, plotLeft, plotTop, best)

	svg.WriteString(`</svg>` + "\n")

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(svg.String()), 0644)
}

func drawAxes(svg *strings.Builder, plotLeft, plotTop int) {
	axisStyle := `font-family="Verdana, sans-serif" font-size="12" fill="#334155"`
	ticks := []struct {
		label string
		value float64
	}{
		{"0", 0},
		{"pi", math.Pi},
		{"2pi", period},
	}

	for _, tick := range ticks {
		x := float64(plotLeft) + tick.value/period*float64(plotSize)
		y := float64(plotTop) + float64(plotSize) - tick.value/period*float64(plotSize)

		svg.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%d" x2="%.2f" y2="%d" stroke="#111827" stroke-width="1"/>`+"\n", x, plotTop+plotSize, x, plotTop+plotSize+6))
		svg.WriteString(fmt.Sprintf(`<text x="%.2f" y="%d" text-anchor="middle" %s>%s</text>`+"\n", x, plotTop+plotSize+24, axisStyle, tick.label))

		svg.WriteString(fmt.Sprintf(`<line x1="%d" y1="%.2f" x2="%d" y2="%.2f" stroke="#111827" stroke-width="1"/>`+"\n", plotLeft-6, y, plotLeft, y))
		svg.WriteString(fmt.Sprintf(`<text x="%d" y="%.2f" text-anchor="end" dominant-baseline="middle" %s>%s</text>`+"\n", plotLeft-12, y, axisStyle, tick.label))
	}

	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-family="Verdana, sans-serif" font-size="15" font-weight="700" fill="#111827">x</text>`+"\n", plotLeft+plotSize/2, plotTop+plotSize+52))
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" transform="rotate(-90 %d %d)" font-family="Verdana, sans-serif" font-size="15" font-weight="700" fill="#111827">y</text>`+"\n", plotLeft-54, plotTop+plotSize/2, plotLeft-54, plotTop+plotSize/2))
}

func drawLegend(svg *strings.Builder, legendLeft, plotTop int) {
	const (
		legendWidth  = 28
		legendHeight = 240
	)

	for i := 0; i < legendHeight; i++ {
		t := 1 - float64(i)/float64(legendHeight-1)
		value := t * maxFunctionValue
		svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="1" fill="%s"/>`+"\n",
			legendLeft, plotTop+i, legendWidth, heatColor(value)))
	}

	svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#111827" stroke-width="1"/>`+"\n", legendLeft, plotTop, legendWidth, legendHeight))
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="Verdana, sans-serif" font-size="13" font-weight="700" fill="#111827">f(x,y)</text>`+"\n", legendLeft, plotTop-12))
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="Verdana, sans-serif" font-size="12" fill="#334155">8</text>`+"\n", legendLeft+legendWidth+8, plotTop+4))
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="Verdana, sans-serif" font-size="12" fill="#334155">4</text>`+"\n", legendLeft+legendWidth+8, plotTop+legendHeight/2+4))
	svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="Verdana, sans-serif" font-size="12" fill="#334155">0</text>`+"\n", legendLeft+legendWidth+8, plotTop+legendHeight+4))
}

func drawBestPoint(svg *strings.Builder, plotLeft, plotTop int, best individual) {
	px := float64(plotLeft) + best.x/period*float64(plotSize)
	py := float64(plotTop) + float64(plotSize) - best.y/period*float64(plotSize)

	svg.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#ffffff" stroke-width="7" stroke-linecap="round"/>`+"\n", px-14, py, px+14, py))
	svg.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#ffffff" stroke-width="7" stroke-linecap="round"/>`+"\n", px, py-14, px, py+14))
	svg.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#111827" stroke-width="3" stroke-linecap="round"/>`+"\n", px-14, py, px+14, py))
	svg.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#111827" stroke-width="3" stroke-linecap="round"/>`+"\n", px, py-14, px, py+14))
	svg.WriteString(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="7" fill="#f8fafc" stroke="#111827" stroke-width="3"/>`+"\n", px, py))

	labelX := px + 18
	labelY := py - 18
	if labelX > float64(plotLeft+plotSize-210) {
		labelX = px - 220
	}
	if labelY < float64(plotTop+28) {
		labelY = py + 28
	}

	svg.WriteString(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="205" height="72" rx="10" fill="#ffffff" fill-opacity="0.92" stroke="#111827" stroke-width="1"/>`+"\n", labelX, labelY))
	svg.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-family="Verdana, sans-serif" font-size="13" font-weight="700" fill="#111827">Najlepszy punkt</text>`+"\n", labelX+12, labelY+22))
	svg.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-family="Verdana, sans-serif" font-size="12" fill="#334155">x = %.5f, y = %.5f</text>`+"\n", labelX+12, labelY+43, best.x, best.y))
	svg.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-family="Verdana, sans-serif" font-size="12" fill="#334155">f(x,y) = %.5f</text>`+"\n", labelX+12, labelY+61, best.fitness))
}

func convertSVGToPNG(svgPath, pngPath string) error {
	svgAbs, err := filepath.Abs(svgPath)
	if err != nil {
		return err
	}
	pngAbs, err := filepath.Abs(pngPath)
	if err != nil {
		return err
	}
	browserPath, err := findChromiumBrowser()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(pngAbs), 0755); err != nil {
		return err
	}
	if err := os.Remove(pngAbs); err != nil && !os.IsNotExist(err) {
		return err
	}

	profileDir, err := os.MkdirTemp("", "svg-to-png-profile-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(profileDir)

	baseArgs := []string{
		"--disable-gpu",
		"--disable-software-rasterizer",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		"--hide-scrollbars",
		"--allow-file-access-from-files",
		"--user-data-dir=" + profileDir,
		"--screenshot=" + pngAbs,
		fmt.Sprintf("--window-size=%d,%d", canvasWidth, canvasHeight),
		fileURL(svgAbs),
	}

	var lastErr error
	for _, headlessFlag := range []string{"--headless=new", "--headless"} {
		args := append([]string{headlessFlag}, baseArgs...)
		cmd := exec.Command(browserPath, args...)
		output, err := cmd.CombinedOutput()
		if err == nil && fileExists(pngAbs) {
			return nil
		}
		if err == nil {
			err = fmt.Errorf("plik PNG nie zostal utworzony")
		}
		lastErr = fmt.Errorf("%s: %w; %s", browserPath, err, strings.TrimSpace(string(output)))
	}

	return lastErr
}

func findChromiumBrowser() (string, error) {
	for _, name := range []string{"chrome", "google-chrome", "chromium", "chromium-browser", "msedge"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	candidates := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
	}
	for _, candidate := range candidates {
		if fileExists(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("nie znaleziono Edge/Chrome do konwersji SVG na PNG")
}

func fileURL(path string) string {
	slashPath := filepath.ToSlash(path)
	if !strings.HasPrefix(slashPath, "/") {
		slashPath = "/" + slashPath
	}
	return (&url.URL{Scheme: "file", Path: slashPath}).String()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func heatColor(value float64) string {
	stops := []struct {
		at      float64
		r, g, b int
	}{
		{0.00, 15, 23, 42},
		{0.35, 14, 116, 144},
		{0.65, 234, 179, 8},
		{1.00, 220, 38, 38},
	}

	t := value / maxFunctionValue
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	for i := 0; i < len(stops)-1; i++ {
		left := stops[i]
		right := stops[i+1]
		if t >= left.at && t <= right.at {
			localT := (t - left.at) / (right.at - left.at)
			r := interpolate(left.r, right.r, localT)
			g := interpolate(left.g, right.g, localT)
			b := interpolate(left.b, right.b, localT)
			return fmt.Sprintf("#%02x%02x%02x", r, g, b)
		}
	}

	last := stops[len(stops)-1]
	return fmt.Sprintf("#%02x%02x%02x", last.r, last.g, last.b)
}

func interpolate(left, right int, t float64) int {
	return int(math.Round(float64(left) + (float64(right)-float64(left))*t))
}
