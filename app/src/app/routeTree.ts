import { createRootRoute } from "@tanstack/react-router";
import { createGrafanaRoute } from "~/features/grafana/routes";
import { createAppLayoutRoute } from "~/features/layout/routes";
import { createUploadRoute } from "~/features/upload/routes";

const RootRoute = createRootRoute();

const AppLayoutRoute = createAppLayoutRoute(RootRoute);

const UploadRoute = createUploadRoute(AppLayoutRoute);
const GrafanaRoute = createGrafanaRoute(AppLayoutRoute);

export const routeTree = RootRoute.addChildren([
  AppLayoutRoute.addChildren([
      UploadRoute,
      GrafanaRoute
  ]),
]);
