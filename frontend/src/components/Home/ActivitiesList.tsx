import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
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
} from "@mui/material";
import {
  Search as SearchIcon,
  LocationOn as LocationIcon,
  AttachMoney as MoneyIcon,
  People as PeopleIcon,
  AccessTime as TimeIcon,
  FilterList as FilterIcon,
} from "@mui/icons-material";
import { activitiesService } from "../../services/activitiesService";
import { Activity, ActivityCategory, DifficultyLevel } from "../../types";

function ActivitiesList(): JSX.Element {
  const navigate = useNavigate();
  const [activities, setActivities] = useState<Activity[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [categoryFilter, setCategoryFilter] = useState<string>("all");
  const [difficultyFilter, setDifficultyFilter] = useState<string>("all");

  useEffect(() => {
    loadActivities();
  }, [categoryFilter]);

  const loadActivities = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      let response;
      if (categoryFilter === "all") {
        response = await activitiesService.getActivities();
      } else {
        response = await activitiesService.getActivitiesByCategory(
          categoryFilter
        );
      }

      // Apply client-side filtering for difficulty and search
      let filtered = response.activities;

      if (difficultyFilter !== "all") {
        filtered = filtered.filter(
          (activity) => activity.difficulty === difficultyFilter
        );
      }

      if (searchQuery.trim()) {
        const query = searchQuery.toLowerCase();
        filtered = filtered.filter(
          (activity) =>
            activity.name.toLowerCase().includes(query) ||
            activity.description.toLowerCase().includes(query) ||
            activity.location.toLowerCase().includes(query) ||
            activity.category.toLowerCase().includes(query)
        );
      }

      setActivities(filtered);
    } catch (err: any) {
      setError(err.message || "Failed to load activities");
      console.error("Error loading activities:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (): void => {
    loadActivities();
  };

  const handleActivityClick = (id: string): void => {
    navigate(`/activities/${id}`);
  };

  const getCategoryColor = (category: string): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
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

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom>
        🏃‍♀️ Sports Activities
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Discover and join amazing sports activities in your area
      </Typography>

      {/* Search and Filters */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Stack spacing={2}>
          <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
            <TextField
              fullWidth
              placeholder="Search activities..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
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
              sx={{ flexGrow: 1, minWidth: 200 }}
            />
            <Button
              variant="contained"
              onClick={handleSearch}
              startIcon={<SearchIcon />}
            >
              Search
            </Button>
          </Box>

          <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap", alignItems: "center" }}>
            <FilterIcon sx={{ color: "text.secondary" }} />
            <FormControl sx={{ minWidth: 150 }}>
              <InputLabel>Category</InputLabel>
              <Select
                value={categoryFilter}
                label="Category"
                onChange={(e) => setCategoryFilter(e.target.value)}
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

            <FormControl sx={{ minWidth: 150 }}>
              <InputLabel>Difficulty</InputLabel>
              <Select
                value={difficultyFilter}
                label="Difficulty"
                onChange={(e) => {
                  setDifficultyFilter(e.target.value);
                  loadActivities();
                }}
              >
                <MenuItem value="all">All Levels</MenuItem>
                <MenuItem value="beginner">Beginner</MenuItem>
                <MenuItem value="intermediate">Intermediate</MenuItem>
                <MenuItem value="advanced">Advanced</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </Stack>
      </Paper>

      {/* Error Message */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Activities Grid */}
      {activities.length === 0 ? (
        <Paper sx={{ p: 4, textAlign: "center" }}>
          <Typography variant="h6" color="textSecondary">
            No activities found
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
            Try adjusting your search or filters
          </Typography>
        </Paper>
      ) : (
        <Grid container spacing={3}>
          {activities.map((activity) => (
            <Grid item xs={12} sm={6} md={4} key={activity.id}>
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
      )}

      {/* Results Count */}
      {activities.length > 0 && (
        <Typography variant="body2" color="textSecondary" sx={{ mt: 3 }}>
          Showing {activities.length} activity{activities.length !== 1 ? "ies" : ""}
        </Typography>
      )}
    </Box>
  );
}

export default ActivitiesList;

