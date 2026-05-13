package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type task struct {
	Name       string
	Directory  string
	ReadmePath string
	CodePath   string
	CodeLang   string
}

var sourceCandidates = []struct{ name, lang string }{
	{"main.go", "go"},
	{"main.py", "python"},
	{"main.ts", "typescript"},
	{"main.js", "javascript"},
}

var extToLang = map[string]string{
	".go": "go",
	".py": "python",
	".ts": "typescript",
	".js": "javascript",
}

func findSourceFile(dir string) (path, lang string, found bool) {
	for _, c := range sourceCandidates {
		p := filepath.Join(dir, c.name)
		if fileExists(p) {
			return p, c.lang, true
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", "", false
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if l, ok := extToLang[strings.ToLower(filepath.Ext(e.Name()))]; ok {
			return filepath.Join(dir, e.Name()), l, true
		}
	}
	return "", "", false
}

type fontSet struct {
	DocumentPath string
}

type markdownImage struct {
	AltText string
	Path    string
}

const (
	fontFamilyRegular    = "Regular"
	fontFamilyMono       = "Mono"
	resultSectionHeading = "## Screen dzialania programu / wynik skryptu"
	resultSectionAltText = "Wynik"
)

func main() {
	outputPath := flag.String("out", "sprawozdanie_algorytmy.pdf", "Sciezka pliku PDF wyjsciowego")
	flag.Parse()

	args := flag.Args()
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "Uzycie: generate_report -out plik.pdf <tytul> <imie_nazwisko> <numer_indeksu>")
		fmt.Fprintln(os.Stderr, "Przyklad: generate_report -out raport.pdf \"Sprawozdanie\" \"Jan Kowalski\" 123456")
		os.Exit(1)
	}
	reportTitle := args[0]
	reportAuthor := args[1]
	reportIndex := args[2]

	tasks, err := discoverTasks(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Blad:", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(os.Stderr, "Blad: nie znaleziono zadnych katalogow z README.md i main.go")
		os.Exit(1)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.SetTitle(reportTitle, false)
	pdf.SetAuthor(reportAuthor, false)
	pdf.SetSubject(reportTitle, false)
	pdf.SetCreator("generate_report.go", false)

	fonts, err := resolveFonts()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Blad:", err)
		os.Exit(1)
	}

	if err := buildPDF(pdf, tasks, fonts, *outputPath, reportTitle, reportAuthor, reportIndex); err != nil {
		fmt.Fprintln(os.Stderr, "Blad:", err)
		os.Exit(1)
	}

	fmt.Printf("Utworzono raport: %s\n", *outputPath)
}

func discoverTasks(root string) ([]task, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	tasks := make([]task, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		readmePath := filepath.Join(root, name, "README.md")
		if !fileExists(readmePath) {
			continue
		}

		codePath, codeLang, codeFound := findSourceFile(filepath.Join(root, name))
		if !codeFound {
			continue
		}

		tasks = append(tasks, task{
			Name:       prettifyTaskName(name),
			Directory:  name,
			ReadmePath: readmePath,
			CodePath:   codePath,
			CodeLang:   codeLang,
		})
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Directory < tasks[j].Directory
	})

	return tasks, nil
}

