import React, { useState, useEffect } from 'react';
import './Reports.css';

const Reports = () => {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  // Form state for generating reports
  const [reportForm, setReportForm] = useState({
    type: 'daily',
    selectedDate: '',      // dla daily
    selectedMonth: '',     // dla monthly
    period_start: '',      // auto-generated
    period_end: ''         // auto-generated
  });

  // Load quick stats on component mount
  useEffect(() => {
    loadQuickStats();
  }, []);

  const loadQuickStats = async (periodType = '30days', customEnd = null) => {
    try {
      setLoading(true);
      setError('');
      
      let url = '/api/reports/stats';
      let startDate, endDate;
      
      // Calculate dates based on period type
      if (periodType === 'all') {
        // Wszystkie zamówienia - bardzo szeroki zakres
        startDate = '2020-01-01';
        endDate = new Date().toISOString().split('T')[0];
      } else if (periodType === 'year') {
        // Aktualny rok
        const currentYear = new Date().getFullYear();
        startDate = `${currentYear}-01-01`;
        endDate = `${currentYear}-12-31`;
      } else if (periodType === '30days') {
        // Ostatnie 30 dni
        const end = new Date();
        const start = new Date();
        start.setDate(end.getDate() - 30);
        
        startDate = start.toISOString().split('T')[0];
        endDate = end.toISOString().split('T')[0];
      } else {
        // Legacy support - gdy podano konkretne daty
        startDate = periodType;
        endDate = customEnd;
      }
      
      url += `?start=${startDate}&end=${endDate}`;
      
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
      
      // Informacja o wybranym okresie
      let periodInfo = '';
      if (periodType === 'all') {
        periodInfo = 'wszystkich zamówień w historii';
      } else if (periodType === 'year') {
        periodInfo = `zamówień z ${new Date().getFullYear()} roku`;
      } else if (periodType === '30days') {
        periodInfo = 'zamówień z ostatnich 30 dni';
      } else {
        periodInfo = `zamówień z okresu ${startDate} - ${endDate}`;
      }
      
      // Jeśli brak danych, pokaż informację
      if (data.total_orders === 0) {
        setError(`Brak ${periodInfo}. Sprawdź czy w bazie są jakiekolwiek zamówienia.`);
      } else {
        // Pokaż info o okresie jako success
        setSuccess(`Znaleziono ${data.total_orders} ${periodInfo} na łączną kwotę ${data.total_amount.toFixed(2)} zł`);
        setTimeout(() => setSuccess(''), 3000); // Auto-hide po 3 sekundach
      }
      
    } catch (err) {
      console.error('Błąd statystyk:', err);
      setError(`Nie udało się pobrać statystyk: ${err.message}. Sprawdź czy raport-service działa na porcie 8083.`);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    
    setReportForm(prev => {
      const newForm = { ...prev, [name]: value };
      
      // Auto-generate period_start and period_end based on type and selection
      if (name === 'type' || name === 'selectedDate' || name === 'selectedMonth') {
        const { start, end } = calculateDateRange(
          name === 'type' ? value : newForm.type,
          name === 'selectedDate' ? value : newForm.selectedDate,
          name === 'selectedMonth' ? value : newForm.selectedMonth
        );
        
        newForm.period_start = start;
        newForm.period_end = end;
      }
      
      return newForm;
    });
  };

  // Calculate start and end dates based on type and selection
  const calculateDateRange = (type, selectedDate, selectedMonth) => {
    let start = '';
    let end = '';
    
    if (type === 'daily' && selectedDate) {
      start = selectedDate;
      end = selectedDate;
    } else if (type === 'monthly' && selectedMonth) {
      // selectedMonth format: "2025-10"
      const [year, month] = selectedMonth.split('-');
      start = `${year}-${month}-01`;
      
      // Last day of month
      const lastDay = new Date(year, month, 0).getDate();
      end = `${year}-${month}-${lastDay.toString().padStart(2, '0')}`;
    }
    
    return { start, end };
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

  // Set default values
  useEffect(() => {
    const today = new Date();
    const currentMonth = today.toISOString().slice(0, 7); // "2025-11"
    const currentDate = today.toISOString().split('T')[0];  // "2025-11-09"
    
    setReportForm(prev => ({
      ...prev,
      selectedDate: currentDate,
      selectedMonth: currentMonth, // "2025-11" będzie pasować do opcji "Listopad 2025"
      period_start: currentDate,
      period_end: currentDate
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
              onClick={() => loadQuickStats('all', null)}
            >
              📊 Wszystkie
            </button>
            <button 
              className="period-btn" 
              onClick={() => loadQuickStats('year', null)}
            >
              📅 W tym roku
            </button>
            <button 
              className="period-btn" 
              onClick={() => loadQuickStats('30days', null)}
            >
              📈 Ostatnie 30 dni
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
                <option value="monthly">Miesięczny</option>
              </select>
            </div>

            {/* Conditional inputs based on report type */}
            {reportForm.type === 'daily' && (
              <div className="form-group">
                <label htmlFor="selected-date">Wybierz dzień:</label>
                <input
                  type="date"
                  id="selected-date"
                  name="selectedDate"
                  value={reportForm.selectedDate}
                  onChange={handleInputChange}
                  required
                  disabled={loading}
                  className="clean-date-input"
                />
              </div>
            )}

            {reportForm.type === 'monthly' && (
              <div className="form-group">
                <label htmlFor="selected-month">Wybierz miesiąc:</label>
                <input
                  type="month"
                  id="selected-month"
                  name="selectedMonth"
                  value={reportForm.selectedMonth}
                  onChange={handleInputChange}
                  required
                  disabled={loading}
                  className="clean-date-input"
                />
              </div>
            )}

            {/* Show calculated period for clarity */}
            {reportForm.period_start && reportForm.period_end && (
              <div className="calculated-period">
                <small>
                  📅 Okres raportu: <strong>
                    {reportForm.period_start === reportForm.period_end 
                      ? reportForm.period_start
                      : `${reportForm.period_start} - ${reportForm.period_end}`
                    }
                  </strong>
                </small>
              </div>
            )}

            <button 
              type="submit" 
              className="generate-btn"
              disabled={loading || !reportForm.period_start}
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