import React, { ReactNode } from "react";
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
import { AdminDashboard, ActivityManagement, CreateActivity } from "./components/Admin";
import { AuthProvider, useAuth } from "./contexts/AuthContext";

// Tema personalizado para la app
const theme = createTheme({
  palette: {
    primary: {
      main: "#2196f3",
    },
    secondary: {
      main: "#f50057",
    },
    background: {
      default: "#f5f5f5",
    },
  },
});

// Props para rutas protegidas
interface RouteProps {
  children: ReactNode;
}

// Componente para rutas protegidas
function ProtectedRoute({ children }: RouteProps): JSX.Element {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
}

// Componente para rutas públicas (solo accesibles si NO estás logueado)
function PublicRoute({ children }: RouteProps): JSX.Element {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  return !isAuthenticated ? <>{children}</> : <Navigate to="/dashboard" />;
}

function AppContent(): JSX.Element {
  const { isAuthenticated } = useAuth();

  return (
    <Box sx={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
      {isAuthenticated && <Navbar />}
      <Box component="main" sx={{ flexGrow: 1, p: isAuthenticated ? 3 : 0 }}>
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
          <Route
            path="/admin"
            element={
              <ProtectedRoute>
                <AdminDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/activities"
            element={
              <ProtectedRoute>
                <ActivityManagement />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/activities/new"
            element={
              <ProtectedRoute>
                <CreateActivity />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/activities/:id/edit"
            element={
              <ProtectedRoute>
                <CreateActivity />
              </ProtectedRoute>
            }
          />

          {/* Redirecciones */}
          <Route
            path="/"
            element={
              <Navigate to={isAuthenticated ? "/activities" : "/login"} />
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
