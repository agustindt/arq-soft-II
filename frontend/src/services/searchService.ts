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
});

// Search Service
export const searchService = {
  // Search activities with filters
  async searchActivities(
    filters: SearchFilters = {}
  ): Promise<SearchResult> {
    const params = new URLSearchParams();

    if (filters.query) {
      params.append("query", filters.query);
    } else {
      params.append("query", "*");
    }

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

    if (filters.size) {
      params.append("size", filters.size.toString());
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
};

export default searchService;

