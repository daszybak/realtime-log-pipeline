import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import { Link } from "@tanstack/react-router";

export function Menu() {
  return (
    <Stack
      sx={{
        width: "256px",
        flexShrink: 0,
        alignItems: "flex-start",
        bgcolor: "background.paper",
        boxShadow: 1,
        height: "100vh",
        overflow: "auto",
      }}
    >
      <List component="nav" sx={{ width: "100%" }}>
        <ListItemButton component={Link} to="/grafana">
          <ListItemText primary="Grafana" />
        </ListItemButton>
        <ListItemButton component={Link} to="/upload">
          <ListItemText primary="Upload File" />
        </ListItemButton>
      </List>
    </Stack>
  );
}