func prettifyTaskName(directory string) string {
	name := stripTaskPrefix(directory)
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func stripTaskPrefix(directory string) string {
	const prefix = "zadanie"
	if !strings.HasPrefix(directory, prefix) {
		return directory
	}
	idx := len(prefix)
	for idx < len(directory) {
		c := directory[idx]
		if c < '0' || c > '9' {
			break
		}
		idx++
	}
	if idx < len(directory) && directory[idx] == '_' {
		return directory[idx+1:]
	}
	return directory
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func resolveFonts() (fontSet, error) {
	localDocument := filepath.Join("fonts", "DejaVuSans.ttf")
	documentCandidates := []string{
		`C:\Windows\Fonts\l_10646.ttf`,
		localDocument,
		`C:\Windows\Fonts\segoeui.ttf`,
		`C:\Windows\Fonts\calibri.ttf`,
		`C:\Windows\Fonts\arial.ttf`,
		`C:\Windows\Fonts\consola.ttf`,
		`C:\Windows\Fonts\CascadiaMono.ttf`,
		`C:\Windows\Fonts\lucon.ttf`,
	}

	document, err := firstExistingPath(documentCandidates)
	if err != nil {
		return fontSet{}, fmt.Errorf("nie udalo sie znalezc fontu dokumentu: %w", err)
	}

	return fontSet{
		DocumentPath: document,
	}, nil
}

func firstExistingPath(candidates []string) (string, error) {
	for _, candidate := range candidates {
		if fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", errors.New("brak odpowiedniego pliku fontu na tym komputerze")
}

func buildPDF(pdf *gofpdf.Fpdf, tasks []task, fonts fontSet, outputPath, reportTitle, reportAuthor, reportIndex string) error {
	pdf.AddUTF8Font(fontFamilyRegular, "", fonts.DocumentPath)
	pdf.AddUTF8Font(fontFamilyRegular, "B", fonts.DocumentPath)
	pdf.AddUTF8Font(fontFamilyRegular, "I", fonts.DocumentPath)
	pdf.AddUTF8Font(fontFamilyMono, "", fonts.DocumentPath)
	if pdf.Err() {
		return pdf.Error()
	}

	addFooter := func() {
		pdf.SetY(-12)
		pdf.SetFont(fontFamilyRegular, "", 8)
		pdf.SetTextColor(110, 110, 110)
		pdf.CellFormat(0, 5, fmt.Sprintf("Strona %d", pdf.PageNo()), "", 0, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}
	pdf.SetFooterFunc(addFooter)

	// przygotuj linki do spisu tresci
	links := make([]int, len(tasks))
	for i := range links {
		links[i] = pdf.AddLink()
	}

	addTitleAndOverviewPage(pdf, tasks, links, reportTitle, reportAuthor, reportIndex)

	for i, task := range tasks {
		content, err := prepareTaskContent(task)
		if err != nil {
			return err
		}

		pdf.Ln(4)
		// ustaw cel linku na poczatek tej sekcji
		pdf.SetLink(links[i], pdf.GetY(), pdf.PageNo())
		if err := addTaskSection(pdf, task, content); err != nil {
			return err
		}
	}

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return err
	}

	return nil
}

func addTitleAndOverviewPage(pdf *gofpdf.Fpdf, tasks []task, links []int, reportTitle, reportAuthor, reportIndex string) {
	pdf.AddPage()
	pdf.SetFont(fontFamilyRegular, "", 13)
	pdf.MultiCell(0, 7, reportAuthor, "", "C", false)
	pdf.SetFont(fontFamilyRegular, "", 11)
	pdf.MultiCell(0, 7, "Nr indeksu: "+reportIndex, "", "C", false)
	pdf.Ln(4)
	pdf.SetFont(fontFamilyRegular, "B", 23)
	pdf.MultiCell(0, 12, reportTitle, "", "C", false)
	pdf.Ln(4)
	pdf.SetFont(fontFamilyRegular, "B", 17)
	pdf.CellFormat(0, 10, "Spis tresci", "", 1, "L", false, 0, "")
	pdf.Ln(2)
	pdf.SetFont(fontFamilyRegular, "", 11)
	for i, task := range tasks {
		// wpis spisu jako link do strony zadania
		pdf.CellFormat(0, 7, fmt.Sprintf("%d. %s", i+1, task.Name), "", 1, "L", false, links[i], "")
	}
	pdf.Ln(4)
}

func addTaskSection(pdf *gofpdf.Fpdf, task task, readme string) error {
	// Jeśli README zaczyna się od nagłówka H1 (# ), nie rysujemy
	// dodatkowego tytułu, żeby nie dublować nazwy zadania.
	first := firstNonEmptyLine(readme)
	if !strings.HasPrefix(strings.TrimSpace(first), "# ") {
		pdf.SetFont(fontFamilyRegular, "B", 17)
		pdf.CellFormat(0, 10, task.Name, "", 1, "L", false, 0, "")
	}
	if err := renderMarkdown(pdf, readme, filepath.Dir(task.ReadmePath)); err != nil {
		return fmt.Errorf("nie mozna wyrenderowac sekcji %q: %w", task.Directory, err)
	}
	pdf.Ln(3)
	return nil
}

func prepareTaskContent(task task) (string, error) {
	readmeBytes, err := os.ReadFile(task.ReadmePath)
	if err != nil {
		return "", fmt.Errorf("nie mozna odczytac %s: %w", task.ReadmePath, err)
	}

	readme := string(readmeBytes)
	imageMarkdown, err := buildGeneratedImageSection(filepath.Dir(task.ReadmePath))
	if err != nil {
		return "", err
	}
	if imageMarkdown != "" && !hasGeneratedResultSection(readme) {
		readme = insertSectionBeforeSourceCode(readme, imageMarkdown)
	}
	if hasSourceCodeSection(readme) {
		return readme, nil
	}

	codeBytes, err := os.ReadFile(task.CodePath)
	if err != nil {
		return "", fmt.Errorf("nie mozna odczytac %s: %w", task.CodePath, err)
	}

	return appendGeneratedSourceCodeSection(readme, string(codeBytes), task.CodeLang), nil
}

func buildGeneratedImageSection(taskDir string) (string, error) {
	imageFile, err := findFirstTaskImage(taskDir)
	if err != nil {
		return "", err
	}
	if imageFile == "" {
		return "", nil
	}

	return resultSectionHeading + "\n\n" +
		fmt.Sprintf("![%s](%s)\n", resultSectionAltText, filepath.ToSlash(imageFile)), nil
}

func findFirstTaskImage(taskDir string) (string, error) {
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return "", fmt.Errorf("nie mozna odczytac katalogu %s: %w", taskDir, err)
	}

	imageFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if isSupportedImageExtension(filepath.Ext(name)) {
			imageFiles = append(imageFiles, name)
		}
	}

	sort.Strings(imageFiles)
	if len(imageFiles) == 0 {
		return "", nil
	}

	return imageFiles[0], nil
}

func insertSectionBeforeSourceCode(readme, section string) string {
	lines := splitLines(readme)
	for i, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if !strings.HasPrefix(line, "#") {
			continue
		}

		heading := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if normalizeHeadingKey(heading) != "kod zrodlowy" {
			continue
		}

		before := strings.Join(lines[:i], "\n")
		after := strings.Join(lines[i:], "\n")

		if strings.TrimSpace(before) == "" {
			return section + "\n" + after
		}

		return strings.TrimRight(before, "\n") + "\n\n" + section + "\n" + after
	}

	readme = strings.TrimRight(readme, "\r\n")
	if readme == "" {
		return section + "\n"
	}

	return readme + "\n\n" + section + "\n"
}

