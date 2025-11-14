import axios, { AxiosInstance } from "axios";
import {
  Reservation,
  ReservationsResponse,
  ReservationResponse,
  CreateReservationRequest,
  ApiResponse,
} from "../types";

// Base URL for Reservations API
const RESERVATIONS_API_URL =
  process.env.REACT_APP_RESERVATIONS_API_URL || "http://localhost:8080";

// Create axios instance for Reservations API
const reservationsApi: AxiosInstance = axios.create({
  baseURL: RESERVATIONS_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor to add auth token
reservationsApi.interceptors.request.use(
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

// Reservations Service
export const reservationsService = {
  // Get all reservations (public, but filtered by user in backend)
  async getReservations(): Promise<Reservation[]> {
    const response = await reservationsApi.get<ReservationsResponse>(
      "/reservas"
    );
    return response.data.reservas;
  },

  // Get reservation by ID
  async getReservationById(id: string): Promise<Reservation> {
    const response = await reservationsApi.get<ReservationResponse>(
      `/reservas/${id}`
    );
    return response.data.Reserva;
  },

  // Create reservation (admin only based on backend, but we'll try for all authenticated users)
  async createReservation(
    reservationData: CreateReservationRequest
  ): Promise<Reservation> {
    // Map frontend format to backend format
    // Backend expects: users_id (array), actividad (string), cupo (number), date (ISO string), status (string)
    const backendFormat = {
      users_id: [], // Will be populated by backend from token
      actividad: reservationData.activityId,
      cupo: reservationData.participants,
      date: reservationData.date,
      status: "Pendiente",
    };

    const response = await reservationsApi.post<ReservationResponse>(
      "/reservas",
      backendFormat
    );
    return response.data.Reserva;
  },

  // Update reservation (admin only)
  async updateReservation(
    id: string,
    reservationData: Partial<Reservation>
  ): Promise<Reservation> {
    const response = await reservationsApi.put<ReservationResponse>(
      `/reservas/${id}`,
      reservationData
    );
    return response.data.Reserva;
  },

  // Delete/Cancel reservation (admin only)
  async deleteReservation(id: string): Promise<void> {
    await reservationsApi.delete(`/reservas/${id}`);
  },

  // Health check
  async healthCheck(): Promise<{ status: string; message: string; service: string }> {
    const response = await reservationsApi.get<{ status: string; message: string; service: string }>("/healthz");
    return response.data;
  },
};

export default reservationsService;

