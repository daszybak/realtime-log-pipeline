import { type AnyRoute, createRoute } from "@tanstack/react-router";
import { Upload } from "./Upload";

export function createUploadRoute<P extends AnyRoute>(parent: P) {
  return createRoute({
    getParentRoute: () => parent,
    path: "/upload",
    component: Upload,
    head: () => ({
      meta: [
        {
          title: "Upload",
        },
      ],
    }),
  });
}
