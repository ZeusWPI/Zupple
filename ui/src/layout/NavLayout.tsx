import { LinkButton } from "@/components/atoms/LinkButton";
import { AppShell, Burger, Container, Group, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LinkProps } from "@tanstack/react-router";
import { PropsWithChildren } from "react";

type Route = {
  title: string;
  link: LinkProps;
}

const routes: Route[] = [
  {
    title: "Home",
    link: { to: "/" },
  },
]

const NavLink = ({ route, closeNavbar }: { route: Route, closeNavbar?: () => void }) => {
  return (
    <div onClick={closeNavbar}>
      <LinkButton to={route.link.to} activeProps={{ c: "red" }} c="black" variant="subtle" tt="uppercase" radius="xl" className="font-black tracking-wide">
        {route.title}
      </LinkButton>
    </div>
  )
};

export const NavLayout = ({ children }: PropsWithChildren) => {
  const [opened, { close, toggle }] = useDisclosure();

  return (
    <AppShell header={{ height: 80 }} navbar={{ width: 300, breakpoint: "lg", collapsed: { desktop: true, mobile: !opened } }} padding="md">
      <AppShell.Header className="shadow-md">
        <Group justify="space-between" h="100%" px="md" wrap="nowrap" align="center">
          <Burger opened={opened} onClick={toggle} hiddenFrom="lg" size="sm" />
          <Group gap="xs" visibleFrom="lg" align="center">
            {routes.map(route => <NavLink key={route.title} route={route} />)}
          </Group>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <Stack align="flex-start">
          {routes.map(route => <NavLink key={route.title} route={route} closeNavbar={close} />)}
        </Stack>
      </AppShell.Navbar>

      <AppShell.Main className="pt-20 pb-0 overflow-hidden">
        <Container fluid className="container mx-auto">
          {children}
        </Container>
      </AppShell.Main>

    </AppShell>
  )
}
