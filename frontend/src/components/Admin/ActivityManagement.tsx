import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Box,
  Typography,
  Button,
  CircularProgress,
  Alert,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Visibility as ViewIcon,
  ToggleOn as ToggleOnIcon,
  ToggleOff as ToggleOffIcon,
  ArrowBack as ArrowBackIcon,
} from "@mui/icons-material";
import { activitiesService } from "../../services/activitiesService";
import { Activity } from "../../types";
import { useAuth } from "../../contexts/AuthContext";

function ActivityManagement(): JSX.Element {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [activities, setActivities] = useState<Activity[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);
  const [selectedActivity, setSelectedActivity] = useState<Activity | null>(null);
  const [deleting, setDeleting] = useState<boolean>(false);
  const [toggling, setToggling] = useState<string | null>(null);

  useEffect(() => {
    if (user?.role !== "admin") {
      navigate("/dashboard");
      return;
    }
    loadActivities();
  }, [user, navigate]);

  const loadActivities = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const response = await activitiesService.getAllActivities();
      setActivities(response.activities);
    } catch (err: any) {
      setError(err.message || "Failed to load activities");
      console.error("Error loading activities:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteClick = (activity: Activity): void => {
    setSelectedActivity(activity);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async (): Promise<void> => {
    if (!selectedActivity) return;

    setDeleting(true);
    try {
      await activitiesService.deleteActivity(selectedActivity.id);
      setDeleteDialogOpen(false);
      setSelectedActivity(null);
      await loadActivities();
    } catch (err: any) {
      console.error("Error deleting activity:", err);
      alert(err.message || "Failed to delete activity");
    } finally {
      setDeleting(false);
    }
  };

  const handleToggleStatus = async (activity: Activity): Promise<void> => {
    setToggling(activity.id);
    try {
      await activitiesService.toggleActivityStatus(activity.id);
      await loadActivities();
    } catch (err: any) {
      console.error("Error toggling activity status:", err);
      alert(err.message || "Failed to toggle activity status");
    } finally {
      setToggling(null);
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
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 3 }}>
        <Box>
          <Button
            startIcon={<ArrowBackIcon />}
            onClick={() => navigate("/admin")}
            sx={{ mb: 1 }}
          >
            Back to Dashboard
          </Button>
          <Typography variant="h4" gutterBottom>
            Activity Management
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate("/admin/activities/new")}
        >
          Create Activity
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Category</TableCell>
              <TableCell>Difficulty</TableCell>
              <TableCell>Location</TableCell>
              <TableCell>Price</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {activities.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography variant="body2" color="textSecondary" sx={{ py: 3 }}>
                    No activities found. Create your first activity!
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              activities.map((activity) => (
                <TableRow key={activity.id} hover>
                  <TableCell>
                    <Typography variant="body2" fontWeight="medium">
                      {activity.name}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip label={activity.category} size="small" />
                  </TableCell>
                  <TableCell>
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
                  </TableCell>
                  <TableCell>{activity.location}</TableCell>
                  <TableCell>${activity.price.toFixed(2)}</TableCell>
                  <TableCell>
                    <Chip
                      label={activity.is_active ? "Active" : "Inactive"}
                      color={activity.is_active ? "success" : "default"}
                      size="small"
                    />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/activities/${activity.id}`)}
                      title="View"
                    >
                      <ViewIcon fontSize="small" />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/admin/activities/${activity.id}/edit`)}
                      title="Edit"
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleToggleStatus(activity)}
                      disabled={toggling === activity.id}
                      title={activity.is_active ? "Deactivate" : "Activate"}
                    >
                      {toggling === activity.id ? (
                        <CircularProgress size={16} />
                      ) : activity.is_active ? (
                        <ToggleOnIcon fontSize="small" color="success" />
                      ) : (
                        <ToggleOffIcon fontSize="small" />
                      )}
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDeleteClick(activity)}
                      color="error"
                      title="Delete"
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
      >
        <DialogTitle>Delete Activity</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "{selectedActivity?.name}"? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleDeleteConfirm}
            disabled={deleting}
          >
            {deleting ? "Deleting..." : "Delete"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default ActivityManagement;

