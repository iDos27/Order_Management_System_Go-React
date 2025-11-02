import axios from 'axios';

export const API_BASE_URL = 'http://localhost:8080/api';
export const AUTH_BASE_URL = 'http://localhost:8081/api/v1';

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
    updateOrderStatus: (id, status) => api.patch(`/orders/${id}/status`, { status }),
};

export default api;