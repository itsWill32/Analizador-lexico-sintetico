import { useState } from 'react';
import './App.css';

// --- CÓDIGOS DE EJEMPLO PARA PRUEBAS ---
const EXAMPLE_CODE_OK = `import { NextPage } from 'next';
// Ejemplo válido para demostrar optimización
const HomePage: NextPage = () => {
  console.log("Renderizando..."); // Esta línea será eliminada
  return <h1>Análisis Exitoso</h1>;
};
export default HomePage;`;

const EXAMPLE_LEXICAL_ERROR = `const x = 1 # 2; // Error: '#' es un carácter inválido en JS/TSX`;

const EXAMPLE_SYNTACTIC_ERROR = `import { NextPage } from 'next';
// Error: falta la llave de cierre '}' en la función
const HomePage: NextPage = () => {
  return <h1>Error Sintáctico`;

const EXAMPLE_SEMANTIC_ERROR = `import { NextPage } from 'next';
// Error: no se puede asignar un número a un string
const HomePage: NextPage = () => {
  const user: string = 123;
  return <h1>Hola, {user}</h1>;
};
export default HomePage;`;


function App() {
  const [code, setCode] = useState(EXAMPLE_CODE_OK);
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(null);

  const analyzeCode = async () => {
    setIsLoading(true);
    setResult(null);
    try {
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code }),
      });

      if (!response.ok) throw new Error(`Error del servidor: ${response.statusText}`);
      const data = await response.json();
      setResult(data);
    } catch (error) {
      setResult({ isValid: false, message: 'Error de Conexión', errorDetail: error.message, errorType: 'CONNECTION' });
    } finally {
      setIsLoading(false);
    }
  };

  const loadExample = (exampleCode) => {
    setCode(exampleCode);
    setResult(null);
  };

  return (
    <div className="app">
      <div className="container">
        <h1>Analizador Robusto con Diagnóstico Avanzado</h1>
        
        <div className="input-section">
          <div className="example-buttons">
            <button onClick={() => loadExample(EXAMPLE_CODE_OK)}>Ejemplo Válido</button>
            {/* <button onClick={() => loadExample(EXAMPLE_LEXICAL_ERROR)}>Error Léxico</button> */}
            <button onClick={() => loadExample(EXAMPLE_SYNTACTIC_ERROR)}>Error Sintáctico</button>
            {/* <button onClick={() => loadExample(EXAMPLE_SEMANTIC_ERROR)}>Error Semántico</button> */}
          </div>
          <textarea
            id="code-input"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            rows={15}
            placeholder="Pega aquí tu código Next.js/TSX..."
          />
          <button onClick={analyzeCode} disabled={isLoading || !code.trim()} className="analyze-button">
            {isLoading ? 'Analizando...' : 'Analizar y Optimizar'}
          </button>
        </div>

        {isLoading && <div className="loading">Analizando...</div>}

        {result && !isLoading && (
          <div className="results-container">
            {result.isValid ? (
              <>
                <div className="result-card success">
                  <h2>✅ {result.message}</h2>
                </div>
                
                <div className="result-card metrics">
                  <h2>📊 Métricas de Optimización (Código)</h2>
                  <div className="metrics-grid">
                    <div className="metric-item"><span>Tamaño Original</span><p>{result.originalSize} Bytes</p></div>
                    <div className="metric-item"><span>Tamaño Optimizado</span><p>{result.optimizedSize} Bytes</p></div>
                    <div className="metric-item reduction"><span>Reducción ⬇️</span><p>{result.reductionPercentage.toFixed(2)}%</p></div>
                  </div>
                </div>

                <div className="result-card server-metrics">
                    <h2>⚙️ Métricas del Servidor (Backend)</h2>
                    <div className="metrics-grid">
                        <div className="metric-item">
                            <span>Uso de Memoria</span>
                            <p>{result.serverMemoryUsage}</p>
                        </div>
                    </div>
                    <small>Nota: Se muestra el uso de memoria por ser una métrica estable. Medir el uso de CPU para una petición tan rápida no es un indicador fiable.</small>
                </div>
                
                <div className="result-card">
                  <h2>Código Optimizado</h2>
                  <pre className="code-block">{result.optimizedCode}</pre>
                </div>
              </>
            ) : (
              <div className="error-panels">
                <div className={`error-card lexical ${result.errorType === 'LEXICAL' ? 'active' : ''}`}>
                  <h2>🚫 Error Léxico</h2>
                  {result.errorType === 'LEXICAL' && <pre className="error-detail">{result.errorDetail}</pre>}
                </div>
                <div className={`error-card syntactic ${result.errorType === 'SYNTACTIC' ? 'active' : ''}`}>
                  <h2>📐 Error Sintáctico</h2>
                  {result.errorType === 'SYNTACTIC' && <pre className="error-detail">{result.errorDetail}</pre>}
                </div>
                <div className={`error-card semantic ${result.errorType === 'SEMANTIC' ? 'active' : ''}`}>
                  <h2>🧠 Error Semántico</h2>
                  {result.errorType === 'SEMANTIC' && <pre className="error-detail">{result.errorDetail}</pre>}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
