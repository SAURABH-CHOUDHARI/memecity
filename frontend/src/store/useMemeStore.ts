// store/useMemeStore.ts
import { create } from 'zustand'
import { Meme } from '@/types/meme'

interface MemeStore {
    memes: Meme[]
    setMemes: (memes: Meme[]) => void
    updateVotes: (id: string, type: 'up' | 'down', action: 'created' | 'flipped' | 'removed') => void
    updateVoteData: (id: string, voteData: { upvotes: number, downvotes: number, userVote: 'up' | 'down' | null }) => void
    addMeme: (meme: Meme) => void
}

export const useMemeStore = create<MemeStore>((set) => ({
    memes: [],

    setMemes: (memes) => set({ memes }),

    updateVotes: (id, type, action) => set((state) => {
        const updated = state.memes.map((meme) => {
            if (meme.ID !== id) return meme

            const delta = action === 'created' ? 1 : action === 'flipped' ? 1 : -1

            if (type === 'up') {
                return {
                    ...meme,
                    upvotes: action === 'flipped' ? meme.upvotes + 1 : meme.upvotes + delta,
                    downvotes: action === 'flipped' ? meme.downvotes - 1 : meme.downvotes,
                }
            } else {
                return {
                    ...meme,
                    downvotes: action === 'flipped' ? meme.downvotes + 1 : meme.downvotes + delta,
                    upvotes: action === 'flipped' ? meme.upvotes - 1 : meme.upvotes,
                }
            }
        })

        return { memes: updated }
    }),

    // New method to directly update vote data from backend
    updateVoteData: (id, voteData) => set((state) => {
        const updated = state.memes.map((meme) => {
            if (meme.ID !== id) return meme

            return {
                ...meme,
                upvotes: voteData.upvotes,
                downvotes: voteData.downvotes,
                userVote: voteData.userVote,
            }
        })

        return { memes: updated }
    }),

    addMeme: (meme) => set((state) => ({ memes: [meme, ...state.memes] }))
}))