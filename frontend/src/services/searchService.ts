import axios, { AxiosInstance } from "axios";
import { SearchResult, SearchFilters } from "../types";

// Base URL for Search API
const SEARCH_API_URL =
  process.env.REACT_APP_SEARCH_API_URL || "http://localhost:8083";

// Create axios instance for Search API
const searchApi: AxiosInstance = axios.create({
  baseURL: SEARCH_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000, // 10 second timeout
});

// Response interceptor for better error handling
searchApi.interceptors.response.use(
  (response) => response,
  (error) => {
    // Enhance error message
    if (error.code === "ECONNABORTED") {
      error.message =
        "Request timeout: The search service took too long to respond";
    } else if (error.code === "ERR_NETWORK") {
      error.message = "Network error: Unable to connect to search service";
    }
    return Promise.reject(error);
  }
);

// Search Service
export const searchService = {
  // Search activities with filters
  async searchActivities(filters: SearchFilters = {}): Promise<SearchResult> {
    const params = new URLSearchParams();

    // Only add query parameter if it's not empty
    if (filters.query && filters.query.trim()) {
      params.append("query", filters.query.trim());
    }
    // If no query, let the backend handle it (it will use *:* for all documents)

    if (filters.category) {
      params.append("category", filters.category);
    }

    if (filters.difficulty) {
      params.append("difficulty", filters.difficulty);
    }

    if (filters.price_min !== undefined) {
      params.append("price_min", filters.price_min.toString());
    }

    if (filters.price_max !== undefined) {
      params.append("price_max", filters.price_max.toString());
    }

    if (filters.page) {
      params.append("page", filters.page.toString());
    }

    if (filters.limit) {
      params.append("limit", filters.limit.toString());
    }

    if (filters.sort) {
      params.append("sort", filters.sort);
    }

    const queryString = params.toString();
    const url = queryString ? `/search?${queryString}` : "/search";

    const response = await searchApi.get<SearchResult>(url);
    return response.data;
  },

  // Get available categories (helper function)
  getCategories(): string[] {
    return [
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
    ];
  },

  // Get available difficulty levels (helper function)
  getDifficulties(): string[] {
    return ["beginner", "intermediate", "advanced"];
  },

  // Health check
  async healthCheck(): Promise<string> {
    const response = await searchApi.get<string>("/health");
    return response.data;
  },
};

export default searchService;
