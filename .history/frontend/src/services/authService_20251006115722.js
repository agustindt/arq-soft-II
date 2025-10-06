import axios from 'axios';

// Configuración base de Axios
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api/v1';

// Crear instancia de axios con configuración base
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para añadir token automáticamente
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor para manejar respuestas y errores
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Si recibimos 401, eliminar token y redireccionar
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Servicio de autenticación
export const authService = {
  // Login
  async login(email, password) {
    const response = await api.post('/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  // Register
  async register(userData) {
    const response = await api.post('/auth/register', {
      email: userData.email,
      username: userData.username,
      password: userData.password,
      first_name: userData.firstName,
      last_name: userData.lastName,
    });
    return response.data;
  },

  // Refresh Token
  async refreshToken() {
    const response = await api.post('/auth/refresh');
    return response.data;
  },

  // Get Profile
  async getProfile() {
    const response = await api.get('/profile');
    return response.data.data;
  },

  // Update Profile
  async updateProfile(profileData) {
    const response = await api.put('/profile', {
      first_name: profileData.firstName,
      last_name: profileData.lastName,
    });
    return response.data;
  },

  // Change Password
  async changePassword(passwordData) {
    const response = await api.put('/profile/password', {
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword,
    });
    return response.data;
  },

  // Health Check
  async healthCheck() {
    const response = await api.get('/health');
    return response.data;
  },
};

// Servicio de usuarios (endpoints públicos)
export const userService = {
  // Get all users
  async getUsers(page = 1, limit = 10) {
    const response = await api.get(`/users?page=${page}&limit=${limit}`);
    return response.data;
  },

  // Get user by ID
  async getUserById(id) {
    const response = await api.get(`/users/${id}`);
    return response.data;
  },
};

export default api;
