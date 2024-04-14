import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import './App.css';
import { Home } from '../components/Home';
import { Login } from '../components/Login';
import { UserRegistration } from '../components/UserRegistration'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<UserRegistration />} />
        <Route path="*" element={<Home />} /> // Catch-all route
      </Routes>
    </Router>
  );
}


export default App;
