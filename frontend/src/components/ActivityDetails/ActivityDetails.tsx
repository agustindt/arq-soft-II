import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  Chip,
  Button,
  CircularProgress,
  Alert,
  Paper,
  Stack,
  Divider,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Breadcrumbs,
  Link,
} from "@mui/material";
import {
  LocationOn as LocationIcon,
  AttachMoney as MoneyIcon,
  People as PeopleIcon,
  AccessTime as TimeIcon,
  ArrowBack as ArrowBackIcon,
  CalendarToday as CalendarIcon,
  BookOnline as BookIcon,
  FitnessCenter as FitnessIcon,
} from "@mui/icons-material";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import dayjs, { Dayjs } from "dayjs";
import { activitiesService } from "../../services/activitiesService";
import { reservationsService } from "../../services/reservationsService";
import { Activity } from "../../types";
import { useAuth } from "../../contexts/AuthContext";
import { useApiStatus } from "../../hooks/useApiStatus";

function ActivityDetails(): JSX.Element {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const apiStatus = useApiStatus(() => activitiesService.healthCheck(), "Activities API");
  const [activity, setActivity] = useState<Activity | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [reservationDialogOpen, setReservationDialogOpen] = useState<boolean>(false);
  const [reservationDate, setReservationDate] = useState<Dayjs | null>(dayjs());
  const [participants, setParticipants] = useState<number>(1);
  const [reserving, setReserving] = useState<boolean>(false);
  const [reservationError, setReservationError] = useState<string | null>(null);
  const [reservationSuccess, setReservationSuccess] = useState<boolean>(false);

  useEffect(() => {
    if (id) {
      loadActivity();
    }
  }, [id]);

  const loadActivity = async (): Promise<void> => {
    if (!id) return;

    setLoading(true);
    setError(null);
    try {
      const data = await activitiesService.getActivityById(id);
      setActivity(data);
    } catch (err: any) {
      setError(err.message || "Failed to load activity");
      console.error("Error loading activity:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleReservation = async (): Promise<void> => {
    if (!activity || !reservationDate) return;

    setReserving(true);
    setReservationError(null);
    setReservationSuccess(false);

    try {
      await reservationsService.createReservation({
        activityId: activity.id,
        date: reservationDate.toISOString(),
        participants: participants,
      });

      setReservationSuccess(true);
      setTimeout(() => {
        setReservationDialogOpen(false);
        setReservationSuccess(false);
        navigate("/my-activities");
      }, 2000);
    } catch (err: any) {
      setReservationError(
        err.response?.data?.message || err.message || "Failed to create reservation"
      );
    } finally {
      setReserving(false);
    }
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

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error || !activity) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error || "Activity not found"}
        </Alert>
        <Button startIcon={<ArrowBackIcon />} onClick={() => navigate("/")}>
          Back to Activities
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 2 }}>
        <Link
          component="button"
          variant="body1"
          onClick={() => navigate("/")}
          sx={{ cursor: "pointer" }}
        >
          Activities
        </Link>
        <Typography color="text.primary">{activity.name}</Typography>
      </Breadcrumbs>

      <Button
        startIcon={<ArrowBackIcon />}
        onClick={() => navigate("/")}
        sx={{ mb: 2 }}
      >
        Back
      </Button>

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
        {/* Main Content */}
        <Grid item xs={12} md={8}>
          <Card>
            {activity.image_url && (
              <Box
                component="img"
                src={activity.image_url}
                alt={activity.name}
                sx={{
                  width: "100%",
                  height: 400,
                  objectFit: "cover",
                }}
              />
            )}

            <CardContent>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "start",
                  mb: 2,
                }}
              >
                <Typography variant="h4" component="h1">
                  {activity.name}
                </Typography>
                {!activity.is_active && (
                  <Chip label="Inactive" color="default" />
                )}
              </Box>

              <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1, mb: 3 }}>
                <Chip
                  label={activity.category}
                  color={getCategoryColor(activity.category)}
                />
                <Chip
                  label={activity.difficulty}
                  color={getDifficultyColor(activity.difficulty)}
                />
              </Box>

              <Typography variant="body1" paragraph>
                {activity.description}
              </Typography>

              <Divider sx={{ my: 3 }} />

              {/* Details Grid */}
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2 }}>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <LocationIcon color="primary" />
                      <Box>
                        <Typography variant="caption" color="textSecondary">
                          Location
                        </Typography>
                        <Typography variant="body1" fontWeight="medium">
                          {activity.location}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2 }}>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <MoneyIcon color="primary" />
                      <Box>
                        <Typography variant="caption" color="textSecondary">
                          Price
                        </Typography>
                        <Typography variant="body1" fontWeight="medium">
                          ${activity.price.toFixed(2)}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2 }}>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <PeopleIcon color="primary" />
                      <Box>
                        <Typography variant="caption" color="textSecondary">
                          Max Capacity
                        </Typography>
                        <Typography variant="body1" fontWeight="medium">
                          {activity.max_capacity} participants
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2 }}>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <TimeIcon color="primary" />
                      <Box>
                        <Typography variant="caption" color="textSecondary">
                          Duration
                        </Typography>
                        <Typography variant="body1" fontWeight="medium">
                          {activity.duration} minutes
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>

                {activity.instructor && (
                  <Grid item xs={12} sm={6}>
                    <Paper sx={{ p: 2 }}>
                      <Stack direction="row" spacing={1} alignItems="center">
                        <FitnessIcon color="primary" />
                        <Box>
                          <Typography variant="caption" color="textSecondary">
                            Instructor
                          </Typography>
                          <Typography variant="body1" fontWeight="medium">
                            {activity.instructor}
                          </Typography>
                        </Box>
                      </Stack>
                    </Paper>
                  </Grid>
                )}
              </Grid>

              {activity.schedule && activity.schedule.length > 0 && (
                <>
                  <Divider sx={{ my: 3 }} />
                  <Typography variant="h6" gutterBottom>
                    Schedule
                  </Typography>
                  <Stack spacing={1}>
                    {activity.schedule.map((schedule, index) => (
                      <Chip key={index} label={schedule} variant="outlined" />
                    ))}
                  </Stack>
                </>
              )}

              {activity.equipment && activity.equipment.length > 0 && (
                <>
                  <Divider sx={{ my: 3 }} />
                  <Typography variant="h6" gutterBottom>
                    Required Equipment
                  </Typography>
                  <Stack direction="row" spacing={1} flexWrap="wrap">
                    {activity.equipment.map((equipment, index) => (
                      <Chip key={index} label={equipment} size="small" />
                    ))}
                  </Stack>
                </>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Reservation Sidebar */}
        <Grid item xs={12} md={4}>
          <Card sx={{ position: "sticky", top: 20 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Book This Activity
              </Typography>

              <Box sx={{ mb: 3 }}>
                <Typography variant="h4" color="primary" fontWeight="bold">
                  ${activity.price.toFixed(2)}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  per participant
                </Typography>
              </Box>

              <Button
                variant="contained"
                fullWidth
                size="large"
                startIcon={<BookIcon />}
                onClick={() => setReservationDialogOpen(true)}
                disabled={!activity.is_active}
                sx={{ mb: 2 }}
              >
                {activity.is_active ? "Make Reservation" : "Not Available"}
              </Button>

              {!activity.is_active && (
                <Alert severity="warning" sx={{ mt: 2 }}>
                  This activity is currently inactive
                </Alert>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Reservation Dialog */}
      <Dialog
        open={reservationDialogOpen}
        onClose={() => setReservationDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Make a Reservation</DialogTitle>
        <DialogContent>
          <Stack spacing={3} sx={{ mt: 1 }}>
            {reservationSuccess && (
              <Alert severity="success">
                Reservation created successfully! Redirecting...
              </Alert>
            )}

            {reservationError && (
              <Alert severity="error">{reservationError}</Alert>
            )}

            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DatePicker
                label="Select Date"
                value={reservationDate}
                onChange={(newValue) => setReservationDate(newValue)}
                minDate={dayjs()}
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />
            </LocalizationProvider>

            <TextField
              label="Number of Participants"
              type="number"
              value={participants}
              onChange={(e) => {
                const value = parseInt(e.target.value, 10);
                if (value > 0 && value <= activity.max_capacity) {
                  setParticipants(value);
                }
              }}
              inputProps={{
                min: 1,
                max: activity.max_capacity,
              }}
              fullWidth
              helperText={`Max: ${activity.max_capacity} participants`}
            />

            <Paper sx={{ p: 2, bgcolor: "grey.50" }}>
              <Typography variant="subtitle2" gutterBottom>
                Reservation Summary
              </Typography>
              <Box sx={{ display: "flex", justifyContent: "space-between", mb: 1 }}>
                <Typography variant="body2">Price per participant:</Typography>
                <Typography variant="body2" fontWeight="bold">
                  ${activity.price.toFixed(2)}
                </Typography>
              </Box>
              <Box sx={{ display: "flex", justifyContent: "space-between", mb: 1 }}>
                <Typography variant="body2">Participants:</Typography>
                <Typography variant="body2" fontWeight="bold">
                  {participants}
                </Typography>
              </Box>
              <Divider sx={{ my: 1 }} />
              <Box sx={{ display: "flex", justifyContent: "space-between" }}>
                <Typography variant="h6">Total:</Typography>
                <Typography variant="h6" color="primary" fontWeight="bold">
                  ${(activity.price * participants).toFixed(2)}
                </Typography>
              </Box>
            </Paper>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setReservationDialogOpen(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleReservation}
            disabled={reserving || !reservationDate || participants < 1}
            startIcon={reserving ? <CircularProgress size={20} /> : <BookIcon />}
          >
            {reserving ? "Processing..." : "Confirm Reservation"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default ActivityDetails;

