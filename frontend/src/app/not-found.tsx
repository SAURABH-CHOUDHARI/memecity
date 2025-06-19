"use client"

import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Home, ArrowLeft } from "lucide-react"

export default function NotFound() {
    const router = useRouter()

    return (
        <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20 flex items-center justify-center p-4 relative overflow-hidden">
            {/* Background decoration */}
            <div className="absolute inset-0 pointer-events-none opacity-30">
                <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-primary/10 rounded-full blur-3xl"></div>
                <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-secondary/10 rounded-full blur-3xl"></div>
            </div>

            <div className="text-center space-y-8 max-w-2xl mx-auto relative z-10">
                {/* 404 Number */}
                <div className="relative">
                    <h1 className="text-6xl sm:text-8xl md:text-9xl font-bold text-primary/80 select-none">
                        404
                    </h1>
                    <div className="absolute inset-0 text-6xl sm:text-8xl md:text-9xl font-bold text-primary/20">
                        404
                    </div>
                </div>

                {/* Title */}
                <h2 className="text-xl sm:text-2xl md:text-3xl font-semibold text-foreground">
                    Page Not Found
                </h2>

                {/* Description Card */}
                <Card className="border-muted bg-card/50 backdrop-blur-sm">
                    <CardContent className="p-4 sm:p-6">
                        <p className="text-muted-foreground text-sm sm:text-base md:text-lg leading-relaxed">
                            Oops! The page youre looking for seems to have wandered off into the digital void. 
                            It might have been moved, deleted, or perhaps it never existed at all.
                        </p>
                    </CardContent>
                </Card>

                {/* Action Buttons */}
                <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center">
                    <Button
                        onClick={() => router.push("/")}
                        size="lg"
                        className="w-full sm:w-auto"
                    >
                        <Home className="w-4 h-4 mr-2" />
                        Go Home
                    </Button>

                    <Button
                        onClick={() => router.back()}
                        variant="outline"
                        size="lg"
                        className="w-full sm:w-auto"
                    >
                        <ArrowLeft className="w-4 h-4 mr-2" />
                        Go Back
                    </Button>

                </div>

                {/* Static dots */}
                <div className="flex justify-center space-x-2 mt-8">
                    <div className="w-2 h-2 bg-primary rounded-full"></div>
                    <div className="w-2 h-2 bg-primary/60 rounded-full"></div>
                    <div className="w-2 h-2 bg-primary/30 rounded-full"></div>
                </div>
            </div>


        </div>
    )
}