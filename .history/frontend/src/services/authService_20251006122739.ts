import axios, { AxiosInstance, AxiosResponse } from "axios";
import {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User,
  UpdateProfileRequest,
  ChangePasswordRequest,
  ApiResponse,
  UsersResponse,
  HealthResponse,
} from "../types";

// Configuración base de Axios
const API_BASE_URL =
  process.env.REACT_APP_API_URL || "http://localhost:8081/api/v1";

// Crear instancia de axios con configuración base
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor para añadir token automáticamente
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token && config.headers) {
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
  (response: AxiosResponse) => {
    return response;
  },
  (error) => {
    // Si recibimos 401, eliminar token y redireccionar
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

// Servicio de autenticación
export const authService = {
  // Login
  async login(
    email: string,
    password: string
  ): Promise<ApiResponse<LoginResponse>> {
    const response = await api.post<ApiResponse<LoginResponse>>("/auth/login", {
      email,
      password,
    });
    return response.data;
  },

  // Register
  async register(
    userData: RegisterRequest
  ): Promise<ApiResponse<LoginResponse>> {
    const response = await api.post<ApiResponse<LoginResponse>>(
      "/auth/register",
      {
        email: userData.email,
        username: userData.username,
        password: userData.password,
        first_name: userData.firstName,
        last_name: userData.lastName,
      }
    );
    return response.data;
  },

  // Refresh Token
  async refreshToken(): Promise<ApiResponse<{ token: string }>> {
    const response = await api.post<ApiResponse<{ token: string }>>(
      "/auth/refresh"
    );
    return response.data;
  },

  // Get Profile
  async getProfile(): Promise<User> {
    const response = await api.get<ApiResponse<User>>("/profile");
    return response.data.data;
  },

  // Update Profile
  async updateProfile(
    profileData: UpdateProfileRequest
  ): Promise<ApiResponse<User>> {
    const response = await api.put<ApiResponse<User>>("/profile", {
      first_name: profileData.firstName,
      last_name: profileData.lastName,
    });
    return response.data;
  },

  // Change Password
  async changePassword(
    passwordData: ChangePasswordRequest
  ): Promise<ApiResponse<null>> {
    const response = await api.put<ApiResponse<null>>("/profile/password", {
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword,
    });
    return response.data;
  },

  // Health Check
  async healthCheck(): Promise<HealthResponse> {
    const response = await api.get<HealthResponse>("/health");
    return response.data;
  },
};

// Servicio de usuarios (endpoints públicos)
export const userService = {
  // Get all users
  async getUsers(
    page: number = 1,
    limit: number = 10
  ): Promise<ApiResponse<UsersResponse>> {
    const response = await api.get<ApiResponse<UsersResponse>>(
      `/users?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  // Get user by ID
  async getUserById(id: number): Promise<ApiResponse<User>> {
    const response = await api.get<ApiResponse<User>>(`/users/${id}`);
    return response.data;
  },
};

export default api;
