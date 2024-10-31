package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Función para remover caracteres especiales
func onlyCharacters(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

// Función para remover las mayúsculas
func toLowerCase(s string) string {
	return strings.ToLower(s)
}

// Leer archivos y almacenar sus contenidos como un String "limpio"
func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		os.Exit(1)
	}
	processedContent := toLowerCase(onlyCharacters(string(content)))
	return processedContent
}

type SuffixArray struct {
	text     string
	suffixes []int
}

// Build a suffix array for a given text
func buildSuffixArray(text string) SuffixArray {
	n := len(text)
	suffixes := make([]int, n)

	// Generate all suffix indices
	for i := 0; i < n; i++ {
		suffixes[i] = i
	}

	// Sort suffixes alphabetically based on the suffix substring
	sort.Slice(suffixes, func(i, j int) bool {
		return text[suffixes[i]:] < text[suffixes[j]:]
	})

	return SuffixArray{
		text:     text,
		suffixes: suffixes,
	}
}

func longestCommonSubstring(sa1, sa2 SuffixArray) (int, string) {
	maxLen := 0
	longestSubstr := ""

	for _, idx1 := range sa1.suffixes {
		for _, idx2 := range sa2.suffixes {
			length := 0
			for idx1+length < len(sa1.text) && idx2+length < len(sa2.text) && sa1.text[idx1+length] == sa2.text[idx2+length] {
				length++
			}
			if length > maxLen {
				maxLen = length
				longestSubstr = sa1.text[idx1 : idx1+length]
			}
		}
	}
	return maxLen, longestSubstr
}

// Levenshtein distance function
func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Create the matrix
	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
	}

	// Initialize the first row and column
	for i := 0; i <= la; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = min(matrix[i-1][j]+1, min(matrix[i][j-1]+1, matrix[i-1][j-1]+cost))
		}
	}

	return matrix[la][lb]
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Calcular porcentaje de similitud
func calculateSimilarity(content1, content2 string, commonLength int) float64 {
	shorterLength := len(content1)
	if len(content2) < len(content1) {
		shorterLength = len(content2)
	}

	percentage := (float64(commonLength) / float64(shorterLength)) * 100
	return percentage
}

// Comparación & guardado de archivos
func compareFiles(original string, others []string, comparisons *[][]string) {
	content1 := readFile(original)

	for _, other := range others {
		content2 := readFile(other)

		// Build suffix arrays for both contents
		sa1 := buildSuffixArray(content1)
		sa2 := buildSuffixArray(content2)

		// Find the longest common substring
		commonLength, commonSubstring := longestCommonSubstring(sa1, sa2)

		if commonLength > 0 {
			percentage := calculateSimilarity(content1, content2, commonLength)
			if percentage >= 20 {
				highlighted1, highlighted2 := highlightSimilarities(content1, content2, commonSubstring)
				*comparisons = append(*comparisons, []string{
					original,
					other,
					strconv.FormatFloat(percentage, 'f', 2, 64),
					highlighted1,
					highlighted2,
				})
			}
		}
	}
}

// Resaltar coincidencias en amarillo en el contenido de cada archivo
func highlightSimilarities(content1, content2, commonSubstring string) (string, string) {
	highlighted1 := strings.Replace(content1, commonSubstring, `<span class="highlight">`+commonSubstring+`</span>`, -1)
	highlighted2 := strings.Replace(content2, commonSubstring, `<span class="highlight">`+commonSubstring+`</span>`, -1)
	return highlighted1, highlighted2
}

func insertLineBreaks(text string) string {
	var result strings.Builder
	for i, r := range text {
		if i > 0 && i%100 == 0 {
			result.WriteString("<br>")
		}
		result.WriteRune(r)
	}
	return result.String()
}

