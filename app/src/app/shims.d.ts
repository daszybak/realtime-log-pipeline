/// <reference types="vite/client" />

import "@tanstack/react-query";
import "@tanstack/router-core";

import type { ErrorDto } from "@gen/app-query";
import type { AxiosError } from "axios";

declare module "*.module.css" {
  const classes: { readonly [key: string]: string };
  export default classes;
}

export interface ToastMeta extends Record<string, unknown> {
  success?: string;
}

declare module "@tanstack/react-query" {
  interface Register {
    mutationMeta: ToastMeta;
    queryMeta: ToastMeta;
    defaultError: Omit<AxiosError, "response"> & { response: ErrorDto };
  }
}

declare module "@tanstack/router-core" {
  interface UpdatableRouteOptionsExtensions {
    meta?: {
      title?: string;
    };
  }
}
