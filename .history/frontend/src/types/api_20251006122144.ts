import { User } from './auth';

// Pagination Types
export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

// Users API Types
export interface UsersResponse {
  users: User[];
  pagination: Pagination;
}

// Health Check Response
export interface HealthResponse {
  status: string;
  message: string;
  service?: string;
}

// Dashboard Stats
export interface DashboardStats {
  totalUsers: number;
  recentUsers: User[];
}

// API Status
export interface ApiStatus {
  status: 'online' | 'error';
  message: string;
}

// Form Errors
export interface FormErrors {
  [key: string]: string;
}

// Message Type (for alerts)
export interface Message {
  type: 'success' | 'error' | 'info' | 'warning' | '';
  text: string;
}
