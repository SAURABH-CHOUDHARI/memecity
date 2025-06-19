"use client"

import { useEffect, useRef } from "react"

export default function MatrixBackground() {
    const canvasRef = useRef<HTMLCanvasElement>(null)

    useEffect(() => {
        const canvas = canvasRef.current
        const ctx = canvas?.getContext("2d")
        if (!canvas || !ctx) return

        // Full screen canvas
        canvas.width = window.innerWidth
        canvas.height = window.innerHeight

        const cols = Math.floor(canvas.width / 20)
        const chars = "アァカサタナハマヤャラワン0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        const matrix: number[] = Array(cols).fill(0)

        const draw = () => {
            ctx.fillStyle = "rgba(0, 0, 0, 0.05)"
            ctx.fillRect(0, 0, canvas.width, canvas.height)

            ctx.fillStyle = "#00FF00"
            ctx.font = "16px monospace"

            for (let i = 0; i < matrix.length; i++) {
                const char = chars.charAt(Math.floor(Math.random() * chars.length))
                ctx.fillText(char, i * 20, matrix[i] * 20)

                if (matrix[i] * 20 > canvas.height && Math.random() > 0.975) {
                    matrix[i] = 0
                }
                matrix[i]++
            }
        }

        const interval = setInterval(draw, 50)
        return () => clearInterval(interval)
    }, [])

    return (
        <canvas
            ref={canvasRef}
            className="fixed top-0 left-0 w-full h-full -z-10 opacity-50"
        />
    )
}
