import { LinkButton } from "@/components/atoms/LinkButton"

export const Index = () => {
  return (
    <div className="flex flex-col items-center pt-20 px-4">
      <h1 className="text-4xl font-bold mb-6 text-center">Hi!</h1>

      <div className="max-w-2xl text-left space-y-4">
        <p className="text-lg">Good luck with the new project!</p>

        <p className="text-lg">
          To test whether your permissions in Auth are working properly,
          you can press the button below to visit a page that requires you to have the <code className="font-mono bg-gray-100 px-1 rounded">example scope</code>.
        </p>

        <p className="text-lg">
          If you don’t have it, a <strong>403 Forbidden</strong> error will be shown.
        </p>

        <p className="text-lg font-semibold mt-6">
          There are two ways to give yourself the permission:
        </p>

        <ol className="list-decimal list-inside space-y-4 text-lg pl-4">
          <li>
            <strong>Backend approach:</strong> Add the scope to the permissions module in the backend.
            After restarting, the backend will push the new scope to Auth.
            If you also assign yourself to the correct role, log in and out again then you’ll have access to the page.
            Om uit te loggen kan je bij de 403 error op <i>Naar de startpagina</i> klikken.
          </li>
          <li>
            <strong>Manual approach:</strong> Go to <code className="font-mono bg-gray-100 px-1 rounded">http://localhost:3003</code> and create it manually in the UI.
            I strongly recommend the first approach so you get a better feel for how the system works.
          </li>
        </ol>
      </div>

      <LinkButton to="/permission" className="mt-12">
        Go to the permission example page
      </LinkButton>
    </div>
  )
}

