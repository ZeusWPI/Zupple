import { PropsWithChildren } from "react"
import { useUserPermissions } from "../api/user"
import { PermissionContext } from "../contexts/permissionContext"
import { useAuth } from "../hooks/useAuth"

export const PermissionProvider = ({ children }: PropsWithChildren) => {
  const { user } = useAuth()
  const { data, isLoading } = useUserPermissions(user)

  return (
    <PermissionContext
      value={{
        permissions: data ?? null,
        isLoading,
      }}
    >
      {children}
    </PermissionContext>
  )
}
