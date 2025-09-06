import { createRootRoute } from "@tanstack/react-router";

import { createGrafanaRoute } from "~/features/grafana";
import { createAppLayoutRoute } from "~/features/layout";
import { createUploadRoute } from "~/features/upload";

const RootRoute = createRootRoute();

const AppLayoutRoute = createAppLayoutRoute(RootRoute);

const UploadRoute = createUploadRoute(AppLayoutRoute);
const GrafanaRoute = createGrafanaRoute(AppLayoutRoute);

export const routeTree = RootRoute.addChildren([
  AppLayoutRoute.addChildren([UploadRoute, GrafanaRoute]),
]);