func firstNonEmptyLine(s string) string {
	for _, line := range splitLines(s) {
		if strings.TrimSpace(line) != "" {
			return line
		}
	}
	return ""
}

func renderMarkdown(pdf *gofpdf.Fpdf, content, assetBaseDir string) error {
	lines := splitLines(content)
	inCode := false
	codeLines := make([]string, 0)

	flushCode := func() {
		if len(codeLines) == 0 {
			return
		}
		addCodeLines(pdf, codeLines, 7.5)
		codeLines = codeLines[:0]
	}

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				flushCode()
				inCode = false
				pdf.Ln(2)
			} else {
				inCode = true
				pdf.Ln(1)
			}
			continue
		}

		if inCode {
			codeLines = append(codeLines, expandTabs(line))
			continue
		}

		if image, ok := parseMarkdownImage(trimmed); ok {
			if err := addMarkdownImage(pdf, assetBaseDir, image); err != nil {
				return err
			}
			continue
		}

		if trimmed == "" {
			pdf.Ln(2)
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "# "):
			setBodyFont(pdf, "B", 15)
			pdf.MultiCell(0, 8, normalizeInline(trimmed[2:]), "", "L", false)
			pdf.Ln(1)
		case strings.HasPrefix(trimmed, "## "):
			setBodyFont(pdf, "B", 13)
			pdf.MultiCell(0, 7, normalizeInline(trimmed[3:]), "", "L", false)
			pdf.Ln(1)
		case strings.HasPrefix(trimmed, "### "):
			setBodyFont(pdf, "B", 11)
			pdf.MultiCell(0, 7, normalizeInline(trimmed[4:]), "", "L", false)
			pdf.Ln(1)
		case strings.HasPrefix(trimmed, "- "):
			setBodyFont(pdf, "", 10)
			pdf.MultiCell(0, 6, "- "+normalizeInline(trimmed[2:]), "", "L", false)
		case isNumberedListItem(trimmed):
			setBodyFont(pdf, "", 10)
			pdf.MultiCell(0, 6, normalizeInline(trimmed), "", "L", false)
		default:
			setBodyFont(pdf, "", 10)
			pdf.MultiCell(0, 6, normalizeInline(trimmed), "", "L", false)
		}
	}

	if inCode {
		flushCode()
	}

	return nil
}

