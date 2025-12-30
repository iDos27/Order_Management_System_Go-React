import { useState } from "react"
import { useAuth } from "../context/AuthContext"
import "./NewOrderForm.css"

const NewOrderForm = ({ onClose }) => {
  const { getAuthHeaders } = useAuth()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  
  const [formData, setFormData] = useState({
    customer_name: "",
    customer_email: "",
    total_amount: "",
    source: "manual"
  })

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/orders', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          ...formData,
          total_amount: parseFloat(formData.total_amount)
        })
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }

      const newOrder = await response.json()
      console.log('Nowe zamówienie utworzone:', newOrder)
      onClose()
    } catch (err) {
      setError('Błąd podczas tworzenia zamówienia: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>+ Nowe zamówienie</h2>
          <button className="close-button" onClick={onClose}>✕</button>
        </div>

        <form onSubmit={handleSubmit} className="order-form">
          {error && <div className="error-message">{error}</div>}
          
          <div className="form-group">
            <label htmlFor="customer_name">Imię i nazwisko</label>
            <input
              type="text"
              id="customer_name"
              name="customer_name"
              value={formData.customer_name}
              onChange={handleChange}
              required
              placeholder="Jan Kowalski"
            />
          </div>

          <div className="form-group">
            <label htmlFor="customer_email">Email</label>
            <input
              type="email"
              id="customer_email"
              name="customer_email"
              value={formData.customer_email}
              onChange={handleChange}
              required
              placeholder="jan@example.com"
            />
          </div>

          <div className="form-group">
            <label htmlFor="total_amount">Kwota (zł)</label>
            <input
              type="number"
              id="total_amount"
              name="total_amount"
              value={formData.total_amount}
              onChange={handleChange}
              required
              step="0.01"
              min="0"
              placeholder="199.99"
            />
          </div>

          <div className="form-group">
            <label htmlFor="source">Źródło</label>
            <input
              type="text"
              id="source"
              name="source"
              value="Ręczne"
              disabled
              style={{ 
                backgroundColor: '#f5f5f5', 
                cursor: 'not-allowed',
                color: '#666'
              }}
            />
          </div>

          <div className="form-actions">
            <button type="button" onClick={onClose} className="btn-secondary">
              Anuluj
            </button>
            <button type="submit" disabled={loading} className="btn-primary">
              {loading ? 'Tworzenie...' : 'Utwórz zamówienie'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default NewOrderForm
