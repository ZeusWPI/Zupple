import { PermissionLayout } from "@/layout/PermissionLayout"

export const PermissionExample = () => {
  return (
    <PermissionLayout permission="example">
      <div className="items-center pt-20">
        <p className="text-center">
          You have the correct permissions!
        </p>
      </div>
    </PermissionLayout>
  )
}