func addCodeBlock(pdf *gofpdf.Fpdf, code string) {
	lines := splitLines(code)
	if len(lines) == 0 {
		pdf.SetFont(fontFamilyRegular, "I", 9)
		pdf.MultiCell(0, 6, "Brak kodu do wyswietlenia.", "", "L", false)
		return
	}

	addCodeLines(pdf, lines, 7.5)
}

func addCodeLines(pdf *gofpdf.Fpdf, lines []string, fontSize float64) {
	pdf.SetFillColor(245, 245, 245)
	pdf.SetTextColor(30, 30, 30)
	for _, rawLine := range lines {
		line := expandTabs(strings.TrimRight(rawLine, "\r"))
		setCodeFont(pdf, fontSize)
		if line == "" {
			pdf.MultiCell(0, 4.5, " ", "", "L", false)
			continue
		}
		pdf.MultiCell(0, 4.5, line, "", "L", false)
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(2)
}

func parseMarkdownImage(line string) (markdownImage, bool) {
	if !strings.HasPrefix(line, "![") || !strings.HasSuffix(line, ")") {
		return markdownImage{}, false
	}

	separatorIndex := strings.Index(line, "](")
	if separatorIndex <= 1 {
		return markdownImage{}, false
	}

	path := strings.TrimSpace(line[separatorIndex+2 : len(line)-1])
	if path == "" {
		return markdownImage{}, false
	}
	if strings.HasPrefix(path, "<") && strings.HasSuffix(path, ">") {
		path = strings.TrimSpace(path[1 : len(path)-1])
	}
	if path == "" {
		return markdownImage{}, false
	}

	return markdownImage{
		AltText: line[2:separatorIndex],
		Path:    path,
	}, true
}

func addMarkdownImage(pdf *gofpdf.Fpdf, assetBaseDir string, image markdownImage) error {
	imagePath := resolveMarkdownAssetPath(assetBaseDir, image.Path)
	if !fileExists(imagePath) {
		return fmt.Errorf("nie znaleziono pliku obrazka %q", image.Path)
	}

	imageOptions, err := imageOptionsFromPath(imagePath)
	if err != nil {
		return err
	}

	info := pdf.RegisterImageOptions(imagePath, imageOptions)
	if pdf.Err() {
		return pdf.Error()
	}

	pageWidth, _ := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	maxWidth := pageWidth - left - right
	width := info.Width()
	if width <= 0 || width > maxWidth {
		width = maxWidth
	}

	pdf.Ln(1)
	pdf.ImageOptions(imagePath, -1, pdf.GetY(), width, 0, true, imageOptions, 0, "")
	if pdf.Err() {
		return pdf.Error()
	}
	pdf.Ln(2)
	return nil
}

func resolveMarkdownAssetPath(assetBaseDir, assetPath string) string {
	normalizedPath := filepath.FromSlash(assetPath)
	if filepath.IsAbs(normalizedPath) {
		return filepath.Clean(normalizedPath)
	}
	return filepath.Clean(filepath.Join(assetBaseDir, normalizedPath))
}

func imageOptionsFromPath(path string) (gofpdf.ImageOptions, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, nil
	case ".jpg", ".jpeg":
		return gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, nil
	case ".gif":
		return gofpdf.ImageOptions{ImageType: "GIF", ReadDpi: true}, nil
	default:
		return gofpdf.ImageOptions{}, fmt.Errorf("nieobslugiwany format obrazka %q", path)
	}
}

