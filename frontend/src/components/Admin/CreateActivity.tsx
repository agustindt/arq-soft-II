import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Alert,
  Paper,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Stack,
  IconButton,
} from "@mui/material";
import {
  ArrowBack as ArrowBackIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
} from "@mui/icons-material";
import { activitiesService } from "../../services/activitiesService";
import { Activity, ACTIVITY_CATEGORIES, DIFFICULTY_LEVELS } from "../../types";
import { useAuth } from "../../contexts/AuthContext";
import { useApiStatus } from "../../hooks/useApiStatus";

function CreateActivity(): JSX.Element {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const apiStatus = useApiStatus(() => activitiesService.healthCheck(), "Activities API");
  const isEditMode = !!id;

  const [loading, setLoading] = useState<boolean>(false);
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const [formData, setFormData] = useState<Partial<Activity>>({
    name: "",
    description: "",
    category: "",
    difficulty: "beginner",
    location: "",
    price: 0,
    duration: 60,
    max_capacity: 10,
    instructor: "",
    schedule: [],
    equipment: [],
    image_url: "",
    is_active: true,
  });

  const [scheduleItem, setScheduleItem] = useState<string>("");
  const [equipmentItem, setEquipmentItem] = useState<string>("");

  useEffect(() => {
    if (user?.role !== "admin") {
      navigate("/dashboard");
      return;
    }

    if (isEditMode && id) {
      loadActivity();
    }
  }, [id, user, navigate, isEditMode]);

  const loadActivity = async (): Promise<void> => {
    if (!id) return;

    setLoading(true);
    setError(null);
    try {
      const activity = await activitiesService.getActivityById(id);
      setFormData({
        name: activity.name,
        description: activity.description,
        category: activity.category,
        difficulty: activity.difficulty,
        location: activity.location,
        price: activity.price,
        duration: activity.duration,
        max_capacity: activity.max_capacity,
        instructor: activity.instructor || "",
        schedule: activity.schedule || [],
        equipment: activity.equipment || [],
        image_url: activity.image_url || "",
        is_active: activity.is_active,
      });
    } catch (err: any) {
      setError(err.message || "Failed to load activity");
      console.error("Error loading activity:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSuccess(false);

    try {
      if (isEditMode && id) {
        await activitiesService.updateActivity(id, formData);
      } else {
        await activitiesService.createActivity(formData);
      }

      setSuccess(true);
      setTimeout(() => {
        navigate("/admin/activities");
      }, 1500);
    } catch (err: any) {
      setError(err.response?.data?.details || err.message || "Failed to save activity");
      console.error("Error saving activity:", err);
    } finally {
      setSaving(false);
    }
  };

  const handleAddSchedule = (): void => {
    if (scheduleItem.trim()) {
      setFormData({
        ...formData,
        schedule: [...(formData.schedule || []), scheduleItem.trim()],
      });
      setScheduleItem("");
    }
  };

  const handleRemoveSchedule = (index: number): void => {
    const newSchedule = [...(formData.schedule || [])];
    newSchedule.splice(index, 1);
    setFormData({ ...formData, schedule: newSchedule });
  };

  const handleAddEquipment = (): void => {
    if (equipmentItem.trim()) {
      setFormData({
        ...formData,
        equipment: [...(formData.equipment || []), equipmentItem.trim()],
      });
      setEquipmentItem("");
    }
  };

  const handleRemoveEquipment = (index: number): void => {
    const newEquipment = [...(formData.equipment || [])];
    newEquipment.splice(index, 1);
    setFormData({ ...formData, equipment: newEquipment });
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
    <Box sx={{ flexGrow: 1, maxWidth: 1200, mx: "auto" }}>
      <Button
        startIcon={<ArrowBackIcon />}
        onClick={() => navigate("/admin/activities")}
        sx={{ mb: 2 }}
      >
        Back to Activities
      </Button>

      <Typography variant="h4" gutterBottom>
        {isEditMode ? "Edit Activity" : "Create New Activity"}
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

      {success && (
        <Alert severity="success" sx={{ mb: 3 }}>
          Activity {isEditMode ? "updated" : "created"} successfully! Redirecting...
        </Alert>
      )}

      <Paper sx={{ p: 3 }}>
        <form onSubmit={handleSubmit}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Activity Name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <FormControl fullWidth required>
                <InputLabel>Category</InputLabel>
                <Select
                  value={formData.category}
                  label="Category"
                  onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                >
                  {ACTIVITY_CATEGORIES.map((cat) => (
                    <MenuItem key={cat} value={cat}>
                      {cat.charAt(0).toUpperCase() + cat.slice(1)}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                multiline
                rows={4}
                required
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <FormControl fullWidth required>
                <InputLabel>Difficulty</InputLabel>
                <Select
                  value={formData.difficulty}
                  label="Difficulty"
                  onChange={(e) => setFormData({ ...formData, difficulty: e.target.value as any })}
                >
                  {DIFFICULTY_LEVELS.map((diff) => (
                    <MenuItem key={diff} value={diff}>
                      {diff.charAt(0).toUpperCase() + diff.slice(1)}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Location"
                value={formData.location}
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                required
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                label="Price"
                type="number"
                value={formData.price}
                onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
                inputProps={{ min: 0, step: 0.01 }}
                required
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                label="Duration (minutes)"
                type="number"
                value={formData.duration}
                onChange={(e) => setFormData({ ...formData, duration: parseInt(e.target.value) || 0 })}
                inputProps={{ min: 1 }}
                required
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                label="Max Capacity"
                type="number"
                value={formData.max_capacity}
                onChange={(e) => setFormData({ ...formData, max_capacity: parseInt(e.target.value) || 0 })}
                inputProps={{ min: 1 }}
                required
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Instructor (optional)"
                value={formData.instructor}
                onChange={(e) => setFormData({ ...formData, instructor: e.target.value })}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Image URL (optional)"
                value={formData.image_url}
                onChange={(e) => setFormData({ ...formData, image_url: e.target.value })}
              />
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle2" gutterBottom>
                Schedule
              </Typography>
              <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
                <TextField
                  placeholder="e.g., Monday 18:00"
                  value={scheduleItem}
                  onChange={(e) => setScheduleItem(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === "Enter") {
                      e.preventDefault();
                      handleAddSchedule();
                    }
                  }}
                  sx={{ flexGrow: 1 }}
                />
                <Button onClick={handleAddSchedule} startIcon={<AddIcon />}>
                  Add
                </Button>
              </Stack>
              <Stack direction="row" spacing={1} flexWrap="wrap">
                {(formData.schedule || []).map((item, index) => (
                  <Chip
                    key={index}
                    label={item}
                    onDelete={() => handleRemoveSchedule(index)}
                  />
                ))}
              </Stack>
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle2" gutterBottom>
                Required Equipment
              </Typography>
              <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
                <TextField
                  placeholder="e.g., Tennis racket"
                  value={equipmentItem}
                  onChange={(e) => setEquipmentItem(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === "Enter") {
                      e.preventDefault();
                      handleAddEquipment();
                    }
                  }}
                  sx={{ flexGrow: 1 }}
                />
                <Button onClick={handleAddEquipment} startIcon={<AddIcon />}>
                  Add
                </Button>
              </Stack>
              <Stack direction="row" spacing={1} flexWrap="wrap">
                {(formData.equipment || []).map((item, index) => (
                  <Chip
                    key={index}
                    label={item}
                    onDelete={() => handleRemoveEquipment(index)}
                  />
                ))}
              </Stack>
            </Grid>

            <Grid item xs={12}>
              <Stack direction="row" spacing={2}>
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  disabled={saving}
                  startIcon={saving ? <CircularProgress size={20} /> : null}
                >
                  {saving ? "Saving..." : isEditMode ? "Update Activity" : "Create Activity"}
                </Button>
                <Button
                  variant="outlined"
                  onClick={() => navigate("/admin/activities")}
                >
                  Cancel
                </Button>
              </Stack>
            </Grid>
          </Grid>
        </form>
      </Paper>
    </Box>
  );
}

export default CreateActivity;

