import React from 'react';
import { Link as RouterLink, Outlet } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Container, Box } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Percent Back
          </Typography>
          <Button color="inherit" component={RouterLink} to="/dashboard/input">Input</Button>
          <Button color="inherit" component={RouterLink} to="/dashboard/table">Table</Button>
          <Button color="inherit" component={RouterLink} to="/dashboard/chart">Chart</Button>
          <Button color="inherit" onClick={handleLogout}>Logout</Button>
        </Toolbar>
      </AppBar>
      <Container>
        <Outlet />
      </Container>
    </Box>
  );
};

export default Dashboard;
