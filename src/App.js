import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [longUrl, setLongUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setShortUrl('');

    axios.post('http://localhost:8080/shorten', { long_url: longUrl })
      .then(response => {
        setShortUrl(`http://localhost:8080/${response.data.short_code}`);
      })
      .catch(error => {
        setError('Error shortening URL. Please try again.');
        console.error('Error shortening URL:', error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>URL Shortener</h1>
        <form onSubmit={handleSubmit}>
          <input
            type="url"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            placeholder="Enter a long URL"
            required
          />
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Shortening...' : 'Shorten'}
          </button>
        </form>
        {error && <p className="error">{error}</p>}
        {shortUrl && (
          <div className="result">
            <h2>Shortened URL:</h2>
            <a href={shortUrl} target="_blank" rel="noopener noreferrer">
              {shortUrl}
            </a>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
