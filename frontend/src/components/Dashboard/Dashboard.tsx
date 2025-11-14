import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  Avatar,
  Chip,
  Button,
  Alert,
  CircularProgress,
  Paper,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Divider,
} from "@mui/material";
import {
  Person as PersonIcon,
  Email as EmailIcon,
  CalendarToday as CalendarIcon,
  Group as GroupIcon,
  DirectionsRun as ActivityIcon,
} from "@mui/icons-material";
import { useAuth } from "../../contexts/AuthContext";
import { userService, authService } from "../../services/authService";
import { activitiesService } from "../../services/activitiesService";
import { formatDate } from "../../utils/dateUtils";
import { DashboardStats, ApiStatus, User, Activity } from "../../types";
import { useNavigate } from "react-router-dom";

function Dashboard(): JSX.Element {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    recentUsers: [],
  });
  const [activities, setActivities] = useState<Activity[]>([]);
  const [totalActivities, setTotalActivities] = useState<number>(0);
  const [apiStatus, setApiStatus] = useState<ApiStatus | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async (): Promise<void> => {
    setLoading(true);
    try {
      // Cargar estadísticas de usuarios
      const usersResponse = await userService.getUsers(1, 5);

      setStats({
        totalUsers: usersResponse.data.pagination.total,
        recentUsers: usersResponse.data.users,
      });

      // Cargar actividades
      try {
        const activitiesResponse = await activitiesService.getActivities();
        setActivities(activitiesResponse.activities.slice(0, 3)); // Mostrar solo las primeras 3
        setTotalActivities(activitiesResponse.count);
      } catch (activitiesError) {
        console.error("Error loading activities:", activitiesError);
        // No fallar si las actividades no se pueden cargar
      }

      // Verificar estado de la API (Users API para Dashboard)
      const healthResponse = await authService.healthCheck();
      setApiStatus({
        status: "online",
        message: `Users API: ${healthResponse.message}`,
      });
    } catch (error: any) {
      console.error("Error loading dashboard data:", error);
      setApiStatus({
        status: "error",
        message: "Failed to connect to API",
      });
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
          Welcome back, {user?.first_name}! 👋
        </Typography>
        <Typography variant="h6" color="textSecondary" sx={{ fontWeight: 400 }}>
          Here's what's happening with your sports activities platform
        </Typography>
      </Box>

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
        {/* User Info Card */}
        <Grid item xs={12} md={6} lg={4}>
          <Card
            sx={{
              background: "linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)",
              border: "1px solid rgba(99, 102, 241, 0.1)",
            }}
          >
            <CardContent>
              <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
                <Avatar
                  sx={{
                    width: 64,
                    height: 64,
                    mr: 2,
                    background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                    fontSize: "1.5rem",
                    fontWeight: 700,
                    boxShadow: "0 4px 12px rgba(99, 102, 241, 0.3)",
                  }}
                >
                  {user?.first_name?.charAt(0)?.toUpperCase()}
                </Avatar>
                <Box>
                  <Typography variant="h6">
                    {user?.first_name} {user?.last_name}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    @{user?.username}
                  </Typography>
                  <Chip
                    label={user?.is_active ? "Active" : "Inactive"}
                    color={user?.is_active ? "success" : "default"}
                    size="small"
                    sx={{ mt: 1 }}
                  />
                </Box>
              </Box>

              <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                <EmailIcon sx={{ mr: 1, color: "text.secondary" }} />
                <Typography variant="body2">{user?.email}</Typography>
              </Box>

              <Box sx={{ display: "flex", alignItems: "center" }}>
                <CalendarIcon sx={{ mr: 1, color: "text.secondary" }} />
                <Typography variant="body2">
                  Member since {formatDate(user?.created_at || "")}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Platform Stats */}
        <Grid item xs={12} md={6} lg={4}>
          <Card
            sx={{
              background: "linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)",
              border: "1px solid rgba(99, 102, 241, 0.1)",
            }}
          >
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 600, mb: 3 }}>
                Platform Stats
              </Typography>

              <Box
                sx={{
                  display: "flex",
                  alignItems: "center",
                  mb: 3,
                  p: 2,
                  borderRadius: 2,
                  background: "linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, rgba(139, 92, 246, 0.05) 100%)",
                }}
              >
                <Box
                  sx={{
                    width: 48,
                    height: 48,
                    borderRadius: 2,
                    background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    mr: 2,
                  }}
                >
                  <GroupIcon sx={{ color: "white" }} />
                </Box>
                <Box>
                  <Typography variant="h4" sx={{ fontWeight: 700, color: "primary.main" }}>
                    {stats.totalUsers}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Total Users
                  </Typography>
                </Box>
              </Box>

              <Box
                sx={{
                  display: "flex",
                  alignItems: "center",
                  p: 2,
                  borderRadius: 2,
                  background: "linear-gradient(135deg, rgba(236, 72, 153, 0.05) 0%, rgba(219, 39, 119, 0.05) 100%)",
                }}
              >
                <Box
                  sx={{
                    width: 48,
                    height: 48,
                    borderRadius: 2,
                    background: "linear-gradient(135deg, #ec4899 0%, #db2777 100%)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    mr: 2,
                  }}
                >
                  <ActivityIcon sx={{ color: "white" }} />
                </Box>
                <Box>
                  <Typography variant="h4" sx={{ fontWeight: 700, color: "secondary.main" }}>
                    {totalActivities}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Total Activities
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Quick Actions */}
        <Grid item xs={12} md={12} lg={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Quick Actions
              </Typography>

              <Button
                fullWidth
                variant="outlined"
                sx={{ mb: 1 }}
                onClick={() => navigate("/profile")}
              >
                Edit Profile
              </Button>

              <Button
                fullWidth
                variant="contained"
                sx={{ mb: 1 }}
                onClick={() => navigate("/activities")}
              >
                Browse Activities
              </Button>

              <Button
                fullWidth
                variant="outlined"
                sx={{ mb: 1 }}
                onClick={() => navigate("/search")}
              >
                Search Activities
              </Button>

              <Button
                fullWidth
                variant="outlined"
                onClick={() => navigate("/my-activities")}
              >
                My Reservations
              </Button>

              {user?.role === "admin" && (
                <Button
                  fullWidth
                  variant="contained"
                  color="secondary"
                  sx={{ mt: 1 }}
                  onClick={() => navigate("/admin")}
                >
                  Admin Panel
                </Button>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Featured Activities */}
        {activities.length > 0 && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    mb: 2,
                  }}
                >
                  <Typography variant="h6">Featured Activities</Typography>
                  <Button
                    size="small"
                    onClick={() => navigate("/activities")}
                  >
                    View All
                  </Button>
                </Box>

                <Grid container spacing={2}>
                  {activities.map((activity) => (
                    <Grid item xs={12} sm={6} md={4} key={activity.id}>
                      <Card
                        sx={{
                          cursor: "pointer",
                          transition: "transform 0.2s",
                          "&:hover": {
                            transform: "translateY(-2px)",
                            boxShadow: 3,
                          },
                        }}
                        onClick={() => navigate(`/activities/${activity.id}`)}
                      >
                        <CardContent>
                          <Typography variant="h6" noWrap>
                            {activity.name}
                          </Typography>
                          <Typography
                            variant="body2"
                            color="textSecondary"
                            sx={{
                              mb: 1,
                              display: "-webkit-box",
                              WebkitLineClamp: 2,
                              WebkitBoxOrient: "vertical",
                              overflow: "hidden",
                            }}
                          >
                            {activity.description}
                          </Typography>
                          <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                            <Chip
                              label={activity.category}
                              size="small"
                              color="primary"
                            />
                            <Chip
                              label={activity.difficulty}
                              size="small"
                              color={
                                activity.difficulty === "beginner"
                                  ? "success"
                                  : activity.difficulty === "intermediate"
                                  ? "warning"
                                  : "error"
                              }
                            />
                          </Box>
                          <Typography
                            variant="h6"
                            color="primary"
                            sx={{ mt: 1 }}
                          >
                            ${activity.price.toFixed(2)}
                          </Typography>
                        </CardContent>
                      </Card>
                    </Grid>
                  ))}
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Recent Users */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Recent Users
              </Typography>

              {stats.recentUsers.length > 0 ? (
                <Paper variant="outlined">
                  <List>
                    {stats.recentUsers.map(
                      (recentUser: User, index: number) => (
                        <React.Fragment key={recentUser.id}>
                          <ListItem>
                            <ListItemAvatar>
                              <Avatar sx={{ bgcolor: "secondary.main" }}>
                                {recentUser.first_name
                                  ?.charAt(0)
                                  ?.toUpperCase() ||
                                  recentUser.username?.charAt(0)?.toUpperCase()}
                              </Avatar>
                            </ListItemAvatar>
                            <ListItemText
                              primary={`${recentUser.first_name} ${recentUser.last_name}`}
                              secondary={
                                <Box>
                                  <Typography
                                    variant="body2"
                                    color="textSecondary"
                                  >
                                    @{recentUser.username} • {recentUser.email}
                                  </Typography>
                                  <Typography
                                    variant="caption"
                                    color="textSecondary"
                                  >
                                    Joined {formatDate(recentUser.created_at)}
                                  </Typography>
                                </Box>
                              }
                            />
                            <Chip
                              label={
                                recentUser.is_active ? "Active" : "Inactive"
                              }
                              color={
                                recentUser.is_active ? "success" : "default"
                              }
                              size="small"
                            />
                          </ListItem>
                          {index < stats.recentUsers.length - 1 && <Divider />}
                        </React.Fragment>
                      )
                    )}
                  </List>
                </Paper>
              ) : (
                <Typography color="textSecondary">No users found</Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}

export default Dashboard;
