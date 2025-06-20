"use client"

import { useState } from "react"
import axios from "axios"
import { Button } from "@/components/ui/button"
import { useMemeStore } from "@/store/useMemeStore"
import { toast } from "sonner"

type Props = {
    memeId: string
    upvotes: number
    downvotes: number
    userVote?: "up" | "down" | null
    type: "up" | "down"
}

export default function MemeVoteButtons({
    memeId,
    upvotes,
    downvotes,
    userVote,
    type,
}: Props) {
    const [isLoading, setIsLoading] = useState(false)
    const { memes, updateVoteData } = useMemeStore()
    
    // Try to get updated data from store, fallback to props
    const storeMeme = memes.find(m => m.ID === memeId)
    const currentUpvotes = storeMeme?.upvotes ?? upvotes
    const currentDownvotes = storeMeme?.downvotes ?? downvotes
    const currentUserVote = storeMeme?.userVote ?? userVote

    const vote = async () => {
        setIsLoading(true)
        const API_URL = process.env.NEXT_PUBLIC_API_URL
        const token = localStorage.getItem("token")

        try {
            const res = await axios.post(
                `${API_URL}/auth/memes/${memeId}/vote`,
                { type },
                { headers: { Authorization: `Bearer ${token}` } }
            )

            if (res.status === 200) {
                const { meme: updatedMeme } = res.data
                
                // Use the new updateVoteData method
                updateVoteData(memeId, {
                    upvotes: updatedMeme.upvotes,
                    downvotes: updatedMeme.downvotes,
                    userVote: updatedMeme.userVote
                })

                toast.success(`${type === "up" ? "Upvoted" : "Downvoted"} successfully!`)
            }
        } catch (err) {
            console.error("Vote failed", err)
            toast.error("Vote failed. Please try again.")
        } finally {
            setIsLoading(false)
        }
    }

    const isActive = currentUserVote === type
    const count = type === "up" ? currentUpvotes : currentDownvotes

    return (
        <Button 
            onClick={vote} 
            variant={isActive ? "default" : "outline"} 
            size="sm"
            disabled={isLoading}
        >
            {type === "up" ? "👍 Upvote" : "👎 Downvote"} ({count})
        </Button>
    )
}