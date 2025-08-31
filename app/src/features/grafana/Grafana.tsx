import { Box, Typography } from '@mui/material';

export function Grafana () {
  const grafanaUrl = 'http://localhost:3000/d/abc123/my-dashboard?orgId=1&theme=light&kiosk';

  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="h4" gutterBottom>
        Monitoring Dashboard
      </Typography>
      <iframe
        src={grafanaUrl}
        title="Grafana Dashboard"
        width="100%"
        height="700px"
      />
    </Box>
  );
};

