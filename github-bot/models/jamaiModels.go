package models


type TableType string

// Define constants for the enum values
const (
	ActionTable    TableType = "action"
	ChatTable      TableType = "chat"
	KnowledgeTable TableType = "knowledge"
)

type RagParams struct {
	K              int    `json:"k"`
	TableID        string `json:"table_id,omitempty"`
	RerankingModel string `json:"reranking_model,omitempty"`
}

type GenConfig struct {
	EmbeddingModel string  `json:"embedding_model,omitempty"`
	Model          string     `json:"model"`
	Messages       []Message  `json:"messages"`
	Temperature    float64    `json:"temperature"`
	MaxTokens      int        `json:"max_tokens"`
	TopP          float64    `json:"top_p"`
	RagParams   *RagParams  `json:"rag_params,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Col struct {
	ID        string    `json:"id"`
	Dtype     string    `json:"dtype"`
	Vlen      int       `json:"vlen,omitempty"`
	GenConfig *GenConfig `json:"gen_config,omitempty"`
}

type CreateAgentChatTableRequest struct {
	ID   string `json:"id"`
	Cols []Col  `json:"cols"`
}

type ConfigureAgentChatTableRequest struct {
	TableID   string              `json:"table_id"`
	ColumnMap map[string]GenConfig `json:"column_map"`
}

type Agent struct{
	ColumnID string
	Messages []Message 
}

type CreateAgentConversationTableRequest struct {
	TableID        string `json:"table_id"`
	ConversationID string `json:"conversation_id"`
	Deploy         bool   `json:"deploy"`
}

type CreateAgentKnowledgeTableRequest struct {
	ID             string `json:"id"`
	Cols           []Col  `json:"cols"`
	EmbeddingModel string `json:"embedding_model"`
}

type AddRowRequest struct {
	TableID string            `json:"table_id"`
	Data    []map[string]string `json:"data"`
	Stream  bool              `json:"stream"`
}

type JamaiAuth struct {
	Authorization string
	XProjectID    string
}

type CreateIssueResponse struct {
	Labels   []string `json:"labels"`
	Priority string   `json:"priority"`
	Response string   `json:"response"`
}

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
		Name    *string `json:"name,omitempty"`
	} `json:"message"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamResponse struct{
	ID               string   `json:"id"`
	Object           string   `json:"object"`
	Created          int64    `json:"created"`
	Model            string   `json:"model"`
	Usage            Usage    `json:"usage"`
	Choices          []Choice `json:"choices"`
	References       *string  `json:"references,omitempty"`
	OutputColumnName string   `json:"output_column_name"`
	RowID            string   `json:"row_id"`
}