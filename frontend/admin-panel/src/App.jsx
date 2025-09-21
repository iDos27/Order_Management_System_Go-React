import { useState, useEffect } from "react"
import { ordersAPI } from "./services/api"
import './App.css'

function App() {
  const [orders, setOrders] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

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

  if (loading) return <div>Ładowanie zamówień...</div>
  if (error) return <div style={{color: 'red'}}>{error}</div>

  return (
    <div className="admin-panel">
      <h1>Admin Panel - Zarządzanie Zamówieniami</h1>
      
      <div className="orders-section">
        <h2>Lista Zamówień ({orders.length})</h2>
        
        {orders.length === 0 ? (
          <p>Brak zamówień</p>
        ) : (
          <div className="orders-list">
            {orders.map(order => (
              <div key={order.id} className="order-card">
                <h3>Zamówienie #{order.id}</h3>
                <p><strong>Klient:</strong> {order.customer_name}</p>
                <p><strong>Email:</strong> {order.customer_email}</p>
                <p><strong>Status:</strong> {order.status}</p>
                <p><strong>Kwota:</strong> {order.total_amount} zł</p>
                <p><strong>Data:</strong> {new Date(order.created_at).toLocaleString()}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default App