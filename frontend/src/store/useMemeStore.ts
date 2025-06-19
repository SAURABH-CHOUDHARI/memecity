// store/useMemeStore.ts
import { create } from "zustand"
import { Meme } from "@/types/meme"

type MemeStore = {
    memes: Meme[]
    setMemes: (data: Meme[]) => void
    updateVotes: (
        id: string, 
        update: {
            upvotes: number
            downvotes: number
            userVote?: "up" | "down" | null
        }
    ) => void
}

export const useMemeStore = create<MemeStore>((set) => ({
    memes: [],
    setMemes: (data) => set({ memes: data }),
    updateVotes: (id, update) =>
        set((state) => ({
            memes: state.memes.map((m) =>
                m.ID === id
                    ? {
                        ...m,
                        upvotes: update.upvotes,
                        downvotes: update.downvotes,
                        userVote: update.userVote
                    }
                    : m
            ),
        })),
}))