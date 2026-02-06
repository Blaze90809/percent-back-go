import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Dashboard from './components/Dashboard';
import PrivateRoute from './components/PrivateRoute';
import RaceInput from './components/RaceInput';
import RaceTable from './components/RaceTable';
import RaceChart from './components/RaceChart';
import ForgotPassword from './components/ForgotPassword';
import ResetPassword from './components/ResetPassword';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password/:token" element={<ResetPassword />} />
        <Route 
          path="/dashboard"
          element={<PrivateRoute><Dashboard /></PrivateRoute>}
        >
          <Route path="input" element={<RaceInput />} />
          <Route path="table" element={<RaceTable />} />
          <Route path="chart" element={<RaceChart />} />
        </Route>
        <Route path="/" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;
