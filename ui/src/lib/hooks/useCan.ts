import { useContext } from "react"
import { PermissionContext } from "../contexts/permissionContext";

export const useCan = (permission: string): { can: boolean, isLoading: boolean } => {
  const context = useContext(PermissionContext);
  if (!context) {
    throw new Error("useCan must be used within a PermissionProvider");
  }

  if (permission === "") {
    return { can: true, isLoading: false }
  }

  const { permissions, isLoading } = context;

  if (isLoading || !permissions) {
    return { can: false, isLoading }
  }

  return { can: permissions.includes(permission), isLoading }
}
