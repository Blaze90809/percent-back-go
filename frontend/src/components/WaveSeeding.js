import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {
  Container,
  Typography,
  Box,
  Card,
  CardContent,
  Grid,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  CircularProgress,
  Alert,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Button
} from '@mui/material';
import { useNavigate } from 'react-router-dom';

const eventLabels = {
  birkie_skate: { name: 'Birkie Skate 50K', short: 'Birkie Skate' },
  birkie_classic: { name: 'Birkie Classic 53K', short: 'Birkie Classic' },
  korte_skate: { name: 'Kortelopet Skate 29K', short: 'Korte Skate' },
  korte_classic: { name: 'Kortelopet Classic 29K', short: 'Korte Classic' },
};

const getWaveChipColor = (wave) => {
  if (wave === 1) return 'primary';
  if (wave === 2 || wave === 3) return 'secondary';
  if (wave > 3 && wave <= 5) return 'info';
  if (wave > 5) return 'default';
  return 'default';
};

const WaveSeeding = () => {
  const [seasons, setSeasons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const fetchWaveSeeding = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get('/api/seeding/my-waves', {
        headers: { Authorization: `Bearer ${token}` }
      });
      setSeasons(response.data || []);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load wave seeding qualifications');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchWaveSeeding();
  }, []);

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 6, display: 'flex', justifyContent: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  const latestSeason = seasons.length > 0 ? seasons[0] : null;

  return (
    <Container maxWidth="lg" sx={{ my: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
          Wave Seeding & Qualifications
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Track which American Birkebeiner and Kortelopet waves your race results qualify you for across each ski season (July 1 – June 30).
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {seasons.length === 0 ? (
        <Card sx={{ p: 4, textAlign: 'center', mt: 3 }}>
          <Typography sx={{ fontSize: 48, mb: 1 }}>🎿</Typography>
          <Typography variant="h5" gutterBottom>
            No Races Logged Yet
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            Add your race results with finish time and first place time to see which Birkie and Kortelopet waves you qualify for.
          </Typography>
          <Button variant="contained" color="primary" onClick={() => navigate('/dashboard/input')}>
            Add Your First Race
          </Button>
        </Card>
      ) : (
        <>
          {/* Latest Season Highlight Banner */}
          {latestSeason && (
            <Box sx={{ mb: 4 }}>
              <Typography variant="h5" component="h2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <span>🏆</span> Current Season Best ({latestSeason.season})
              </Typography>
              <Grid container spacing={2}>
                {Object.keys(eventLabels).map((eventKey) => {
                  const waveNum = latestSeason.best_waves?.[eventKey] || 0;
                  const waveLabel = latestSeason.best_wave_labels?.[eventKey] || 'Not Qualified';
                  const earnedRace = latestSeason.best_races?.[eventKey];

                  return (
                    <Grid item xs={12} sm={6} md={3} key={eventKey}>
                      <Card sx={{ height: '100%', borderTop: 4, borderColor: waveNum > 0 ? 'primary.main' : 'divider' }}>
                        <CardContent>
                          <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                            {eventLabels[eventKey].name}
                          </Typography>
                          <Typography variant="h4" component="div" fontWeight="bold" sx={{ my: 1 }}>
                            {waveNum > 0 ? (
                              <Chip
                                label={waveLabel}
                                color={getWaveChipColor(waveNum)}
                                sx={{ fontSize: '1rem', py: 2, px: 1, fontWeight: 'bold' }}
                              />
                            ) : (
                              <Typography variant="body2" color="text.secondary" sx={{ py: 0.5 }}>
                                No official qualification
                              </Typography>
                            )}
                          </Typography>
                          {earnedRace ? (
                            <Typography variant="caption" color="text.secondary" display="block">
                              Earned at: {earnedRace}
                            </Typography>
                          ) : (
                            <Typography variant="caption" color="text.disabled" display="block">
                              Requires official qualifier
                            </Typography>
                          )}
                        </CardContent>
                      </Card>
                    </Grid>
                  );
                })}
              </Grid>
            </Box>
          )}

          {/* Season-by-Season Details */}
          <Typography variant="h5" component="h2" gutterBottom sx={{ mt: 4 }}>
            Season History & Breakdown
          </Typography>

          {seasons.map((seasonData, index) => (
            <Accordion key={seasonData.season} defaultExpanded={index === 0} sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<Typography sx={{ fontSize: '0.85rem' }}>▼</Typography>}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%', pr: 2, flexWrap: 'wrap', gap: 1 }}>
                  <Typography variant="h6" fontWeight="bold">
                    {seasonData.season} Season ({seasonData.races.length} {seasonData.races.length === 1 ? 'race' : 'races'})
                  </Typography>
                  <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                    {Object.keys(eventLabels).map((eventKey) => {
                      const waveNum = seasonData.best_waves?.[eventKey];
                      if (!waveNum || waveNum <= 0) return null;
                      return (
                        <Chip
                          key={eventKey}
                          size="small"
                          label={`${eventLabels[eventKey].short}: Wave ${waveNum}`}
                          color={getWaveChipColor(waveNum)}
                        />
                      );
                    })}
                  </Box>
                </Box>
              </AccordionSummary>
              <AccordionDetails>
                <TableContainer component={Paper} variant="outlined">
                  <Table size="small">
                    <TableHead>
                      <TableRow sx={{ bgcolor: 'action.hover' }}>
                        <TableCell><strong>Race Name</strong></TableCell>
                        <TableCell><strong>Date</strong></TableCell>
                        <TableCell align="right"><strong>Distance</strong></TableCell>
                        <TableCell align="right"><strong>% Back</strong></TableCell>
                        <TableCell><strong>Matched Qualifier</strong></TableCell>
                        <TableCell align="center"><strong>Birkie Skate</strong></TableCell>
                        <TableCell align="center"><strong>Birkie Classic</strong></TableCell>
                        <TableCell align="center"><strong>Korte Skate</strong></TableCell>
                        <TableCell align="center"><strong>Korte Classic</strong></TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {seasonData.races.map((race) => {
                        const qInfo = race.matched_qualifier;
                        const isOfficial = race.is_official_qualifier;
                        const hasAdj = qInfo && qInfo.id !== 'none' && qInfo.adjustment !== 0;

                        return (
                          <TableRow key={race.id} hover sx={!isOfficial ? { bgcolor: 'action.hover', opacity: 0.9 } : {}}>
                            <TableCell component="th" scope="row">
                              <strong>{race.race_name}</strong>
                              {!isOfficial && (
                                <Typography variant="caption" color="text.secondary" display="block">
                                  Unofficial Event
                                </Typography>
                              )}
                            </TableCell>
                            <TableCell>{race.race_date}</TableCell>
                            <TableCell align="right">{race.race_distance ? `${race.race_distance} km` : '-'}</TableCell>
                            <TableCell align="right">{Number(race.percent_back).toFixed(2)}%</TableCell>
                            <TableCell>
                              {isOfficial ? (
                                hasAdj ? (
                                  <Chip
                                    size="small"
                                    label={`${qInfo.name} (${qInfo.adjustment > 0 ? '+' : ''}${qInfo.adjustment}%)`}
                                    color="info"
                                    variant="outlined"
                                  />
                                ) : qInfo && qInfo.id && qInfo.id.startsWith('birkie') ? (
                                  <Chip
                                    size="small"
                                    label="Birkie Direct (0%)"
                                    color="success"
                                    variant="outlined"
                                  />
                                ) : qInfo && qInfo.id && qInfo.id.startsWith('korte') ? (
                                  <Chip
                                    size="small"
                                    label="Korte Direct (0%)"
                                    color="success"
                                    variant="outlined"
                                  />
                                ) : (
                                  <Chip
                                    size="small"
                                    label="Direct (0%)"
                                    color="default"
                                    variant="outlined"
                                  />
                                )
                              ) : (
                                <Chip
                                  size="small"
                                  label={qInfo?.name || "Non-Qualifier"}
                                  color="default"
                                  variant="outlined"
                                  sx={{ color: 'text.secondary', borderStyle: 'dashed' }}
                                />
                              )}
                            </TableCell>
                            {['birkie_skate', 'birkie_classic', 'korte_skate', 'korte_classic'].map((eventKey) => {
                              const qual = race.qualifications?.[eventKey];
                              if (!qual) return <TableCell key={eventKey} align="center">-</TableCell>;

                              return (
                                <TableCell key={eventKey} align="center">
                                  <Chip
                                    size="small"
                                    label={qual.qualified_wave_label}
                                    color={isOfficial ? getWaveChipColor(qual.qualified_wave) : 'default'}
                                    variant={isOfficial ? 'filled' : 'outlined'}
                                    sx={{
                                      fontWeight: isOfficial ? 'bold' : 'normal',
                                      ...(isOfficial ? {} : { borderStyle: 'dashed', color: 'text.secondary' })
                                    }}
                                  />
                                  <Typography variant="caption" color="text.secondary" display="block">
                                    PBE: {Number(qual.final_pbe).toFixed(2)}%
                                    {isOfficial && qual.adjustment !== 0 && (
                                      <span style={{ fontSize: '0.7rem', color: '#666' }}>
                                        {' '}({qual.adjustment > 0 ? '+' : ''}{qual.adjustment}%)
                                      </span>
                                    )}
                                  </Typography>
                                  {isOfficial && qual.has_next_wave && qual.distance_to_next > 0 && (
                                    <Typography variant="caption" color="text.secondary" display="block">
                                      +{qual.distance_to_next}% to W{qual.qualified_wave - 1}
                                    </Typography>
                                  )}
                                </TableCell>
                              );
                            })}
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                </TableContainer>
                <Box sx={{ mt: 1.5, px: 1 }}>
                  <Typography variant="caption" color="text.secondary" sx={{ display: 'block', fontStyle: 'italic' }}>
                    * Dashed "Est. Wave" rows indicate unofficial results (e.g. rollerski races, events &lt; 15km, or unlisted races). While your percent back is calculated, only official qualifiers count toward your official Birkie &amp; Kortelopet wave placements and season bests.
                  </Typography>
                </Box>
              </AccordionDetails>
            </Accordion>
          ))}
        </>
      )}
    </Container>
  );
};

export default WaveSeeding;
