import React, { useState, useEffect } from 'react';
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
} from '@mui/material';
import {
  Person as PersonIcon,
  Email as EmailIcon,
  CalendarToday as CalendarIcon,
  Group as GroupIcon,
  DirectionsRun as ActivityIcon,
} from '@mui/icons-material';
import { useAuth } from '../../contexts/AuthContext';
import { userService, authService } from '../../services/authService';
import { DashboardStats, ApiStatus, User } from '../../types';

function Dashboard(): JSX.Element {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    recentUsers: [],
  });
  const [apiStatus, setApiStatus] = useState<ApiStatus | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async (): Promise<void> => {
    setLoading(true);
    try {
      // Cargar estadísticas
      const usersResponse = await userService.getUsers(1, 5);

      setStats({
        totalUsers: usersResponse.data.pagination.total,
        recentUsers: usersResponse.data.users,
      });

      // Verificar estado de la API
      const healthResponse = await authService.healthCheck();
      setApiStatus({
        status: 'online',
        message: healthResponse.message,
      });
    } catch (error: any) {
      console.error('Error loading dashboard data:', error);
      setApiStatus({
        status: 'error',
        message: 'Failed to connect to API',
      });
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Typography variant="h4" gutterBottom>
        Welcome back, {user?.first_name}! 👋
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom>
        Here's what's happening with your sports activities platform
      </Typography>

      {/* API Status */}
      {apiStatus && (
        <Alert
          severity={apiStatus.status === 'online' ? 'success' : 'error'}
          sx={{ mb: 3 }}
        >
          API Status: {apiStatus.message}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* User Info Card */}
        <Grid item xs={12} md={6} lg={4}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Avatar
                  sx={{ width: 56, height: 56, mr: 2, bgcolor: 'primary.main' }}
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
                    label={user?.is_active ? 'Active' : 'Inactive'}
                    color={user?.is_active ? 'success' : 'default'}
                    size="small"
                    sx={{ mt: 1 }}
                  />
                </Box>
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <EmailIcon sx={{ mr: 1, color: 'text.secondary' }} />
                <Typography variant="body2">{user?.email}</Typography>
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <CalendarIcon sx={{ mr: 1, color: 'text.secondary' }} />
                <Typography variant="body2">
                  Member since {formatDate(user?.created_at || '')}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Platform Stats */}
        <Grid item xs={12} md={6} lg={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Platform Stats
              </Typography>

              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <GroupIcon sx={{ mr: 2, color: 'primary.main' }} />
                <Box>
                  <Typography variant="h4">{stats.totalUsers}</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Total Users
                  </Typography>
                </Box>
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <ActivityIcon sx={{ mr: 2, color: 'secondary.main' }} />
                <Box>
                  <Typography variant="h4">0</Typography>
                  <Typography variant="body2" color="textSecondary">
                    Activities (Coming Soon)
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
                onClick={() => (window.location.href = '/profile')}
              >
                Edit Profile
              </Button>

              <Button fullWidth variant="outlined" disabled sx={{ mb: 1 }}>
                Create Activity (Soon)
              </Button>

              <Button fullWidth variant="outlined" disabled>
                View Analytics (Soon)
              </Button>
            </CardContent>
          </Card>
        </Grid>

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
                    {stats.recentUsers.map((recentUser: User, index: number) => (
                      <React.Fragment key={recentUser.id}>
                        <ListItem>
                          <ListItemAvatar>
                            <Avatar sx={{ bgcolor: 'secondary.main' }}>
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
                            label={recentUser.is_active ? 'Active' : 'Inactive'}
                            color={recentUser.is_active ? 'success' : 'default'}
                            size="small"
                          />
                        </ListItem>
                        {index < stats.recentUsers.length - 1 && <Divider />}
                      </React.Fragment>
                    ))}
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
