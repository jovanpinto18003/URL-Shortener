import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

const App = () => {
    const [shortURL, setShortURL] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (url) => {
        setError('');
        try {
            const response = await axios.post('http://localhost:8080/shorten', { url }, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            setShortURL(response.data.short_url);
        } catch (error) {
            console.error('Error shortening URL:', error);
            setError('Failed to shorten URL. Please try again.');
        }
    };

    return (
        <div className="App">
            <header>
                <h1>URL Shortener</h1>
            </header>
            <form
                onSubmit={(e) => {
                    e.preventDefault();
                    const url = e.target.elements.url.value;
                    handleSubmit(url);
                }}
            >
                <input
                    type="text"
                    name="url"
                    placeholder="Enter URL"
                    required
                />
                <button type="submit">Shorten</button>
            </form>
            {error && <p className="error">{error}</p>}
            {shortURL && (
                <div className="result">
                    <a href={shortURL} target="_blank" rel="noopener noreferrer">
                        <input type="text" value={shortURL} readOnly />
                    </a>
                </div>
            )}
        </div>
    );
};

export default App;
