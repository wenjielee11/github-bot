package models

// TableType defines a type for different kinds of tables used in the system.
type TableType string

// Define constants for the enum values of TableType
const (
	ActionTable    TableType = "action"
	ChatTable      TableType = "chat"
	KnowledgeTable TableType = "knowledge"
)

// RagParams defines the parameters for the RAG (Retrieval-Augmented Generation) model.
type RagParams struct {
	K              int    `json:"k"`                         // The number of documents to retrieve.
	TableID        string `json:"table_id,omitempty"`        // The ID of the table to retrieve documents from.
	RerankingModel string `json:"reranking_model,omitempty"` // The model used for reranking the retrieved documents.
}

// GenConfig defines the configuration for generating responses.
type GenConfig struct {
	EmbeddingModel string     `json:"embedding_model,omitempty"` // The model used for embeddings.
	Model          string     `json:"model"`                     // The model used for generation.
	Messages       []Message  `json:"messages"`                  // The messages to be included in the generation.
	Temperature    float64    `json:"temperature"`               // The temperature setting for generation.
	MaxTokens      int        `json:"max_tokens"`                // The maximum number of tokens for the generated response.
	TopP           float64    `json:"top_p"`                     // The nucleus sampling parameter.
	RagParams      *RagParams `json:"rag_params,omitempty"`      // The RAG parameters for retrieval-augmented generation.
}

// Message defines the structure of a message used in generation.
type Message struct {
	Role    string `json:"role"`    // The role of the message sender (e.g., "user", "system").
	Content string `json:"content"` // The content of the message.
}

// Col defines a column in a table, including its ID, data type, length, and generation configuration.
type Col struct {
	ID        string     `json:"id"`                   // The ID of the column.
	Dtype     string     `json:"dtype"`                // The data type of the column.
	Vlen      int        `json:"vlen,omitempty"`       // The variable length of the column, if applicable.
	GenConfig *GenConfig `json:"gen_config,omitempty"` // The generation configuration for the column.
}

// CreateAgentChatTableRequest defines the request structure for creating an agent chat table.
type CreateAgentChatTableRequest struct {
	ID   string `json:"id"`  // The ID of the table to be created.
	Cols []Col  `json:"cols"` // The columns of the table.
}

// ConfigureAgentChatTableRequest defines the request structure for configuring an agent chat table.
type ConfigureAgentChatTableRequest struct {
	TableID   string               `json:"table_id"`   // The ID of the table to be configured.
	ColumnMap map[string]GenConfig `json:"column_map"` // The map of column IDs to their generation configurations.
}

// Agent defines the structure of an agent, including its column ID and messages.
type Agent struct {
	ColumnID string    // The ID of the column the agent is associated with.
	Messages []Message // The messages associated with the agent.
}

// CreateAgentConversationTableRequest defines the request structure for creating an agent conversation table.
type CreateAgentConversationTableRequest struct {
	TableID        string `json:"table_id"`        // The ID of the table to be created.
	ConversationID string `json:"conversation_id"` // The ID of the conversation.
	Deploy         bool   `json:"deploy"`          // Whether the table should be deployed immediately.
}

// CreateAgentKnowledgeTableRequest defines the request structure for creating an agent knowledge table.
type CreateAgentKnowledgeTableRequest struct {
	ID             string `json:"id"`             // The ID of the table to be created.
	Cols           []Col  `json:"cols"`           // The columns of the table.
	EmbeddingModel string `json:"embedding_model"` // The embedding model to be used.
}

// AddRowRequest defines the request structure for adding rows to a table.
type AddRowRequest struct {
	TableID string              `json:"table_id"` // The ID of the table to add rows to.
	Data    []map[string]string `json:"data"`     // The data for the rows to be added.
	Stream  bool                `json:"stream"`   // Whether to stream the data.
}

// JamaiAuth defines the authentication structure for JAM.AI.
type JamaiAuth struct {
	Authorization string // The authorization token.
	XProjectID    string // The project ID.
}

// CreateIssueResponse defines the structure of the response when creating an issue.
type CreateIssueResponse struct {
	Labels   []string `json:"labels"` // The labels assigned to the issue.
	Priority string   `json:"priority"` // The priority of the issue.
	Response string   `json:"response"` // The response message for the issue.
}

// Choice defines the structure of a choice in the response.
type Choice struct {
	Message struct {
		Role    string  `json:"role"`    // The role of the message sender.
		Content string  `json:"content"` // The content of the message.
		Name    *string `json:"name,omitempty"` // The name of the sender, if applicable.
	} `json:"message"`
	Index        int         `json:"index"`        // The index of the choice.
	FinishReason interface{} `json:"finish_reason"` // The reason the choice finished.
}

// Usage defines the structure of token usage in the response.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // The number of tokens in the prompt.
	CompletionTokens int `json:"completion_tokens"` // The number of tokens in the completion.
	TotalTokens      int `json:"total_tokens"`      // The total number of tokens used.
}

// StreamResponse defines the structure of a streaming response.
type StreamResponse struct {
	ID               string   `json:"id"`                // The ID of the response.
	Object           string   `json:"object"`            // The object type.
	Created          int64    `json:"created"`           // The creation timestamp.
	Model            string   `json:"model"`             // The model used for generation.
	Usage            Usage    `json:"usage"`             // The usage details of the response.
	Choices          []Choice `json:"choices"`           // The choices in the response.
	References       *string  `json:"references,omitempty"` // The references, if any.
	OutputColumnName string   `json:"output_column_name"` // The name of the output column.
	RowID            string   `json:"row_id"`            // The ID of the row.
}

// CreatePullReqSecretResponse defines the structure of the response when checking for secret key leakage in a pull request.
type CreatePullReqSecretResponse struct {
	Leak   bool   `json:"leak"`   // Whether a leak was detected.
	Commit string `json:"commit"` // The commit hash.
	Response string `json:"response"` // The response message.
}
