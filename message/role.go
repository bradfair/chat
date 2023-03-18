package message

// Role is a type that represents the roles in OpenAI's chat API.
type Role string

const (
	// RoleAssistant represents the chatbot.
	RoleAssistant Role = "assistant"
	// RoleSystem represents the entity managing the conversation. This is typically used for prompts and other content that is not directly from the user, though such content could also be sent as if from the user.
	RoleSystem Role = "system"
	// RoleUser represents the user interacting with the chatbot.
	RoleUser Role = "user"
)
