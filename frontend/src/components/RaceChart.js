import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Container, Typography, Box, Card, CardContent, Select, MenuItem, FormControl, InputLabel, Button, TextField } from '@mui/material';

const RaceChart = () => {
  const [originalRaces, setOriginalRaces] = useState([]);
  const [filteredRaces, setFilteredRaces] = useState([]);
  const [chartData, setChartData] = useState([]);
  const [years, setYears] = useState([]);
  const [selectedYear, setSelectedYear] = useState('');
  const [averagePercentBack, setAveragePercentBack] = useState(null);
  const [filterRaceName, setFilterRaceName] = useState('');
  const [error, setError] = useState('');
  const [windowWidth, setWindowWidth] = useState(window.innerWidth);

  const getSeason = (dateStr) => {
    const date = new Date(dateStr);
    // If month is July (6) or later, it belongs to the next year's season
    return date.getMonth() >= 6 ? date.getFullYear() + 1 : date.getFullYear();
  };

  useEffect(() => {
    const handleResize = () => setWindowWidth(window.innerWidth);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    const fetchRaces = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('/api/races', { headers: { Authorization: `Bearer ${token}` } });
        const races = response.data || [];
        const sortedRaces = races.sort((a, b) => new Date(a.RaceDate) - new Date(b.RaceDate));
        setOriginalRaces(sortedRaces);
        setFilteredRaces(sortedRaces);

        // Initialize filters
        setFilterRaceName('');
        setSelectedYear('');

        const uniqueYears = [...new Set(sortedRaces.map(race => getSeason(race.RaceDate)))].sort((a, b) => b - a);
        setYears(uniqueYears);
      } catch (err) {
        setError('Failed to fetch races');
      }
    };
    fetchRaces();
  }, []);

  const handleRaceNameChange = (event) => {
    setFilterRaceName(event.target.value);
  };

  useEffect(() => {
    let currentRaces = originalRaces;
    // 1. Filter by Year
    if (selectedYear) {
      currentRaces = currentRaces.filter(race => getSeason(race.RaceDate) === selectedYear);
    }
    // 2. Filter by Race Name
    if (filterRaceName) {
      const lowerCaseFilter = filterRaceName.toLowerCase();
      currentRaces = currentRaces.filter(race => race.RaceName.toLowerCase().includes(lowerCaseFilter));
    }
    setFilteredRaces(currentRaces);
  }, [originalRaces, selectedYear, filterRaceName]);

  useEffect(() => {
    if (filteredRaces.length > 0) {
      // 1. Calculate stats (sum and count) for each season
      const seasonStats = {};
      filteredRaces.forEach(race => {
        const season = getSeason(race.RaceDate);
        if (!seasonStats[season]) {
          seasonStats[season] = { sum: 0, count: 0 };
        }
        seasonStats[season].sum += race.PercentBack;
        seasonStats[season].count += 1;
      });

      // 2. Assign the season average as the Trend value for each race
      const dataWithTrend = filteredRaces.map((race) => {
        const season = getSeason(race.RaceDate);
        const stats = seasonStats[season];
        const seasonAvg = stats.sum / stats.count;

        return {
          ...race,
          PercentBack: parseFloat(Number(race.PercentBack).toFixed(2)),
          Trend: parseFloat(seasonAvg.toFixed(2))
        };
      });
      setChartData(dataWithTrend);
      calculateAverage(filteredRaces); // Calculate overall average when filters change
    } else {
      setChartData([]);
      setAveragePercentBack(null); // Reset average if no races are filtered
    }
  }, [filteredRaces]);

  const calculateAverage = (races) => {
    if (races.length === 0) {
      setAveragePercentBack(null);
      return;
    }
    const total = races.reduce((acc, race) => acc + race.PercentBack, 0);
    setAveragePercentBack((total / races.length).toFixed(2));
  };

  const handleYearChange = (event) => {
    const year = event.target.value;
    setSelectedYear(year);
  };

  const handleReset = () => {
    setSelectedYear('');
    setFilterRaceName('');
  };

  const formatDate = (dateStr) => {
    try {
      const date = new Date(dateStr);
      if (isNaN(date.getTime())) return dateStr;
      return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
    } catch {
      return dateStr;
    }
  };

  const chartHeight = windowWidth < 600 ? 280 : 400;

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Card>
          <CardContent>
            <Typography variant="h4" component="h1" gutterBottom>
              Percent Back Over Time
            </Typography>
            {error && <Typography color="error">{error}</Typography>}
            
            {/* Responsive Filter Box */}
            <Box sx={{ 
              display: 'flex', 
              flexDirection: { xs: 'column', sm: 'row' }, 
              gap: 2, 
              mb: 3
            }}>
              <TextField
                label="Filter by Race Name"
                value={filterRaceName}
                onChange={handleRaceNameChange}
                sx={{ flexGrow: 1 }}
                fullWidth
              />
              <Box sx={{ display: 'flex', gap: 2, minWidth: { xs: '100%', sm: '280px' } }}>
                <FormControl sx={{ flexGrow: 1, minWidth: 120 }}>
                  <InputLabel id="year-select-label">Season</InputLabel>
                  <Select
                    labelId="year-select-label"
                    id="year-select"
                    value={selectedYear}
                    label="Season"
                    onChange={handleYearChange}
                  >
                    <MenuItem value="">
                      <em>All</em>
                    </MenuItem>
                    {years.map(year => (
                      <MenuItem key={year} value={year}>{year}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <Button variant="outlined" onClick={handleReset} sx={{ height: '56px', minWidth: '90px' }}>
                  Reset
                </Button>
              </Box>
            </Box>

            {averagePercentBack !== null && (
              <Typography variant="h6" gutterBottom>
                Average Percent Back: {averagePercentBack}%
              </Typography>
            )}

            {/* Mobile-optimized Responsive Chart Container */}
            <ResponsiveContainer width="100%" height={chartHeight}>
              <LineChart data={chartData} margin={{ top: 10, right: 10, left: -25, bottom: 20 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="RaceDate" 
                  tickFormatter={formatDate} 
                  tick={{ fontSize: 10, angle: -45, textAnchor: 'end' }} 
                  height={50} 
                />
                <YAxis tick={{ fontSize: 10 }} />
                <Tooltip />
                <Legend wrapperStyle={{ fontSize: 12, paddingTop: 10 }} />
                <Line type="monotone" dataKey="PercentBack" stroke="#8884d8" name="Percent Back" strokeWidth={2} />
                <Line type="monotone" dataKey="Trend" stroke="#82ca9d" dot={false} strokeDasharray="5 5" name="Trend Line" strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
};

export default RaceChart;
