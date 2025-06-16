import { useState } from 'react'
import './App.css'

const EXAMPLE_CODE = `import React, { useState, useEffect } from 'react';
import { NextPage } from 'next';


const HomePage: NextPage = () => {
  return (
    <>
      <h1>Añade tu codigo a analizar</h1>
    </>
  );
};

export default HomePage;`;

function App() {
  const [code, setCode] = useState(EXAMPLE_CODE);
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(null);

  const analyzeCode = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code }),
      });

      if (!response.ok) {
        throw new Error('Error en la comunicación con el servidor');
      }

      const data = await response.json();
      setResult(data);
    } catch (error) {
      setResult({
        isValid: false,
        message: 'Error de conexión',
        errorDetail: error.message,
        tokens: []
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="app">
      <div className="container">
        <h1>Analizador Léxico y Sintáctico para Next.js/TSX</h1>
        
        <div className="input-section">
          <label htmlFor="code-input">Código Next.js/TSX:</label>
          <textarea
            id="code-input"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="Pega aquí tu código Next.js/TSX..."
            rows={20}
            cols={80}
          />
          
          <button 
            onClick={analyzeCode} 
            disabled={isLoading || !code.trim()}
            className="analyze-button"
          >
            {isLoading ? 'Analizando...' : 'Analizar Código'}
          </button>
        </div>

        {isLoading && (
          <div className="loading">
            <div className="loading-spinner"></div>
            <p>Analizando código...</p>
          </div>
        )}

        {result && !isLoading && (
          <div className="results-section">
            {result.isValid ? (
              <div className="success-section">
                <div className="success-message">
                  <h2>✅ {result.message}</h2>
                  <p>Se encontraron {result.tokens.length} tokens</p>
                </div>
                
                <div className="tokens-table-container">
                  <table className="tokens-table">
                    <thead>
                      <tr>
                        <th>Línea</th>
                        <th>Tipo de Token</th>
                        <th>Valor</th>
                      </tr>
                    </thead>
                    <tbody>
                      {result.tokens.map((token, index) => (
                        <tr key={index}>
                          <td>{token.line}</td>
                          <td className="token-type">{token.type}</td>
                          <td className="token-value">{token.value}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            ) : (
              <div className="error-section">
                <div className="error-message">
                  <h2>❌ {result.message}</h2>
                  <div className="error-detail">
                    <strong>Detalle del error:</strong>
                    <p>{result.errorDetail}</p>
                  </div>
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