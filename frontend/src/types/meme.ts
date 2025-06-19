// types/meme.ts
export type Meme = {
    ID: string
    Title: string
    ImageURL: string
    Tags: string[]
    Caption: string
    Price: number
    OnSale: boolean
    OwnerID: string
    CreatedAt: string
    upvotes: number
    downvotes: number
    userVote?: "up" | "down" | null // Add this optional field
    Owner: {
        ID: string
        Username: string
        Email: string
        Credits: number
        Password: string
        ProfilePic: string
        CreatedAt: string
    }
}