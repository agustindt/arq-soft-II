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

  const getCategoryColor = (
    category: string
  ): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
    const colors: Record<
      string,
      "default" | "primary" | "secondary" | "success" | "warning" | "error"
    > = {
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
    const colors: Record<
      string,
      "default" | "primary" | "secondary" | "success" | "warning" | "error"
    > = {
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
      <Box sx={{ mb: 4 }}>
        <Typography
          variant="h4"
          gutterBottom
          sx={{
            fontWeight: 700,
            background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
            backgroundClip: "text",
            mb: 1,
          }}
        >
          🏃‍♀️ Sports Activities
        </Typography>
        <Typography variant="h6" color="textSecondary" sx={{ fontWeight: 400 }}>
          Discover and join amazing sports activities in your area
        </Typography>
      </Box>

      {/* Search and Filters */}
      <Paper
        sx={{
          p: 3,
          mb: 4,
          borderRadius: 3,
          background: "linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)",
          border: "1px solid rgba(99, 102, 241, 0.1)",
        }}
      >
        <Stack spacing={3}>
          {/* Search Row */}
          <Box
            sx={{
              display: "flex",
              gap: 2,
              alignItems: "stretch",
              flexWrap: { xs: "wrap", sm: "nowrap" },
            }}
          >
            <TextField
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
              sx={{
                flex: { xs: "1 1 100%", sm: "1 1 auto" },
                minWidth: { xs: "100%", sm: 300 },
              }}
            />
            <Button
              variant="contained"
              onClick={handleSearch}
              startIcon={<SearchIcon />}
              sx={{
                minWidth: { xs: "100%", sm: 140 },
                px: 3,
                py: 1.5,
                fontWeight: 600,
                background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                "&:hover": {
                  background:
                    "linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%)",
                  transform: "translateY(-2px)",
                  boxShadow: "0 6px 16px rgba(99, 102, 241, 0.4)",
                },
                transition: "all 0.3s ease",
              }}
            >
              Search
            </Button>
          </Box>

          {/* Filters Row */}
          <Box
            sx={{
              display: "flex",
              gap: 2,
              alignItems: "center",
              flexWrap: { xs: "wrap", sm: "nowrap" },
            }}
          >
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: 1,
                minWidth: { xs: "100%", sm: "auto" },
                mr: { xs: 0, sm: 1 },
              }}
            >
              <FilterIcon
                sx={{
                  color: "primary.main",
                  fontSize: 24,
                }}
              />
              <Typography
                variant="body2"
                sx={{
                  fontWeight: 600,
                  color: "text.secondary",
                  display: { xs: "none", sm: "block" },
                }}
              >
                Filters:
              </Typography>
            </Box>
            <FormControl
              sx={{
                minWidth: { xs: "100%", sm: 180 },
                flex: { xs: "1 1 100%", sm: "0 1 auto" },
              }}
            >
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

            <FormControl
              sx={{
                minWidth: { xs: "100%", sm: 180 },
                flex: { xs: "1 1 100%", sm: "0 1 auto" },
              }}
            >
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
                  transition: "all 0.3s cubic-bezier(0.4, 0, 0.2, 1)",
                  border: "1px solid rgba(0, 0, 0, 0.05)",
                  "&:hover": {
                    transform: "translateY(-8px)",
                    boxShadow: "0 12px 24px rgba(99, 102, 241, 0.15)",
                    borderColor: "rgba(99, 102, 241, 0.2)",
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
                      background:
                        "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                    }}
                  >
                    <Typography
                      variant="h3"
                      sx={{
                        color: "white",
                        fontWeight: 700,
                        textShadow: "0 2px 4px rgba(0,0,0,0.2)",
                      }}
                    >
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

                  <Box
                    sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 2 }}
                  >
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
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 0.5 }}
                      >
                        <PeopleIcon fontSize="small" color="action" />
                        <Typography variant="body2" color="textSecondary">
                          Max: {activity.max_capacity}
                        </Typography>
                      </Box>

                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 0.5 }}
                      >
                        <TimeIcon fontSize="small" color="action" />
                        <Typography variant="body2" color="textSecondary">
                          {activity.duration} min
                        </Typography>
                      </Box>
                    </Box>
                  </Stack>
                </CardContent>

                <CardActions sx={{ p: 2, pt: 0 }}>
                  <Button
                    size="medium"
                    variant="contained"
                    fullWidth
                    onClick={(e) => {
                      e.stopPropagation();
                      handleActivityClick(activity.id);
                    }}
                    sx={{
                      py: 1.2,
                      fontWeight: 600,
                      background:
                        "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                      "&:hover": {
                        background:
                          "linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%)",
                        transform: "translateY(-2px)",
                        boxShadow: "0 6px 16px rgba(99, 102, 241, 0.4)",
                      },
                      transition: "all 0.3s ease",
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
          Showing {activities.length} activity
          {activities.length !== 1 ? "ies" : ""}
        </Typography>
      )}
    </Box>
  );
}

export default ActivitiesList;
