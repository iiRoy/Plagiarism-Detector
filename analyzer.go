package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// Función para eliminar los caractéres especiales de un texto
func onlyCharacters(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func contains(slice []string, item string) bool {
	item = filepath.Clean(item)
	for _, str := range slice {
		if filepath.Clean(str) == item {
			return true
		}
	}
	return false
}

// Función para cambiar a minúsculas los elementos en el archivo
func toLowerCase(s string) string {
	var result string
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}

// Función para leer un archivo y devolver su contenido
func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		os.Exit(1)
	}

	processedContent := toLowerCase(onlyCharacters(string(content)))

	return processedContent
}

// Función para encontrar la mayor substring común entre dos textos
func longestCommonSubstring(text1, text2 string) string {
	maxLen := 0
	longestSubstr := ""

	// Creamos una tabla para almacenar la longitud de substrings comunes
	table := make([][]int, len(text1)+1)
	for i := range table {
		table[i] = make([]int, len(text2)+1)
	}

	// Llenamos la tabla comparando cada carácter de text1 y text2
	for i := 1; i <= len(text1); i++ {
		for j := 1; j <= len(text2); j++ {
			if text1[i-1] == text2[j-1] {
				table[i][j] = table[i-1][j-1] + 1
				if table[i][j] > maxLen {
					maxLen = table[i][j]
					longestSubstr = text1[i-maxLen : i]
				}
			}
		}
	}

	return longestSubstr
}

// Función para generar n-gramas de un texto
func nGrams(text string, n int) []string {
	words := strings.Fields(text)
	grams := []string{}

	for i := 0; i <= len(words)-n; i++ {
		gram := strings.Join(words[i:i+n], " ")
		grams = append(grams, gram)
	}

	return grams
}

func intersection(ngrams1, ngrams2 []string) []string {
	common := []string{}
	ngramSet := make(map[string]bool)

	// Añadir todos los n-gramas del primer texto a un set
	for _, ngram := range ngrams1 {
		ngramSet[ngram] = true
	}

	// Verificar si los n-gramas del segundo texto están en el set
	for _, ngram := range ngrams2 {
		if ngramSet[ngram] {
			common = append(common, ngram)
		}
	}

	return common
}

// Función para comparar archivos y determinar similitud basada en más del 50% de contenido común
func compareFiles(original string, others []string) {
	// Leer el archivo original
	content1 := readFile(original)
	ngrams1 := nGrams(content1, 3)

	// Comparar el archivo original con los otros archivos
	similarFiles := make([][2]string, 0)

	for _, other := range others {
		content2 := readFile(other)
		ngrams2 := nGrams(content2, 3)

		commonNgrams := intersection(ngrams1, ngrams2)

		// Verificar que la similitud sea mayor al umbral
		if float64(len(commonNgrams))/float64(min(len(ngrams1), len(ngrams2))) > 0.2 { // Cambiado al 20%
			similarFiles = append(similarFiles, [2]string{original, other})
		}
	}

	// Mostrar resultados
	fmt.Println("Similar files (share more than 25% of content):")
	for _, pair := range similarFiles {
		fmt.Printf("File %s and File %s are similar.\n", pair[0], pair[1])
	}
}

// Función auxiliar para encontrar el mínimo entre dos números
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	// Loop through the original files and compare with corresponding task files
	for _, original := range Originals {
		fmt.Printf("\nComparing %s:\n", original)
		compareFiles(original, taskFilesMap[original])
	}
}
