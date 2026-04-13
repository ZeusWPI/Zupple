import { Outlet, Scripts } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { NavLayout } from "./layout/NavLayout";
import { AuthLayout } from "./layout/AuthLayout";

// <Scripts /> executes any relevant scripts defined in the router

export const App = () => {
  return (
    <>
      <Scripts />
      <AuthLayout>
        <NavLayout>
          <Outlet />
        </NavLayout>
      </AuthLayout>
      <TanStackRouterDevtools />
      <ReactQueryDevtools />
    </>
  )
}