func normalizeInline(text string) string {
	text = strings.ReplaceAll(text, "`", "")
	text = strings.ReplaceAll(text, "\t", "    ")
	return text
}

func hasGeneratedResultSection(readme string) bool {
	for _, rawLine := range splitLines(readme) {
		line := strings.TrimSpace(rawLine)
		if !strings.HasPrefix(line, "#") {
			continue
		}

		heading := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if normalizeHeadingKey(heading) == "screen dzialania programu / wynik skryptu" {
			return true
		}
	}

	return false
}

func hasSourceCodeSection(readme string) bool {
	for _, rawLine := range splitLines(readme) {
		line := strings.TrimSpace(rawLine)
		if !strings.HasPrefix(line, "#") {
			continue
		}
		heading := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if normalizeHeadingKey(heading) == "kod zrodlowy" {
			return true
		}
	}

	return false
}

func isSupportedImageExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif":
		return true
	default:
		return false
	}
}

func appendGeneratedSourceCodeSection(readme, code, lang string) string {
	readme = strings.TrimRight(readme, "\r\n")
	code = strings.TrimRight(code, "\r\n")

	section := "## Kod źródłowy\n\n```" + lang + "\n" + code + "\n```\n"
	if readme == "" {
		return section
	}

	return readme + "\n\n" + section
}

func normalizeHeadingKey(text string) string {
	replacer := strings.NewReplacer(
		"ą", "a",
		"ć", "c",
		"ę", "e",
		"ł", "l",
		"ń", "n",
		"ó", "o",
		"ś", "s",
		"ź", "z",
		"ż", "z",
		"Ą", "a",
		"Ć", "c",
		"Ę", "e",
		"Ł", "l",
		"Ń", "n",
		"Ó", "o",
		"Ś", "s",
		"Ź", "z",
		"Ż", "z",
	)

	normalized := replacer.Replace(text)
	normalized = strings.ToLower(normalized)
	return strings.Join(strings.Fields(normalized), " ")
}

func setBodyFont(pdf *gofpdf.Fpdf, style string, size float64) {
	pdf.SetFont(fontFamilyRegular, style, size)
}

func setCodeFont(pdf *gofpdf.Fpdf, size float64) {
	pdf.SetFont(fontFamilyMono, "", size)
}

func expandTabs(text string) string {
	return strings.ReplaceAll(text, "\t", "    ")
}

func splitLines(text string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{text}
	}
	if len(lines) == 0 {
		return []string{}
	}
	return lines
}

func isNumberedListItem(line string) bool {
	if len(line) < 3 {
		return false
	}
	separatorIndex := strings.Index(line, ". ")
	if separatorIndex <= 0 {
		return false
	}
	for _, r := range line[:separatorIndex] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
