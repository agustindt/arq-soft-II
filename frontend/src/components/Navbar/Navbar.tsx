import React, { useState, MouseEvent } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Menu,
  MenuItem,
  Avatar,
  Box,
  IconButton,
  Divider,
} from "@mui/material";
import {
  AccountCircle as AccountCircleIcon,
  Dashboard as DashboardIcon,
  Logout as LogoutIcon,
  Person as PersonIcon,
  DirectionsRun as ActivitiesIcon,
  Search as SearchIcon,
  BookOnline as BookIcon,
  AdminPanelSettings as AdminIcon,
} from "@mui/icons-material";
import { useAuth } from "../../contexts/AuthContext";
import { getUserRoleFromToken } from "../../utils/jwtUtils";

function Navbar(): JSX.Element {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const { user, logout, token } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const roleFromToken = token ? getUserRoleFromToken(token) : null;
  const isAdmin =
    roleFromToken === "admin" ||
    roleFromToken === "root" ||
    roleFromToken === "super_admin" ||
    user?.role === "admin";

  const handleMenuOpen = (event: MouseEvent<HTMLButtonElement>): void => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = (): void => {
    setAnchorEl(null);
  };

  const handleLogout = (): void => {
    logout();
    handleMenuClose();
    navigate("/login");
  };

  const handleNavigation = (path: string): void => {
    navigate(path);
    handleMenuClose();
  };

  const isMenuOpen = Boolean(anchorEl);

  return (
    <AppBar position="sticky" elevation={0}>
      <Toolbar sx={{ px: { xs: 2, sm: 3 } }}>
        <Typography
          variant="h6"
          sx={{ flexGrow: 1, fontWeight: "bold", cursor: "pointer" }}
          onClick={() => navigate("/dashboard")}
        >
          FitLife
        </Typography>

        <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
          <Button
            color="inherit"
            startIcon={<ActivitiesIcon />}
            onClick={() => navigate("/activities")}
            sx={{
              backgroundColor:
                location.pathname === "/activities"
                  ? "rgba(255,255,255,0.2)"
                  : "transparent",
              borderRadius: 2,
              px: 2,
              py: 1,
              transition: "all 0.2s",
              "&:hover": {
                backgroundColor: "rgba(255,255,255,0.15)",
                transform: "translateY(-1px)",
              },
            }}
          >
            Activities
          </Button>

          <Button
            color="inherit"
            startIcon={<SearchIcon />}
            onClick={() => navigate("/search")}
            sx={{
              backgroundColor:
                location.pathname === "/search"
                  ? "rgba(255,255,255,0.2)"
                  : "transparent",
              borderRadius: 2,
              px: 2,
              py: 1,
              transition: "all 0.2s",
              "&:hover": {
                backgroundColor: "rgba(255,255,255,0.15)",
                transform: "translateY(-1px)",
              },
            }}
          >
            Search
          </Button>

          <Button
            color="inherit"
            startIcon={<BookIcon />}
            onClick={() => navigate("/my-activities")}
            sx={{
              backgroundColor:
                location.pathname === "/my-activities"
                  ? "rgba(255,255,255,0.2)"
                  : "transparent",
              borderRadius: 2,
              px: 2,
              py: 1,
              transition: "all 0.2s",
              "&:hover": {
                backgroundColor: "rgba(255,255,255,0.15)",
                transform: "translateY(-1px)",
              },
            }}
          >
            My Activities
          </Button>

          {isAdmin && (
            <Button
              color="inherit"
              startIcon={<AdminIcon />}
              onClick={() => navigate("/admin")}
              sx={{
                backgroundColor:
                  location.pathname.startsWith("/admin")
                    ? "rgba(255,255,255,0.2)"
                    : "transparent",
                borderRadius: 2,
                px: 2,
                py: 1,
                transition: "all 0.2s",
                "&:hover": {
                  backgroundColor: "rgba(255,255,255,0.15)",
                  transform: "translateY(-1px)",
                },
              }}
            >
              Admin
            </Button>
          )}

          <Button
            color="inherit"
            startIcon={<DashboardIcon />}
            onClick={() => navigate("/dashboard")}
            sx={{
              backgroundColor:
                location.pathname === "/dashboard"
                  ? "rgba(255,255,255,0.2)"
                  : "transparent",
              borderRadius: 2,
              px: 2,
              py: 1,
              transition: "all 0.2s",
              "&:hover": {
                backgroundColor: "rgba(255,255,255,0.15)",
                transform: "translateY(-1px)",
              },
            }}
          >
            Dashboard
          </Button>

          {/* Perfil */}
          <IconButton
            color="inherit"
            onClick={handleMenuOpen}
            size="large"
            sx={{ ml: 1 }}
          >
            <Avatar sx={{ bgcolor: "secondary.main" }}>
              <AccountCircleIcon />
            </Avatar>
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={isMenuOpen}
            onClose={handleMenuClose}
            anchorOrigin={{
              vertical: "bottom",
              horizontal: "right",
            }}
            transformOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            sx={{ mt: 1 }}
          >
            <MenuItem onClick={() => handleNavigation("/profile")}>
              <PersonIcon fontSize="small" sx={{ mr: 1 }} />
              Profile
            </MenuItem>
            <MenuItem onClick={() => handleNavigation("/extended-profile")}>
              <PersonIcon fontSize="small" sx={{ mr: 1 }} />
              Extended Profile
            </MenuItem>

            <Divider />

            <MenuItem onClick={handleLogout}>
              <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  );
}

export default Navbar;
