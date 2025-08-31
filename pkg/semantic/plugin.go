// Package semantic provides a plugin architecture for language-specific semantic analysis
package semantic

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// SemanticChange represents a semantic change detected in code
type SemanticChange struct {
	Type           string            `json:"type"`           // feat, fix, refactor, etc.
	Scope          string            `json:"scope"`          // api, auth, validation, etc.
	Description    string            `json:"description"`    // Human-readable change summary
	Intent         string            `json:"intent"`         // Why this change was made
	Impact         string            `json:"impact"`         // What this affects
	BreakingChange bool              `json:"breaking"`       // Is this a breaking change?
	Files          []string          `json:"files"`          // Affected files
	Confidence     float64           `json:"confidence"`     // 0-1 confidence score
	Reasoning      string            `json:"reasoning"`      // Explanation of analysis
	Metadata       map[string]string `json:"metadata"`       // Plugin-specific metadata
}

// FileChange represents a change to a single file
type FileChange struct {
	Path         string
	Language     string
	BeforeContent string
	AfterContent  string
	DiffContent   string
	ChangeType    string // "added", "modified", "deleted"
}

// AnalysisContext provides context for semantic analysis
type AnalysisContext struct {
	Repository  string
	Branch      string
	Files       []FileChange
	ProjectType string            // detected project type
	Config      map[string]string // plugin-specific config
}

// SemanticPlugin defines the interface for language-specific semantic analyzers
type SemanticPlugin interface {
	// Metadata
	Name() string
	Version() string
	SupportedExtensions() []string
	SupportedFilePatterns() []string
	
	// Analysis capabilities
	CanAnalyze(file FileChange) bool
	AnalyzeFile(ctx context.Context, file FileChange, context AnalysisContext) (*SemanticChange, error)
	AnalyzeProject(ctx context.Context, context AnalysisContext) (*SemanticChange, error)
	
	// Configuration
	DefaultConfig() map[string]string
	ValidateConfig(config map[string]string) error
}

// PluginRegistry manages available semantic analysis plugins
type PluginRegistry struct {
	plugins map[string]SemanticPlugin
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]SemanticPlugin),
	}
}

// Register registers a semantic analysis plugin
func (r *PluginRegistry) Register(plugin SemanticPlugin) error {
	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}
	
	r.plugins[name] = plugin
	return nil
}

// GetPlugin returns a plugin by name
func (r *PluginRegistry) GetPlugin(name string) (SemanticPlugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// GetPluginForFile returns the most appropriate plugin for a file
func (r *PluginRegistry) GetPluginForFile(file FileChange) SemanticPlugin {
	// Try extension matching first
	ext := strings.ToLower(filepath.Ext(file.Path))
	for _, plugin := range r.plugins {
		for _, supportedExt := range plugin.SupportedExtensions() {
			if ext == supportedExt {
				return plugin
			}
		}
	}
	
	// Try pattern matching
	for _, plugin := range r.plugins {
		for _, pattern := range plugin.SupportedFilePatterns() {
			if matched, _ := filepath.Match(pattern, file.Path); matched {
				return plugin
			}
		}
	}
	
	// Try plugin-specific analysis
	for _, plugin := range r.plugins {
		if plugin.CanAnalyze(file) {
			return plugin
		}
	}
	
	return nil
}

// ListPlugins returns all registered plugins
func (r *PluginRegistry) ListPlugins() []SemanticPlugin {
	plugins := make([]SemanticPlugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// SemanticAnalyzer orchestrates semantic analysis using plugins
type SemanticAnalyzer struct {
	registry *PluginRegistry
	config   map[string]map[string]string // plugin-name -> config
}

// NewSemanticAnalyzer creates a new semantic analyzer
func NewSemanticAnalyzer(registry *PluginRegistry) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		registry: registry,
		config:   make(map[string]map[string]string),
	}
}

// SetPluginConfig sets configuration for a specific plugin
func (s *SemanticAnalyzer) SetPluginConfig(pluginName string, config map[string]string) error {
	plugin, exists := s.registry.GetPlugin(pluginName)
	if !exists {
		return fmt.Errorf("plugin %s not found", pluginName)
	}
	
	if err := plugin.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid config for plugin %s: %w", pluginName, err)
	}
	
	s.config[pluginName] = config
	return nil
}

