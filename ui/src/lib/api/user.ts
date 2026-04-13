import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, NO_CONVERTER, NO_DATA, NO_FILES } from "./query";
import { convertUserToModel, User, UserPermissions } from "../types/user";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

const STALE_30_MIN = 30 * 60 * 1000;
const STALE_5_MIN = 5 * 60 * 1000;

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUserToModel)).data,
    retry: 0,
    staleTime: STALE_30_MIN,
  })
}

export const useUserLogin = () => {
  window.location.href = `/api/${ENDPOINT_AUTH}/login/azureadv2`
}

export const useUserLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => (await apiPost(`${ENDPOINT_AUTH}/logout`, NO_DATA, NO_CONVERTER, NO_FILES)).data,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["user"] })
  })
}

export const useUserPermissions = (user: User | null) => {
  return useQuery({
    queryKey: ["permissions"],
    queryFn: async () => (await apiGet<UserPermissions>(`${ENDPOINT_USER}/permissions`, NO_CONVERTER)).data,
    staleTime: STALE_5_MIN,
    enabled: !!user,
    throwOnError: true,
  })
}
