import React, { useState, useContext } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { AUTH_BASE_URL } from '../services/api';
import './Dashboard.css';

const Dashboard = ({ children }) => {
  const { user, logout } = useContext(AuthContext);
  const navigate = useNavigate();
  const location = useLocation();
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const [registerForm, setRegisterForm] = useState({
    email: '',
    password: '',
    role: 'employee'
  });
  const [registerLoading, setRegisterLoading] = useState(false);
  const [registerError, setRegisterError] = useState('');
  const [registerSuccess, setRegisterSuccess] = useState('');

  const handleRegisterChange = (e) => {
    setRegisterForm({
      ...registerForm,
      [e.target.name]: e.target.value
    });
  };

  const handleRegisterSubmit = async (e) => {
    e.preventDefault();
    setRegisterLoading(true);
    setRegisterError('');
    setRegisterSuccess('');

    try {
      const response = await fetch(`${AUTH_BASE_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(registerForm)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Błąd rejestracji');
      }

      setRegisterSuccess('Użytkownik został pomyślnie zarejestrowany!');
      setRegisterForm({
        email: '',
        password: '',
        role: 'employee'
      });
      
      // Auto-close modal after success
      setTimeout(() => {
        setShowRegisterModal(false);
        setRegisterSuccess('');
      }, 2000);

    } catch (error) {
      setRegisterError(error.message);
    } finally {
      setRegisterLoading(false);
    }
  };

  const closeModal = () => {
    setShowRegisterModal(false);
    setRegisterError('');
    setRegisterSuccess('');
    setRegisterForm({
      email: '',
      password: '',
      role: 'employee'
    });
  };

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="header-left">
          <h1>System zarządzania zamówieniami z wykorzystaniem technologii React.js oraz Go</h1>
        </div>
        
        <div className="header-right">
          <div className="user-info">
            <span className="welcome-text">Witaj, </span>
            <span className="user-email">{user?.email}</span>
            <span className={`role-badge role-${user?.role}`}>
              {user?.role === 'admin' ? 'Administrator' : 
               user?.role === 'employee' ? 'Pracownik' : user?.role}
            </span>
          </div>
          
          <div className="header-actions">
            {user?.role === 'admin' && (
              <button 
                className="register-btn"
                onClick={() => setShowRegisterModal(true)}
              >
                + Dodaj użytkownika
              </button>
            )}
            <button 
              className="logout-btn"
              onClick={logout}
            >
              Wyloguj
            </button>
          </div>
        </div>
      </header>

      {/* Navigation - Only for admins */}
      {user?.role === 'admin' && (
        <nav className="dashboard-navigation">
          <button 
            className={`nav-btn ${location.pathname === '/' ? 'active' : ''}`}
            onClick={() => navigate('/')}
          >
            📋 Zamówienia
          </button>
          <button 
            className={`nav-btn ${location.pathname === '/reports' ? 'active' : ''}`}
            onClick={() => navigate('/reports')}
          >
            📊 Raporty
          </button>
        </nav>
      )}

      <main className="dashboard-content">
        {children}
      </main>

      {/* Register Modal - Only for admins */}
      {showRegisterModal && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Dodaj nowego użytkownika</h3>
              <button className="close-btn" onClick={closeModal}>×</button>
            </div>

            {registerError && (
              <div className="error-message">
                {registerError}
              </div>
            )}

            {registerSuccess && (
              <div className="success-message">
                {registerSuccess}
              </div>
            )}

            <form onSubmit={handleRegisterSubmit} className="register-form">
              <div className="form-group">
                <label htmlFor="register-email">Email:</label>
                <input
                  type="email"
                  id="register-email"
                  name="email"
                  value={registerForm.email}
                  onChange={handleRegisterChange}
                  required
                  disabled={registerLoading}
                  placeholder="user@example.com"
                />
              </div>

              <div className="form-group">
                <label htmlFor="register-password">Hasło:</label>
                <input
                  type="password"
                  id="register-password"
                  name="password"
                  value={registerForm.password}
                  onChange={handleRegisterChange}
                  required
                  disabled={registerLoading}
                  placeholder="Minimum 6 znaków"
                  minLength="6"
                />
              </div>

              <div className="form-group">
                <label htmlFor="register-role">Rola:</label>
                <select
                  id="register-role"
                  name="role"
                  value={registerForm.role}
                  onChange={handleRegisterChange}
                  disabled={registerLoading}
                >
                  <option value="employee">Pracownik</option>
                  <option value="admin">Administrator</option>
                  <option value="customer">Klient</option>
                </select>
              </div>

              <div className="modal-actions">
                <button 
                  type="button" 
                  className="cancel-btn"
                  onClick={closeModal}
                  disabled={registerLoading}
                >
                  Anuluj
                </button>
                <button 
                  type="submit" 
                  className="submit-btn"
                  disabled={registerLoading}
                >
                  {registerLoading ? 'Rejestracja...' : 'Zarejestruj użytkownika'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;