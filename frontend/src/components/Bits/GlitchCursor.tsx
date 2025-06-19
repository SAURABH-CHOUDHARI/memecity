// components/GlitchCursor.tsx
"use client"
import { useEffect, useRef } from "react"

export default function GlitchCursor() {
    const canvasRef = useRef<HTMLCanvasElement>(null)

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas) return

        const ctx = canvas.getContext("2d")
        if (!ctx) return

        let mouseX = 0
        let mouseY = 0

        const glitchColors = ["#0ff", "#f0f", "#0f0"]
        const trail: { x: number; y: number; life: number }[] = []

        const draw = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height)

            trail.forEach((p, i) => {
                if (p.life > 0) {
                    ctx.fillStyle = glitchColors[i % glitchColors.length]
                    ctx.fillRect(p.x, p.y, 4, 4)
                    p.life -= 0.1
                }
            })

            trail.push({ x: mouseX, y: mouseY, life: 1.0 })
            if (trail.length > 10) trail.shift()

            requestAnimationFrame(draw)
        }

        const handleMouseMove = (e: MouseEvent) => {
            mouseX = e.clientX
            mouseY = e.clientY
        }

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight
        window.addEventListener("mousemove", handleMouseMove)
        draw()

        return () => {
            window.removeEventListener("mousemove", handleMouseMove)
        }
    }, [])

    return (
        <canvas
            ref={canvasRef}
            style={{
                position: "fixed",
                top: 0,
                left: 0,
                pointerEvents: "none",
                zIndex: 9999,
                mixBlendMode: "difference",
            }}
        />
    )
}
