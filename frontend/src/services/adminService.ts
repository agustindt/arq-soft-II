import api from "./authService";
import { ApiResponse, User } from "../types";

/**
 * Admin Service
 * Servicios de administración para gestionar usuarios y estadísticas del sistema
 */

export interface AdminUser {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AdminUsersListResponse {
  users: AdminUser[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface SystemStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  users_by_role: {
    user: number;
    admin: number;
    root: number;
  };
  recent_registrations: number; // últimos 7 días
}

export interface CreateUserRequest {
  email: string;
  username: string;
  password: string;
  first_name: string;
  last_name: string;
  role: string;
}

export interface UpdateRoleRequest {
  role: string;
}

export interface UpdateStatusRequest {
  is_active: boolean;
}

export const adminService = {
  /**
   * Obtiene lista completa de usuarios (admin only)
   */
  async getAllUsers(
    page: number = 1,
    limit: number = 20
  ): Promise<ApiResponse<AdminUsersListResponse>> {
    const response = await api.get<ApiResponse<AdminUsersListResponse>>(
      `/admin/users?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  /**
   * Crea un nuevo usuario (admin only)
   */
  async createUser(
    userData: CreateUserRequest
  ): Promise<ApiResponse<AdminUser>> {
    const response = await api.post<ApiResponse<AdminUser>>(
      "/admin/users",
      userData
    );
    return response.data;
  },

  /**
   * Actualiza el rol de un usuario (admin only)
   */
  async updateUserRole(
    userId: number,
    role: string
  ): Promise<ApiResponse<AdminUser>> {
    const response = await api.put<ApiResponse<AdminUser>>(
      `/admin/users/${userId}/role`,
      { role }
    );
    return response.data;
  },

  /**
   * Actualiza el estado (activo/inactivo) de un usuario (admin only)
   */
  async updateUserStatus(
    userId: number,
    isActive: boolean
  ): Promise<ApiResponse<AdminUser>> {
    const response = await api.put<ApiResponse<AdminUser>>(
      `/admin/users/${userId}/status`,
      { is_active: isActive }
    );
    return response.data;
  },

  /**
   * Elimina un usuario (root only)
   */
  async deleteUser(userId: number): Promise<ApiResponse<null>> {
    const response = await api.delete<ApiResponse<null>>(
      `/admin/users/${userId}`
    );
    return response.data;
  },

  /**
   * Obtiene estadísticas del sistema (admin only)
   */
  async getSystemStats(): Promise<ApiResponse<SystemStats>> {
    const response = await api.get<ApiResponse<SystemStats>>("/admin/stats");
    return response.data;
  },
};

export default adminService;
