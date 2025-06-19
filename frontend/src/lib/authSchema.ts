import { z } from "zod"

export const SignupSchema = z.object({
    email: z.string().email(),
    username: z.string().min(3, "Username too short"),
    password: z.string().min(6, "Password too short"),
})

export const LoginSchema = z.object({
    email: z.string().email(),
    password: z.string().min(6),
})

export type SignupForm = z.infer<typeof SignupSchema>
export type LoginForm = z.infer<typeof LoginSchema>
