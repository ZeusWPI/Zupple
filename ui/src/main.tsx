import "./lib/instrument";
import React from "react";
import ReactDOM from "react-dom/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./lib/api/query";
import { someTheme } from "./lib/theme/theme";
import { MantineProvider } from "@mantine/core";

import "./index.css";
import "@mantine/core/styles.layer.css";
import "@mantine/notifications/styles.layer.css";
import "@mantine/dates/styles.css";
import { router } from "./router";
import { RouterProvider } from "@tanstack/react-router";
import { AuthProvider } from "./lib/providers/AuthProvider";
import { Notifications } from "@mantine/notifications";
import { PermissionProvider } from "./lib/providers/PermissionProvider";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={someTheme}>
        <AuthProvider>
          <PermissionProvider>
            <Notifications />
            <RouterProvider router={router} />
          </PermissionProvider>
        </AuthProvider>
      </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>
);
