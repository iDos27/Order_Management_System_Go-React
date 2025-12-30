import { useState, useEffect } from "react"
import { useAuth } from "../context/AuthContext"
import OrderCard from "./OrderCard"
import useWebSocket from "../hooks/useWebSocket"
import NewOrderForm from "./NewOrderForm"

const OrderManagement = () => {
  const [orders, setOrders] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [currentView, setCurrentView] = useState('list')
  const [selectedOrder, setSelectedOrder] = useState(null)
  const [showNewOrderForm, setShowNewOrderForm] = useState(false)

  // Pobierz token i helper z AuthContext
  const { token, isAuthenticated, getAuthHeaders } = useAuth()

  const { lastMessage } = useWebSocket(`ws://${window.location.host}/ws`)

  useEffect(() => { 
    fetchOrders()
  }, [])

  useEffect(() => {
    if (lastMessage && lastMessage.type === 'order_update') {
      const { order_id, new_status } = lastMessage.payload;
      
      // Sprawdzamy czy zamówienie już istnieje
      setOrders(prevOrders => {
        const orderExists = prevOrders.find(order => order.id === order_id);
        
        if (orderExists) {
          // Aktualizujemy istniejące zamówienie
          return prevOrders.map(order => 
            order.id === order_id 
              ? { ...order, status: new_status }
              : order
          );
        } else {
          // Nowe zamówienie - pobieramy z API w następnym cyklu
          console.log(`Nowe zamówienie #${order_id}, pobieram z API...`);
          setTimeout(() => fetchOrders(), 100);
          return prevOrders;
        }
      });
      
      // Aktualizujemy selectedOrder jeśli to to samo zamówienie
      setSelectedOrder(prev => {
        if (prev && prev.id === order_id) {
          return { ...prev, status: new_status };
        }
        return prev;
      });
    }
  }, [lastMessage]);

  const fetchOrders = async () => {
    try {
      setLoading(true)
      
      // Sprawdź autoryzację
      if (!isAuthenticated || !token) {
        setError('Wymagane logowanie')
        return
      }

      // Wykonaj request z tokenem z AuthContext
      const response = await fetch('/api/orders', {
        headers: getAuthHeaders(),
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }

      const ordersData = await response.json()
      setOrders(ordersData)
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

  const handleStatusChange = async (orderId, newStatus) => {
    try {
      // Sprawdź autoryzację
      if (!isAuthenticated || !token) {
        alert('Wymagane logowanie')
        return
      }

      // Wykonaj request z tokenem z AuthContext
      const response = await fetch(`/api/orders/${orderId}/status`, {
        method: 'PATCH',
        headers: getAuthHeaders(),
        body: JSON.stringify({ status: newStatus }),
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }

      console.log(`Status zamówienia ${orderId} zmieniony na ${newStatus}`)
    } catch (error) {
      console.error('Błąd podczas zmiany statusu:', error)
      alert('Błąd podczas zmiany statusu zamówienia')
    }
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
      <div className="order-management">
        <div className="management-header">
          <h2>Zarządzanie Zamówieniami</h2>
          <button 
            className="add-order-button" 
            onClick={() => setShowNewOrderForm(true)}
          >
            + Dodaj zamówienie
          </button>
        </div>

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
    <div className="order-management">
      <div className="details-header">
        <button className="back-button" onClick={handleBackToList}>
          ← Powrót do listy
        </button>
        <h2>Zamówienie #{selectedOrder.id}</h2>
      </div>
      
      <div className="order-details-container">
        <div className="status-card">
          <h3>Status zamówienia</h3>
          <div className="status-controls">
            <div 
              className="current-status-badge"
              style={{ backgroundColor: getStatusColor(selectedOrder.status) }}
            >
              {getStatusText(selectedOrder.status)}
            </div>

            <div className="action-buttons">
              {getAvailableActions(selectedOrder.status).map(action => (
                <button
                  key={action.status}
                  className={`action-button ${action.type}`}
                  onClick={() => handleStatusChange(selectedOrder.id, action.status)}
                >
                  {action.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      
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

  const getAvailableActions = (currentStatus) => {
    switch(currentStatus) {
      case 'new':
        return [
          { status: 'confirmed', label: '✅ Potwierdź', type: 'primary' },
          { status: 'cancelled', label: '❌ Anuluj', type: 'danger' }
        ]
      case 'confirmed':
        return [
          { status: 'shipped', label: '🚚 Wyślij', type: 'primary' },
          { status: 'new', label: '⬅️ Cofnij', type: 'secondary' },
          { status: 'cancelled', label: '❌ Anuluj', type: 'danger' }
        ]
      case 'shipped':
        return [
          { status: 'delivered', label: '📦 Dostarczone', type: 'primary' },
          { status: 'confirmed', label: '⬅️ Cofnij', type: 'secondary' },
          { status: 'cancelled', label: '❌ Anuluj', type: 'danger' }
        ]
      case 'delivered':
        return [
          { status: 'shipped', label: '⬅️ Cofnij do wysłane', type: 'secondary' }
        ]
      case 'cancelled':
        return [
          { status: 'new', label: '🔄 Przywróć', type: 'secondary' }
        ]
      default:
        return []
    }
  }

  if (loading) return <div className="loading">Ładowanie zamówień...</div>
  if (error) return <div className="error">{error}</div>

  return (
    <>
      {currentView === 'list' ? renderOrdersList() : renderOrderDetails()}
      {showNewOrderForm && (
        <NewOrderForm
          onClose={() => setShowNewOrderForm(false)}
        />
      )}
    </>
  )
}

export default OrderManagement