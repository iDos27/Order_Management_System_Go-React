import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const ordersAPI = {
    getAllOrders: () => api.get('/orders'),
    getOrderById: (id) => api.get(`/orders/${id}`),
    createOrder: (orderData) => api.post('/orders', orderData),
    updateOrder: (id, orderData) => api.put(`/orders/${id}`, orderData),
};

export default api;