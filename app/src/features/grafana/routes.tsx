import { type AnyRoute, createRoute } from "@tanstack/react-router";
import { Grafana } from "./Grafana";

export function createGrafanaRoute<P extends AnyRoute>(parent: P) {
  return createRoute({
    getParentRoute: () => parent,
    path: "/grafana",
    component: Grafana,
    head: () => ({
      meta: [
        {
          title: "Grafana",
        },
      ],
    }),
  });
}
