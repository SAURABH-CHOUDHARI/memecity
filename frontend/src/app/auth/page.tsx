"use client"

import axios, { AxiosError } from "axios"
import { useRouter } from "next/navigation"
import { toast } from "sonner"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card"
import {
    SignupSchema,
    LoginSchema,
    SignupForm,
    LoginForm,
} from "@/lib/authSchema"
import { useUserStore } from "@/store/useUserStore"


export default function AuthPage() {
    const router = useRouter()
    const [tab, setTab] = useState<"login" | "signup">("login")
    const [loading, setLoading] = useState(false)

    const signupForm = useForm<SignupForm>({
        resolver: zodResolver(SignupSchema),
    })

    const loginForm = useForm<LoginForm>({
        resolver: zodResolver(LoginSchema),
    })

    const API_URL = process.env.NEXT_PUBLIC_API_URL

    const onSignup = async (data: SignupForm) => {
        try {
            setLoading(true)
            const res = await axios.post(`${API_URL}/users`, data)

            localStorage.setItem("token", res.data.token)
            localStorage.setItem("user", JSON.stringify(res.data.user))

            useUserStore.getState().setUser(res.data.user)

            toast.success("Account created! Welcome 🎉")
            router.push("/")
        } catch (err) {
            const error = err as AxiosError<{ error: string }>
            const message = error.response?.data?.error || "Signup failed"
            toast.error(`Signup failed: ${message}`)
        } finally {
            setLoading(false)
        }
    }


    const onLogin = async (data: LoginForm) => {
        try {
            setLoading(true)
            const res = await axios.post(`${API_URL}/users/login`, data)

            localStorage.setItem("token", res.data.token)
            localStorage.setItem("user", JSON.stringify(res.data.user))

            useUserStore.getState().setUser(res.data.user)

            toast.success("Welcome back! 👋")
            router.push("/")
        } catch (err) {
            const error = err as AxiosError<{ error: string }>
            const message = error.response?.data?.error || "Login failed"
            toast.error(`Login failed: ${message}`)
        } finally {
            setLoading(false)
        }
    }



    return (
        <div className="min-h-screen flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle className="text-center">MemeCity</CardTitle>
                </CardHeader>
                <CardContent>
                    <Tabs
                        value={tab}
                        onValueChange={(val: string) => setTab(val as "login" | "signup")}
                        className="w-full"
                    >
                        <TabsList className="grid grid-cols-2 mb-4">
                            <TabsTrigger value="login">Login</TabsTrigger>
                            <TabsTrigger value="signup">Signup</TabsTrigger>
                        </TabsList>

                        <TabsContent value="login">
                            <form
                                onSubmit={loginForm.handleSubmit(onLogin)}
                                className="space-y-4"
                            >
                                <Input
                                    placeholder="Email"
                                    {...loginForm.register("email")}
                                />
                                <Input
                                    placeholder="Password"
                                    type="password"
                                    {...loginForm.register("password")}
                                />
                                <Button className="w-full" disabled={loading}>
                                    {loading ? "Logging in..." : "Login"}
                                </Button>
                            </form>
                        </TabsContent>

                        <TabsContent value="signup">
                            <form
                                onSubmit={signupForm.handleSubmit(onSignup)}
                                className="space-y-4"
                            >
                                <Input
                                    placeholder="Email"
                                    {...signupForm.register("email")}
                                />
                                <Input
                                    placeholder="Username"
                                    {...signupForm.register("username")}
                                />
                                <Input
                                    placeholder="Password"
                                    type="password"
                                    {...signupForm.register("password")}
                                />
                                <Button className="w-full" disabled={loading}>
                                    {loading ? "Signing up..." : "Signup"}
                                </Button>
                            </form>
                        </TabsContent>
                    </Tabs>
                </CardContent>
            </Card>
        </div>
    )
}
