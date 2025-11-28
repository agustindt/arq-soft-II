import { User } from "./auth";

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
  status: "online" | "error";
  message: string;
}

// Form Errors
export interface FormErrors {
  [key: string]: string;
}

// Message Type (for alerts)
export interface Message {
  type: "success" | "error" | "info" | "warning" | "";
  text: string;
}

// Activity Types
export interface Activity {
  id: string;
  name: string;
  description: string;
  category: string;
  difficulty: "beginner" | "intermediate" | "advanced";
  location: string;
  price: number;
  duration: number; // in minutes
  max_capacity: number;
  instructor?: string;
  schedule?: string[];
  equipment?: string[];
  image_url?: string;
  is_active: boolean;
  created_by?: number;
  created_at: string;
  updated_at: string;
}

export interface ActivitiesResponse {
  activities: Activity[];
  count: number;
}

export interface ActivityResponse {
  activity: Activity;
}

// Reservation Types
export interface Reservation {
  id: string;
  users_id: number[];
  cupo: number;
  actividad: string; // Activity ID
  schedule: string; // Specific schedule slot (e.g., "Lunes 20:00")
  date: string; // ISO date string
  status: "Pendiente" | "confirmada" | "cancelada";
  created_at: string;
  updated_at: string;
}

export interface ReservationsResponse {
  reservas: Reservation[];
  count: number;
}

export interface ReservationResponse {
  Reserva: Reservation;
}

export interface CreateReservationRequest {
  activityId: string;
  schedule: string; // Specific schedule slot (e.g., "Lunes 20:00")
  date: string; // ISO date string
  participants: number;
}

export interface ScheduleAvailability {
  activity_id: string;
  date: string;
  availability: {
    [schedule: string]: number; // "Lunes 20:00" -> 7 (available spots)
  };
}

// Search Types
export interface SearchFilters {
  query?: string;
  category?: string;
  difficulty?: string;
  price_min?: number;
  price_max?: number;
  page?: number;
  limit?: number;
  sort?: string;
}

export interface SearchResult {
  query: string;
  results: Activity[];
  total_found: number;
  page?: number;
  limit?: number;
  timestamp: string;
}

// Activity Categories
export const ACTIVITY_CATEGORIES = [
  "football",
  "basketball",
  "tennis",
  "swimming",
  "running",
  "cycling",
  "yoga",
  "fitness",
  "volleyball",
  "paddle",
] as const;

export type ActivityCategory = typeof ACTIVITY_CATEGORIES[number];

// Difficulty Levels
export const DIFFICULTY_LEVELS = [
  "beginner",
  "intermediate",
  "advanced",
] as const;

export type DifficultyLevel = typeof DIFFICULTY_LEVELS[number];
