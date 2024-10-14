package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
)

// Función para leer un archivo y devolver su contenido
func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}
	return string(content)
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
func compareFiles(files []string) {
	fileContents := make(map[string]string)

	// Leer todos los archivos
	for _, file := range files {
		content := readFile(file)
		fileContents[file] = content
	}

	// Comparar los archivos
	similarFiles := make([][2]string, 0)

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			// Obtener el contenido de los archivos
			content1 := fileContents[files[i]]
			content2 := fileContents[files[j]]

			// Encontrar la mayor substring común
			longestCommon := longestCommonSubstring(content1, content2)

			// Determinar el tamaño de la substring común
			lenCommon := len(longestCommon)
			lenText1 := len(content1)
			lenText2 := len(content2)

			// Verificar si la substring común es mayor al 50% del texto más corto
			if float64(lenCommon) >= 0.3*float64(min(lenText1, lenText2)) {
				similarFiles = append(similarFiles, [2]string{files[i], files[j]})
			}
		}
	}

	// Mostrar resultados finales
	fmt.Println("Archivos similares (comparten más del 50% de contenido):")
	for _, pair := range similarFiles {
		fmt.Printf("Archivo %s y Archivo %s son similares.\n", pair[0], pair[1])
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
	// Leer todos los archivos de la carpeta "data"
	files, err := ioutil.ReadDir("data")
	if err != nil {
		fmt.Println("Error al leer la carpeta data:", err)
		os.Exit(1)
	}

	// Filtrar solo los archivos .txt
	var textFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" {
			textFiles = append(textFiles, filepath.Join("data", file.Name()))
		}
	}

	// Comparar los archivos para determinar similitud
	compareFiles(textFiles)
}