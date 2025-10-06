import React, { useState, ChangeEvent, FormEvent, useRef } from "react";
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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  IconButton,
  Paper,
  Tabs,
  Tab,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Tooltip,
} from "@mui/material";
import {
  Edit as EditIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
  PhotoCamera as PhotoCameraIcon,
  Delete as DeleteIcon,
  Person as PersonIcon,
  FitnessCenter as FitnessCenterIcon,
  SportsBaseball as SportsIcon,
  Link as LinkIcon,
  ExpandMore as ExpandMoreIcon,
  Height as HeightIcon,
  Scale as WeightIcon,
  LocationOn as LocationIcon,
  Phone as PhoneIcon,
  Cake as CakeIcon,
} from "@mui/icons-material";
import { useAuth } from "../../contexts/AuthContext";
import BirthDatePicker from "../BirthDatePicker";
import {
  UpdateProfileRequest,
  SocialLinks,
  Message,
} from "../../types";

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`profile-tabpanel-${index}`}
      aria-labelledby={`profile-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

// Lista de deportes predefinidos
const SPORTS_OPTIONS = [
  'Football', 'Basketball', 'Tennis', 'Swimming', 'Running', 'Cycling',
  'Boxing', 'Wrestling', 'Martial Arts', 'Volleyball', 'Baseball', 'Soccer',
  'Golf', 'Hiking', 'Rock Climbing', 'Surfing', 'Skiing', 'Snowboarding',
  'Skateboarding', 'Yoga', 'Pilates', 'CrossFit', 'Weightlifting', 'Gymnastics'
];

function ExtendedProfile(): JSX.Element {
  const { user, updateProfile, uploadAvatar, deleteAvatar } = useAuth();
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const [tabValue, setTabValue] = useState(0);
  const [editMode, setEditMode] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<Message>({ type: "", text: "" });
  const [avatarUploading, setAvatarUploading] = useState<boolean>(false);

  // Parse sports interests from JSON string
  const parseSportsInterests = (sportsString?: string | null): string[] => {
    if (!sportsString) return [];
    try {
      return JSON.parse(sportsString);
    } catch {
      return [];
    }
  };

  const [profileData, setProfileData] = useState<UpdateProfileRequest>({
    first_name: user?.first_name || "",
    last_name: user?.last_name || "",
    bio: user?.bio || "",
    phone: user?.phone || "",
    birth_date: user?.birth_date ? user.birth_date.split('T')[0] : "",
    location: user?.location || "",
    gender: user?.gender || undefined,
    height: user?.height || undefined,
    weight: user?.weight || undefined,
    fitness_level: user?.fitness_level || undefined,
    social_links: user?.social_links || {
      instagram: "",
      twitter: "",
      facebook: "",
      linkedin: "",
      youtube: "",
      website: "",
    },
  });

  const [selectedSports, setSelectedSports] = useState<string[]>(
    parseSportsInterests(user?.sports_interests)
  );

  // Handle tab change
  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // Handle profile field changes
  const handleProfileChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    setProfileData((prev) => ({
      ...prev,
      [name]: value === "" ? undefined : value,
    }));
  };

  // Handle birth date change
  const handleBirthDateChange = (date: string | null): void => {
    setProfileData((prev) => ({
      ...prev,
      birth_date: date || undefined,
    }));
  };

  // Handle social links changes
  const handleSocialLinksChange = (platform: keyof SocialLinks, value: string): void => {
    setProfileData((prev) => ({
      ...prev,
      social_links: {
        ...prev.social_links,
        [platform]: value,
      },
    }));
  };

  // Handle sports selection
  const handleSportAdd = (sport: string): void => {
    if (!selectedSports.includes(sport)) {
      setSelectedSports((prev) => [...prev, sport]);
    }
  };

  const handleSportRemove = (sport: string): void => {
    setSelectedSports((prev) => prev.filter((s) => s !== sport));
  };

  // Handle avatar upload
  const handleAvatarUpload = async (e: ChangeEvent<HTMLInputElement>): Promise<void> => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith('image/')) {
      setMessage({
        type: "error",
        text: "Please select a valid image file",
      });
      return;
    }

    // Validate file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
      setMessage({
        type: "error",
        text: "Image must be smaller than 5MB",
      });
      return;
    }

    setAvatarUploading(true);
    setMessage({ type: "", text: "" });

    try {
      const result = await uploadAvatar(file);
      if (result.success) {
        setMessage({
          type: "success",
          text: "Avatar updated successfully!",
        });
      } else {
        setMessage({
          type: "error",
          text: result.error || "Failed to upload avatar",
        });
      }
    } catch (error) {
      setMessage({
        type: "error",
        text: "Failed to upload avatar",
      });
    } finally {
      setAvatarUploading(false);
      // Reset file input
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  // Handle avatar deletion
  const handleAvatarDelete = async (): Promise<void> => {
    setAvatarUploading(true);
    setMessage({ type: "", text: "" });

    try {
      const result = await deleteAvatar();
      if (result.success) {
        setMessage({
          type: "success",
          text: "Avatar removed successfully!",
        });
      } else {
        setMessage({
          type: "error",
          text: result.error || "Failed to remove avatar",
        });
      }
    } catch (error) {
      setMessage({
        type: "error",
        text: "Failed to remove avatar",
      });
    } finally {
      setAvatarUploading(false);
    }
  };

  // Handle form submission
  const handleProfileSubmit = async (e: FormEvent): Promise<void> => {
    e.preventDefault();
    setLoading(true);
    setMessage({ type: "", text: "" });

    try {
      // Include sports interests as JSON string
      const dataToSubmit = {
        ...profileData,
        sports_interests: selectedSports.length > 0 ? JSON.stringify(selectedSports) : undefined,
      };

      const result = await updateProfile(dataToSubmit);
      if (result.success) {
        setMessage({
          type: "success",
          text: "Profile updated successfully!",
        });
        setEditMode(false);
      } else {
        setMessage({
          type: "error",
          text: result.error || "Failed to update profile",
        });
      }
    } catch (error) {
      setMessage({
        type: "error",
        text: "Failed to update profile",
      });
    } finally {
      setLoading(false);
    }
  };

  // Cancel edit mode
  const handleCancel = (): void => {
    setEditMode(false);
    setMessage({ type: "", text: "" });
    // Reset data
    setProfileData({
      first_name: user?.first_name || "",
      last_name: user?.last_name || "",
      bio: user?.bio || "",
      phone: user?.phone || "",
      birth_date: user?.birth_date ? user.birth_date.split('T')[0] : "",
      location: user?.location || "",
      gender: user?.gender || undefined,
      height: user?.height || undefined,
      weight: user?.weight || undefined,
      fitness_level: user?.fitness_level || undefined,
      social_links: user?.social_links || {
        instagram: "",
        twitter: "",
        facebook: "",
        linkedin: "",
        youtube: "",
        website: "",
      },
    });
    setSelectedSports(parseSportsInterests(user?.sports_interests));
  };

  if (!user) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  const avatarUrl = user.avatar_url 
    ? `${process.env.REACT_APP_API_URL?.replace('/api/v1', '') || 'http://localhost:8081'}${user.avatar_url}`
    : undefined;

  return (
    <Box sx={{ maxWidth: 1200, margin: "0 auto", padding: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
        🏃‍♀️ Extended Profile
      </Typography>

      {message.text && (
        <Alert 
          severity={message.type as any} 
          sx={{ mb: 3 }}
          onClose={() => setMessage({ type: "", text: "" })}
        >
          {message.text}
        </Alert>
      )}

      {/* Avatar Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box display="flex" alignItems="center" gap={3}>
            <Box position="relative">
              <Avatar
                src={avatarUrl}
                sx={{ 
                  width: 120, 
                  height: 120,
                  border: 3,
                  borderColor: 'primary.main'
                }}
              >
                {user.first_name?.[0]}{user.last_name?.[0]}
              </Avatar>
              {avatarUploading && (
                <Box
                  position="absolute"
                  top={0}
                  left={0}
                  right={0}
                  bottom={0}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  bgcolor="rgba(0,0,0,0.5)"
                  borderRadius="50%"
                >
                  <CircularProgress size={24} />
                </Box>
              )}
            </Box>
            
            <Box>
              <Typography variant="h5" gutterBottom>
                {user.first_name} {user.last_name}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                @{user.username} • {user.role}
              </Typography>
              {user.email_verified ? (
                <Chip label="Email Verified" color="success" size="small" />
              ) : (
                <Chip label="Email Not Verified" color="warning" size="small" />
              )}
              
              <Box sx={{ mt: 2 }}>
                <input
                  type="file"
                  accept="image/*"
                  ref={fileInputRef}
                  onChange={handleAvatarUpload}
                  style={{ display: 'none' }}
                />
                <Button
                  variant="outlined"
                  startIcon={<PhotoCameraIcon />}
                  onClick={() => fileInputRef.current?.click()}
                  disabled={avatarUploading}
                  sx={{ mr: 1 }}
                >
                  Upload Avatar
                </Button>
                {user.avatar_url && (
                  <Button
                    variant="outlined"
                    color="error"
                    startIcon={<DeleteIcon />}
                    onClick={handleAvatarDelete}
                    disabled={avatarUploading}
                  >
                    Remove
                  </Button>
                )}
              </Box>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {/* Profile Tabs */}
      <Card>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={tabValue} onChange={handleTabChange} aria-label="profile tabs">
            <Tab icon={<PersonIcon />} label="Personal Info" />
            <Tab icon={<FitnessCenterIcon />} label="Fitness" />
            <Tab icon={<SportsIcon />} label="Sports" />
            <Tab icon={<LinkIcon />} label="Social Links" />
          </Tabs>
        </Box>

        <form onSubmit={handleProfileSubmit}>
          {/* Personal Info Tab */}
          <TabPanel value={tabValue} index={0}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="First Name"
                  name="first_name"
                  value={profileData.first_name || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  InputProps={{
                    startAdornment: <PersonIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Last Name"
                  name="last_name"
                  value={profileData.last_name || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Bio"
                  name="bio"
                  value={profileData.bio || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  multiline
                  rows={3}
                  placeholder="Tell us about yourself..."
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Phone"
                  name="phone"
                  value={profileData.phone || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  InputProps={{
                    startAdornment: <PhoneIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <BirthDatePicker
                  value={profileData.birth_date || null}
                  onChange={handleBirthDateChange}
                  disabled={!editMode}
                  label="Birth Date"
                  fullWidth
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Location"
                  name="location"
                  value={profileData.location || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  InputProps={{
                    startAdornment: <LocationIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth disabled={!editMode}>
                  <InputLabel>Gender</InputLabel>
                  <Select
                    name="gender"
                    value={profileData.gender || ""}
                    label="Gender"
                    onChange={(e) => setProfileData(prev => ({ 
                      ...prev, 
                      gender: e.target.value as any || undefined 
                    }))}
                  >
                    <MenuItem value="">
                      <em>Prefer not to say</em>
                    </MenuItem>
                    <MenuItem value="male">Male</MenuItem>
                    <MenuItem value="female">Female</MenuItem>
                    <MenuItem value="other">Other</MenuItem>
                    <MenuItem value="prefer_not_to_say">Prefer not to say</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </TabPanel>

          {/* Fitness Tab */}
          <TabPanel value={tabValue} index={1}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={4}>
                <TextField
                  fullWidth
                  label="Height (cm)"
                  name="height"
                  type="number"
                  value={profileData.height || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  InputProps={{
                    startAdornment: <HeightIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  fullWidth
                  label="Weight (kg)"
                  name="weight"
                  type="number"
                  value={profileData.weight || ""}
                  onChange={handleProfileChange}
                  disabled={!editMode}
                  InputProps={{
                    startAdornment: <WeightIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <FormControl fullWidth disabled={!editMode}>
                  <InputLabel>Fitness Level</InputLabel>
                  <Select
                    name="fitness_level"
                    value={profileData.fitness_level || ""}
                    label="Fitness Level"
                    onChange={(e) => setProfileData(prev => ({ 
                      ...prev, 
                      fitness_level: e.target.value as any || undefined 
                    }))}
                  >
                    <MenuItem value="">
                      <em>Not specified</em>
                    </MenuItem>
                    <MenuItem value="beginner">Beginner</MenuItem>
                    <MenuItem value="intermediate">Intermediate</MenuItem>
                    <MenuItem value="advanced">Advanced</MenuItem>
                    <MenuItem value="professional">Professional</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </TabPanel>

          {/* Sports Tab */}
          <TabPanel value={tabValue} index={2}>
            <Typography variant="h6" gutterBottom>
              Sports Interests
            </Typography>
            
            {editMode && (
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Add Sport</InputLabel>
                <Select
                  value=""
                  label="Add Sport"
                  onChange={(e) => handleSportAdd(e.target.value)}
                >
                  {SPORTS_OPTIONS.filter(sport => !selectedSports.includes(sport)).map((sport) => (
                    <MenuItem key={sport} value={sport}>
                      {sport}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}

            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
              {selectedSports.map((sport) => (
                <Chip
                  key={sport}
                  label={sport}
                  onDelete={editMode ? () => handleSportRemove(sport) : undefined}
                  color="primary"
                  variant="outlined"
                />
              ))}
              {selectedSports.length === 0 && (
                <Typography color="text.secondary">
                  No sports selected
                </Typography>
              )}
            </Box>
          </TabPanel>

          {/* Social Links Tab */}
          <TabPanel value={tabValue} index={3}>
            <Grid container spacing={3}>
              {Object.entries(profileData.social_links || {}).map(([platform, url]) => (
                <Grid item xs={12} md={6} key={platform}>
                  <TextField
                    fullWidth
                    label={platform.charAt(0).toUpperCase() + platform.slice(1)}
                    value={url || ""}
                    onChange={(e) => handleSocialLinksChange(platform as keyof SocialLinks, e.target.value)}
                    disabled={!editMode}
                    placeholder={`https://${platform}.com/username`}
                  />
                </Grid>
              ))}
            </Grid>
          </TabPanel>

          {/* Action Buttons */}
          <Box sx={{ p: 3, pt: 0 }}>
            <Divider sx={{ mb: 3 }} />
            {editMode ? (
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Button
                  type="submit"
                  variant="contained"
                  startIcon={loading ? <CircularProgress size={20} /> : <SaveIcon />}
                  disabled={loading}
                >
                  {loading ? "Saving..." : "Save Changes"}
                </Button>
                <Button
                  variant="outlined"
                  startIcon={<CancelIcon />}
                  onClick={handleCancel}
                  disabled={loading}
                >
                  Cancel
                </Button>
              </Box>
            ) : (
              <Button
                variant="contained"
                startIcon={<EditIcon />}
                onClick={() => setEditMode(true)}
              >
                Edit Profile
              </Button>
            )}
          </Box>
        </form>
      </Card>
    </Box>
  );
}

export default ExtendedProfile;
