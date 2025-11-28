import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Box, Typography, Grid, Card, CardContent, Stack, Alert } from "@mui/material";
import {
  People as PeopleIcon,
  DirectionsRun as ActivityIcon,
  Security as SecurityIcon,
  Settings as SettingsIcon,
} from "@mui/icons-material";
import { useAuth } from "../contexts/AuthContext";
import { adminService, SystemStats } from "../services/adminService";
import UserManagement from "../components/Admin/UserManagement";
import { activitiesService } from "../services/activitiesService";

function RootDashboard(): JSX.Element {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [userStats, setUserStats] = useState<SystemStats | null>(null);
  const [activityStats, setActivityStats] = useState({ total: 0, active: 0 });
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    // Check if user is root
    if (user?.role !== "root") {
      navigate("/dashboard");
      return;
    }

    loadStats();
  }, [user, navigate]);

  const loadStats = async (): Promise<void> => {
    setLoading(true);
    try {
      // Load user statistics
      const userStatsResponse = await adminService.getSystemStats();
      setUserStats(userStatsResponse.data);

      // Load activity statistics
      const allActivities = await activitiesService.getAllActivities();
      const activeActivities = allActivities.activities.filter(
        (a) => a.is_active
      );
      setActivityStats({
        total: allActivities.count,
        active: activeActivities.length,
      });
    } catch (err) {
      console.error("Error loading stats:", err);
    } finally {
      setLoading(false);
    }
  };

  if (user?.role !== "root") {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">Access Denied - Root privileges required</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom sx={{ color: "error.main", fontWeight: "bold" }}>
        🔒 Root Control Panel
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Full system administration and control
      </Typography>

      {/* System Statistics */}
      {userStats && (
        <>
          <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
            📊 System Statistics
          </Typography>
          <Grid container spacing={3} sx={{ mb: 4 }}>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ border: "1px solid", borderColor: "error.main" }}>
                <CardContent>
                  <Stack direction="row" spacing={2} alignItems="center">
                    <PeopleIcon sx={{ fontSize: 40, color: "error.main" }} />
                    <Box>
                      <Typography variant="h4">{userStats.total_users}</Typography>
                      <Typography variant="body2" color="textSecondary">
                        Total Users
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
                    <ActivityIcon sx={{ fontSize: 40, color: "primary.main" }} />
                    <Box>
                      <Typography variant="h4">{activityStats.total}</Typography>
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
                    <SecurityIcon sx={{ fontSize: 40, color: "warning.main" }} />
                    <Box>
                      <Typography variant="h4">
                        {userStats.users_by_role.admin + userStats.users_by_role.root}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        Admin + Root
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
                    <SettingsIcon sx={{ fontSize: 40, color: "info.main" }} />
                    <Box>
                      <Typography variant="h4">{userStats.recent_registrations}</Typography>
                      <Typography variant="body2" color="textSecondary">
                        New Users (7d)
                      </Typography>
                    </Box>
                  </Stack>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </>
      )}

      {/* User Management Section */}
      <Typography variant="h6" gutterBottom sx={{ mb: 2, mt: 4 }}>
        👥 User Management
      </Typography>
      <UserManagement />
    </Box>
  );
}

export default RootDashboard;
