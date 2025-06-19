import { create } from "zustand"

interface User {
    id: string
    email: string
    username: string
    profilePic?: string
}

interface UserStore {
    user: User | null
    setUser: (user: User | null) => void
    loadUserFromStorage: () => void
}

export const useUserStore = create<UserStore>((set) => ({
    user: null,
    setUser: (user) => {
        if (user) {
            localStorage.setItem("user", JSON.stringify(user))
        } else {
            localStorage.removeItem("user")
        }
        set({ user })
    },
    loadUserFromStorage: () => {
        try {
            const userJSON = localStorage.getItem("user")
            if (userJSON) {
                const user = JSON.parse(userJSON)
                set({ user })
            }
        } catch (err) {
            console.error("Failed to load user from storage", err)
            set({ user: null })
        }
    },
}))
