package main

import (
	"fmt"
	"io/ioutil"
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

// Construcción del Suffix Array
func buildSuffixArray(text string) []int {
	n := len(text)
	suffixArr := make([]int, n)
	for i := range suffixArr {
		suffixArr[i] = i
	}
	sort.Slice(suffixArr, func(i, j int) bool {
		return text[suffixArr[i]:] < text[suffixArr[j]:]
	})
	return suffixArr
}

// Encontrar el prefijo común más largo entre 2 Strings
func longestCommonPrefix(text string, i, j int) int {
	n := len(text)
	length := 0
	for i+length < n && j+length < n && text[i+length] == text[j+length] {
		length++
	}
	return length
}

// Búsqueda del Substring más largo en común usadno Suffix Array
func longestCommonSubstring(text1, text2 string) int {
	combined := text1 + "#" + text2
	suffixArr := buildSuffixArray(combined)

	n1 := len(text1)
	longestLen := 0

	for i := 1; i < len(suffixArr); i++ {
		if (suffixArr[i-1] < n1 && suffixArr[i] > n1) || (suffixArr[i-1] > n1 && suffixArr[i] < n1) {
			lcp := longestCommonPrefix(combined, suffixArr[i-1], suffixArr[i])
			if lcp > longestLen {
				longestLen = lcp
			}
		}
	}
	return longestLen
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
		commonLength := longestCommonSubstring(content1, content2)

		if commonLength > 0 {
			percentage := calculateSimilarity(content1, content2, commonLength)
			// Guardar si la comparación es al menos 20% similar
			if percentage >= 20 {
				*comparisons = append(*comparisons, []string{
					original,
					other,
					strconv.FormatFloat(percentage, 'f', 2, 64),
				})
			}
		}
	}
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
}
