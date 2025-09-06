import { ClientOnly, RouterProvider } from "@tanstack/react-router";

import { router } from "./router";

export function App() {
  return (
    <ClientOnly>
      <RouterProvider router={router} />
    </ClientOnly>
  );
}
