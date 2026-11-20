import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
  TextField,
  Button,
  Container,
  Typography,
  Box,
  Card,
  CardContent,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Divider,
  Grid,
  Alert
} from '@mui/material';

const eventLabels = {
  birkie_skate: 'Birkie Skate 50K',
  birkie_classic: 'Birkie Classic 53K',
  korte_skate: 'Korte Skate 29K',
  korte_classic: 'Korte Classic 29K',
};

const getWaveChipColor = (wave) => {
  if (wave === 1) return 'primary';
  if (wave === 2 || wave === 3) return 'secondary';
  if (wave > 3 && wave <= 5) return 'info';
  return 'default';
};

const RaceInput = () => {
  const [raceName, setRaceName] = useState('');
  const [raceDate, setRaceDate] = useState('');
  const [raceDistance, setRaceDistance] = useState('');
  const [userRaceTime, setUserRaceTime] = useState('');
  const [firstPlaceTime, setFirstPlaceTime] = useState('');
  const [qualifiers, setQualifiers] = useState([]);
  const [selectedQualifier, setSelectedQualifier] = useState('auto');
  const [seedingPreview, setSeedingPreview] = useState(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const fetchQualifiers = async () => {
      try {
        const response = await axios.get('/api/seeding/qualifiers');
        setQualifiers(response.data || []);
      } catch (err) {
        console.error('Failed to load qualifiers', err);
      }
    };
    fetchQualifiers();
  }, []);

  const formatTime = (value) => {
    const digits = value.replace(/\D/g, '').slice(0, 6);
    if (digits.length <= 2) {
      return digits;
    } else if (digits.length <= 4) {
      return `${digits.slice(0, digits.length - 2)}:${digits.slice(digits.length - 2)}`;
    } else {
      return `${digits.slice(0, digits.length - 4)}:${digits.slice(digits.length - 4, digits.length - 2)}:${digits.slice(digits.length - 2)}`;
    }
  };

  const parseTimeToSeconds = (timeString) => {
    const parts = timeString.split(':');
    if (parts.length === 2) {
      return parseInt(parts[0], 10) * 60 + parseInt(parts[1], 10);
    } else if (parts.length === 3) {
      return parseInt(parts[0], 10) * 3600 + parseInt(parts[1], 10) * 60 + parseInt(parts[2], 10);
    }
    return 0;
  };

  const calculatePercentBack = () => {
    const userTimeSeconds = parseTimeToSeconds(userRaceTime);
    const firstPlaceTimeSeconds = parseTimeToSeconds(firstPlaceTime);

    if (userTimeSeconds > 0 && firstPlaceTimeSeconds > 0 && userTimeSeconds >= firstPlaceTimeSeconds) {
      return ((userTimeSeconds - firstPlaceTimeSeconds) / firstPlaceTimeSeconds) * 100;
    }
    return 0;
  };

  // Update live seeding preview when times or qualifier changes
  useEffect(() => {
    const pct = calculatePercentBack();
    if (pct > 0) {
      const qId = selectedQualifier === 'auto' ? '' : selectedQualifier;
      axios
        .post('/api/seeding/evaluate', {
          raceName,
          raceDate,
          raceDistance: parseFloat(raceDistance) || 0,
          percentBack: pct,
          qualifierId: qId,
        })
        .then((res) => {
          setSeedingPreview(res.data);
        })
        .catch(() => {
          setSeedingPreview(null);
        });
    } else {
      setSeedingPreview(null);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [raceName, raceDate, raceDistance, userRaceTime, firstPlaceTime, selectedQualifier]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const calculatedPercentBack = calculatePercentBack();

    if (calculatedPercentBack === 0 && (userRaceTime !== '' || firstPlaceTime !== '')) {
      setError('Please enter valid race times (MM:SS or HH:MM:SS) where your time is greater than or equal to first place.');
      return;
    }

    try {
      const token = localStorage.getItem('token');
      await axios.post(
        '/api/races/create',
        {
          raceName,
          raceDate,
          raceDistance: parseFloat(raceDistance),
          percentBack: calculatedPercentBack,
        },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setSuccess('Race created successfully!');
      setRaceName('');
      setRaceDate('');
      setRaceDistance('');
      setUserRaceTime('');
      setFirstPlaceTime('');
      setSelectedQualifier('auto');
      setSeedingPreview(null);
      setError('');
    } catch (err) {
      setError('Failed to create race');
      setSuccess('');
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ marginTop: 4, mb: 6, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Card sx={{ width: '100%' }}>
          <CardContent sx={{ p: 3 }}>
            <Typography component="h1" variant="h5" fontWeight="bold" gutterBottom>
              Add New Race
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Enter race results to compute your percent back and preview Birkie/Kortelopet wave qualifications.
            </Typography>

            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="raceName"
                label="Race Name"
                name="raceName"
                placeholder="e.g. Seeley Hills Classic, Mora Vasaloppet, Birkie Skate"
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

              <FormControl fullWidth margin="normal">
                <InputLabel id="qualifier-label">Qualifier Adjustment</InputLabel>
                <Select
                  labelId="qualifier-label"
                  id="qualifier"
                  value={selectedQualifier}
                  label="Qualifier Adjustment"
                  onChange={(e) => setSelectedQualifier(e.target.value)}
                >
                  <MenuItem value="auto">
                    <em>Auto-detect from Race Name</em>
                  </MenuItem>
                  {qualifiers.map((q) => (
                    <MenuItem key={q.id} value={q.id}>
                      {q.name} ({q.adjustment > 0 ? `+${q.adjustment}%` : `${q.adjustment}%`})
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

              <TextField
                margin="normal"
                required
                fullWidth
                id="userRaceTime"
                label="Your Race Time (MM:SS or HH:MM:SS)"
                name="userRaceTime"
                value={userRaceTime}
                onChange={(e) => setUserRaceTime(formatTime(e.target.value))}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                id="firstPlaceTime"
                label="First Place Finisher Time (MM:SS or HH:MM:SS)"
                name="firstPlaceTime"
                value={firstPlaceTime}
                onChange={(e) => setFirstPlaceTime(formatTime(e.target.value))}
              />

              {/* Real-time Wave Seeding Preview Card */}
              {seedingPreview && (
                <Card variant="outlined" sx={{ mt: 2, mb: 2, bgcolor: 'background.paper', borderColor: seedingPreview.is_official_qualifier ? 'primary.main' : 'divider' }}>
                  <CardContent sx={{ p: 2 }}>
                    <Typography variant="subtitle2" color={seedingPreview.is_official_qualifier ? 'primary' : 'text.secondary'} fontWeight="bold" gutterBottom>
                      Live Seeding Preview ({seedingPreview.season} Season)
                    </Typography>

                    <Typography variant="body2" sx={{ mb: 1 }}>
                      Raw Percent Back: <strong>{Number(seedingPreview.percent_back).toFixed(2)}%</strong>
                      {seedingPreview.matched_qualifier?.id !== 'none' && seedingPreview.is_official_qualifier && (
                        <span>
                          {' '}| Qualifier: <em>{seedingPreview.matched_qualifier.name}</em> (
                          {seedingPreview.matched_qualifier.adjustment > 0 ? '+' : ''}
                          {seedingPreview.matched_qualifier.adjustment}%)
                        </span>
                      )}
                    </Typography>

                    {!seedingPreview.is_official_qualifier && (
                      <Alert severity="info" sx={{ my: 1, py: 0.5, fontSize: '0.8rem' }}>
                        This event is not recognized as an official Birkie qualifier (e.g. rollerski, &lt; 15km, or unlisted). Wave placements below are estimates only and do not count toward official start waves.
                      </Alert>
                    )}

                    <Divider sx={{ my: 1 }} />

                    <Grid container spacing={1} sx={{ mt: 0.5 }}>
                      {Object.keys(eventLabels).map((eventKey) => {
                        const q = seedingPreview.qualifications?.[eventKey];
                        if (!q) return null;
                        return (
                          <Grid item xs={6} key={eventKey}>
                            <Box sx={{ p: 1, border: '1px solid', borderColor: 'divider', borderRadius: 1 }}>
                              <Typography variant="caption" color="text.secondary" display="block">
                                {eventLabels[eventKey]}
                              </Typography>
                              <Chip
                                size="small"
                                label={q.qualified_wave_label}
                                color={seedingPreview.is_official_qualifier ? getWaveChipColor(q.qualified_wave) : 'default'}
                                variant={seedingPreview.is_official_qualifier ? 'filled' : 'outlined'}
                                sx={{
                                  my: 0.5,
                                  fontWeight: seedingPreview.is_official_qualifier ? 'bold' : 'normal',
                                  ...(seedingPreview.is_official_qualifier ? {} : { borderStyle: 'dashed' })
                                }}
                              />
                              <Typography variant="caption" display="block" color="text.secondary">
                                PBE: {Number(q.final_pbe).toFixed(2)}%
                              </Typography>
                            </Box>
                          </Grid>
                        );
                      })}
                    </Grid>
                  </CardContent>
                </Card>
              )}

              {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
              {success && <Alert severity="success" sx={{ mt: 2 }}>{success}</Alert>}

              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                sx={{ mt: 3, mb: 1 }}
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
