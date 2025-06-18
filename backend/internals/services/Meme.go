package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

func CreateMeme(conn storage.Repository, title, imageURL, caption string, tags []string, ownerID uuid.UUID) error {
	if title == "" {
		return errors.New("title is required")
	}

if caption == "" {
	ctx := context.Background()

	// System instruction for meme caption generation
	systemInstruction := `You are a witty meme caption generator. Your task is to create funny, relatable, and engaging captions for memes based on the provided tags. 

Guidelines:
- Keep captions short and punchy (under 100 characters when possible)
- Use internet humor, trending phrases, and relatable situations
- Make it appropriate for general audiences
- Be creative and match the vibe of the tags
- Don't explain the joke - just deliver it
- Use popular meme formats when relevant (e.g., "When you...", "POV:", "Nobody:", "Me:", etc.)

Generate only the caption text, nothing else.`

	// Use tags directly (assuming tags is []string)
	var tagStrings []string
	if len(tags) > 0 {
		tagStrings = tags
	}

	// Create focused prompt
	var prompt string
	if len(tagStrings) > 0 {
		prompt = fmt.Sprintf("Generate a funny meme caption for tags: %v", tagStrings)
		if title != "" {
			prompt += fmt.Sprintf(" (Meme title: '%s')", title)
		}
	} else if title != "" {
		prompt = fmt.Sprintf("Generate a funny meme caption for a meme titled: '%s'", title)
	} else {
		prompt = "Generate a funny, generic meme caption"
	}

	// Generate with system instruction
	result, err := conn.GeminiClient.Models.GenerateContent(
		ctx, 
		"gemini-2.5-flash", 
		genai.Text(systemInstruction+"\n\nUser request: "+prompt), 
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to generate caption with Gemini: %w", err)
	}

	generatedCaption := strings.TrimSpace(result.Text())
	
	// Fallback if generation fails or returns empty
	if generatedCaption == "" {
		caption = "When the meme speaks for itself 😂"
	} else {
		caption = generatedCaption
	}
}

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	meme := models.Meme{
		ID:       uuid.New(),
		Title:    title,
		ImageURL: imageURL,
		Tags:     tagsJSON,
		Caption:  caption,
		OwnerID:  ownerID,
	}

	return conn.DB.Create(&meme).Error
}
