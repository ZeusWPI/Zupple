
import { createContext } from "react";
import { UserPermissions } from "../types/user";

interface PermissionContextType {
  permissions: UserPermissions | null;
  isLoading: boolean;
}

export const PermissionContext = createContext<PermissionContextType | undefined>(undefined)
