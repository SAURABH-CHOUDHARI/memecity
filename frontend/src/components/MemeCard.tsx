// components/MemeCard.tsx
import { Meme } from "@/types/meme"
import Image from "next/image"
import MemeVoteButtons from "@/components/MemeVoteButtons"
import BidButton from "@/components/BidButton"

interface MemeCardProps {
    meme: Meme
    showVoting?: boolean
    showBidding?: boolean
}

export default function MemeCard({ meme, showVoting = true, showBidding = true }: MemeCardProps) {
    return (
        <div className="bg-black/80 border border-green-500 rounded-xl shadow-[0_0_10px_#0f0] hover:shadow-[0_0_15px_#0f0] transition-shadow duration-200 overflow-hidden text-green-400 font-mono">
            {/* Header */}
            <div className="flex items-center justify-between p-4 pb-2">
                <h2 className="text-lg font-bold truncate flex-1 mr-2 text-green-300">
                    {meme.Title}
                </h2>
                <div className="flex items-center gap-3 text-sm">
                    <span className="flex items-center gap-1 text-lime-400 font-bold">
                        {meme.upvotes} 🔥
                    </span>
                    {meme.downvotes > 0 && (
                        <span className="flex items-center gap-1 text-red-500 font-bold">
                            {meme.downvotes} 👎
                        </span>
                    )}
                </div>
            </div>

            {/* Image */}
            <div className="px-4">
                <div className="relative rounded-lg overflow-hidden bg-zinc-800">
                    <Image
                        src={meme.ImageURL}
                        alt={meme.Caption}
                        width={400}
                        height={300}
                        className="object-cover w-full h-auto max-h-80 contrast-125 saturate-150"
                        priority={false}
                    />
                </div>
            </div>

            {/* Content */}
            <div className="p-4 pt-3">
                {/* Caption */}
                <p className="text-sm text-green-500 mb-3 leading-relaxed italic">
                    {meme.Caption}
                </p>

                {/* Tags */}
                {meme.Tags && meme.Tags.length > 0 && (
                    <div className="flex gap-1.5 mb-4 flex-wrap">
                        {meme.Tags.map(tag => (
                            <span
                                key={tag}
                                className="px-2 py-1 bg-green-900 text-xs rounded-md border border-green-500/50 text-green-300 hover:bg-green-800 transition-colors"
                            >
                                #{tag}
                            </span>
                        ))}
                    </div>
                )}

                {/* Action Buttons */}
                <div className="flex items-center justify-between">
                    {/* Voting Buttons */}
                    {showVoting ? (
                        <div className="flex gap-2">
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
                        </div>
                    ) : (
                        <div className="flex gap-2">
                            <span className="text-xs text-green-600 px-3 py-2">
                                {!showVoting && meme.OwnerID ? 'Your meme' : 'Login to interact'}
                            </span>
                        </div>
                    )}

                    {/* Bid Button */}
                    {showBidding && meme.OnSale && (
                        <BidButton memeId={meme.ID} />
                    )}

                    {/* Price/Sale Status */}
                    {meme.OnSale && (
                        <div className="text-right">
                            <span className="text-xs bg-green-700 text-green-200 px-2 py-1 rounded-md font-medium">
                                {meme.Price > 0 ? `₹${meme.Price}` : 'Free'}
                            </span>
                        </div>
                    )}
                </div>

                {/* Metadata */}
                <div className="flex justify-between items-center mt-3 pt-3 border-t border-green-600/40">
                    <span className="text-xs text-green-600">
                        {new Date(meme.CreatedAt).toLocaleDateString('en-IN', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric'
                        })}
                    </span>
                    <span className="text-xs text-green-600">
                        ID: {meme.ID.slice(0, 8)}...
                    </span>
                </div>
            </div>
        </div>
    )
}
