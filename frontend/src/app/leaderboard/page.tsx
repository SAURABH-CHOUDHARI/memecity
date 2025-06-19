// app/leaderboard/page.tsx
import axios from "axios"
import { Meme } from "@/types/meme"
import MemeCard from "@/components/MemeCard"
import { useMemeStore } from "@/store/useMemeStore"

export const dynamic = "force-dynamic"

async function fetchLeaderboardMemes(): Promise<Meme[]> {
    const API_URL = process.env.NEXT_PUBLIC_API_URL
    const res = await axios.get(`${API_URL}/memes/leaderboard?limit=10&offset=0`)
    return res.data
}

export default async function LeaderboardPage() {
    const memes = await fetchLeaderboardMemes()

    // Hydrate Zustand store on server-render
    useMemeStore.getState().setMemes(memes)

    return (
        <div className="max-w-4xl mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4">🏆 Leaderboard</h1>
            <div className="grid md:grid-cols-2 gap-4">
                {memes.map((meme) => (
                    <MemeCard key={meme.ID} meme={meme} />
                ))}
            </div>
        </div>
    )
}
