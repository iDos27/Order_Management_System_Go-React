import axios from 'axios';

// Ścieżki dla NGINX proxy
export const API_BASE_URL = '/api/orders';
export const AUTH_BASE_URL = '/api/auth';

const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const ordersAPI = {
    getAllOrders: () => axiosInstance.get(''),
    getOrderById: (id) => axiosInstance.get(`/${id}`),
    createOrder: (orderData) => axiosInstance.post('', orderData),
    updateOrder: (id, orderData) => axiosInstance.put(`/${id}`, orderData),
    updateOrderStatus: (id, status) => axiosInstance.patch(`/${id}/status`, { status }),
};

// API object z funkcjami dla AuthContext
const api = {
    login: async (credentials) => {
        try {
            const response = await axios.post(`${AUTH_BASE_URL}/login`, credentials, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            return response.data;
        } catch (error) {
            if (error.response) {
                throw new Error(error.response.data.error || `HTTP ${error.response.status}`);
            }
            throw new Error('Błąd połączenia z serwerem');
        }
    },

    register: async (userData) => {
        try {
            const response = await axios.post(`${AUTH_BASE_URL}/register`, userData, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            return response.data;
        } catch (error) {
            if (error.response) {
                throw new Error(error.response.data.error || `HTTP ${error.response.status}`);
            }
            throw new Error('Błąd połączenia z serwerem');
        }
    },

    getOrders: async (token) => {
        try {
            const response = await axios.get(`${API_BASE_URL}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });
            return response.data;
        } catch (error) {
            if (error.response) {
                throw new Error(error.response.data.error || `HTTP ${error.response.status}`);
            }
            throw new Error('Błąd połączenia z serwerem');
        }
    },

    createOrder: async (orderData, token) => {
        try {
            const response = await axios.post(`${API_BASE_URL}`, orderData, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });
            return response.data;
        } catch (error) {
            if (error.response) {
                throw new Error(error.response.data.error || `HTTP ${error.response.status}`);
            }
            throw new Error('Błąd połączenia z serwerem');
        }
    },

    updateOrderStatus: async (orderId, status, token) => {
        try {
            const response = await axios.patch(`${API_BASE_URL}/${orderId}/status`, 
                { status }, 
                {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                }
            );
            return response.data;
        } catch (error) {
            if (error.response) {
                throw new Error(error.response.data.error || `HTTP ${error.response.status}`);
            }
            throw new Error('Błąd połączenia z serwerem');
        }
    }
};

export default api;