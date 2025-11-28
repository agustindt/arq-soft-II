import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import {
  CalendarToday as CalendarIcon,
  LocationOn as LocationIcon,
  AttachMoney as MoneyIcon,
  People as PeopleIcon,
  Cancel as CancelIcon,
  Visibility as ViewIcon,
  EventBusy as EventBusyIcon,
  AccessTime as TimeIcon,
} from "@mui/icons-material";
import { reservationsService } from "../../services/reservationsService";
import { activitiesService } from "../../services/activitiesService";
import { formatDateTime } from "../../utils/dateUtils";
import { Reservation, Activity } from "../../types";
import { useAuth } from "../../contexts/AuthContext";
import { useApiStatus } from "../../hooks/useApiStatus";

interface ReservationWithActivity extends Reservation {
  activity?: Activity | null;
  activityDeleted?: boolean;
}

function MyReservations(): JSX.Element {
  const navigate = useNavigate();
  const { user } = useAuth();
  const healthCheckFn = useCallback(() => reservationsService.healthCheck(), []);
  const apiStatus = useApiStatus(healthCheckFn, "Reservations API");
  const [reservations, setReservations] = useState<ReservationWithActivity[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [cancelDialogOpen, setCancelDialogOpen] = useState<boolean>(false);
  const [selectedReservation, setSelectedReservation] = useState<Reservation | null>(null);
  const [cancelling, setCancelling] = useState<boolean>(false);

  useEffect(() => {
    loadReservations();
  }, []);

  const loadReservations = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      // Use 'mine' scope to always get only the current user's reservations
      // This ensures admins see only their own reservations in "My Activities"
      const reservas = await reservationsService.getReservations('mine');
      
      // Load activity details for each reservation
      const reservationsWithActivities = await Promise.all(
        reservas.map(async (reservation) => {
          try {
            const activity = await activitiesService.getActivityById(
              reservation.actividad
            );
            return { ...reservation, activity };
          } catch (err: any) {
            console.error(`Failed to load activity ${reservation.actividad}:`, err);
            // If activity not found or inactive, mark it in the reservation
            return { 
              ...reservation, 
              activity: null,
              activityDeleted: true 
            };
          }
        })
      );

      setReservations(reservationsWithActivities);
    } catch (err: any) {
      setError(err.message || "Failed to load reservations");
      console.error("Error loading reservations:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelClick = (reservation: Reservation): void => {
    setSelectedReservation(reservation);
    setCancelDialogOpen(true);
  };

  const handleCancelConfirm = async (): Promise<void> => {
    if (!selectedReservation) return;

    setCancelling(true);
    try {
      await reservationsService.deleteReservation(selectedReservation.id);
      setCancelDialogOpen(false);
      setSelectedReservation(null);
      await loadReservations();
    } catch (err: any) {
      console.error("Error cancelling reservation:", err);
      alert(err.message || "Failed to cancel reservation");
    } finally {
      setCancelling(false);
    }
  };

  const getStatusColor = (
    status: string
  ): "default" | "primary" | "secondary" | "success" | "warning" | "error" => {
    const colors: Record<string, "default" | "primary" | "secondary" | "success" | "warning" | "error"> = {
      Pendiente: "warning",
      confirmada: "success",
      cancelada: "error",
    };
    return colors[status] || "default";
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
        📅 My Reservations
      </Typography>
      <Typography variant="subtitle1" color="textSecondary" gutterBottom sx={{ mb: 3 }}>
        Manage your activity reservations
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

      {reservations.length === 0 ? (
        <Paper sx={{ p: 4, textAlign: "center" }}>
          <EventBusyIcon sx={{ fontSize: 64, color: "text.secondary", mb: 2 }} />
          <Typography variant="h6" color="textSecondary">
            No reservations found
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 1, mb: 3 }}>
            Start exploring activities and make your first reservation!
          </Typography>
          <Button variant="contained" onClick={() => navigate("/")}>
            Browse Activities
          </Button>
        </Paper>
      ) : (
        <Grid container spacing={3}>
          {reservations.map((reservation) => (
            <Grid item xs={12} md={6} key={reservation.id}>
              <Card>
                <CardContent>
                  <Box
                    sx={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "start",
                      mb: 2,
                    }}
                  >
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="h6">
                        {reservation.activity?.name || `Activity ${reservation.actividad}`}
                      </Typography>
                      {reservation.activityDeleted && (
                        <Chip
                          label="Activity Cancelled"
                          color="error"
                          size="small"
                          sx={{ mt: 1 }}
                        />
                      )}
                    </Box>
                    <Chip
                      label={reservation.status}
                      color={getStatusColor(reservation.status)}
                      size="small"
                    />
                  </Box>

                  {reservation.activity && (
                    <>
                      <Typography
                        variant="body2"
                        color="textSecondary"
                        sx={{ mb: 2 }}
                      >
                        {reservation.activity.description}
                      </Typography>

                      <Stack spacing={1.5}>
                        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                          <CalendarIcon fontSize="small" color="action" />
                          <Typography variant="body2">
                            <strong>Date:</strong> {formatDateTime(reservation.date)}
                          </Typography>
                        </Box>

                        {reservation.schedule && (
                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <TimeIcon fontSize="small" color="action" />
                            <Typography variant="body2">
                              <strong>Schedule:</strong> {reservation.schedule}
                            </Typography>
                          </Box>
                        )}

                        {reservation.activity.location && (
                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <LocationIcon fontSize="small" color="action" />
                            <Typography variant="body2">
                              <strong>Location:</strong> {reservation.activity.location}
                            </Typography>
                          </Box>
                        )}

                        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                          <PeopleIcon fontSize="small" color="action" />
                          <Typography variant="body2">
                            <strong>Participants:</strong> {reservation.cupo}
                          </Typography>
                        </Box>

                        {reservation.activity.price && (
                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <MoneyIcon fontSize="small" color="action" />
                            <Typography variant="body2">
                              <strong>Total:</strong> $
                              {(reservation.activity.price * reservation.cupo).toFixed(2)}
                            </Typography>
                          </Box>
                        )}

                        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                          <Typography variant="caption" color="textSecondary">
                            Reserved on: {formatDateTime(reservation.created_at)}
                          </Typography>
                        </Box>
                      </Stack>
                    </>
                  )}

                  {!reservation.activity && (
                    <Alert 
                      severity={reservation.activityDeleted ? "warning" : "info"} 
                      sx={{ mt: 2 }}
                    >
                      {reservation.activityDeleted 
                        ? "⚠️ This activity has been cancelled or removed by the administrator" 
                        : "Activity details not available"}
                    </Alert>
                  )}

                  <Box sx={{ display: "flex", gap: 1, mt: 3 }}>
                    {reservation.activity && (
                      <Button
                        size="small"
                        variant="outlined"
                        startIcon={<ViewIcon />}
                        onClick={() => navigate(`/activities/${reservation.activity?.id}`)}
                      >
                        View Activity
                      </Button>
                    )}
                    {reservation.status !== "cancelada" && (
                      <Button
                        size="small"
                        variant="outlined"
                        color="error"
                        startIcon={<CancelIcon />}
                        onClick={() => handleCancelClick(reservation)}
                      >
                        Cancel
                      </Button>
                    )}
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Cancel Confirmation Dialog */}
      <Dialog
        open={cancelDialogOpen}
        onClose={() => setCancelDialogOpen(false)}
      >
        <DialogTitle>Cancel Reservation</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to cancel this reservation? This action cannot be undone.
          </Typography>
          {selectedReservation && (
            <Box sx={{ mt: 2 }}>
              <Typography variant="body2" color="textSecondary">
                Activity ID: {selectedReservation.actividad}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                Date: {formatDateTime(selectedReservation.date)}
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)}>Keep Reservation</Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleCancelConfirm}
            disabled={cancelling}
          >
            {cancelling ? "Cancelling..." : "Cancel Reservation"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default MyReservations;

