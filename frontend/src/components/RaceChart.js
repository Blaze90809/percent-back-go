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

  const getSeason = (dateStr) => {
    const date = new Date(dateStr);
    // If month is July (6) or later, it belongs to the next year's season
    return date.getMonth() >= 6 ? date.getFullYear() + 1 : date.getFullYear();
  };

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

        if (sortedRaces.length > 0) {
          // Initial calculation
          calculateAverage(sortedRaces);
        }
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
    } else {
      setChartData([]);
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
    // The filtering logic is now handled by the useEffect hook above,
    // which depends on selectedYear, so we just update the state.
  };

  const handleReset = () => {
    setSelectedYear('');
    setFilterRaceName('');
    // The filtering logic is now handled by the useEffect hook above,
    // which depends on selectedYear and filterRaceName, so we just update the state.
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Card>
          <CardContent>
            <Typography variant="h4" component="h1" gutterBottom>
              Percent Back Over Time
            </Typography>
            {error && <Typography color="error">{error}</Typography>}
            
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <TextField
                label="Filter by Race Name"
                value={filterRaceName}
                onChange={(e) => handleRaceNameChange(e)}
                sx={{ minWidth: 250 }}
              />
              <FormControl sx={{ minWidth: 120, mr: 2 }}>
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
              <Button variant="outlined" onClick={handleReset}>Reset</Button>
            </Box>

            {averagePercentBack !== null && (
              <Typography variant="h6" gutterBottom>
                Average Percent Back: {averagePercentBack}%
              </Typography>
            )}

            <ResponsiveContainer width="100%" height={400}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="RaceDate" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="PercentBack" stroke="#8884d8" name="Percent Back" />
                <Line type="monotone" dataKey="Trend" stroke="#82ca9d" dot={false} strokeDasharray="5 5" name="Trend Line" />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
};

export default RaceChart;
