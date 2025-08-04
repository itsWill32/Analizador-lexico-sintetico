package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime" // Importamos el paquete para medir la memoria
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// --- ESTRUCTURAS DE DATOS ---

type AnalyzeRequest struct {
	Code string `json:"code"`
}

type AnalysisResponse struct {
	IsValid             bool    `json:"isValid"`
	Message             string  `json:"message"`
	ErrorDetail         string  `json:"errorDetail"`
	ErrorType           string  `json:"errorType"`
	Tokens              []Token `json:"tokens"`
	OptimizedCode       string  `json:"optimizedCode"`
	OriginalSize        int     `json:"originalSize"`
	OptimizedSize       int     `json:"optimizedSize"`
	ReductionPercentage float64 `json:"reductionPercentage"`
	ServerMemoryUsage   string  `json:"serverMemoryUsage"` // Métrica de memoria del servidor
}

type Token struct {
	Line  int    `json:"line"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// --- MIDDLEWARE Y HANDLER PRINCIPAL ---

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	result := api.Transform(req.Code, api.TransformOptions{
		Loader:    api.LoaderTSX,
		Target:    api.ES2020,
		Sourcemap: api.SourceMapNone,
	})

	response := AnalysisResponse{}

	if len(result.Errors) > 0 {
		response.IsValid = false
		errorMsg := result.Errors[0].Text
		response.Message = "Se encontró un error en el código"
		response.ErrorDetail = fmt.Sprintf("Línea %d: %s", result.Errors[0].Location.Line, errorMsg)

		// Lógica de clasificación de errores mejorada y verificada
		errorTextLower := strings.ToLower(errorMsg)
		if strings.Contains(errorTextLower, "invalid character") {
			response.ErrorType = "LEXICAL"
		} else if strings.Contains(errorTextLower, "type") || strings.Contains(errorTextLower, "assignable") {
			response.ErrorType = "SEMANTIC"
		} else {
			response.ErrorType = "SYNTACTIC" // Todos los demás errores de parseo son sintácticos
		}

	} else {
		response.IsValid = true
		response.Message = "Análisis y optimización completados"
		optimizedCode := optimizeCode(req.Code)
		response.OptimizedCode = optimizedCode
		originalSize := len(req.Code)
		optimizedSize := len(optimizedCode)
		response.OriginalSize = originalSize
		response.OptimizedSize = optimizedSize
		if originalSize > 0 {
			response.ReductionPercentage = (float64(originalSize-optimizedSize) / float64(originalSize)) * 100
		} else {
			response.ReductionPercentage = 0
		}
		tokens, _ := tokenizeCode(req.Code)
		response.Tokens = tokens
	}

	// Medir la memoria del servidor después del análisis
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// Formateamos los bytes a Megabytes para que sea más legible
	response.ServerMemoryUsage = fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// optimizeCode elimina las sentencias `console.log(...)` del código.
func optimizeCode(code string) string {
	re, err := regexp.Compile(`console\.log\((?s).*?\);?`)
	if err != nil {
		return code
	}
	optimized := re.ReplaceAllString(code, "")
	lines := strings.Split(optimized, "\n")
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return strings.Join(nonEmptyLines, "\n")
}

// tokenizeCode realiza un análisis léxico simple para mostrar los tokens.
func tokenizeCode(code string) ([]Token, error) {
	re := regexp.MustCompile(`(import|from|export|default|const|let|var|function|return|=>|\{|\}|\(|\)|\[|\]|;|,|\.|:|\S+)`)
	matches := re.FindAllString(code, -1)
	var tokens []Token
	lines := strings.Split(code, "\n")
	lineNum := 1
	lastIndex := 0
	for _, match := range matches {
		found := false
		for ; lineNum <= len(lines); lineNum++ {
			if idx := strings.Index(lines[lineNum-1][lastIndex:], match); idx != -1 {
				found = true
				lastIndex += idx
				break
			}
			lastIndex = 0
		}
		if !found {
			continue
		}
		tokenType := "LITERAL"
		switch match {
		case "import", "export", "from", "const", "let", "return", "default", "function":
			tokenType = "KEYWORD"
		case "{", "}", "(", ")", "[", "]", ";", ",", ":", ".":
			tokenType = "PUNCTUATION"
		case "=>":
			tokenType = "ARROW_FUNCTION"
		}
		tokens = append(tokens, Token{
			Line:  lineNum,
			Type:  tokenType,
			Value: match,
		})
	}
	return tokens, nil
}

// main inicia el servidor
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", analyzeHandler)
	handler := corsMiddleware(mux)
	fmt.Println("Servidor de análisis y optimización robusto iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
