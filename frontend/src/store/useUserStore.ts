import { create } from "zustand"
import { persist } from "zustand/middleware"

interface User {
    id: string
    email: string
    username: string
    profilePic: string
    credits: number
}

interface UserStore {
    user: User | null
    setUser: (user: User | null) => void
    updateCredits: (delta: number) => void
    logout: () => void
}

export const useUserStore = create<UserStore>()(
    persist(
        (set) => ({
            user: null,
            setUser: (user) => set({ user }),
            updateCredits: (delta) => set((state) => {
                if (!state.user) return state
                return {
                    user: {
                        ...state.user,
                        credits: state.user.credits + delta,
                    },
                }
            }),
            logout: () => {
                localStorage.removeItem("token")
                set({ user: null})
            }
        }),
        {
            name: "user-storage", // localStorage key
            partialize: (state) => ({
                user: state.user,
            }),
        }
    )
)