func generateHTML(comparisons [][]string) string {
	htmlContent := `<!DOCTYPE html>
	<html lang="es">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Similitudes entre Archivos</title>
		<style>
			.highlight { background-color: yellow; font-weight: bold; }
		</style>
	</head>
	<body>
		<h1>Comparaciones de Archivos</h1>
	`

	for _, comp := range comparisons {
		htmlContent += `
		<div>
			<h2>Archivo base: ` + comp[0] + `, Comparado con: ` + comp[1] + `</h2>
			<p><strong>Similitud:</strong> ` + comp[2] + `%</p>
			<h3>Contenido del Archivo Base:</h3>
			<p>` + insertLineBreaks(comp[3]) + `</p>
			<h3>Contenido del Archivo Comparado:</h3>
			<p>` + insertLineBreaks(comp[4]) + `</p>
		</div>
		<hr>
		`
	}

	htmlContent += `
	</body>
	</html>
	`
	return htmlContent
}

func handler(w http.ResponseWriter, r *http.Request) {
	comparisons := [][]string{}
	compareFiles("data/orig_taska.txt", []string{"data/g0pE_taska.txt"}, &comparisons)
	compareFiles("data/orig_taske.txt", []string{"data/g0pE_taske.txt"}, &comparisons)
	compareFiles("data/orig_taskb.txt", []string{"data/g0pE_taskb.txt"}, &comparisons)
	compareFiles("data/orig_taska.txt", []string{"data/g4pC_taska.txt"}, &comparisons)
	compareFiles("data/orig_taska.txt", []string{"data/g3pC_taska.txt"}, &comparisons)
	compareFiles("data/orig_taske.txt", []string{"data/g2pB_taske.txt"}, &comparisons)
	compareFiles("data/orig_taskd.txt", []string{"data/g3pA_taskd.txt"}, &comparisons)
	compareFiles("data/orig_taska.txt", []string{"data/g2pC_taska.txt"}, &comparisons)
	compareFiles("data/orig_taskd.txt", []string{"data/g4pC_taskd.txt"}, &comparisons)
	compareFiles("data/orig_taske.txt", []string{"data/g4pB_taske.txt"}, &comparisons)

	sort.Slice(comparisons, func(i, j int) bool {
		similarityI, _ := strconv.ParseFloat(comparisons[i][2], 64)
		similarityJ, _ := strconv.ParseFloat(comparisons[j][2], 64)
		return similarityI > similarityJ
	})
	if len(comparisons) > 10 {
		comparisons = comparisons[:10]
	}

	htmlContent := generateHTML(comparisons)
	fmt.Fprint(w, htmlContent)
}

