import './OrderCard.css';

function OrderCard({ order, onClick }) {
    const getStatusColor = (status) => {
        switch (status) {
            case 'new': return '#ff9800';
            case 'confirmed': return '#2196f3'
            case 'shipped': return '#4caf50'
            case 'delivered': return '#8bc34a'
            case 'canceled': return '#f44336'
            default: return '#757575'
        }
    }

    const getStatusText = (status) => {
        switch (status) {
            case 'new': return 'NOWE'
            case 'confirmed': return 'POTWIERDZONE'
            case 'shipped': return 'WYSŁANE'
            case 'delivered': return 'DOSTARCZONE';
            case 'canceled': return 'ANULOWANE'
            default: return status.toUpperCase()
        }
    }

    return (
        <div className="order-card" onClick={() => onClick(order)}>
          <div className="order-header">
            <h3>Zamówienie #{order.id}</h3>
            <div 
              className="status-badge"
              style={{ backgroundColor: getStatusColor(order.status) }}
            >
              {getStatusText(order.status)}
            </div>
          </div>
          <div className="order-info">
            <p className="customer-name">{order.customer_name}</p>
            <p className="order-amount">{order.total_amount} zł</p>
          </div>
        </div>
    )
}

export default OrderCard;