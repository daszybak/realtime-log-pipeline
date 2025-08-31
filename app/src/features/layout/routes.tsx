import { type AnyRoute, createRoute, redirect } from "@tanstack/react-router";

import { AppLayout } from "./AppLayout";

export function createAppLayoutRoute<P extends AnyRoute>(parent: P) {
  return createRoute({
    getParentRoute: () => parent,
    path: "/",
    component: AppLayout,
    beforeLoad: ({ location }) => {
      if (location.pathname === "/") {
        throw redirect({
          to: "/grafana",
          replace: true,
        });
      }
    },
  });
}
