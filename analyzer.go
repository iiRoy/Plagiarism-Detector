package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

// Función para eliminar los caracteres especiales de un texto
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

// SuffixArray es una estructura que contiene los sufijos de un texto y permite buscar substrings
type SuffixArray struct {
	text    string
	suffixes []int
}

// Crear un sufijo array para un texto
func buildSuffixArray(text string) SuffixArray {
	n := len(text)
	suffixes := make([]int, n)

	// Generar todos los sufijos
	for i := 0; i < n; i++ {
		suffixes[i] = i
	}

	// Ordenar los sufijos alfabéticamente
	sort.Slice(suffixes, func(i, j int) bool {
		return text[suffixes[i]:] < text[suffixes[j]:]
	})

	return SuffixArray{
		text:    text,
		suffixes: suffixes,
	}
}

// Buscar el substring más largo entre dos textos utilizando Suffix Array
func longestCommonSubstringSA(text1, text2 string) string {
	sa1 := buildSuffixArray(text1)
	sa2 := buildSuffixArray(text2)

	maxLen := 0
	longestSubstr := ""

	// Buscar el substring más largo común entre los dos sufijo arrays
	for _, idx1 := range sa1.suffixes {
		for _, idx2 := range sa2.suffixes {
			length := 0
			for idx1+length < len(text1) && idx2+length < len(text2) && text1[idx1+length] == text2[idx2+length] {
				length++            // Comparar caracteres en ambos textos
			}
			if length > maxLen {
				maxLen = length
				longestSubstr = text1[idx1 : idx1+length]
			}
		}
	}

	return longestSubstr
}

// Función para calcular la distancia de Levenshtein entre dos cadenas
func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Crear la matriz
	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
	}

	// Inicializar la primera fila y columna
	for i := 0; i <= la; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	// Llenar la matriz
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

// Función para comparar archivos y determinar similitud basada en más del 25% de contenido común
func compareFiles(original string, others []string) {
    // Leer el archivo original
    content1 := readFile(original)

    // Crear una lista para almacenar archivos similares
    similarFiles := make([][2]string, 0)

    // Comparar el archivo original con los otros archivos
    for _, other := range others {
        content2 := readFile(other)

        // Encontrar el Substring más Largo usando Suffix Array
        longestCommon := longestCommonSubstringSA(content1, content2)

        // Determinar el Substring común más largo
        lenCommon := len(longestCommon)
        lenText1 := len(content1)
        lenText2 := len(content2)

        // Verificar que la similitud sea mayor al 25% usando el Suffix Array
        if float64(lenCommon) >= 0.25*float64(min(lenText1, lenText2)) {
            similarFiles = append(similarFiles, [2]string{original, other})
        } else {
            // Si no cumple el criterio del Suffix Array, calcular Levenshtein
            distance := levenshtein(content1, content2)
            similarity := 1 - float64(distance)/float64(max(lenText1, lenText2))

            // Verificar si la similitud basada en Levenshtein es mayor al 75%
            if similarity >= 0.75 {
                similarFiles = append(similarFiles, [2]string{original, other})
            }
        }
    }

    // Mostrar resultados
    if len(similarFiles) > 0 {
        fmt.Println("Archivos similares (comparten más del 25% de contenido o 75% de similitud con Levenshtein):")
        for _, pair := range similarFiles {
            fmt.Printf("File %s y File %s son similares.\n", pair[0], pair[1])
        }
    } else {
        fmt.Println("No se encontraron archivos similares para", original)
    }
}

// Función auxiliar para encontrar el mínimo entre dos números
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Función auxiliar para encontrar el máximo entre dos números
func max(a, b int) int {
	if a > b {
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

	// Mapa de los archivos de tareas
	taskFilesMap := map[string][]string{
		"data/orig_taska.txt": TaskA,
		"data/orig_taskb.txt": TaskB,
		"data/orig_taskc.txt": TaskC,

		"data/orig_taskd.txt": TaskD,

		"data/orig_taske.txt": TaskE,

		// Añade más tareas...
	}

	// Comparar archivos
	for _, original := range Originals {
		fmt.Printf("\nComparando %s:\n", original)
		compareFiles(original, taskFilesMap[original])
	}
}
