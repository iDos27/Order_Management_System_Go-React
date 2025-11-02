import { createContext, useContext, useState, useEffect } from 'react';
import api from '../services/api';

export const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Inicjalizacja stanu autoryzacji z localStorage
  useEffect(() => {
    const initializeAuth = () => {
      try {
        // Sprawdź oba możliwe klucze tokenów dla kompatybilności
        const storedToken = localStorage.getItem('auth_token') || localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');

        if (storedToken && storedUser) {
          const userData = JSON.parse(storedUser);
          setToken(storedToken);
          setUser(userData);
          setIsAuthenticated(true);
          
          // Upewnij się, że oba klucze są ustawione dla spójności
          localStorage.setItem('auth_token', storedToken);
          localStorage.setItem('token', storedToken);
          
          console.log('Autoryzacja przywrócona z localStorage');
        } else {
          console.log('Nie znaleziono zapisanej autoryzacji');
        }
      } catch (error) {
        console.error('Błąd inicjalizacji autoryzacji:', error);
        // Wyczyść uszkodzone dane
        localStorage.removeItem('auth_token');
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      } finally {
        setIsLoading(false);
      }
    };

    initializeAuth();
  }, []);

  const login = async (credentials) => {
    try {
      setIsLoading(true);
      console.log('Próba logowania...');
      
      const response = await api.login(credentials);
      console.log('Odpowiedź logowania:', response);

      if (response.token && response.user) {
        // Zapisz token pod oboma kluczami dla kompatybilności
        localStorage.setItem('auth_token', response.token);
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));

        // Aktualizuj stan
        setToken(response.token);
        setUser(response.user);
        setIsAuthenticated(true);

        console.log('Logowanie udane, token zapisany');
        return { success: true, user: response.user };
      } else {
        console.error('Nieprawidłowy format odpowiedzi logowania:', response);
        throw new Error('Nieprawidłowy format odpowiedzi');
      }
    } catch (error) {
      console.error('Logowanie nieudane:', error);
      
      // Wyczyść częściowy stan
      logout();
      
      return { 
        success: false, 
        error: error.message || 'Logowanie nieudane' 
      };
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    console.log('Wylogowywanie...');
    
    // Wyczyść localStorage
    localStorage.removeItem('auth_token');
    localStorage.removeItem('token');
    localStorage.removeItem('user');

    // Wyczyść stan
    setToken(null);
    setUser(null);
    setIsAuthenticated(false);
    
    console.log('Wylogowanie zakończone');
  };

  const register = async (userData) => {
    try {
      setIsLoading(true);
      console.log('Próba rejestracji...');
      
      const response = await api.register(userData);
      console.log('Odpowiedź rejestracji:', response);

      if (response.token && response.user) {
        // Automatyczne logowanie po udanej rejestracji
        localStorage.setItem('auth_token', response.token);
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));

        setToken(response.token);
        setUser(response.user);
        setIsAuthenticated(true);

        console.log('Rejestracja i automatyczne logowanie udane');
        return { success: true, user: response.user };
      } else {
        console.log('Rejestracja udana, ale bez automatycznego logowania');
        return { success: true, message: 'Rejestracja udana' };
      }
    } catch (error) {
      console.error('Rejestracja nieudana:', error);
      return { 
        success: false, 
        error: error.message || 'Rejestracja nieudana' 
      };
    } finally {
      setIsLoading(false);
    }
  };

  // Walidacja tokenu
  const validateToken = async () => {
    if (!token) return false;

    try {
      const response = await fetch('/api/v1/validate', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        console.log('Token jest prawidłowy');
        return true;
      } else {
        console.log('Token jest nieprawidłowy, wylogowywanie');
        logout();
        return false;
      }
    } catch (error) {
      console.error('Błąd walidacji tokenu:', error);
      logout();
      return false;
    }
  };

  // Funkcja pomocnicza do pobierania nagłówków autoryzacji dla wywołań API
  const getAuthHeaders = () => {
    if (!token) {
      throw new Error('Brak dostępnego tokenu autoryzacji');
    }
    
    return {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    };
  };

  const value = {
    // Stan
    user,
    token,
    isAuthenticated,
    isLoading,
    
    // Akcje
    login,
    logout,
    register,
    validateToken,
    getAuthHeaders,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;