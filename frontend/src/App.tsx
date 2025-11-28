import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { Box } from "@mui/material";

import Login from "./components/Login/Login";
import Register from "./components/Register/Register";
import Dashboard from "./components/Dashboard/Dashboard";
import Profile from "./components/Profile/Profile";
import ExtendedProfile from "./components/ExtendedProfile/ExtendedProfile";
import Navbar from "./components/Navbar/Navbar";
import { ActivitiesList } from "./components/Home";
import { ActivityDetails } from "./components/ActivityDetails";
import { SearchPage } from "./components/Search";
import { MyReservations } from "./components/MyActivities";
import {
  AdminDashboard,
  ActivityManagement,
  CreateActivity,
  UserManagement,
} from "./components/Admin";
import { AuthProvider, useAuth } from "./contexts/AuthContext";
import { getUserRoleFromToken } from "./utils/jwtUtils";
import { ReactNode } from "react";

// Tema personalizado moderno para la app
const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#6366f1",
      light: "#818cf8",
      dark: "#4f46e5",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#ec4899",
      light: "#f472b6",
      dark: "#db2777",
      contrastText: "#ffffff",
    },
    background: {
      default: "#f8fafc",
      paper: "#ffffff",
    },
    text: {
      primary: "#1e293b",
      secondary: "#64748b",
    },
  },
});

interface RouteProps {
  children: ReactNode;
}

// Rutas protegidas
function ProtectedRoute({ children }: RouteProps): JSX.Element {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
}

// Rutas públicas
function PublicRoute({ children }: RouteProps): JSX.Element {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  return !isAuthenticated ? <>{children}</> : <Navigate to="/my-activities" />;
}

// Nueva protección para ADMIN
function RequireAdmin({ children }: RouteProps): JSX.Element {
  const { isAuthenticated, loading, token } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  const role = token ? getUserRoleFromToken(token) : null;
  const isAdmin =
    role === "admin" || role === "root" || role === "super_admin";

  if (!isAuthenticated || !isAdmin) {
    return <Navigate to="/login" />;
  }

  return <>{children}</>;
}

function AppContent(): JSX.Element {
  const { isAuthenticated } = useAuth();

  return (
    <Box sx={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
      {isAuthenticated && <Navbar />}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: isAuthenticated ? { xs: 2, sm: 3, md: 4 } : 0,
          backgroundColor: "background.default",
        }}
      >
        <Routes>
          {/* Rutas públicas */}
          <Route
            path="/login"
            element={
              <PublicRoute>
                <Login />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute>
                <Register />
              </PublicRoute>
            }
          />

          {/* Rutas protegidas */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/extended-profile"
            element={
              <ProtectedRoute>
                <ExtendedProfile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/activities"
            element={
              <ProtectedRoute>
                <ActivitiesList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/activities/:id"
            element={
              <ProtectedRoute>
                <ActivityDetails />
              </ProtectedRoute>
            }
          />
          <Route
            path="/search"
            element={
              <ProtectedRoute>
                <SearchPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/my-activities"
            element={
              <ProtectedRoute>
                <MyReservations />
              </ProtectedRoute>
            }
          />

          {/* ADMIN */}
          <Route
            path="/admin"
            element={
              <RequireAdmin>
                <AdminDashboard />
              </RequireAdmin>
            }
          />
          <Route
            path="/admin/activities"
            element={
              <RequireAdmin>
                <ActivityManagement />
              </RequireAdmin>
            }
          />
          <Route
            path="/admin/activities/new"
            element={
              <RequireAdmin>
                <CreateActivity />
              </RequireAdmin>
            }
          />
          <Route
            path="/admin/activities/:id/edit"
            element={
              <RequireAdmin>
                <CreateActivity />
              </RequireAdmin>
            }
          />
          <Route
            path="/admin/users"
            element={
              <RequireAdmin>
                <UserManagement />
              </RequireAdmin>
            }
          />

          {/* Redirecciones */}
          <Route
            path="/"
            element={
              <Navigate to={isAuthenticated ? "/my-activities" : "/login"} />
            }
          />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </Box>
    </Box>
  );
}

function App(): JSX.Element {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </Router>
    </ThemeProvider>
  );
}

export default App;
