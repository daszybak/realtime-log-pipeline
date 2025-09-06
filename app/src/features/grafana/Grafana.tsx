import { Box, Typography } from "@mui/material";

export function Grafana() {
  // Updated to use the actual dashboard we created
  const grafanaUrl =
    "http://localhost:3000/d/realtime-pipeline-streamer/realtime-log-pipeline-streamer-service?orgId=1&theme=light&kiosk";

  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="h4" gutterBottom>
        Realtime Log Pipeline - Monitoring Dashboard
      </Typography>
      <Typography variant="body1" sx={{ mb: 2 }}>
        Real-time metrics, logs, and traces for the cryptocurrency streaming
        pipeline
      </Typography>
      <iframe
        src={grafanaUrl}
        title="Grafana Dashboard - Streamer Service"
        width="100%"
        height="800px"
        style={{ border: "none", borderRadius: "8px" }}
      />
    </Box>
  );
}
