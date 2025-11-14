import React, { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActions,
  Button,
  Chip,
  CircularProgress,
  Alert,
  TextField,
  InputAdornment,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Paper,
  Stack,
  Slider,
  Divider,
} from "@mui/material";
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

  // Search filters state
  const [query, setQuery] = useState<string>(searchParams.get("q") || "");
  const [category, setCategory] = useState<string>(searchParams.get("category") || "all");
  const [difficulty, setDifficulty] = useState<string>(searchParams.get("difficulty") || "all");
  const [priceRange, setPriceRange] = useState<number[]>([0, 1000]);
  const [page, setPage] = useState<number>(1);
  const [pageSize] = useState<number>(12);

  useEffect(() => {
    // Load initial search if query params exist
    if (searchParams.get("q") || searchParams.get("category") || searchParams.get("difficulty")) {
      performSearch();
    }
  }, []);

  const performSearch = async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const filters: SearchFilters = {
        query: query && query.trim() ? query.trim() : undefined,
        category: category !== "all" ? category : undefined,
        difficulty: difficulty !== "all" ? difficulty : undefined,
        price_min: priceRange[0] > 0 ? priceRange[0] : undefined,
        price_max: priceRange[1] < 1000 ? priceRange[1] : undefined,
        page: page,
        size: pageSize,
      };

      const result = await searchService.searchActivities(filters);
      setActivities(result.results || []);
      setTotalFound(result.total_found || 0);

      // Update URL params
      const newParams = new URLSearchParams();
      if (filters.query) newParams.set("q", filters.query);
      if (filters.category) newParams.set("category", filters.category);
      if (filters.difficulty) newParams.set("difficulty", filters.difficulty);
      setSearchParams(newParams);
    } catch (err: any) {
      console.error("Error searching activities:", err);
      
      // Better error handling
      if (err.response) {
        // Server responded with error status
        setError(err.response.data?.message || `Server error: ${err.response.status}`);
      } else if (err.request) {
        // Request was made but no response received
        setError("Network error: Unable to connect to search service. Please check your connection.");
      } else {
        // Something else happened
        setError(err.message || "Failed to search activities. Please try again.");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (): void => {
    setPage(1);
    performSearch();
  };

  const handleActivityClick = (id: string): void => {
    navigate(`/activities/${id}`);
  };

  const getCategoryColor = (
    category: string
  ): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
    const colors: Record<string, "default" | "primary" | "secondary" | "success" | "warning" | "error"> = {
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
    const colors: Record<string, "default" | "primary" | "secondary" | "success" | "warning" | "error"> = {
      beginner: "success",
      intermediate: "warning",
      advanced: "error",
    };
    return colors[difficulty] || "default";
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom>
        🔍 Search Activities
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Find the perfect activity for you
      </Typography>

      {/* API Status */}
      {apiStatus && (
        <Alert
          severity={apiStatus.status === "online" ? "success" : "error"}
          sx={{ mb: 3 }}
        >
          API Status: {apiStatus.message}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* Filters Sidebar */}
        <Grid item xs={12} md={3}>
          <Paper sx={{ p: 3, position: "sticky", top: 20 }}>
            <Typography variant="h6" gutterBottom>
              <FilterIcon sx={{ verticalAlign: "middle", mr: 1 }} />
              Filters
            </Typography>

            <Stack spacing={3}>
              {/* Search Query */}
              <TextField
                fullWidth
                label="Search"
                placeholder="Search activities..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyPress={(e) => {
                  if (e.key === "Enter") {
                    handleSearch();
                  }
                }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
              />

              {/* Category Filter */}
              <FormControl fullWidth>
                <InputLabel>Category</InputLabel>
                <Select
                  value={category}
                  label="Category"
                  onChange={(e) => {
                    setCategory(e.target.value);
                    // Trigger search automatically when filter changes
                    setTimeout(() => {
                      setPage(1);
                      performSearch();
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

              {/* Difficulty Filter */}
              <FormControl fullWidth>
                <InputLabel>Difficulty</InputLabel>
                <Select
                  value={difficulty}
                  label="Difficulty"
                  onChange={(e) => {
                    setDifficulty(e.target.value);
                    // Trigger search automatically when filter changes
                    setTimeout(() => {
                      setPage(1);
                      performSearch();
                    }, 100);
                  }}
                >
                  <MenuItem value="all">All Levels</MenuItem>
                  <MenuItem value="beginner">Beginner</MenuItem>
                  <MenuItem value="intermediate">Intermediate</MenuItem>
                  <MenuItem value="advanced">Advanced</MenuItem>
                </Select>
              </FormControl>

              {/* Price Range */}
              <Box>
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
                <Box sx={{ display: "flex", justifyContent: "space-between", mt: 1 }}>
                  <Typography variant="caption">${priceRange[0]}</Typography>
                  <Typography variant="caption">${priceRange[1]}</Typography>
                </Box>
              </Box>

              <Button
                variant="contained"
                fullWidth
                onClick={handleSearch}
                disabled={loading}
                startIcon={loading ? <CircularProgress size={20} /> : <SearchIcon />}
              >
                {loading ? "Searching..." : "Search"}
              </Button>
            </Stack>
          </Paper>
        </Grid>

        {/* Results */}
        <Grid item xs={12} md={9}>
          {loading && activities.length === 0 ? (
            <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error">{error}</Alert>
          ) : activities.length === 0 ? (
            <Paper sx={{ p: 4, textAlign: "center" }}>
              <Typography variant="h6" color="textSecondary">
                No activities found
              </Typography>
              <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
                Try adjusting your search criteria or filters
              </Typography>
            </Paper>
          ) : (
            <>
              <Box sx={{ mb: 2, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <Typography variant="body1" color="textSecondary">
                  Found {totalFound} result{totalFound !== 1 ? "s" : ""}
                </Typography>
              </Box>

              <Grid container spacing={3}>
                {activities.map((activity) => (
                  <Grid item xs={12} sm={6} lg={4} key={activity.id}>
                    <Card
                      sx={{
                        height: "100%",
                        display: "flex",
                        flexDirection: "column",
                        cursor: "pointer",
                        transition: "transform 0.2s, box-shadow 0.2s",
                        "&:hover": {
                          transform: "translateY(-4px)",
                          boxShadow: 6,
                        },
                      }}
                      onClick={() => handleActivityClick(activity.id)}
                    >
                      {activity.image_url ? (
                        <CardMedia
                          component="img"
                          height="200"
                          image={activity.image_url}
                          alt={activity.name}
                        />
                      ) : (
                        <Box
                          sx={{
                            height: 200,
                            bgcolor: "primary.light",
                            display: "flex",
                            alignItems: "center",
                            justifyContent: "center",
                          }}
                        >
                          <Typography variant="h4" color="primary.contrastText">
                            {activity.name.charAt(0)}
                          </Typography>
                        </Box>
                      )}

                      <CardContent sx={{ flexGrow: 1 }}>
                        <Box
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "start",
                            mb: 1,
                          }}
                        >
                          <Typography variant="h6" component="h2" noWrap>
                            {activity.name}
                          </Typography>
                          {!activity.is_active && (
                            <Chip label="Inactive" size="small" color="default" />
                          )}
                        </Box>

                        <Typography
                          variant="body2"
                          color="textSecondary"
                          sx={{
                            mb: 2,
                            display: "-webkit-box",
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: "vertical",
                            overflow: "hidden",
                          }}
                        >
                          {activity.description}
                        </Typography>

                        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 2 }}>
                          <Chip
                            label={activity.category}
                            size="small"
                            color={getCategoryColor(activity.category)}
                          />
                          <Chip
                            label={activity.difficulty}
                            size="small"
                            color={getDifficultyColor(activity.difficulty)}
                          />
                        </Box>

                        <Stack spacing={1}>
                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <LocationIcon fontSize="small" color="action" />
                            <Typography variant="body2" color="textSecondary">
                              {activity.location}
                            </Typography>
                          </Box>

                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <MoneyIcon fontSize="small" color="action" />
                            <Typography variant="body2" fontWeight="bold">
                              ${activity.price.toFixed(2)}
                            </Typography>
                          </Box>

                          <Box
                            sx={{
                              display: "flex",
                              alignItems: "center",
                              gap: 2,
                            }}
                          >
                            <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                              <PeopleIcon fontSize="small" color="action" />
                              <Typography variant="body2" color="textSecondary">
                                Max: {activity.max_capacity}
                              </Typography>
                            </Box>

                            <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                              <TimeIcon fontSize="small" color="action" />
                              <Typography variant="body2" color="textSecondary">
                                {activity.duration} min
                              </Typography>
                            </Box>
                          </Box>
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
            </>
          )}
        </Grid>
      </Grid>
    </Box>
  );
}

export default SearchPage;

