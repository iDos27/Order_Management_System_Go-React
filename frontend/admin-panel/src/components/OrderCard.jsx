import './OrderCard.css';

function OrderCard({ order, onClick }) {
    const getSourceColor = (source) => {
        switch (source) {
            case 'źródło_jeden': return '#e48650ff';
            case 'źródło_dwa': return '#d36b8aff';
            case 'website': return '#2196f3';
            case 'manual': return '#4caf50';
            default: return '#757575'
        }
    }

    const getSourceText = (source) => {
        switch (source) {
            case 'źródło_jeden': return 'Źródło 1';
            case 'źródło_dwa': return 'Źródło 2';
            case 'website': return 'Strona WWW';
            case 'manual': return 'Ręczne';
            default: return source ? source.toUpperCase() : 'NIEZNANE';
        }
    }

    return (
        <div className="order-card" onClick={() => onClick(order)}>
          <div className="order-header">
            <h3>Zamówienie #{order.id}</h3>
            <div 
              className="source-badge"
              style={{ backgroundColor: getSourceColor(order.source) }}
            >
              {getSourceText(order.source)}
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