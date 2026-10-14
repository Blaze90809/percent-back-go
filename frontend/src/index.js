import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import axios from 'axios';
import { Capacitor } from '@capacitor/core';

// Set Axios base URL dynamically for Capacitor mobile viewports
if (Capacitor.isNativePlatform()) {
  // If a custom API URL is provided in the build environment (e.g., ngrok for local dev), use it.
  // Otherwise, default to the production backend.
  axios.defaults.baseURL = process.env.REACT_APP_API_URL || 'https://nordicracetrack.com'; 
} else {
  // On standard web browser, use relative paths
  axios.defaults.baseURL = window.location.origin;
}

// Bypass ngrok free tier browser warning screen for local development API requests
axios.defaults.headers.common['ngrok-skip-browser-warning'] = 'true';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
