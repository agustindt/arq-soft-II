import React, { useState, ChangeEvent, FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  Container,
  CircularProgress,
  Divider,
} from "@mui/material";
import { useAuth } from "../../contexts/AuthContext";
import { FormErrors } from "../../types";

interface FormData {
  email: string;
  password: string;
}

function Login(): JSX.Element {
  const [formData, setFormData] = useState<FormData>({
    email: "",
    password: "",
  });
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { login, loading, error, clearError } = useAuth();
  const navigate = useNavigate();

  // Manejar cambios en los inputs
  const handleChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));

    // Limpiar errores cuando el usuario empiece a escribir
    if (formErrors[name]) {
      setFormErrors((prev) => ({
        ...prev,
        [name]: "",
      }));
    }

    // Limpiar error general
    if (error) {
      clearError();
    }
  };

  // Validar formulario
  const validateForm = (): boolean => {
    const errors: FormErrors = {};

    if (!formData.email.trim()) {
      errors.email = "Email is required";
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      errors.email = "Email is invalid";
    }

    if (!formData.password) {
      errors.password = "Password is required";
    } else if (formData.password.length < 6) {
      errors.password = "Password must be at least 6 characters";
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Manejar envío del formulario
  const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const result = await login(formData.email, formData.password);

    if (result.success) {
      navigate("/my-activities");
    }
  };

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
        position: "relative",
        overflow: "hidden",
        "&::before": {
          content: '""',
          position: "absolute",
          width: "100%",
          height: "100%",
          background: "radial-gradient(circle at 20% 50%, rgba(120, 119, 198, 0.3) 0%, transparent 50%), radial-gradient(circle at 80% 80%, rgba(139, 92, 246, 0.3) 0%, transparent 50%)",
          pointerEvents: "none",
        },
      }}
    >
      <Container component="main" maxWidth="sm" sx={{ position: "relative", zIndex: 1 }}>
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            py: 4,
          }}
        >
          <Card
            sx={{
              width: "100%",
              borderRadius: 4,
              boxShadow: "0 20px 60px rgba(0, 0, 0, 0.3)",
              overflow: "hidden",
            }}
          >
            <Box
              sx={{
                background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                p: 3,
                textAlign: "center",
                color: "white",
              }}
            >
              <Typography
                component="h1"
                variant="h4"
                gutterBottom
                sx={{ fontWeight: 700, mb: 1 }}
              >
                🏃‍♀️ Sports Activities
              </Typography>
              <Typography variant="h6" sx={{ opacity: 0.9, fontWeight: 400 }}>
                Welcome back! Sign in to continue
              </Typography>
            </Box>
            <CardContent sx={{ p: 4 }}>

            {error && (
              <Alert severity="error" sx={{ mt: 2, mb: 2 }}>
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 3 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                autoFocus
                value={formData.email}
                onChange={handleChange}
                error={Boolean(formErrors.email)}
                helperText={formErrors.email}
                disabled={loading}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
                value={formData.password}
                onChange={handleChange}
                error={Boolean(formErrors.password)}
                helperText={formErrors.password}
                disabled={loading}
              />

              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{
                  mt: 3,
                  mb: 2,
                  py: 1.5,
                  fontSize: "1rem",
                  fontWeight: 600,
                  background: "linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)",
                  "&:hover": {
                    background: "linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%)",
                    transform: "translateY(-2px)",
                    boxShadow: "0 8px 20px rgba(99, 102, 241, 0.4)",
                  },
                  transition: "all 0.3s ease",
                }}
                disabled={loading}
                size="large"
              >
                {loading ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  "Sign In"
                )}
              </Button>

              <Divider sx={{ my: 3 }}>
                <Typography variant="body2" color="textSecondary">
                  OR
                </Typography>
              </Divider>

              <Box sx={{ textAlign: "center" }}>
                <Typography variant="body2" color="textSecondary">
                  Don't have an account?{" "}
                  <Link
                    to="/register"
                    style={{
                      textDecoration: "none",
                      color: "#6366f1",
                      fontWeight: 600,
                      transition: "color 0.2s",
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.color = "#4f46e5")}
                    onMouseLeave={(e) => (e.currentTarget.style.color = "#6366f1")}
                  >
                    Sign Up
                  </Link>
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>

        {/* Demo credentials */}
        <Card
          sx={{
            width: "100%",
            mt: 3,
            background: "rgba(255, 255, 255, 0.98)",
            backdropFilter: "blur(10px)",
            borderRadius: 3,
            boxShadow: "0 8px 32px rgba(0, 0, 0, 0.15)",
            border: "2px solid rgba(251, 191, 36, 0.3)",
          }}
        >
          <CardContent sx={{ p: 3 }}>
            <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
              <Box
                sx={{
                  width: 48,
                  height: 48,
                  borderRadius: 2,
                  background: "linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  mr: 2,
                  boxShadow: "0 4px 12px rgba(251, 191, 36, 0.4)",
                }}
              >
                <Typography variant="h4">💡</Typography>
              </Box>
              <Typography variant="h6" sx={{ fontWeight: 700, color: "#1e293b" }}>
                Demo Credentials
              </Typography>
            </Box>
            <Box
              sx={{
                p: 3,
                borderRadius: 2,
                background: "linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%)",
                mb: 2,
                border: "1px solid rgba(251, 191, 36, 0.2)",
              }}
            >
              <Typography variant="subtitle1" sx={{ fontWeight: 700, mb: 2, color: "#92400e" }}>
                👤 Admin User
              </Typography>
              <Box sx={{ display: "flex", flexDirection: "column", gap: 1.5 }}>
                <Box>
                  <Typography variant="body2" sx={{ fontWeight: 600, mb: 0.5, color: "#78350f" }}>
                    Email:
                  </Typography>
                  <Typography
                    variant="body1"
                    sx={{
                      fontFamily: "monospace",
                      fontWeight: 600,
                      color: "#1e293b",
                      p: 1,
                      bgcolor: "white",
                      borderRadius: 1,
                      border: "1px solid rgba(0,0,0,0.1)",
                    }}
                  >
                    admin@example.com
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="body2" sx={{ fontWeight: 600, mb: 0.5, color: "#78350f" }}>
                    Password:
                  </Typography>
                  <Typography
                    variant="body1"
                    sx={{
                      fontFamily: "monospace",
                      fontWeight: 600,
                      color: "#1e293b",
                      p: 1,
                      bgcolor: "white",
                      borderRadius: 1,
                      border: "1px solid rgba(0,0,0,0.1)",
                    }}
                  >
                    password
                  </Typography>
                </Box>
              </Box>
            </Box>
            <Typography
              variant="body2"
              color="textSecondary"
              sx={{ fontSize: "0.875rem", textAlign: "center", fontStyle: "italic" }}
            >
              Or create a new account using the "Sign Up" button above
            </Typography>
          </CardContent>
        </Card>
      </Box>
    </Container>
    </Box>
  );
}

export default Login;
