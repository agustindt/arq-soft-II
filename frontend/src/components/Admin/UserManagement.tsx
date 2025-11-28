import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  CircularProgress,
  Alert,
  Stack,
  Pagination,
  Tooltip,
} from "@mui/material";
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  CheckCircle as ActiveIcon,
  Cancel as InactiveIcon,
  Add as AddIcon,
} from "@mui/icons-material";
import { adminService, AdminUser, CreateUserRequest } from "../../services/adminService";
import { useAuth } from "../../contexts/AuthContext";

const ROLES = ["user", "moderator", "admin", "super_admin", "root"];

function UserManagement(): JSX.Element {
  const { user } = useAuth();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [total, setTotal] = useState<number>(0);

  // Estados para diálogos
  const [editDialogOpen, setEditDialogOpen] = useState<boolean>(false);
  const [createDialogOpen, setCreateDialogOpen] = useState<boolean>(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);
  const [selectedUser, setSelectedUser] = useState<AdminUser | null>(null);

  // Estados para formularios
  const [newUser, setNewUser] = useState<CreateUserRequest>({
    email: "",
    username: "",
    password: "",
    first_name: "",
    last_name: "",
    role: "user",
  });
  const [editRole, setEditRole] = useState<string>("");

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const response = await adminService.getAllUsers(page, 10);
      setUsers(response.data.users);
      setTotal(response.data.total);
      setTotalPages(response.data.total_pages);
    } catch (err: any) {
      setError(err.response?.data?.message || "Error al cargar usuarios");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateUser = async (): Promise<void> => {
    try {
      await adminService.createUser(newUser);
      setSuccess("Usuario creado exitosamente");
      setCreateDialogOpen(false);
      setNewUser({
        email: "",
        username: "",
        password: "",
        first_name: "",
        last_name: "",
        role: "user",
      });
      loadUsers();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || "Error al crear usuario");
      setTimeout(() => setError(null), 5000);
    }
  };

  const handleUpdateRole = async (): Promise<void> => {
    if (!selectedUser) return;
    try {
      await adminService.updateUserRole(selectedUser.id, editRole);
      setSuccess("Rol actualizado exitosamente");
      setEditDialogOpen(false);
      loadUsers();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || "Error al actualizar rol");
      setTimeout(() => setError(null), 5000);
    }
  };

  const handleToggleStatus = async (userId: number, currentStatus: boolean): Promise<void> => {
    try {
      await adminService.updateUserStatus(userId, !currentStatus);
      setSuccess(`Usuario ${!currentStatus ? "activado" : "desactivado"} exitosamente`);
      loadUsers();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || "Error al actualizar estado");
      setTimeout(() => setError(null), 5000);
    }
  };

  const handleDeleteUser = async (): Promise<void> => {
    if (!selectedUser) return;
    try {
      await adminService.deleteUser(selectedUser.id);
      setSuccess("Usuario eliminado exitosamente");
      setDeleteDialogOpen(false);
      loadUsers();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || "Error al eliminar usuario");
      setTimeout(() => setError(null), 5000);
    }
  };

  const openEditDialog = (usr: AdminUser): void => {
    setSelectedUser(usr);
    setEditRole(usr.role);
    setEditDialogOpen(true);
  };

  const openDeleteDialog = (usr: AdminUser): void => {
    setSelectedUser(usr);
    setDeleteDialogOpen(true);
  };

  const getRoleColor = (role: string): "default" | "primary" | "secondary" | "error" | "warning" | "info" | "success" => {
    switch (role) {
      case "root":
        return "error";
      case "super_admin":
        return "warning";
      case "admin":
        return "secondary";
      case "moderator":
        return "info";
      default:
        return "default";
    }
  };

  if (loading && users.length === 0) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Typography variant="h4">Gestión de Usuarios</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setCreateDialogOpen(true)}
        >
          Crear Usuario
        </Button>
      </Stack>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {success && (
        <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
          {success}
        </Alert>
      )}

      <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
        Total: {total} usuarios
      </Typography>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Usuario</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Nombre</TableCell>
              <TableCell>Rol</TableCell>
              <TableCell>Estado</TableCell>
              <TableCell>Acciones</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((usr) => (
              <TableRow key={usr.id}>
                <TableCell>{usr.id}</TableCell>
                <TableCell>{usr.username}</TableCell>
                <TableCell>{usr.email}</TableCell>
                <TableCell>
                  {usr.first_name} {usr.last_name}
                </TableCell>
                <TableCell>
                  <Chip label={usr.role} color={getRoleColor(usr.role)} size="small" />
                </TableCell>
                <TableCell>
                  <Chip
                    icon={usr.is_active ? <ActiveIcon /> : <InactiveIcon />}
                    label={usr.is_active ? "Activo" : "Inactivo"}
                    color={usr.is_active ? "success" : "default"}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Stack direction="row" spacing={1}>
                    <Tooltip title="Editar rol">
                      <IconButton
                        size="small"
                        color="primary"
                        onClick={() => openEditDialog(usr)}
                        disabled={usr.id === user?.id}
                      >
                        <EditIcon />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title={usr.is_active ? "Desactivar" : "Activar"}>
                      <IconButton
                        size="small"
                        color={usr.is_active ? "warning" : "success"}
                        onClick={() => handleToggleStatus(usr.id, usr.is_active)}
                        disabled={usr.id === user?.id}
                      >
                        {usr.is_active ? <InactiveIcon /> : <ActiveIcon />}
                      </IconButton>
                    </Tooltip>
                    {user?.role === "root" && (
                      <Tooltip title="Eliminar">
                        <IconButton
                          size="small"
                          color="error"
                          onClick={() => openDeleteDialog(usr)}
                          disabled={usr.id === user?.id}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                  </Stack>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
        <Pagination
          count={totalPages}
          page={page}
          onChange={(_, value) => setPage(value)}
          color="primary"
        />
      </Box>

      {/* Diálogo de crear usuario */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Crear Nuevo Usuario</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            <TextField
              label="Email"
              type="email"
              fullWidth
              value={newUser.email}
              onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
            />
            <TextField
              label="Nombre de usuario"
              fullWidth
              value={newUser.username}
              onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
            />
            <TextField
              label="Contraseña"
              type="password"
              fullWidth
              value={newUser.password}
              onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
            />
            <TextField
              label="Nombre"
              fullWidth
              value={newUser.first_name}
              onChange={(e) => setNewUser({ ...newUser, first_name: e.target.value })}
            />
            <TextField
              label="Apellido"
              fullWidth
              value={newUser.last_name}
              onChange={(e) => setNewUser({ ...newUser, last_name: e.target.value })}
            />
            <FormControl fullWidth>
              <InputLabel>Rol</InputLabel>
              <Select
                value={newUser.role}
                label="Rol"
                onChange={(e) => setNewUser({ ...newUser, role: e.target.value })}
              >
                {ROLES.map((role) => (
                  <MenuItem key={role} value={role}>
                    {role}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancelar</Button>
          <Button onClick={handleCreateUser} variant="contained">
            Crear
          </Button>
        </DialogActions>
      </Dialog>

      {/* Diálogo de editar rol */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)}>
        <DialogTitle>Editar Rol de Usuario</DialogTitle>
        <DialogContent>
          <FormControl fullWidth sx={{ mt: 2 }}>
            <InputLabel>Rol</InputLabel>
            <Select
              value={editRole}
              label="Rol"
              onChange={(e) => setEditRole(e.target.value)}
            >
              {ROLES.map((role) => (
                <MenuItem key={role} value={role}>
                  {role}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Cancelar</Button>
          <Button onClick={handleUpdateRole} variant="contained">
            Actualizar
          </Button>
        </DialogActions>
      </Dialog>

      {/* Diálogo de eliminar usuario */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Eliminar Usuario</DialogTitle>
        <DialogContent>
          <Typography>
            ¿Estás seguro de que deseas eliminar al usuario {selectedUser?.username}?
            Esta acción no se puede deshacer.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancelar</Button>
          <Button onClick={handleDeleteUser} color="error" variant="contained">
            Eliminar
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default UserManagement;