// AnalyzeChanges analyzes a set of file changes using appropriate plugins
func (s *SemanticAnalyzer) AnalyzeChanges(ctx context.Context, files []FileChange) ([]*SemanticChange, error) {
	context := AnalysisContext{
		Files:       files,
		ProjectType: s.detectProjectType(files),
	}
	
	var changes []*SemanticChange
	
	// Analyze individual files
	for _, file := range files {
		plugin := s.registry.GetPluginForFile(file)
		if plugin == nil {
			continue // Skip files without appropriate plugins
		}
		
		// Get plugin config
		pluginConfig := s.config[plugin.Name()]
		if pluginConfig == nil {
			pluginConfig = plugin.DefaultConfig()
		}
		context.Config = pluginConfig
		
		change, err := plugin.AnalyzeFile(ctx, file, context)
		if err != nil {
			continue // Log error but continue with other files
		}
		
		if change != nil {
			changes = append(changes, change)
		}
	}
	
	// Try project-level analysis
	projectChanges := s.analyzeProjectLevel(ctx, context)
	changes = append(changes, projectChanges...)
	
	return s.consolidateChanges(changes), nil
}

// detectProjectType attempts to detect the project type from files
func (s *SemanticAnalyzer) detectProjectType(files []FileChange) string {
	for _, file := range files {
		switch {
		case strings.Contains(file.Path, "terraform") || strings.HasSuffix(file.Path, ".tf"):
			return "terraform"
		case strings.Contains(file.Path, "kubernetes") || strings.HasSuffix(file.Path, ".yaml") && strings.Contains(file.AfterContent, "apiVersion"):
			return "kubernetes"
		case strings.HasSuffix(file.Path, "go.mod"):
			return "go"
		case strings.HasSuffix(file.Path, "package.json"):
			return "nodejs"
		case strings.HasSuffix(file.Path, "requirements.txt") || strings.HasSuffix(file.Path, "pyproject.toml"):
			return "python"
		}
	}
	return "generic"
}

// analyzeProjectLevel performs project-level analysis using plugins
func (s *SemanticAnalyzer) analyzeProjectLevel(ctx context.Context, context AnalysisContext) []*SemanticChange {
	var changes []*SemanticChange
	
	for _, plugin := range s.registry.ListPlugins() {
		pluginConfig := s.config[plugin.Name()]
		if pluginConfig == nil {
			pluginConfig = plugin.DefaultConfig()
		}
		context.Config = pluginConfig
		
		change, err := plugin.AnalyzeProject(ctx, context)
		if err != nil || change == nil {
			continue
		}
		
		changes = append(changes, change)
	}
	
	return changes
}

// consolidateChanges merges and prioritizes semantic changes
func (s *SemanticAnalyzer) consolidateChanges(changes []*SemanticChange) []*SemanticChange {
	if len(changes) == 0 {
		return changes
	}
	
	// Group by type and scope
	groups := make(map[string][]*SemanticChange)
	for _, change := range changes {
		key := fmt.Sprintf("%s:%s", change.Type, change.Scope)
		groups[key] = append(groups[key], change)
	}
	
	// Consolidate groups
	var consolidated []*SemanticChange
	for _, group := range groups {
		if len(group) == 1 {
			consolidated = append(consolidated, group[0])
		} else {
			merged := s.mergeChanges(group)
			consolidated = append(consolidated, merged)
		}
	}
	
	return consolidated
}

// mergeChanges merges multiple similar changes into one
func (s *SemanticAnalyzer) mergeChanges(changes []*SemanticChange) *SemanticChange {
	if len(changes) == 0 {
		return nil
	}
	
	primary := changes[0]
	
	// Merge files
	allFiles := make(map[string]bool)
	for _, change := range changes {
		for _, file := range change.Files {
			allFiles[file] = true
		}
	}
	
	files := make([]string, 0, len(allFiles))
	for file := range allFiles {
		files = append(files, file)
	}
	
	// Calculate average confidence
	totalConfidence := 0.0
	for _, change := range changes {
		totalConfidence += change.Confidence
	}
	avgConfidence := totalConfidence / float64(len(changes))
	
	// Merge breaking change (any breaking = breaking)
	breaking := false
	for _, change := range changes {
		if change.BreakingChange {
			breaking = true
			break
		}
	}
	
	return &SemanticChange{
		Type:           primary.Type,
		Scope:          primary.Scope,
		Description:    fmt.Sprintf("%s (%d files)", primary.Description, len(files)),
		Intent:         primary.Intent,
		Impact:         primary.Impact,
		BreakingChange: breaking,
		Files:          files,
		Confidence:     avgConfidence,
		Reasoning:      "Consolidated from multiple similar changes",
		Metadata:       primary.Metadata,
	}
}