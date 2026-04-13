import { API } from "./api";
import { z } from "zod";
import { JSONBody } from "./general";

export interface User {
  id: number;
  name: string;
  email: string;
}

export type UserPermissions = string[]

// Converters

export const convertUserToModel = (user: API.User): User => {
  return userSchema.parse(user);
}

// Schemas

export const userSchema = z.object({
  id: z.number(),
  name: z.string().nonempty(),
  email: z.string(),
})
export type NewsSchema = z.infer<typeof userSchema> & JSONBody
