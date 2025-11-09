import React, { useContext } from 'react';
import { Navigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';

const ProtectedRoute = ({ children, requiredRole = null, requireAdmin = false }) => {
  const { user, isLoading } = useContext(AuthContext);

  // Show loading while checking auth status
  if (isLoading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontSize: '18px',
        color: '#666'
      }}>
        Sprawdzanie uprawnień...
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // Check if admin is required
  if (requireAdmin && user.role !== 'admin') {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        textAlign: 'center',
        padding: '20px'
      }}>
        <h2 style={{ color: '#e74c3c', marginBottom: '10px' }}>
          Brak uprawnień administratora
        </h2>
        <p style={{ color: '#666', fontSize: '16px' }}>
          Dostęp do raportów mają tylko administratorzy.
        </p>
        <p style={{ color: '#666', fontSize: '14px', marginTop: '10px' }}>
          Twoja rola: <strong>{user.role}</strong>
        </p>
      </div>
    );
  }

  // Check role if required (backwards compatibility)
  if (requiredRole && user.role !== requiredRole) {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        textAlign: 'center',
        padding: '20px'
      }}>
        <h2 style={{ color: '#e74c3c', marginBottom: '10px' }}>
          Brak uprawnień
        </h2>
        <p style={{ color: '#666', fontSize: '16px' }}>
          Nie masz uprawnień do tej części systemu.
        </p>
        <p style={{ color: '#666', fontSize: '14px', marginTop: '10px' }}>
          Twoja rola: <strong>{user.role}</strong>
        </p>
        <p style={{ color: '#666', fontSize: '14px' }}>
          Wymagana rola: <strong>{requiredRole}</strong>
        </p>
      </div>
    );
  }

  return children;
};

export default ProtectedRoute;