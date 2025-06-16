package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

type AnalyzeRequest struct {
	Code string `json:"code"`
}

type Token struct {
	Line  int    `json:"line"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AnalysisResponse struct {
	IsValid     bool    `json:"isValid"`
	Message     string  `json:"message"`
	ErrorDetail string  `json:"errorDetail"`
	Tokens      []Token `json:"tokens"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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
		Loader: api.LoaderTSX,
		Target: api.ES2020,
	})

	response := AnalysisResponse{}

	if len(result.Errors) > 0 {
		response.IsValid = false
		response.Message = "Error de sintaxis detectado"
		response.ErrorDetail = fmt.Sprintf("Línea %d: %s",
			result.Errors[0].Location.Line,
			result.Errors[0].Text)
		response.Tokens = []Token{}
	} else {
		tokens, err := tokenizeCode(req.Code)
		if err != nil {
			response.IsValid = false
			response.Message = "Error en el análisis léxico"
			response.ErrorDetail = err.Error()
			response.Tokens = []Token{}
		} else {
			response.IsValid = true
			response.Message = "Análisis completado exitosamente"
			response.ErrorDetail = ""
			response.Tokens = tokens
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func tokenizeCode(code string) ([]Token, error) {
	lexer := js.NewLexer(parse.NewInputString(code))

	var tokens []Token
	currentLine := 1

	for {
		tokenType, text := lexer.Next()
		if tokenType == js.ErrorToken {
			break
		}

		tokenTypeStr := getTokenTypeName(tokenType)
		tokenValue := string(text)

		if tokenType != js.WhitespaceToken && tokenType != js.LineTerminatorToken {
			tokens = append(tokens, Token{
				Line:  currentLine,
				Type:  tokenTypeStr,
				Value: tokenValue,
			})
		}

		currentLine += strings.Count(tokenValue, "\n")
	}

	return tokens, nil
}

func getTokenTypeName(tokenType js.TokenType) string {
	switch tokenType {
	case js.ErrorToken:
		return "ERROR"
	case js.IdentifierToken:
		return "IDENTIFIER"
	case js.NumericToken:
		return "NUMERIC"
	case js.StringToken:
		return "STRING"
	case js.TemplateToken:
		return "TEMPLATE"
	case js.BreakToken:
		return "BREAK"
	case js.CaseToken:
		return "CASE"
	case js.CatchToken:
		return "CATCH"
	case js.ClassToken:
		return "CLASS"
	case js.ConstToken:
		return "CONST"
	case js.ContinueToken:
		return "CONTINUE"
	case js.DebuggerToken:
		return "DEBUGGER"
	case js.DefaultToken:
		return "DEFAULT"
	case js.DeleteToken:
		return "DELETE"
	case js.DoToken:
		return "DO"
	case js.ElseToken:
		return "ELSE"
	case js.EnumToken:
		return "ENUM"
	case js.ExportToken:
		return "EXPORT"
	case js.ExtendsToken:
		return "EXTENDS"
	case js.FalseToken:
		return "FALSE"
	case js.FinallyToken:
		return "FINALLY"
	case js.ForToken:
		return "FOR"
	case js.FunctionToken:
		return "FUNCTION"
	case js.IfToken:
		return "IF"
	case js.ImportToken:
		return "IMPORT"
	case js.InToken:
		return "IN"
	case js.InstanceofToken:
		return "INSTANCEOF"
	case js.InterfaceToken:
		return "INTERFACE"
	case js.LetToken:
		return "LET"
	case js.NewToken:
		return "NEW"
	case js.NullToken:
		return "NULL"
	case js.PackageToken:
		return "PACKAGE"
	case js.PrivateToken:
		return "PRIVATE"
	case js.ProtectedToken:
		return "PROTECTED"
	case js.PublicToken:
		return "PUBLIC"
	case js.ReturnToken:
		return "RETURN"
	case js.StaticToken:
		return "STATIC"
	case js.SuperToken:
		return "SUPER"
	case js.SwitchToken:
		return "SWITCH"
	case js.ThisToken:
		return "THIS"
	case js.ThrowToken:
		return "THROW"
	case js.TrueToken:
		return "TRUE"
	case js.TryToken:
		return "TRY"
	case js.TypeofToken:
		return "TYPEOF"
	case js.VarToken:
		return "VAR"
	case js.VoidToken:
		return "VOID"
	case js.WhileToken:
		return "WHILE"
	case js.WithToken:
		return "WITH"
	case js.YieldToken:
		return "YIELD"
	case js.OpenParenToken:
		return "OPEN_PAREN"
	case js.CloseParenToken:
		return "CLOSE_PAREN"
	case js.OpenBraceToken:
		return "OPEN_BRACE"
	case js.CloseBraceToken:
		return "CLOSE_BRACE"
	case js.OpenBracketToken:
		return "OPEN_BRACKET"
	case js.CloseBracketToken:
		return "CLOSE_BRACKET"
	case js.SemicolonToken:
		return "SEMICOLON"
	case js.CommaToken:
		return "COMMA"
	case js.DotToken:
		return "DOT"
	case js.ColonToken:
		return "COLON"
	case js.OrToken:
		return "OR"
	case js.AndToken:
		return "AND"
	case js.BitOrToken:
		return "BIT_OR"
	case js.BitXorToken:
		return "BIT_XOR"
	case js.BitAndToken:
		return "BIT_AND"
	case js.EqToken:
		return "EQUAL"
	case js.EqEqToken:
		return "STRICT_EQUAL"
	case js.LtToken:
		return "LESS_THAN"
	case js.GtToken:
		return "GREATER_THAN"
	case js.AddToken:
		return "ADD"
	case js.SubToken:
		return "SUBTRACT"
	case js.MulToken:
		return "MULTIPLY"
	case js.DivToken:
		return "DIVIDE"
	case js.ModToken:
		return "MODULO"
	case js.ExpToken:
		return "EXPONENT"
	case js.IncrToken:
		return "INCREMENT"
	case js.DecrToken:
		return "DECREMENT"
	case js.NotToken:
		return "NOT"
	case js.BitNotToken:
		return "BIT_NOT"
	case js.ArrowToken:
		return "ARROW"
	case js.EllipsisToken:
		return "ELLIPSIS"
	case js.LineTerminatorToken:
		return "LINE_TERMINATOR"
	case js.WhitespaceToken:
		return "WHITESPACE"
	default:

		tokenStr := string(tokenType)
		if strings.Contains(tokenStr, "=") {
			return "ASSIGNMENT_OP"
		}
		return fmt.Sprintf("TOKEN_%d", int(tokenType))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", analyzeHandler)

	handler := corsMiddleware(mux)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
