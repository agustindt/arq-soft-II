import React, { useState, ChangeEvent, FormEvent } from "react";
import {
  Box,
  Typography,
  Card,
  CardContent,
  TextField,
  Button,
  Grid,
  Avatar,
  Alert,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
} from "@mui/material";
import {
  Edit as EditIcon,
  Lock as LockIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
} from "@mui/icons-material";
import { useAuth } from "../../contexts/AuthContext";
import { formatDateTime } from "../../utils/dateUtils";
import {
  FormErrors,
  Message,
  UpdateProfileRequest,
  ChangePasswordRequest,
} from "../../types";
import { useApiStatus } from "../../hooks/useApiStatus";
import { authService } from "../../services/authService";

interface PasswordFormData extends ChangePasswordRequest {
  confirmPassword: string;
}

function Profile(): JSX.Element {
  const { user, updateProfile, changePassword } = useAuth();
  const apiStatus = useApiStatus(() => authService.healthCheck(), "Users API");
  const [editMode, setEditMode] = useState<boolean>(false);
  const [passwordDialogOpen, setPasswordDialogOpen] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<Message>({ type: "", text: "" });

  const [profileData, setProfileData] = useState<UpdateProfileRequest>({
    first_name: user?.first_name || "",
    last_name: user?.last_name || "",
  });

  const [passwordData, setPasswordData] = useState<PasswordFormData>({
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  const [passwordErrors, setPasswordErrors] = useState<FormErrors>({});

  // Manejar cambios en el perfil
  const handleProfileChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setProfileData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Manejar cambios en contraseña
  const handlePasswordChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setPasswordData((prev) => ({
      ...prev,
      [name]: value,
    }));

    // Limpiar errores
    if (passwordErrors[name]) {
      setPasswordErrors((prev) => ({
        ...prev,
        [name]: "",
      }));
    }
  };

  // Validar contraseña
  const validatePassword = (): boolean => {
    const errors: FormErrors = {};

    if (!passwordData.currentPassword) {
      errors.currentPassword = "Current password is required";
    }

    if (!passwordData.newPassword) {
      errors.newPassword = "New password is required";
    } else if (passwordData.newPassword.length < 6) {
      errors.newPassword = "Password must be at least 6 characters";
    }

    if (!passwordData.confirmPassword) {
      errors.confirmPassword = "Please confirm your password";
    } else if (passwordData.newPassword !== passwordData.confirmPassword) {
      errors.confirmPassword = "Passwords do not match";
    }

    setPasswordErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Guardar perfil
  const handleSaveProfile = async (): Promise<void> => {
    setLoading(true);
    setMessage({ type: "", text: "" });

    try {
      const result = await updateProfile(profileData);

      if (result.success) {
        setEditMode(false);
        setMessage({ type: "success", text: "Profile updated successfully!" });
      } else {
        setMessage({ type: "error", text: result.error || "Update failed" });
      }
    } catch (error: any) {
      setMessage({ type: "error", text: "Failed to update profile" });
    } finally {
      setLoading(false);
    }
  };

  // Cambiar contraseña
  const handleChangePassword = async (): Promise<void> => {
    if (!validatePassword()) {
      return;
    }

    setLoading(true);

    try {
      const result = await changePassword({
        currentPassword: passwordData.currentPassword,
        newPassword: passwordData.newPassword,
      });

      if (result.success) {
        setPasswordDialogOpen(false);
        setPasswordData({
          currentPassword: "",
          newPassword: "",
          confirmPassword: "",
        });
        setMessage({ type: "success", text: "Password changed successfully!" });
      } else {
        setMessage({
          type: "error",
          text: result.error || "Password change failed",
        });
      }
    } catch (error: any) {
      setMessage({ type: "error", text: "Failed to change password" });
    } finally {
      setLoading(false);
    }
  };

  // Cancelar edición
  const handleCancelEdit = (): void => {
    setProfileData({
      first_name: user?.first_name || "",
      last_name: user?.last_name || "",
    });
    setEditMode(false);
  };

  return (
    <Box sx={{ maxWidth: 1000, mx: "auto" }}>
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
          My Profile
        </Typography>
        <Typography variant="h6" color="textSecondary" sx={{ fontWeight: 400 }}>
          Manage your personal information and account settings
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

      {message.text && (
        <Alert
          severity={message.type as "success" | "error"}
          sx={{ mb: 3 }}
          onClose={() => setMessage({ type: "", text: "" })}
        >
          {message.text}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* Profile Information */}
        <Grid item xs={12} md={8}>
          <Card
            sx={{
              background: "linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)",
              border: "1px solid rgba(99, 102, 241, 0.1)",
            }}
          >
            <CardContent>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  mb: 3,
                }}
              >
                <Typography variant="h6">Profile Information</Typography>
                {!editMode ? (
                  <Button
                    startIcon={<EditIcon />}
                    onClick={() => setEditMode(true)}
                  >
                    Edit
                  </Button>
                ) : (
                  <Box>
                    <Button
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                      disabled={loading}
                      sx={{ mr: 1 }}
                    >
                      {loading ? <CircularProgress size={20} /> : "Save"}
                    </Button>
                    <Button
                      startIcon={<CancelIcon />}
                      onClick={handleCancelEdit}
                      disabled={loading}
                    >
                      Cancel
                    </Button>
                  </Box>
                )}
              </Box>

              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="First Name"
                    name="first_name"
                    value={editMode ? profileData.first_name : user?.first_name}
                    onChange={handleProfileChange}
                    disabled={!editMode}
                    variant={editMode ? "outlined" : "filled"}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Last Name"
                    name="last_name"
                    value={editMode ? profileData.last_name : user?.last_name}
                    onChange={handleProfileChange}
                    disabled={!editMode}
                    variant={editMode ? "outlined" : "filled"}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label="Email"
                    value={user?.email}
                    disabled
                    variant="filled"
                    helperText="Email cannot be changed"
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label="Username"
                    value={user?.username}
                    disabled
                    variant="filled"
                    helperText="Username cannot be changed"
                  />
                </Grid>
              </Grid>

              <Divider sx={{ my: 3 }} />

              <Button
                fullWidth
                startIcon={<LockIcon />}
                onClick={() => setPasswordDialogOpen(true)}
                variant="outlined"
                sx={{
                  py: 1.5,
                  fontWeight: 600,
                  borderColor: "primary.main",
                  color: "primary.main",
                  "&:hover": {
                    borderColor: "primary.dark",
                    backgroundColor: "rgba(99, 102, 241, 0.05)",
                    transform: "translateY(-2px)",
                    boxShadow: "0 4px 12px rgba(99, 102, 241, 0.2)",
                  },
                  transition: "all 0.3s ease",
                }}
              >
                Change Password
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Profile Summary */}
        <Grid item xs={12} md={4}>
          <Card
            sx={{
              background: "linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)",
              border: "1px solid rgba(99, 102, 241, 0.1)",
            }}
          >
            <CardContent sx={{ textAlign: "center" }}>
              <Avatar
                sx={{
                  width: 120,
                  height: 120,
                  mx: "auto",
                  mb: 3,
                  background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                  fontSize: "3rem",
                  fontWeight: 700,
                  boxShadow: "0 8px 24px rgba(99, 102, 241, 0.3)",
                }}
              >
                {user?.first_name?.charAt(0)?.toUpperCase()}
              </Avatar>

              <Typography variant="h6" gutterBottom>
                {user?.first_name} {user?.last_name}
              </Typography>

              <Typography variant="body2" color="textSecondary" gutterBottom>
                @{user?.username}
              </Typography>

              <Typography variant="body2" color="textSecondary" gutterBottom>
                {user?.email}
              </Typography>

              <Divider sx={{ my: 2 }} />

              <Typography variant="caption" color="textSecondary">
                Member since
              </Typography>
              <Typography variant="body2">
                {formatDateTime(user?.created_at || "")}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Change Password Dialog */}
      <Dialog
        open={passwordDialogOpen}
        onClose={() => setPasswordDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Change Password</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 1 }}>
            <TextField
              fullWidth
              label="Current Password"
              name="currentPassword"
              type="password"
              value={passwordData.currentPassword}
              onChange={handlePasswordChange}
              error={Boolean(passwordErrors.currentPassword)}
              helperText={passwordErrors.currentPassword}
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="New Password"
              name="newPassword"
              type="password"
              value={passwordData.newPassword}
              onChange={handlePasswordChange}
              error={Boolean(passwordErrors.newPassword)}
              helperText={passwordErrors.newPassword}
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Confirm New Password"
              name="confirmPassword"
              type="password"
              value={passwordData.confirmPassword}
              onChange={handlePasswordChange}
              error={Boolean(passwordErrors.confirmPassword)}
              helperText={passwordErrors.confirmPassword}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setPasswordDialogOpen(false)}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            onClick={handleChangePassword}
            variant="contained"
            disabled={loading}
          >
            {loading ? <CircularProgress size={20} /> : "Change Password"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default Profile;
