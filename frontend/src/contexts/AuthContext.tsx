import React, {
  createContext,
  useContext,
  useReducer,
  useEffect,
  ReactNode,
} from "react";
import { authService } from "../services/authService";
import {
  AuthState,
  AuthContextType,
  AuthAction,
  RegisterRequest,
  UpdateProfileRequest,
  ChangePasswordRequest,
  AuthResult,
} from "../types";

// Estado inicial
const initialState: AuthState = {
  user: null,
  token: localStorage.getItem("token"),
  isAuthenticated: false,
  loading: true,
  error: null,
};

// Actions
const AUTH_ACTIONS = {
  LOGIN_START: "LOGIN_START",
  LOGIN_SUCCESS: "LOGIN_SUCCESS",
  LOGIN_FAILURE: "LOGIN_FAILURE",
  REGISTER_START: "REGISTER_START",
  REGISTER_SUCCESS: "REGISTER_SUCCESS",
  REGISTER_FAILURE: "REGISTER_FAILURE",
  LOGOUT: "LOGOUT",
  LOAD_USER_START: "LOAD_USER_START",
  LOAD_USER_SUCCESS: "LOAD_USER_SUCCESS",
  LOAD_USER_FAILURE: "LOAD_USER_FAILURE",
  UPDATE_PROFILE_SUCCESS: "UPDATE_PROFILE_SUCCESS",
  CLEAR_ERROR: "CLEAR_ERROR",
} as const;

// Reducer
function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case "LOGIN_START":
    case "REGISTER_START":
    case "LOAD_USER_START":
      return {
        ...state,
        loading: true,
        error: null,
      };

    case "LOGIN_SUCCESS":
    case "REGISTER_SUCCESS":
      localStorage.setItem("token", action.payload.token);
      return {
        ...state,
        token: action.payload.token,
        user: action.payload.user,
        isAuthenticated: true,
        loading: false,
        error: null,
      };

    case "LOAD_USER_SUCCESS":
    case "UPDATE_PROFILE_SUCCESS":
      return {
        ...state,
        user: action.payload,
        isAuthenticated: true,
        loading: false,
        error: null,
      };

    case "LOGIN_FAILURE":
    case "REGISTER_FAILURE":
    case "LOAD_USER_FAILURE":
      return {
        ...state,
        token: null,
        user: null,
        isAuthenticated: false,
        loading: false,
        error: action.payload,
      };

    case "LOGOUT":
      localStorage.removeItem("token");
      return {
        ...state,
        token: null,
        user: null,
        isAuthenticated: false,
        loading: false,
        error: null,
      };

    case "CLEAR_ERROR":
      return {
        ...state,
        error: null,
      };

    default:
      return state;
  }
}

// Context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Provider Props
interface AuthProviderProps {
  children: ReactNode;
}

// Provider
export function AuthProvider({ children }: AuthProviderProps): JSX.Element {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Cargar usuario si hay token al inicializar
  useEffect(() => {
    const loadUser = async (): Promise<void> => {
      const token = localStorage.getItem("token");
      if (token) {
        try {
          dispatch({ type: "LOAD_USER_START" });
          const userData = await authService.getProfile();
          dispatch({
            type: "LOAD_USER_SUCCESS",
            payload: userData,
          });
        } catch (error: any) {
          dispatch({
            type: "LOAD_USER_FAILURE",
            payload: error.message,
          });
          localStorage.removeItem("token");
        }
      } else {
        dispatch({
          type: "LOAD_USER_FAILURE",
          payload: "No token found",
        });
      }
    };

    loadUser();
  }, []);

  // Login
  const login = async (
    email: string,
    password: string
  ): Promise<AuthResult> => {
    try {
      dispatch({ type: "LOGIN_START" });
      const response = await authService.login(email, password);
      dispatch({
        type: "LOGIN_SUCCESS",
        payload: response.data,
      });
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || error.message || "Login failed";
      dispatch({
        type: "LOGIN_FAILURE",
        payload: errorMessage,
      });
      return { success: false, error: errorMessage };
    }
  };

  // Register
  const register = async (userData: RegisterRequest): Promise<AuthResult> => {
    try {
      dispatch({ type: "REGISTER_START" });
      const response = await authService.register(userData);
      dispatch({
        type: "REGISTER_SUCCESS",
        payload: response.data,
      });
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || error.message || "Registration failed";
      dispatch({
        type: "REGISTER_FAILURE",
        payload: errorMessage,
      });
      return { success: false, error: errorMessage };
    }
  };

  // Logout
  const logout = (): void => {
    dispatch({ type: "LOGOUT" });
  };

  // Update Profile
  const updateProfile = async (
    profileData: UpdateProfileRequest
  ): Promise<AuthResult> => {
    try {
      const response = await authService.updateProfile(profileData);
      dispatch({
        type: "UPDATE_PROFILE_SUCCESS",
        payload: response.data,
      });
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || error.message || "Update failed";
      return { success: false, error: errorMessage };
    }
  };

  // Change Password
  const changePassword = async (
    passwordData: ChangePasswordRequest
  ): Promise<AuthResult> => {
    try {
      await authService.changePassword(passwordData);
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message ||
        error.message ||
        "Password change failed";
      return { success: false, error: errorMessage };
    }
  };

  // Upload Avatar
  const uploadAvatar = async (file: File): Promise<AuthResult> => {
    try {
      const response = await authService.uploadAvatar(file);
      dispatch({
        type: "UPDATE_PROFILE_SUCCESS",
        payload: response.data,
      });
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || error.message || "Avatar upload failed";
      return { success: false, error: errorMessage };
    }
  };

  // Delete Avatar
  const deleteAvatar = async (): Promise<AuthResult> => {
    try {
      const response = await authService.deleteAvatar();
      dispatch({
        type: "UPDATE_PROFILE_SUCCESS",
        payload: response.data,
      });
      return { success: true };
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || error.message || "Avatar deletion failed";
      return { success: false, error: errorMessage };
    }
  };

  // Clear Error
  const clearError = (): void => {
    dispatch({ type: "CLEAR_ERROR" });
  };

  const value: AuthContextType = {
    ...state,
    login,
    register,
    logout,
    updateProfile,
    changePassword,
    uploadAvatar,
    deleteAvatar,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// Hook personalizado
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
