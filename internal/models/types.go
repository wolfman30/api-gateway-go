package models

// CreateReelRequest matches the OpenAPI schema from ai-twin-contracts.
type CreateReelRequest struct {
	ProjectID          string              `json:"projectId"`
	ICP                IdealClientProfile  `json:"icp"`
	Idea               string              `json:"idea"`
	FluxModel          FluxModelConfig     `json:"fluxModel"`
	FluxPrompt         FluxPromptRequest   `json:"fluxPrompt"`
	KlingPreferences   *KlingPreferences   `json:"klingPreferences,omitempty"`
	CaptionPreferences *CaptionPreferences `json:"captionPreferences,omitempty"`
}

type IdealClientProfile struct {
	Industry           string   `json:"industry"`
	AudiencePainPoints []string `json:"audiencePainPoints"`
	DesiredOutcome     string   `json:"desiredOutcome,omitempty"`
}

type FluxModelConfig struct {
	LoraURL  string  `json:"loraUrl"`
	CfgScale float64 `json:"cfgScale,omitempty"`
	Steps    int     `json:"steps,omitempty"`
}

type FluxPromptRequest struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negativePrompt,omitempty"`
	AspectRatio    string `json:"aspectRatio,omitempty"`
	BatchSize      int    `json:"batchSize,omitempty"`
}

type KlingPreferences struct {
	StylePreset     string  `json:"stylePreset,omitempty"`
	NegativePrompt  string  `json:"negativePrompt,omitempty"`
	GuidanceScale   float64 `json:"guidanceScale,omitempty"`
	DurationSeconds float64 `json:"durationSeconds,omitempty"`
}

type CaptionPreferences struct {
	HookStyle    string        `json:"hookStyle,omitempty"`
	BodyStyle    string        `json:"bodyStyle,omitempty"`
	CallToAction *CallToAction `json:"callToAction,omitempty"`
}

type CallToAction struct {
	Type    string `json:"type"`
	Keyword string `json:"keyword,omitempty"`
}

// CreateReelResponse is returned after accepting a reel request.
type CreateReelResponse struct {
	RunID string `json:"runId"`
}

// RunStatusResponse describes the current state of a run.
type RunStatusResponse struct {
	RunID  string    `json:"runId"`
	Status string    `json:"status"`
	Steps  []RunStep `json:"steps"`
}

type RunStep struct {
	Name      string   `json:"name"`
	Status    string   `json:"status"`
	UpdatedAt string   `json:"updatedAt"`
	Artifacts []string `json:"artifacts,omitempty"`
}
