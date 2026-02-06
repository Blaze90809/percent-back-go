import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Container, Typography, Box, Button, Card, CardContent, TableSortLabel } from '@mui/material';

const RaceTable = () => {
  const [races, setRaces] = useState([]);
  const [error, setError] = useState('');
  const [sortConfig, setSortConfig] = useState({ key: 'RaceDate', direction: 'desc' });

  const fetchRaces = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get('/api/races', { headers: { Authorization: `Bearer ${token}` } });
      const sortedRaces = response.data.sort((a, b) => new Date(b.RaceDate) - new Date(a.RaceDate));
      setRaces(sortedRaces);
    } catch (err) {
      setError('Failed to fetch races');
    }
  };

  useEffect(() => {
    fetchRaces();
  }, []);

  const handleDelete = async (id) => {
    try {
      const token = localStorage.getItem('token');
      await axios.delete(`/api/races/delete/${id}`, { headers: { Authorization: `Bearer ${token}` } });
      setRaces(races.filter(race => race.ID !== id));
    } catch (err) {
      setError('Failed to delete race');
    }
  };

  const handleSort = (key) => {
    let direction = 'asc';
    if (sortConfig.key === key && sortConfig.direction === 'asc') {
      direction = 'desc';
    }
    setSortConfig({ key, direction });

    const sortedRaces = [...races].sort((a, b) => {
      if (a[key] < b[key]) {
        return direction === 'asc' ? -1 : 1;
      }
      if (a[key] > b[key]) {
        return direction === 'asc' ? 1 : -1;
      }
      return 0;
    });
    setRaces(sortedRaces);
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Card>
          <CardContent>
            <Typography variant="h4" component="h1" gutterBottom>
              Race Data
            </Typography>
            {error && <Typography color="error">{error}</Typography>}
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell onClick={() => handleSort('RaceName')}>
                      <TableSortLabel
                        active={sortConfig.key === 'RaceName'}
                        direction={sortConfig.key === 'RaceName' ? sortConfig.direction : 'asc'}
                      >
                        Race Name
                      </TableSortLabel>
                    </TableCell>
                    <TableCell onClick={() => handleSort('RaceDate')}>
                      <TableSortLabel
                        active={sortConfig.key === 'RaceDate'}
                        direction={sortConfig.key === 'RaceDate' ? sortConfig.direction : 'asc'}
                      >
                        Date
                      </TableSortLabel>
                    </TableCell>
                    <TableCell onClick={() => handleSort('RaceDistance')}>
                      <TableSortLabel
                        active={sortConfig.key === 'RaceDistance'}
                        direction={sortConfig.key === 'RaceDistance' ? sortConfig.direction : 'asc'}
                      >
                        Distance (km)
                      </TableSortLabel>
                    </TableCell>
                    <TableCell onClick={() => handleSort('PercentBack')}>
                      <TableSortLabel
                        active={sortConfig.key === 'PercentBack'}
                        direction={sortConfig.key === 'PercentBack' ? sortConfig.direction : 'asc'}
                      >
                        Percent Back
                      </TableSortLabel>
                    </TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {races.map((race) => (
                    <TableRow key={race.ID}>
                      <TableCell>{race.RaceName}</TableCell>
                      <TableCell>{race.RaceDate}</TableCell>
                      <TableCell>{race.RaceDistance}</TableCell>
                      <TableCell>{race.PercentBack}</TableCell>
                      <TableCell>
                        <Button 
                          variant="outlined" 
                          color="error" 
                          onClick={() => handleDelete(race.ID)}
                        >
                          Delete
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
};

export default RaceTable;
