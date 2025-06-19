// components/BidButton.tsx
"use client";

import { useState } from "react";
import axios from "axios";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function BidButton({ memeId }: { memeId: string }) {
    const [open, setOpen] = useState(false);
    const [bidAmount, setBidAmount] = useState("");
    const [loading, setLoading] = useState(false);

    const placeBid = async () => {
        const API_URL = process.env.NEXT_PUBLIC_API_URL;
        const token = localStorage.getItem("token");

        if (!bidAmount || isNaN(Number(bidAmount))) {
            toast.error("Please enter a valid bid amount");
            return;
        }

        setLoading(true);
        try {
            await axios.post(
                `${API_URL}/auth/memes/${memeId}/bid`,
                { credits: Number(bidAmount) },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );
            toast.success("Bid placed successfully!");
            setOpen(false);
            setBidAmount("");
        } catch (err) {
            console.error("Bid failed", err);
            toast.error("Failed to place bid.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button onClick={() => setOpen(true)} variant="outline" size="sm">
                💰 Place Bid
            </Button>

            <Dialog open={open} onOpenChange={setOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Enter your bid</DialogTitle>
                    </DialogHeader>

                    <Input
                        placeholder="e.g. 100"
                        value={bidAmount}
                        onChange={(e) => setBidAmount(e.target.value)}
                        type="number"
                    />

                    <DialogFooter className="mt-4">
                        <Button onClick={placeBid} disabled={loading}>
                            {loading ? "Placing..." : "Confirm Bid"}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
