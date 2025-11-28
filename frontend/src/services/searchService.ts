import { SearchFilters, SearchResult } from "../types";
import { searchApi } from "./api";

export const searchService = {
  async healthCheck(): Promise<string> {
    try {
      await searchApi.get("/health");
      return "Search API is up";
    } catch (error) {
      return "Search API unreachable";
    }
  },

  async searchActivities(filters: SearchFilters): Promise<SearchResult> {
    const params = new URLSearchParams();

    if (filters.query && filters.query.trim()) {
      params.append("query", filters.query.trim());
    }
    // If no query, backend uses default (*:*)

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
};
