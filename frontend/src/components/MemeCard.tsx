// components/MemeCard.tsx
import { Meme } from "@/types/meme"
import Image from "next/image"
import MemeVoteButtons from "@/components/MemeVoteButtons"
import BidButton from "@/components/BidButton";


export default function MemeCard({ meme }: { meme: Meme }) {

    return (
        <div className="bg-card border rounded-xl shadow p-4">
            <div className="flex items-center justify-between mb-2">
                <h2 className="text-lg font-semibold">{meme.Title}</h2>
                <span className="text-sm text-muted-foreground">{meme.upvotes}🔥</span>
            </div>
            <Image
                src={meme.ImageURL}
                alt={meme.Caption}
                width={400}
                height={300}
                className="rounded-lg object-cover w-full h-auto"
            />
            <p className="mt-2 text-sm text-muted-foreground">{meme.Caption}</p>
            <div className="flex gap-2 mt-2 flex-wrap">
                {meme.Tags.map(tag => (
                    <span
                        key={tag}
                        className="px-2 py-0.5 bg-muted text-xs rounded border border-muted-foreground/20"
                    >
                        #{tag}
                    </span>
                ))}
            </div>
            <div className="flex gap-2 mt-3">
                <MemeVoteButtons
                    memeId={meme.ID}
                    upvotes={meme.upvotes}
                    downvotes={meme.downvotes}
                    userVote={meme.userVote}
                    type="up"
                />
                <MemeVoteButtons
                    memeId={meme.ID}
                    upvotes={meme.upvotes}
                    downvotes={meme.downvotes}
                    userVote={meme.userVote}
                    type="down"
                />
                <BidButton memeId={meme.ID} />
            </div>
        </div>
    )
}
