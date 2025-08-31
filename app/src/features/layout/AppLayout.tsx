import Box from "@mui/material/Box";
import { Outlet, useMatches } from "@tanstack/react-router";
import { useEffect } from "react";

import { Menu } from "./Menu";

export function AppLayout() {
  const matches = useMatches();
  const crumbs = matches
    .filter((match) => match.routeId !== "/")
    .map((match) => {
      let title = "Unknown Title";
      if (
        Array.isArray(match.meta) &&
        match.meta.length > 0 &&
        match.meta[0]?.title
      ) {
        title = match.meta[0].title;
      }

      return {
        id: match.routeId,
        pathname: match.pathname ?? "#",
        label: title,
      };
    });

  const currentTitle =
    crumbs.length > 0 ? crumbs[crumbs.length - 1].label : "Home";

  useEffect(() => {
    if (currentTitle) {
      document.title = `${currentTitle} – Dashboard`;
    }
  }, [currentTitle]);

  return (
    <Box sx={{ display: "flex", height: "100vh", overflow: "hidden" }}>
      {/* TODO Prevent extra space below the `Menu` when Grafana is open.*/}
      <Menu />
      {/* Ensure the main content area can shrink within the flex container to avoid horizontal overflow */}
      {/*
       * Make the main content area a flex column container so that child components
       * (e.g. pages using MUI DataGrid) can grow and shrink properly within the
       * available viewport height. The additional `minHeight: 0` ensures that the
       * flex item is allowed to shrink below its content's intrinsic height which
       * is required for the DataGrid starting from v6.2+.
       */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          display: "flex",
          flexDirection: "column",
          p: 4,
          minWidth: 0,
          minHeight: 0,
          height: "100%",
          overflow: "auto",
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
}
