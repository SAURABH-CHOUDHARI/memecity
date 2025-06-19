// Updated Navbar component
"use client"

import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { useUserStore } from "@/store/useUserStore"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Skeleton } from "@/components/ui/skeleton"
import axios from "axios"
import { navItems } from "@/config/navigation"
import { createUserMenuItems } from "@/config/userMenu"
import { ThemeToggle } from "@/components/theme-toggle"
import { useHydration } from "@/hooks/useHydration"
import FuzzyText from "./Bits/FuzzyText"

export default function Navbar() {
    const pathname = usePathname()
    const router = useRouter()
    const hydrated = useHydration()
    const { user, setUser } = useUserStore()

    const logout = async () => {
        const API_URL = process.env.NEXT_PUBLIC_API_URL
        const token = localStorage.getItem("token")
        const res = await axios.post(
            `${API_URL}/auth/users/logout`,
            {},
            { headers: { Authorization: `Bearer ${token}` } }
        )
        if (res.data.message) {
            localStorage.removeItem("token")
            setUser(null)
            toast.success("Logged out")
            router.push("/auth")
        }
    }

    const userMenuItems = user ? createUserMenuItems(
        user.id,
        (userId) => router.push(`/profile/${userId}`),
        logout
    ) : []

    return (
        <nav className="w-full border-b border-border bg-background px-4 py-3 flex justify-between items-center shadow-sm">
            <Link href="/" className="text-lg font-bold tracking-wide">
                <FuzzyText
                    baseIntensity={0.1}
                    hoverIntensity={0.3}
                    enableHover={true}
                    fontSize={30}
                >
                    MemeCity
                </FuzzyText>
            </Link>

            <div className="flex items-center gap-4">
                {navItems.map(({ href, label, icon }) => (
                    <Link
                        key={href}
                        href={href}
                        className={`text-sm font-medium transition-colors flex items-center gap-1 ${pathname === href
                            ? "text-primary"
                            : "text-muted-foreground hover:text-primary"
                            }`}
                    >
                        {icon && <span className="text-xs">{icon}</span>}
                        {label}
                    </Link>
                ))}

                <ThemeToggle />

                {!hydrated ? (
                    // Loading state during hydration
                    <Skeleton className="h-8 w-8 rounded-full" />
                ) : user ? (
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Avatar className="h-10 w-10 cursor-pointer">
                                <AvatarImage src={user.profilePic || ""} alt={user.username} />
                                <AvatarFallback>{user.username?.[0]?.toUpperCase() || "U"}</AvatarFallback>
                            </Avatar>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-40">
                            {userMenuItems.map(({ label, action, icon, variant }) => (
                                <DropdownMenuItem
                                    key={label}
                                    onClick={action}
                                    className={variant === 'destructive' ? 'text-destructive' : ''}
                                >
                                    {icon && <span className="mr-2 text-xs">{icon}</span>}
                                    {label}
                                </DropdownMenuItem>
                            ))}
                        </DropdownMenuContent>
                    </DropdownMenu>
                ) : (
                    <Button size="sm" onClick={() => router.push("/auth")}>
                        Login
                    </Button>
                )}
            </div>
        </nav>
    )
}