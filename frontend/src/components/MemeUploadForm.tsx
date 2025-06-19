// components/MemeUploadForm.tsx
"use client";

import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import axios from "axios";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";

const MemeSchema = z.object({
    title: z.string().min(1, "Title is required"),
    image_url: z.string().url("Must be a valid URL"),
    tags: z.string().optional(),
    caption: z.string(),
});

type MemeForm = z.infer<typeof MemeSchema>;

export default function MemeUploadForm() {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<MemeForm>({
        resolver: zodResolver(MemeSchema),
    });

    const [loading, setLoading] = useState(false);
    const [msg, setMsg] = useState("");

    const onSubmit = async (data: MemeForm) => {
        setLoading(true);
        setMsg("");
        try {
            const token = localStorage.getItem("token");
            const tags = data.tags
                ? data.tags.split(",").map((t) => t.trim()).filter(Boolean)
                : [];

            await axios.post(
                `${process.env.NEXT_PUBLIC_API_URL}/auth/memes`,
                {
                    title: data.title,
                    image_url: data.image_url,
                    caption: data.caption,
                    tags,
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            setMsg("✅ Meme uploaded successfully!");
            reset();
        } catch (error) {
            console.error("Upload failed", error);
            setMsg("❌ Failed to upload meme.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 bg-black/90 p-6 rounded-lg border border-green-500 shadow-[0_0_12px_#00ff00] text-green-300 font-mono">
            <Input placeholder="Meme Title" {...register("title")} className="bg-zinc-900 border border-green-600 text-green-200 placeholder:text-green-500" />
            {errors.title && <p className="text-sm text-red-500">{errors.title.message}</p>}

            <Textarea placeholder="Give Caption or let AI take care of it..." {...register("caption")} className="bg-zinc-900 border border-green-600 text-green-200 placeholder:text-green-500" />
            {errors.caption && <p className="text-sm text-red-500">{errors.caption.message}</p>}

            <Input placeholder="Image URL" {...register("image_url")} className="bg-zinc-900 border border-green-600 text-green-200 placeholder:text-green-500" />
            {errors.image_url && (
                <p className="text-sm text-red-500">{errors.image_url.message}</p>
            )}

            <Input placeholder="Tags (comma-separated)" {...register("tags")} className="bg-zinc-900 border border-green-600 text-green-200 placeholder:text-green-500" />

            <Button type="submit" disabled={loading} className="bg-green-700 hover:bg-green-600 text-black font-bold">
                {loading ? "Uploading..." : "Upload Meme"}
            </Button>

            {msg && <p className="text-sm text-green-400">{msg}</p>}
        </form>
    );
}
