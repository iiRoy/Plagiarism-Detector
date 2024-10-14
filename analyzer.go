package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// Suffix Estructura para almacenar un sufijo y su índice
type Suffix struct {
	index int
	suff  string
}

// Función para leer un archivo y devolver su contenido como una cadena
func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}
	return string(content)
}

// Función para generar el Suffix Array
func buildSuffixArray(s string) []int {
	var suffixes []Suffix

	// Generar todos los sufijos
	for i := 0; i < len(s); i++ {
		suffixes = append(suffixes, Suffix{i, s[i:]})
	}

	// Ordenar los sufijos basados en la cadena de sufijos
	sort.Slice(suffixes, func(i, j int) bool {
		return suffixes[i].suff < suffixes[j].suff
	})

	// Extraer solo los índices para el Suffix Array
	suffixArray := make([]int, len(s))
	for i := range suffixes {
		suffixArray[i] = suffixes[i].index
	}

	return suffixArray
}

// Función para contar los caracteres en BWT
func countChars(bwt string) map[byte]int {
	charCount := make(map[byte]int)
	for i := 0; i < len(bwt); i++ {
		charCount[bwt[i]]++
	}
	return charCount
}

// Función para construir el arreglo C
func buildC(bwt string) map[byte]int {
	charCount := countChars(bwt)
	var chars []byte
	for ch := range charCount {
		chars = append(chars, ch)
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	C := make(map[byte]int)
	total := 0
	for _, ch := range chars {
		C[ch] = total
		total += charCount[ch]
	}

	fmt.Println("Paso 1: Construir el arreglo C (número de caracteres anteriores):", C)
	return C
}

// Función para construir el arreglo Occur
func buildOccur(bwt string) map[byte][]int {
	charCount := countChars(bwt)
	Occur := make(map[byte][]int)
	for ch := range charCount {
		Occur[ch] = make([]int, len(bwt)+1)
	}

	for i := 0; i < len(bwt); i++ {
		for ch := range charCount {
			if i > 0 {
				Occur[ch][i] = Occur[ch][i-1]
			}
		}
		Occur[bwt[i]][i]++
	}

	fmt.Println("Paso 2: Construir el arreglo Occur (ocurrencias por carácter hasta cada posición):", Occur)
	return Occur
}

// Función para realizar la transformación inversa usando LF-mapping
func inverseBWT(bwt string, C map[byte]int, Occur map[byte][]int) string {
	original := make([]byte, len(bwt))
	var index int

	// Encontrar el índice del símbolo '$' que marca el final de la cadena
	for i := 0; i < len(bwt); i++ {
		if bwt[i] == '$' {
			index = i
			break
		}
	}

	fmt.Println("Paso 3: Encontrar la posición del símbolo '$' (posición inicial en BWT):", index)

	// Reconstrucción de la cadena original
	fmt.Println("Paso 4: Reconstrucción de la cadena original usando LF-mapping:")
	for i := len(bwt) - 1; i >= 0; i-- {
		original[i] = bwt[index]
		fmt.Printf("Posición %d -> Carácter '%c' -> Índice en LF: %d\n", i, bwt[index], index)
		index = C[bwt[index]] + Occur[bwt[index]][index] - 1
	}

	return string(original)
}

// Función que ejecuta el proceso de BWT
func runBWT(bwt string) (map[byte]int, map[byte][]int, string) {
	fmt.Println("Paso 0: BWT de entrada:", bwt)

	// Construir el arreglo C
	c := buildC(bwt)

	// Construir el arreglo Occur
	occur := buildOccur(bwt)

	// Realizar la transformación inversa usando LF-mapping
	original := inverseBWT(bwt, c, occur)

	return c, occur, original
}

// Función para realizar la búsqueda hacia atrás (backwardSearch) usando FM-Index
func backwardSearch(pattern string, bwt string) []int {
	SuffixArray := buildSuffixArray(bwt)
	bwtC, bwtOcc, _ := runBWT(bwt)

	i := len(pattern) - 1
	c := pattern[i]

	// Verificamos si el carácter está en el mapa C
	if _, exists := bwtC[c]; !exists {
		return []int{} // El carácter no está en el texto
	}

	// Inicializamos los límites l y r
	l := bwtC[c]
	r := len(bwt) - 1
	if next, exists := bwtC[c+1]; exists {
		r = next - 1
	}

	// Iteramos sobre el patrón de derecha a izquierda
	for l <= r && i > 0 {
		i--
		c = pattern[i]

		if _, exists := bwtC[c]; !exists {
			return []int{} // Si el carácter no está en C
		}

		// Actualizamos los límites usando el LF-mapping
		l = bwtC[c] + bwtOcc[c][l-1]
		r = bwtC[c] + bwtOcc[c][r] - 1
	}

	// Si l <= r, retornamos el rango en el suffix array
	if l <= r {
		return SuffixArray[l : r+1]
	}
	return []int{}
}

func main() {
	// Llama a "Suffix Array"
	fmt.Println("Ejecutando Suffix Array")

	// Get all files in "data" directory
	files, err := ioutil.ReadDir("data")
	if err != nil {
		fmt.Println("Error al leer la carpeta data:", err)
		os.Exit(1)
	}

	// Iterate through all files in "data" directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" {
			filePath := filepath.Join("data", file.Name())
			fmt.Println("Procesando archivo:", filePath)

			// Read and process each .txt file
			content := readFile(filePath)
			content = content + "$"

			// Construye el Suffix Array
			suffixArray := buildSuffixArray(content)
			fmt.Println("Suffix Array:")
			fmt.Println(suffixArray)

			// Llama a la función "Burrow's Wheeler Transform"
			fmt.Println("Ejecutando Burrow's Wheeler Transform")
			bwt := "annb$aa" // Ejemplo de BWT para la cadena "banana"
			_, _, original := runBWT(bwt)
			fmt.Println("Cadena original reconstruida:")
			fmt.Println(original)

			// Llama a la función "FM-Index"
			fmt.Println("Ejecutando FM-Index")
			pattern := "ban"
			positions := backwardSearch(pattern, bwt)
			fmt.Println("Posiciones del patrón:", positions)
		}
	}
}
