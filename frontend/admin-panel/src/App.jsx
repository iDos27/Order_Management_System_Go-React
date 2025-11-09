import React from "react"
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import LoginPage from './pages/LoginPage'
import ProtectedRoute from './components/ProtectedRoute'
import Dashboard from './components/Dashboard'
import OrderManagement from './components/OrderManagement'
import Reports from './components/Reports'
import './App.css'

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route 
            path="/" 
            element={
              <ProtectedRoute>
                <Dashboard>
                  <OrderManagement />
                </Dashboard>
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/reports" 
            element={
              <ProtectedRoute requireAdmin={true}>
                <Dashboard>
                  <Reports />
                </Dashboard>
              </ProtectedRoute>
            } 
          />
        </Routes>
      </AuthProvider>
    </Router>
  )
}

export default App