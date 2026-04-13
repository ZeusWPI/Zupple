import { isResponseNot200Error } from "@/lib/api/query";
import { Button, Container, Title } from "@mantine/core";
import { ErrorComponentProps, useNavigate } from "@tanstack/react-router";
import { Error404 } from "./404";
import { Forbidden } from "./Forbidden";

export const Error = ({ error, reset }: ErrorComponentProps) => {
  const navigate = useNavigate()

  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404:
        return (
          <Container fluid className="pt-[10%]">
            <Error404 />
          </Container>
        )
      case 403:
        return (
          <Container fluid className="pt-[10%]">
            <Forbidden />
          </Container>
        )
    }
  }

  const handleReturn = () => {
    reset()
    navigate({ to: "/" })
  }

  return (
    <div className="flex flex-col justify-center items-center h-full pt-[10%]">
      <p className="font-semibold text-primary">
        500
      </p>
      <Title order={1} className="mt-4 text-balance font-semibold tracking-tight">
        Server Error
      </Title>
      <p className="flex flex-col items-center mt-6 text-pretty text-lg font-medium text-gray-500 sm:text-xl/8">
        <span>Oh no, er is iets misgelopen!</span>
        <br />
        <span>We zijn op de hoogte gebracht.</span>
      </p>
      <div className="mt-10 flex items-center justify-center gap-x-6">
        <Button onClick={handleReturn}>
          Keer terug naar het startscherm
        </Button>
      </div>
    </div>
  )
}
