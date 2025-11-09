import React, { useState, useEffect } from 'react';
import './Reports.css';

const Reports = () => {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  // Form state for generating reports
  const [reportForm, setReportForm] = useState({
    type: 'weekly',
    period_start: '',
    period_end: ''
  });

  // Load quick stats on component mount
  useEffect(() => {
    loadQuickStats();
  }, []);

  const loadQuickStats = async (customStart = null, customEnd = null) => {
    try {
      setLoading(true);
      setError('');
      
      // Jeśli nie podano dat, użyj ostatnich 30 dni (więcej niż 7)
      let url = '/api/reports/stats';
      if (customStart && customEnd) {
        url += `?start=${customStart}&end=${customEnd}`;
      } else {
        // Ostatnie 30 dni jako fallback
        const endDate = new Date();
        const startDate = new Date();
        startDate.setDate(endDate.getDate() - 30);
        
        url += `?start=${startDate.toISOString().split('T')[0]}&end=${endDate.toISOString().split('T')[0]}`;
      }
      
      const response = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token') || localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      
      const data = await response.json();
      setStats(data);
      
      // Jeśli brak danych, pokaż informację
      if (data.total_orders === 0) {
        setError('Brak zamówień w wybranym okresie. Spróbuj wybrać starsze daty.');
      }
      
    } catch (err) {
      console.error('Błąd statystyk:', err);
      setError(`Nie udało się pobrać statystyk: ${err.message}. Sprawdź czy raport-service działa na porcie 8083.`);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    setReportForm({
      ...reportForm,
      [e.target.name]: e.target.value
    });
  };

  const generateReport = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const response = await fetch('/api/reports/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('auth_token') || localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          type: reportForm.type,
          period_start: new Date(reportForm.period_start).toISOString(),
          period_end: new Date(reportForm.period_end).toISOString()
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Błąd generowania raportu');
      }

      const report = await response.json();
      setSuccess(`Raport wygenerowany pomyślnie! ID: ${report.id}, Zamówienia: ${report.total_orders}, Kwota: ${report.total_amount.toFixed(2)} zł`);
      
      // Refresh stats after generating report
      loadQuickStats();
      
    } catch (err) {
      setError('Nie udało się wygenerować raportu: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  // Set default dates (last 7 days)
  useEffect(() => {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(endDate.getDate() - 7);
    
    setReportForm(prev => ({
      ...prev,
      period_start: startDate.toISOString().split('T')[0],
      period_end: endDate.toISOString().split('T')[0]
    }));
  }, []);

  return (
    <div className="reports-container">
      <div className="reports-header">
        <h2>📊 Raporty sprzedażowe</h2>
        <p>Generowanie i analiza raportów zamówień</p>
      </div>

      {error && (
        <div className="error-message">
          <span className="error-icon">⚠️</span>
          {error}
        </div>
      )}

      {success && (
        <div className="success-message">
          <span className="success-icon">✅</span>
          {success}
        </div>
      )}

      <div className="reports-grid">
        {/* Quick Stats Card */}
        <div className="stats-card">
          <h3>📈 Statystyki zamówień</h3>
          
          {/* Quick period buttons */}
          <div className="period-buttons">
            <button 
              className="period-btn" 
              onClick={() => loadQuickStats('2025-10-01', '2025-10-31')}
            >
              Październik 2025
            </button>
            <button 
              className="period-btn" 
              onClick={() => loadQuickStats('2025-10-20', '2025-10-27')}
            >
              20-27.10.2025
            </button>
            <button 
              className="period-btn" 
              onClick={() => loadQuickStats()}
            >
              Ostatnie 30 dni
            </button>
          </div>
          {loading && !stats ? (
            <div className="loading">Ładowanie...</div>
          ) : stats ? (
            <div className="stats-content">
              <div className="stat-item">
                <span className="stat-label">Łączne zamówienia:</span>
                <span className="stat-value">{stats.total_orders}</span>
              </div>
              <div className="stat-item">
                <span className="stat-label">Łączna kwota:</span>
                <span className="stat-value">{stats.total_amount.toFixed(2)} zł</span>
              </div>
              
              {stats.sources && stats.sources.length > 0 && (
                <div className="sources-breakdown">
                  <h4>Według źródeł:</h4>
                  {stats.sources.map((source, index) => (
                    <div key={index} className="source-item">
                      <span className="source-name">{source.source_name}:</span>
                      <span className="source-stats">
                        {source.count} zamówień ({source.amount.toFixed(2)} zł)
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ) : (
            <div className="no-data">Brak danych do wyświetlenia</div>
          )}
          
          <button 
            className="refresh-btn" 
            onClick={() => loadQuickStats()}
            disabled={loading}
          >
            🔄 Odśwież
          </button>
        </div>

        {/* Generate Report Card */}
        <div className="generate-card">
          <h3>📋 Generowanie raportu</h3>
          
          <form onSubmit={generateReport} className="report-form">
            <div className="form-group">
              <label htmlFor="report-type">Typ raportu:</label>
              <select
                id="report-type"
                name="type"
                value={reportForm.type}
                onChange={handleInputChange}
                disabled={loading}
              >
                <option value="daily">Dzienny</option>
                <option value="weekly">Tygodniowy</option>
                <option value="monthly">Miesięczny</option>
              </select>
            </div>

            <div className="form-group">
              <label htmlFor="period-start">Data początkowa:</label>
              <input
                type="date"
                id="period-start"
                name="period_start"
                value={reportForm.period_start}
                onChange={handleInputChange}
                required
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <label htmlFor="period-end">Data końcowa:</label>
              <input
                type="date"
                id="period-end"
                name="period_end"
                value={reportForm.period_end}
                onChange={handleInputChange}
                required
                disabled={loading}
              />
            </div>

            <button 
              type="submit" 
              className="generate-btn"
              disabled={loading}
            >
              {loading ? 'Generowanie...' : '📄 Generuj raport Excel'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Reports;