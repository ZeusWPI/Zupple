import { LoadingSpinner } from "@/components/molecules/LoadingSpinner";
import { useCan } from "@/lib/hooks/useCan";
import { Forbidden } from "@/pages/Forbidden";
import { PropsWithChildren } from "react";

type Props = PropsWithChildren<{
  permission: string;
}>

export const PermissionLayout = ({ permission, children }: Props) => {
  const { can, isLoading } = useCan(permission);

  if (isLoading) {
    return <LoadingSpinner />
  }

  if (!can) {
    return <Forbidden />
  }

  return children;
}
