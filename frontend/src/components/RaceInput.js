import React, { useState } from 'react';
import axios from 'axios';
import { TextField, Button, Container, Typography, Box, Card, CardContent } from '@mui/material';

const RaceInput = () => {
  const [raceName, setRaceName] = useState('');
  const [raceDate, setRaceDate] = useState('');
  const [raceDistance, setRaceDistance] = useState('');
  const [userRaceTime, setUserRaceTime] = useState('');
  const [firstPlaceTime, setFirstPlaceTime] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const calculatePercentBack = () => {
    const userTimeSeconds = parseTimeToSeconds(userRaceTime);
    const firstPlaceTimeSeconds = parseTimeToSeconds(firstPlaceTime);

    if (userTimeSeconds > 0 && firstPlaceTimeSeconds > 0) {
      return ((userTimeSeconds - firstPlaceTimeSeconds) / firstPlaceTimeSeconds) * 100;
    }
    return 0;
  };

  const parseTimeToSeconds = (timeString) => {
    const parts = timeString.split(':');
    if (parts.length === 2) {
      return parseInt(parts[0]) * 60 + parseInt(parts[1]);
    } else if (parts.length === 3) {
      return parseInt(parts[0]) * 3600 + parseInt(parts[1]) * 60 + parseInt(parts[2]);
    }
    return 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const calculatedPercentBack = calculatePercentBack();

    if (calculatedPercentBack === 0 && (userRaceTime !== '' || firstPlaceTime !== '')) {
      setError('Please enter valid race times (MM:SS or HH:MM:SS).');
      return;
    }

    try {
      const token = localStorage.getItem('token');
      await axios.post('/api/races/create', 
        { 
          raceName,
          raceDate,
          raceDistance: parseFloat(raceDistance),
          percentBack: calculatedPercentBack 
        }, 
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setSuccess('Race created successfully!');
      setRaceName('');
      setRaceDate('');
      setRaceDistance('');
      setUserRaceTime('');
      setFirstPlaceTime('');
      setError('');
    } catch (err) {
      setError('Failed to create race');
      setSuccess('');
    }
  };

  return (
    <Container maxWidth="xs">
      <Box sx={{ marginTop: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Card>
          <CardContent>
            <Typography component="h1" variant="h5">
              Add New Race
            </Typography>
            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="raceName"
                label="Race Name"
                name="raceName"
                autoFocus
                value={raceName}
                onChange={(e) => setRaceName(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="raceDate"
                label="Race Date"
                type="date"
                id="raceDate"
                InputLabelProps={{
                  shrink: true,
                }}
                value={raceDate}
                onChange={(e) => setRaceDate(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="raceDistance"
                label="Race Distance (km)"
                type="number"
                id="raceDistance"
                value={raceDistance}
                onChange={(e) => setRaceDistance(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                id="userRaceTime"
                label="Your Race Time (MM:SS or HH:MM:SS)"
                name="userRaceTime"
                value={userRaceTime}
                onChange={(e) => setUserRaceTime(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                id="firstPlaceTime"
                label="First Place Finisher Time (MM:SS or HH:MM:SS)"
                name="firstPlaceTime"
                value={firstPlaceTime}
                onChange={(e) => setFirstPlaceTime(e.target.value)}
              />
              {error && <Typography color="error">{error}</Typography>}
              {success && <Typography color="success">{success}</Typography>}
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                Add Race
              </Button>
            </Box>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
};

export default RaceInput;
