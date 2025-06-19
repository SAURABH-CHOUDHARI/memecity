// app/leaderboard/page.tsx
"use client"
import axios from "axios"
import { useEffect, useState } from "react"
import { Meme } from "@/types/meme"
import MemeCard from "@/components/MemeCard"
import { useMemeStore } from "@/store/useMemeStore"
import { useUserStore } from "@/store/useUserStore"

export default function LeaderboardPage() {
    const [memes, setMemes] = useState<Meme[]>([])
    const [loading, setLoading] = useState(true)
    const { setMemes: setStoreMemes } = useMemeStore()
    const { user } = useUserStore()

    useEffect(() => {
        async function fetchLeaderboardMemes() {
            try {
                const API_URL = process.env.NEXT_PUBLIC_API_URL
                const res = await axios.get(`${API_URL}/memes/leaderboard?limit=10&offset=0`)
                setMemes(res.data)
                setStoreMemes(res.data)
            } catch (error) {
                console.error("Failed to fetch memes:", error)
            } finally {
                setLoading(false)
            }
        }

        fetchLeaderboardMemes()
    }, [setStoreMemes])

    if (loading) return <div>Loading...</div>

    return (
        <div className="container mx-auto p-4">
            <div className="text-center mb-8">
                <h1 className="text-3xl font-bold">🏆 Leaderboard</h1>
            </div>
            
            <div className="grid grid-cols-3 max-xl:grid-cols-2 max-sm:grid-cols-1 gap-6">
                {memes.map((meme) => {
                    const isOwner = user?.id === meme.OwnerID
                    const isLoggedIn = !!user
                    
                    return (
                        <MemeCard
                            key={meme.ID}
                            meme={meme}
                            showVoting={isLoggedIn && !isOwner}
                            showBidding={isLoggedIn && !isOwner && meme.OnSale}
                        />
                    )
                })}
            </div>
        </div>
    )
}