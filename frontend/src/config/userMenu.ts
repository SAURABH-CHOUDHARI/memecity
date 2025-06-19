// config/userMenu.ts
import { UserMenuItem } from "@/types/navigation"

export const createUserMenuItems = (
    userId: string,
    onProfileClick: (userId: string) => void,
    onLogout: () => void
): UserMenuItem[] => [
    {
        label: "Profile",
        action: () => onProfileClick(userId),
        icon: "👤"
    },
    {
        label: "My Memes",
        action: () => onProfileClick(userId), // You can change this to navigate to user's memes
        icon: "🖼️"
    },
    {
        label: "Settings",
        action: () => {}, // Add settings logic later
        icon: "⚙️"
    },
    {
        label: "Logout",
        action: onLogout,
        icon: "🚪",
        variant: 'destructive' as const
    }
]
