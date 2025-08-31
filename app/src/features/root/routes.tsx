import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  createRootRouteWithContext,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

export function createRootRoute() {
  return createRootRouteWithContext()({
    component: () => (
      <>
        {/* NOTE We statically exclude Devtools 
            from getting in the production bundle. */}
        {import.meta.env.DEV && (
          <>
            <TanStackRouterDevtools />
            <ReactQueryDevtools />
          </>
        )}
        <Outlet />
      </>
    ),
    beforeLoad({ location }) {
      if (location.pathname === "/") {
        throw redirect({
          to: "/grafana",
          replace: true,
        });
      }
    },
    head: () => ({
      meta: [
        {
          title: "Dashboard",
        },
      ],
    }),
  });
}
