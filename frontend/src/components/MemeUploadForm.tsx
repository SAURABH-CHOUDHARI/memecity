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
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input placeholder="Meme Title" {...register("title")} />
            {errors.title && <p className="text-sm text-red-500">{errors.title.message}</p>}

            <Textarea placeholder="Give Caption or let AI take care of it..." {...register("caption")} />
            {errors.caption && <p className="text-sm text-red-500">{errors.caption.message}</p>}

            <Input placeholder="Image URL" {...register("image_url")} />
            {errors.image_url && (
                <p className="text-sm text-red-500">{errors.image_url.message}</p>
            )}

            <Input placeholder="Tags (comma-separated)" {...register("tags")} />

            <Button type="submit" disabled={loading}>
                {loading ? "Uploading..." : "Upload Meme"}
            </Button>

            {msg && <p className="text-sm text-muted-foreground">{msg}</p>}
        </form>
    );
}
