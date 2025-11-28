import {
  Box,
  Typography,
  Grid,
  Paper,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  TextField,
  Slider,
  Button,
  Stack,
  Chip,
  Card,
  CardActions,
  CardContent,
  InputAdornment,
  Alert,
} from "@mui/material";
import { useNavigate, useSearchParams } from "react-router-dom";
import React, { useEffect, useState } from "react";
import {
  Search as SearchIcon,
  LocationOn as LocationIcon,
  AttachMoney as MoneyIcon,
  People as PeopleIcon,
  AccessTime as TimeIcon,
  FilterList as FilterIcon,
} from "@mui/icons-material";
import { searchService } from "../../services/searchService";
import { Activity, SearchFilters } from "../../types";
import { useApiStatus } from "../../hooks/useApiStatus";

function SearchPage(): JSX.Element {
  const navigate = useNavigate();
  const apiStatus = useApiStatus(() => searchService.healthCheck(), "Search API");
  const [searchParams, setSearchParams] = useSearchParams();
  const [activities, setActivities] = useState<Activity[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [totalFound, setTotalFound] = useState<number>(0);

  const [query, setQuery] = useState<string>(searchParams.get("q") || "");
  const [category, setCategory] = useState<string>(searchParams.get("category") || "all");
  const [difficulty, setDifficulty] = useState<string>(searchParams.get("difficulty") || "all");
  const [priceRange, setPriceRange] = useState<number[]>([0, 1000]);

  const initialPage = parseInt(searchParams.get("page") || "1", 10) || 1;
  const initialLimit = parseInt(searchParams.get("limit") || "12", 10) || 12;
  const [page, setPage] = useState<number>(initialPage);
  const [limit, setLimit] = useState<number>(initialLimit);

  useEffect(() => {
    if (searchParams.get("q") || searchParams.get("category") || searchParams.get("difficulty")) {
      performSearch();
    }
  }, []);

  const performSearch = async (targetPage = page): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const filters: SearchFilters = {
        query: query && query.trim() ? query.trim() : undefined,
        category: category !== "all" ? category : undefined,
        difficulty: difficulty !== "all" ? difficulty : undefined,
        price_min: priceRange[0] > 0 ? priceRange[0] : undefined,
        price_max: priceRange[1] < 1000 ? priceRange[1] : undefined,
        page: targetPage,
        limit: limit,
      };

      const result = await searchService.searchActivities(filters);
      setActivities(result.results || []);
      setTotalFound(result.total_found || 0);

      if (result.page) setPage(result.page);
      if (result.limit) setLimit(result.limit);

      const newParams = new URLSearchParams();
      if (filters.query) newParams.set("q", filters.query);
      if (filters.category) newParams.set("category", filters.category);
      if (filters.difficulty) newParams.set("difficulty", filters.difficulty);
      newParams.set("page", String(targetPage));
      newParams.set("limit", String(limit));
      setSearchParams(newParams);
    } catch (err: any) {
      console.error("Error searching activities:", err);

      if (err.response) {
        setError(err.response.data?.message || `Server error: ${err.response.status}`);
      } else if (err.request) {
        setError("Network error: Unable to connect to search service.");
      } else {
        setError(err.message || "Failed to search activities.");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (): void => {
    setPage(1);
    performSearch(1);
  };

  const handleActivityClick = (id: string): void => {
    navigate(`/activities/${id}`);
  };

  const getCategoryColor = (
    category: string
  ): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
    const colors: Record<string, any> = {
      football: "primary",
      basketball: "secondary",
      tennis: "success",
      swimming: "primary",
      running: "warning",
      cycling: "error",
      yoga: "success",
      fitness: "primary",
      volleyball: "secondary",
      paddle: "default",
    };
    return colors[category] || "default";
  };

  const getDifficultyColor = (
    difficulty: string
  ): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
    const colors: Record<string, any> = {
      beginner: "success",
      intermediate: "warning",
      advanced: "error",
    };
    return colors[difficulty] || "default";
  };

  const totalPages = Math.max(1, Math.ceil((totalFound || 0) / (limit || 1)));

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom>
        🔍 Search Activities
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Find the perfect activity for you
      </Typography>

      {apiStatus && (
        <Alert severity={apiStatus.status === "online" ? "success" : "error"} sx={{ mb: 3 }}>
          API Status: {apiStatus.message}
        </Alert>
      )}

      <Grid container spacing={3}>
        <Grid item xs={12} md={3}>
          <Paper sx={{ p: 3, position: "sticky", top: 20 }}>
            <Typography variant="h6" gutterBottom>
              <FilterIcon sx={{ verticalAlign: "middle", mr: 1 }} />
              Filters
            </Typography>

            <TextField
              fullWidth
              label="Search..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === "Enter") handleSearch();
              }}
              sx={{ mb: 3 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />

            <FormControl fullWidth sx={{ mb: 3 }}>
              <InputLabel>Category</InputLabel>
              <Select
                value={category}
                label="Category"
                onChange={(e) => {
                  setCategory(e.target.value);
                  setTimeout(() => {
                    setPage(1);
                    performSearch(1);
                  }, 100);
                }}
              >
                <MenuItem value="all">All Categories</MenuItem>
                <MenuItem value="football">Football</MenuItem>
                <MenuItem value="basketball">Basketball</MenuItem>
                <MenuItem value="tennis">Tennis</MenuItem>
                <MenuItem value="swimming">Swimming</MenuItem>
                <MenuItem value="running">Running</MenuItem>
                <MenuItem value="cycling">Cycling</MenuItem>
                <MenuItem value="yoga">Yoga</MenuItem>
                <MenuItem value="fitness">Fitness</MenuItem>
                <MenuItem value="volleyball">Volleyball</MenuItem>
                <MenuItem value="paddle">Paddle</MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth sx={{ mb: 3 }}>
              <InputLabel>Difficulty</InputLabel>
              <Select
                value={difficulty}
                label="Difficulty"
                onChange={(e) => {
                  setDifficulty(e.target.value);
                  setTimeout(() => {
                    setPage(1);
                    performSearch(1);
                  }, 100);
                }}
              >
                <MenuItem value="all">All Levels</MenuItem>
                <MenuItem value="beginner">Beginner</MenuItem>
                <MenuItem value="intermediate">Intermediate</MenuItem>
                <MenuItem value="advanced">Advanced</MenuItem>
              </Select>
            </FormControl>

            <Box sx={{ mb: 2 }}>
              <Typography gutterBottom>Price Range</Typography>
              <Slider
                value={priceRange}
                onChange={(_, newValue) => setPriceRange(newValue as number[])}
                valueLabelDisplay="auto"
                min={0}
                max={1000}
                step={10}
                marks={[
                  { value: 0, label: "$0" },
                  { value: 500, label: "$500" },
                  { value: 1000, label: "$1000+" },
                ]}
              />
              <Button fullWidth variant="contained" sx={{ mt: 1 }} onClick={handleSearch}>
                Apply Filters
              </Button>
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={12} md={9}>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {!loading && activities.length === 0 && (
            <Typography variant="h6" color="textSecondary">
              No activities found.
            </Typography>
          )}

          {loading ? (
            <Typography>Loading...</Typography>
          ) : (
            <>
              <Grid container spacing={3}>
                {activities.map((activity) => (
                  <Grid item xs={12} sm={6} md={4} key={activity.id}>
                    <Card
                      sx={{
                        cursor: "pointer",
                        transition: "all 0.2s",
                        "&:hover": {
                          transform: "scale(1.03)",
                          boxShadow: 4,
                        },
                      }}
                      onClick={() => handleActivityClick(activity.id)}
                    >
                      <CardContent>
                        <Typography variant="h6" gutterBottom>
                          {activity.name}
                        </Typography>
                        <Stack spacing={1} sx={{ mt: 1 }}>
                          <Chip
                            label={activity.category}
                            color={getCategoryColor(activity.category)}
                            size="small"
                          />
                          <Chip
                            label={activity.difficulty}
                            color={getDifficultyColor(activity.difficulty)}
                            size="small"
                          />

                          <Stack direction="row" spacing={1} alignItems="center">
                            <LocationIcon fontSize="small" />
                            <Typography variant="body2">{activity.location}</Typography>
                          </Stack>

                          <Stack direction="row" spacing={1} alignItems="center">
                            <MoneyIcon fontSize="small" />
                            <Typography variant="body2">${activity.price}</Typography>
                          </Stack>

                          <Stack direction="row" spacing={1} alignItems="center">
                            <PeopleIcon fontSize="small" />
                            <Typography variant="body2">{activity.capacity} seats</Typography>
                          </Stack>

                          <Stack direction="row" spacing={1} alignItems="center">
                            <TimeIcon fontSize="small" />
                            <Typography variant="body2">{activity.duration} min</Typography>
                          </Stack>
                        </Stack>
                      </CardContent>

                      <CardActions>
                        <Button
                          size="small"
                          variant="contained"
                          fullWidth
                          onClick={(e) => {
                            e.stopPropagation();
                            handleActivityClick(activity.id);
                          }}
                        >
                          View Details
                        </Button>
                      </CardActions>
                    </Card>
                  </Grid>
                ))}
              </Grid>

              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  mt: 3,
                  gap: 2,
                  flexWrap: "wrap",
                }}
              >
                <Button
                  variant="outlined"
                  disabled={page <= 1 || loading}
                  onClick={() => {
                    const newPage = Math.max(1, page - 1);
                    setPage(newPage);
                    performSearch(newPage);
                  }}
                >
                  Previous
                </Button>

                <Typography variant="body2" color="textSecondary">
                  Page {page} of {totalPages} • {totalFound} result
                  {totalFound === 1 ? "" : "s"}
                </Typography>

                <Button
                  variant="outlined"
                  disabled={page >= totalPages || loading}
                  onClick={() => {
                    const newPage = Math.min(totalPages, page + 1);
                    setPage(newPage);
                    performSearch(newPage);
                  }}
                >
                  Next
                </Button>
              </Box>
            </>
          )}
        </Grid>
      </Grid>
    </Box>
  );
}

export default SearchPage;
