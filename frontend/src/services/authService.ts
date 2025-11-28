import axios, { AxiosInstance, AxiosResponse, AxiosError, InternalAxiosRequestConfig } from "axios";
import {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User,
  UpdateProfileRequest,
  ChangePasswordRequest,
  AvatarUploadResponse,
  ApiResponse,
  UsersResponse,
  HealthResponse,
} from "../types";
import { isTokenExpired, isTokenExpiringSoon } from "../utils/jwtUtils";

// Configuración base de Axios
const API_BASE_URL =
  process.env.REACT_APP_API_URL || "http://localhost:8081/api/v1";

// Flag para evitar múltiples refreshes simultáneos
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: any) => void;
  reject: (reason?: any) => void;
}> = [];

const processQueue = (error: any = null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

// Crear instancia de axios con configuración base
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor para añadir token automáticamente y verificar expiración
api.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem("token");

    // Si no hay token, continuar sin autenticación
    if (!token) {
      return config;
    }

    // Verificar si el token ha expirado completamente
    if (isTokenExpired(token)) {
      localStorage.removeItem("token");
      window.location.href = "/login";
      return Promise.reject(new Error("Token expired"));
    }

    // Si el token está próximo a expirar y no es una petición de refresh
    if (isTokenExpiringSoon(token) && !config.url?.includes("/auth/refresh")) {
      try {
        // Intentar refrescar el token automáticamente
        const response = await axios.post<ApiResponse<{ token: string }>>(
          `${API_BASE_URL}/auth/refresh`,
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        const newToken = response.data.data.token;
        localStorage.setItem("token", newToken);

        // Usar el nuevo token para esta petición
        if (config.headers) {
          config.headers.Authorization = `Bearer ${newToken}`;
        }
      } catch (error) {
        console.error("Failed to refresh token:", error);
        // Continuar con el token actual si el refresh falla
      }
    }

    // Añadir el token a la petición
    if (config.headers) {
      const currentToken = localStorage.getItem("token");
      if (currentToken) {
        config.headers.Authorization = `Bearer ${currentToken}`;
      }
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
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // Si recibimos 401 y no es una petición de refresh o login
    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !originalRequest.url?.includes("/auth/login") &&
      !originalRequest.url?.includes("/auth/refresh")
    ) {
      if (isRefreshing) {
        // Si ya se está refrescando, agregar a la cola
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            return axios(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const token = localStorage.getItem("token");

      if (!token) {
        localStorage.removeItem("token");
        window.location.href = "/login";
        return Promise.reject(error);
      }

      try {
        // Intentar refrescar el token
        const response = await axios.post<ApiResponse<{ token: string }>>(
          `${API_BASE_URL}/auth/refresh`,
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        const newToken = response.data.data.token;
        localStorage.setItem("token", newToken);

        processQueue(null, newToken);

        // Reintentar la petición original con el nuevo token
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }

        return axios(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        localStorage.removeItem("token");
        window.location.href = "/login";
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
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
    const response = await api.put<ApiResponse<User>>('/profile', profileData);
    return response.data;
  },

  // Upload Avatar
  async uploadAvatar(file: File): Promise<AvatarUploadResponse> {
    const formData = new FormData();
    formData.append('avatar', file);
    
    const response = await api.post<AvatarUploadResponse>('/profile/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Delete Avatar
  async deleteAvatar(): Promise<ApiResponse<User>> {
    const response = await api.delete<ApiResponse<User>>('/profile/avatar');
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
