import axios, { AxiosInstance } from "axios";
import {
  Activity,
  ActivitiesResponse,
  ActivityResponse,
  ApiResponse,
} from "../types";

// Base URL for Activities API
const ACTIVITIES_API_URL =
  process.env.REACT_APP_ACTIVITIES_API_URL || "http://localhost:8082";

// Create axios instance for Activities API
const activitiesApi: AxiosInstance = axios.create({
  baseURL: ACTIVITIES_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor to add auth token
activitiesApi.interceptors.request.use(
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

// Activities Service
export const activitiesService = {
  // Get all active activities (public)
  async getActivities(): Promise<ActivitiesResponse> {
    const response = await activitiesApi.get<ActivitiesResponse>("/activities");
    return response.data;
  },

  // Get all activities including inactive (admin only)
  async getAllActivities(): Promise<ActivitiesResponse> {
    const response = await activitiesApi.get<ActivitiesResponse>(
      "/activities/all"
    );
    return response.data;
  },

  // Get activity by ID
  async getActivityById(id: string): Promise<Activity> {
    const response = await activitiesApi.get<ActivityResponse>(
      `/activities/${id}`
    );
    return response.data.activity;
  },

  // Get activities by category
  async getActivitiesByCategory(
    category: string
  ): Promise<ActivitiesResponse> {
    const response = await activitiesApi.get<ActivitiesResponse>(
      `/activities/category/${category}`
    );
    return response.data;
  },

  // Create activity (admin only)
  async createActivity(
    activityData: Partial<Activity>
  ): Promise<ApiResponse<Activity>> {
    const response = await activitiesApi.post<ApiResponse<Activity>>(
      "/activities",
      activityData
    );
    return response.data;
  },

  // Update activity (admin only)
  async updateActivity(
    id: string,
    activityData: Partial<Activity>
  ): Promise<ApiResponse<Activity>> {
    const response = await activitiesApi.put<ApiResponse<Activity>>(
      `/activities/${id}`,
      activityData
    );
    return response.data;
  },

  // Delete activity (admin only)
  async deleteActivity(id: string): Promise<void> {
    await activitiesApi.delete(`/activities/${id}`);
  },

  // Toggle activity active status (admin only)
  async toggleActivityStatus(id: string): Promise<ApiResponse<Activity>> {
    const response = await activitiesApi.patch<ApiResponse<Activity>>(
      `/activities/${id}/toggle`
    );
    return response.data;
  },

  // Health check
  async healthCheck(): Promise<{ status: string; service: string }> {
    const response = await activitiesApi.get<{ status: string; service: string }>("/healthz");
    return response.data;
  },
};

export default activitiesService;

