import { useAuth } from "@/lib/hooks/useAuth";
import { Paper, Title, Button, Text, Stack, Center } from "@mantine/core";
import { FaMicrosoft } from "react-icons/fa6";

export const Login = () => {
  const { login } = useAuth();

  return (
    <Center h="100vh">
      <Paper shadow="sm" p="xl" className="w-96">
        <Stack align="center">
          <Title>Login</Title>
          <Text c="gray">Zupple</Text>
          <Button onClick={login} color="red" size="lg" className="my-12">
            <FaMicrosoft size={"1.7rem"} className="mr-2" />
            Office 365
          </Button>
        </Stack>
      </Paper>
    </Center>
  );
};
