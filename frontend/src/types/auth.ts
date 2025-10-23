// Social Links Type
export interface SocialLinks {
  instagram?: string;
  twitter?: string;
  facebook?: string;
  linkedin?: string;
  youtube?: string;
  website?: string;
}

// Authentication Types
export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  // Extended profile fields
  avatar_url?: string | null;
  bio?: string | null;
  phone?: string | null;
  birth_date?: string | null;
  location?: string | null;
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say' | null;
  height?: number | null; // in cm
  weight?: number | null; // in kg
  sports_interests?: string | null; // JSON string of sports array
  fitness_level?: 'beginner' | 'intermediate' | 'advanced' | 'professional' | null;
  social_links: SocialLinks;
  // System fields
  role: string;
  email_verified: boolean;
  email_verified_at?: string | null;
  is_active: boolean;
  last_login_at?: string | null;
  created_at: string;
  updated_at: string;
}

// Public User Response (for listings)
export interface PublicUser {
  id: number;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url?: string | null;
  bio?: string | null;
  location?: string | null;
  social_links: SocialLinks;
  fitness_level?: string | null;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

export interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<AuthResult>;
  register: (userData: RegisterRequest) => Promise<AuthResult>;
  logout: () => void;
  updateProfile: (profileData: UpdateProfileRequest) => Promise<AuthResult>;
  changePassword: (passwordData: ChangePasswordRequest) => Promise<AuthResult>;
  uploadAvatar: (file: File) => Promise<AuthResult>;
  deleteAvatar: () => Promise<AuthResult>;
  clearError: () => void;
}

export interface AuthResult {
  success: boolean;
  error?: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  avatar_url?: string;
  bio?: string;
  phone?: string;
  birth_date?: string; // ISO date string
  location?: string;
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say';
  height?: number; // in cm
  weight?: number; // in kg
  sports_interests?: string; // JSON string of sports array
  fitness_level?: 'beginner' | 'intermediate' | 'advanced' | 'professional';
  social_links?: SocialLinks;
}

// Avatar Upload Response
export interface AvatarUploadResponse {
  message: string;
  avatar_url: string;
  data: User;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}

// API Response Types
export interface ApiResponse<T> {
  message: string;
  data: T;
}

export interface ApiError {
  error: string;
  message: string;
}

// Auth Action Types
export type AuthAction =
  | { type: "LOGIN_START" }
  | { type: "LOGIN_SUCCESS"; payload: LoginResponse }
  | { type: "LOGIN_FAILURE"; payload: string }
  | { type: "REGISTER_START" }
  | { type: "REGISTER_SUCCESS"; payload: LoginResponse }
  | { type: "REGISTER_FAILURE"; payload: string }
  | { type: "LOGOUT" }
  | { type: "LOAD_USER_START" }
  | { type: "LOAD_USER_SUCCESS"; payload: User }
  | { type: "LOAD_USER_FAILURE"; payload: string }
  | { type: "UPDATE_PROFILE_SUCCESS"; payload: User }
  | { type: "CLEAR_ERROR" };
