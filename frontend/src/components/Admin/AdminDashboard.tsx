import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  CircularProgress,
  Alert,
  Paper,
  Stack,
} from "@mui/material";
import {
  Add as AddIcon,
  List as ListIcon,
  People as PeopleIcon,
  DirectionsRun as ActivityIcon,
  BookOnline as ReservationIcon,
} from "@mui/icons-material";
import { activitiesService } from "../../services/activitiesService";
import { useAuth } from "../../contexts/AuthContext";
import { useApiStatus } from "../../hooks/useApiStatus";

function AdminDashboard(): JSX.Element {
  const navigate = useNavigate();
  const { user } = useAuth();
  const apiStatus = useApiStatus(() => activitiesService.healthCheck(), "Activities API");
  const [stats, setStats] = useState({
    totalActivities: 0,
    activeActivities: 0,
    inactiveActivities: 0,
  });
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Check if user is admin
    if (user?.role !== "admin") {
      navigate("/dashboard");
      return;
    }

    loadStats();
  }, [user, navigate]);

  const loadStats = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const allActivities = await activitiesService.getAllActivities();
      const activeActivities = allActivities.activities.filter(
        (a) => a.is_active
      );

      setStats({
        totalActivities: allActivities.count,
        activeActivities: activeActivities.length,
        inactiveActivities: allActivities.count - activeActivities.length,
      });
    } catch (err: any) {
      setError(err.message || "Failed to load statistics");
      console.error("Error loading stats:", err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (user?.role !== "admin") {
    return (
      <Alert severity="error">
        You don't have permission to access this page.
      </Alert>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom>
        👨‍💼 Admin Dashboard
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Manage activities, reservations, and platform settings
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

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Quick Actions */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <ActivityIcon sx={{ fontSize: 40, color: "primary.main" }} />
                <Box>
                  <Typography variant="h4">{stats.totalActivities}</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Total Activities
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <ActivityIcon sx={{ fontSize: 40, color: "success.main" }} />
                <Box>
                  <Typography variant="h4">{stats.activeActivities}</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Active Activities
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <ActivityIcon sx={{ fontSize: 40, color: "warning.main" }} />
                <Box>
                  <Typography variant="h4">{stats.inactiveActivities}</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Inactive Activities
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <ReservationIcon sx={{ fontSize: 40, color: "secondary.main" }} />
                <Box>
                  <Typography variant="h4">-</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Total Reservations
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Management Actions */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Activity Management
            </Typography>
            <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
              Create, edit, and manage activities
            </Typography>
            <Stack spacing={2}>
              <Button
                variant="contained"
                fullWidth
                startIcon={<AddIcon />}
                onClick={() => navigate("/admin/activities/new")}
              >
                Create New Activity
              </Button>
              <Button
                variant="outlined"
                fullWidth
                startIcon={<ListIcon />}
                onClick={() => navigate("/admin/activities")}
              >
                Manage Activities
              </Button>
            </Stack>
          </Paper>
        </Grid>

        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              User Management
            </Typography>
            <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
              View and manage users (Coming soon)
            </Typography>
            <Button
              variant="outlined"
              fullWidth
              startIcon={<PeopleIcon />}
              disabled
            >
              Manage Users (Coming Soon)
            </Button>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
}

export default AdminDashboard;

