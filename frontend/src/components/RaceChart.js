import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Container, Typography, Box, Card, CardContent, Select, MenuItem, FormControl, InputLabel, Button } from '@mui/material';

const RaceChart = () => {
  const [originalRaces, setOriginalRaces] = useState([]);
  const [filteredRaces, setFilteredRaces] = useState([]);
  const [years, setYears] = useState([]);
  const [selectedYear, setSelectedYear] = useState('');
  const [averagePercentBack, setAveragePercentBack] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchRaces = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('/api/races', { headers: { Authorization: `Bearer ${token}` } });
        const races = response.data || [];
        const sortedRaces = races.sort((a, b) => new Date(a.RaceDate) - new Date(b.RaceDate));
        setOriginalRaces(sortedRaces);
        setFilteredRaces(sortedRaces);

        const uniqueYears = [...new Set(sortedRaces.map(race => new Date(race.RaceDate).getFullYear()))];
        setYears(uniqueYears);

        if (sortedRaces.length > 0) {
          calculateAverage(sortedRaces);
        }
      } catch (err) {
        setError('Failed to fetch races');
      }
    };
    fetchRaces();
  }, []);

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

    if (year === '') {
      setFilteredRaces(originalRaces);
      calculateAverage(originalRaces);
    } else {
      const filtered = originalRaces.filter(race => new Date(race.RaceDate).getFullYear() === year);
      setFilteredRaces(filtered);
      calculateAverage(filtered);
    }
  };

  const handleReset = () => {
    setSelectedYear('');
    setFilteredRaces(originalRaces);
    calculateAverage(originalRaces);
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
              <FormControl sx={{ minWidth: 120, mr: 2 }}>
                <InputLabel id="year-select-label">Year</InputLabel>
                <Select
                  labelId="year-select-label"
                  id="year-select"
                  value={selectedYear}
                  label="Year"
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
              <LineChart data={filteredRaces}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="RaceDate" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="PercentBack" stroke="#8884d8" />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
};

export default RaceChart;