func main() {
	Originals := []string{
		"data/orig_taska.txt",
		"data/orig_taskb.txt",
		"data/orig_taskc.txt",
		"data/orig_taskd.txt",
		"data/orig_taske.txt",
	}

	TaskA := []string{
		"data/g0pA_taska.txt",
		"data/g0pB_taska.txt",
		"data/g0pC_taska.txt",
		"data/g0pD_taska.txt",
		"data/g0pE_taska.txt",
		"data/g1pA_taska.txt",
		"data/g1pB_taska.txt",
		"data/g1pD_taska.txt",
		"data/g2pA_taska.txt",
		"data/g2pB_taska.txt",
		"data/g2pC_taska.txt",
		"data/g2pE_taska.txt",
		"data/g3pA_taska.txt",
		"data/g3pB_taska.txt",
		"data/g3pC_taska.txt",
		"data/g4pB_taska.txt",
		"data/g4pC_taska.txt",
		"data/g4pD_taska.txt",
		"data/g4pE_taska.txt",
	}

	TaskB := []string{
		"data/g0pA_taskb.txt",
		"data/g0pB_taskb.txt",
		"data/g0pC_taskb.txt",
		"data/g0pD_taskb.txt",
		"data/g0pE_taskb.txt",
		"data/g1pA_taskb.txt",
		"data/g1pB_taskb.txt",
		"data/g1pD_taskb.txt",
		"data/g2pA_taskb.txt",
		"data/g2pB_taskb.txt",
		"data/g2pC_taskb.txt",
		"data/g2pE_taskb.txt",
		"data/g3pA_taskb.txt",
		"data/g3pB_taskb.txt",
		"data/g3pC_taskb.txt",
		"data/g4pB_taskb.txt",
		"data/g4pC_taskb.txt",
		"data/g4pD_taskb.txt",
		"data/g4pE_taskb.txt",
	}

	TaskC := []string{
		"data/g0pA_taskc.txt",
		"data/g0pB_taskc.txt",
		"data/g0pC_taskc.txt",
		"data/g0pD_taskc.txt",
		"data/g0pE_taskc.txt",
		"data/g1pA_taskc.txt",
		"data/g1pB_taskc.txt",
		"data/g1pD_taskc.txt",
		"data/g2pA_taskc.txt",
		"data/g2pB_taskc.txt",
		"data/g2pC_taskc.txt",
		"data/g2pE_taskc.txt",
		"data/g3pA_taskc.txt",
		"data/g3pB_taskc.txt",
		"data/g3pC_taskc.txt",
		"data/g4pB_taskc.txt",
		"data/g4pC_taskc.txt",
		"data/g4pD_taskc.txt",
		"data/g4pE_taskc.txt",
	}

	TaskD := []string{
		"data/g0pA_taskd.txt",
		"data/g0pB_taskd.txt",
		"data/g0pC_taskd.txt",
		"data/g0pD_taskd.txt",
		"data/g0pE_taskd.txt",
		"data/g1pA_taskd.txt",
		"data/g1pB_taskd.txt",
		"data/g1pD_taskd.txt",
		"data/g2pA_taskd.txt",
		"data/g2pB_taskd.txt",
		"data/g2pC_taskd.txt",
		"data/g2pE_taskd.txt",
		"data/g3pA_taskd.txt",
		"data/g3pB_taskd.txt",
		"data/g3pC_taskd.txt",
		"data/g4pB_taskd.txt",
		"data/g4pC_taskd.txt",
		"data/g4pD_taskd.txt",
		"data/g4pE_taskd.txt",
	}

	TaskE := []string{
		"data/g0pA_taske.txt",
		"data/g0pB_taske.txt",
		"data/g0pC_taske.txt",
		"data/g0pD_taske.txt",
		"data/g0pE_taske.txt",
		"data/g1pA_taske.txt",
		"data/g1pB_taske.txt",
		"data/g1pD_taske.txt",
		"data/g2pA_taske.txt",
		"data/g2pB_taske.txt",
		"data/g2pC_taske.txt",
		"data/g2pE_taske.txt",
		"data/g3pA_taske.txt",
		"data/g3pB_taske.txt",
		"data/g3pC_taske.txt",
		"data/g4pB_taske.txt",
		"data/g4pC_taske.txt",
		"data/g4pD_taske.txt",
		"data/g4pE_taske.txt",
	}

	taskFilesMap := map[string][]string{
		"data/orig_taska.txt": TaskA,
		"data/orig_taskb.txt": TaskB,
		"data/orig_taskc.txt": TaskC,
		"data/orig_taskd.txt": TaskD,
		"data/orig_taske.txt": TaskE,
	}

	// Matriz de todas las comparaciones
	var comparisons [][]string

	// Comparar los "Originales" con los "Task" y guardar los datos
	for _, original := range Originals {
		compareFiles(original, taskFilesMap[original], &comparisons)
	}

	// Ordenar en base a porcentaje de similitud
	sort.Slice(comparisons, func(i, j int) bool {
		similarityI, _ := strconv.ParseFloat(comparisons[i][2], 64)
		similarityJ, _ := strconv.ParseFloat(comparisons[j][2], 64)
		return similarityI > similarityJ
	})

	fmt.Println("\nTop 10 Archivos más parecidos:")
	for i, comparison := range comparisons {
		if i >= 10 {
			break
		}
		fmt.Printf("Archivo base: %s, Archivo comparado: %s, Similitud: %s%%\n", comparison[0], comparison[1], comparison[2])
	}

	http.HandleFunc("/", handler)
	fmt.Println("")
	fmt.Println("HTML Generado con éxito")
	fmt.Println("Servidor escuchando en http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
