// types/navigation.ts
export interface NavItem {
    href: string
    label: string
    icon?: string
}

export interface UserMenuItem {
    label: string
    action: () => void
    icon?: string
    variant?: 'default' | 'destructive'
}