// app/add-meme/page.tsx
import MemeUploadForm from "@/components/MemeUploadForm";

export const metadata = {
    title: "Add Meme | MemeCity",
};

export default function AddMemePage() {
    return (
        <div className="max-w-2xl mx-auto p-6">
            <h1 className="text-2xl font-bold mb-6">Upload Your Meme</h1>
            <MemeUploadForm />
        </div>
    );
}
