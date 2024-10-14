package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

// Función para comparar archivos y determinar similitud basada en más del 50% de contenido común
func compareFiles(originals, others []string) {
	fileContents := make(map[string]string)

	// Leer todos los archivos
	for _, file := range originals {
		content := readFile(file)
		fileContents[file] = content
	}

	// Comparar cada original con el resto de los archivos
	similarFiles := make([][2]string, 0)

	for _, original := range originals {
		content1 := fileContents[original]
		for _, other := range others {
			// Contenidos del 2° Archivo
			content2 := readFile(other)

			// Encontrar el Substring más Largo
			longestCommon := longestCommonSubstring(content1, content2)

			// Determinar el Substring común más largo
			lenCommon := len(longestCommon)
			lenText1 := len(content1)
			lenText2 := len(content2)

			// Verificar que la similitud sea mayor al 50%
			if float64(lenCommon) >= 0.5*float64(min(lenText1, lenText2)) {
				similarFiles = append(similarFiles, [2]string{original, other})
			}
		}
	}

	// Mostrar resultados
	fmt.Println("Similar files (share more than 50% of content):")
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

	// Leer todos los archivos de la carpeta "data"
	files, err := ioutil.ReadDir("data")
	if err != nil {
		fmt.Println("Error al leer la carpeta data:", err)
		os.Exit(1)
	}

	// Filtrar solo los archivos .txt
	var otherFiles []string
	for _, file := range files {
		filePath := filepath.Join("data", file.Name())
		if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" && !contains(Originals, filePath) {
			otherFiles = append(otherFiles, filePath)
		}
	}

	// Comparar los archivos para determinar similitud
	compareFiles(Originals, otherFiles)
}
