// config/navigation.ts
import { NavItem } from "@/types/navigation"

export const navItems: NavItem[] = [
    {
        href: "/",
        label: "Home",
        icon: "🏠"
    },
    {
        href: "/leaderboard", 
        label: "Leaderboard",
        icon: "🏆"
    },
    {
        href: "/create",
        label: "Create",
        icon: "➕"
    },
    {
        href: "/marketplace",
        label: "Marketplace", 
        icon: "🛒"
    }
]