import { useState, useEffect } from "react"
import { ordersAPI } from "./services/api"
import OrderCard from "./components/OrderCard"
import './App.css'

function App() {
  const [orders, setOrders] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [currentView, setCurrentView] = useState('list') // 'list' lub 'details'
  const [selectedOrder, setSelectedOrder] = useState(null)

  useEffect(() => { 
    fetchOrders()
  }, [])

  const fetchOrders = async () => {
    try {
      setLoading(true)
      const response = await ordersAPI.getAllOrders()
      setOrders(response.data)
      setError(null)
    } catch (err) {
      setError('Błąd podczas pobierania zamówień: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleOrderClick = (order) => {
    setSelectedOrder(order)
    setCurrentView('details')
  }

  const handleBackToList = () => {
    setCurrentView('list')
    setSelectedOrder(null)
  }

  const renderOrdersList = () => {
    // Grupowanie zamówień według statusu
    const groupedOrders = {
      'new': orders.filter(order => order.status === 'new'),
      'confirmed': orders.filter(order => order.status === 'confirmed'),
      'shipped': orders.filter(order => order.status === 'shipped'),
      'delivered': orders.filter(order => order.status === 'delivered'),
      'cancelled': orders.filter(order => order.status === 'cancelled')
    }

    const statusConfig = {
      'new': { title: 'NOWE', color: '#ff9800', count: groupedOrders.new.length },
      'confirmed': { title: 'POTWIERDZONE', color: '#2196f3', count: groupedOrders.confirmed.length },
      'shipped': { title: 'WYSŁANE', color: '#4caf50', count: groupedOrders.shipped.length },
      'delivered': { title: 'DOSTARCZONE', color: '#8bc34a', count: groupedOrders.delivered.length },
      'cancelled': { title: 'ANULOWANE', color: '#f44336', count: groupedOrders.cancelled.length }
    }

    return (
      <div className="admin-panel">
        <h1>Admin Panel - Zarządzanie Zamówieniami</h1>

        <div className="kanban-board">
          {Object.entries(statusConfig).map(([status, config]) => (
            <div key={status} className="kanban-column">
              <div 
                className="column-header"
                style={{ backgroundColor: config.color }}
              >
                <h3>{config.title}</h3>
                <span className="count-badge">{config.count}</span>
              </div>

              <div className="column-content">
                {groupedOrders[status].length === 0 ? (
                  <div className="empty-column">
                    <p>Brak zamówień</p>
                  </div>
                ) : (
                  groupedOrders[status].map(order => (
                    <OrderCard 
                      key={order.id} 
                      order={order} 
                      onClick={handleOrderClick}
                    />
                  ))
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

const renderOrderDetails = () => (
  <div className="admin-panel">
    <div className="details-header">
      <button className="back-button" onClick={handleBackToList}>
        ← Powrót do listy
      </button>
      <h1>Zamówienie #{selectedOrder.id}</h1>
    </div>
    
    <div className="order-details-container">
      {/* Status Card z możliwością zmiany */}
      <div className="status-card">
        <h3>Status zamówienia</h3>
        <div className="status-display">
          <div 
            className="current-status-badge"
            style={{ backgroundColor: getStatusColor(selectedOrder.status) }}
          >
            {getStatusText(selectedOrder.status)}
          </div>
          {/* TODO: Tutaj będzie dropdown do zmiany statusu */}
        </div>
      </div>

      {/* Info Cards */}
      <div className="details-grid">
        <div className="details-card">
          <div className="card-icon">👤</div>
          <h3>Informacje o kliencie</h3>
          <div className="info-row">
            <span className="label">Imię i nazwisko:</span>
            <span className="value">{selectedOrder.customer_name}</span>
          </div>
          <div className="info-row">
            <span className="label">Email:</span>
            <span className="value">{selectedOrder.customer_email}</span>
          </div>
        </div>
        
        <div className="details-card">
          <div className="card-icon">📦</div>
          <h3>Szczegóły zamówienia</h3>
          <div className="info-row">
            <span className="label">Numer zamówienia:</span>
            <span className="value">#{selectedOrder.id}</span>
          </div>
          <div className="info-row">
            <span className="label">Kwota:</span>
            <span className="value amount">{selectedOrder.total_amount} zł</span>
          </div>
          <div className="info-row">
            <span className="label">Data utworzenia:</span>
            <span className="value">{new Date(selectedOrder.created_at).toLocaleString()}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
)

// Helper functions (dodaj je na początku komponentu)
  const getStatusColor = (status) => {
    switch(status) {
      case 'new': return '#ff9800'
      case 'confirmed': return '#2196f3'  
      case 'shipped': return '#4caf50'
      case 'delivered': return '#8bc34a'
      case 'cancelled': return '#f44336'
      default: return '#757575'
    }
  }
  
  const getStatusText = (status) => {
    switch(status) {
      case 'new': return 'NOWE'
      case 'confirmed': return 'POTWIERDZONE'
      case 'shipped': return 'WYSŁANE'
      case 'delivered': return 'DOSTARCZONE'
      case 'cancelled': return 'ANULOWANE'
      default: return status.toUpperCase()
    }
  }

  if (loading) return <div className="loading">Ładowanie zamówień...</div>
  if (error) return <div className="error">{error}</div>

  return currentView === 'list' ? renderOrdersList() : renderOrderDetails()
}

export default App