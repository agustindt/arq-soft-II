export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  user: User;
}

export interface Activity {
  id: string;
  name: string;
  description: string;
  category: string;
  difficulty: string;
  location: string;
  duration: number;
  price: number;
  maxCapacity: number;
  instructor?: string;
  dateCreated?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface ActivityResponse {
  activity: Activity;
}

export interface ActivityListResponse {
  activities: Activity[];
}

export interface CreateActivityRequest {
  name: string;
  description: string;
  category: string;
  difficulty: string;
  location: string;
  duration: number;
  price: number;
  maxCapacity: number;
}

export interface UpdateActivityRequest extends CreateActivityRequest {}

export interface Reservation {
  id: string;
  userId: number;
  activityId: string;
  participants: number;
  date: string;
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
  date: string; // ISO date string
  participants: number;
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
