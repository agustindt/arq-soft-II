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
  Security as SecurityIcon,
} from "@mui/icons-material";
import { useAuth } from "../../contexts/AuthContext";

function Navbar(): JSX.Element {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

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
          component="div"
          sx={{
            flexGrow: 1,
            cursor: "pointer",
            display: "flex",
            alignItems: "center",
            fontWeight: 700,
            fontSize: "1.25rem",
            transition: "opacity 0.2s",
            "&:hover": {
              opacity: 0.8,
            },
          }}
          onClick={() => navigate("/activities")}
        >
          🏃‍♀️ Sports Activities
        </Typography>

        <Box sx={{ display: "flex", alignItems: "center", gap: { xs: 0.5, sm: 1 } }}>
          {/* Navigation Buttons */}
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
          {/* Admin Panel - Solo para admins */}
          {user?.role === "admin" && (
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
              Admin Panel
            </Button>
          )}

          {/* Root Panel - Solo para root */}
          {user?.role === "root" && (
            <Button
              color="inherit"
              startIcon={<SecurityIcon />}
              onClick={() => navigate("/root")}
              sx={{
                backgroundColor:
                  location.pathname.startsWith("/root")
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
                color: "error.light",
              }}
            >
              Root Panel
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

          {/* User Menu */}
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenuOpen}
            color="inherit"
            sx={{
              ml: 1,
              transition: "transform 0.2s",
              "&:hover": {
                transform: "scale(1.1)",
              },
            }}
          >
            <Avatar
              sx={{
                width: 36,
                height: 36,
                bgcolor: "rgba(255,255,255,0.2)",
                border: "2px solid rgba(255,255,255,0.3)",
                fontWeight: 600,
              }}
            >
              {user?.first_name?.charAt(0)?.toUpperCase() ||
                user?.username?.charAt(0)?.toUpperCase() ||
                "U"}
            </Avatar>
          </IconButton>

          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: "bottom",
              horizontal: "right",
            }}
            keepMounted
            transformOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            open={isMenuOpen}
            onClose={handleMenuClose}
            PaperProps={{
              sx: {
                mt: 1.5,
                minWidth: 220,
                borderRadius: 2,
                boxShadow: "0 8px 24px rgba(0, 0, 0, 0.15)",
                border: "1px solid rgba(0, 0, 0, 0.05)",
              },
            }}
          >
            <Box sx={{ px: 2, py: 1 }}>
              <Typography variant="subtitle2" color="textSecondary">
                {user?.email}
              </Typography>
              <Typography variant="body2" fontWeight="bold">
                {user?.first_name} {user?.last_name}
              </Typography>
            </Box>
            <Divider />
            <MenuItem onClick={() => handleNavigation("/profile")}>
              <AccountCircleIcon sx={{ mr: 1 }} />
              Profile
            </MenuItem>
            <MenuItem onClick={() => handleNavigation("/extended-profile")}>
              <PersonIcon sx={{ mr: 1 }} />
              Extended Profile
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <LogoutIcon sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  );
}

export default Navbar;
