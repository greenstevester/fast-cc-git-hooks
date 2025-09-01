// Package plugins contains language-specific semantic analysis plugins
package plugins

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/semantic"
)

// TerraformPlugin provides semantic analysis for Terraform files
type TerraformPlugin struct {
	version string
}

// NewTerraformPlugin creates a new Terraform semantic analyzer plugin
func NewTerraformPlugin() *TerraformPlugin {
	return &TerraformPlugin{
		version: "1.0.0",
	}
}

// Name returns the plugin name
func (t *TerraformPlugin) Name() string {
	return "terraform"
}

// Version returns the plugin version
func (t *TerraformPlugin) Version() string {
	return t.version
}

// SupportedExtensions returns file extensions this plugin supports
func (t *TerraformPlugin) SupportedExtensions() []string {
	return []string{".tf", ".tfvars", ".tfvars.json"}
}

// SupportedFilePatterns returns file patterns this plugin supports
func (t *TerraformPlugin) SupportedFilePatterns() []string {
	return []string{
		"*.tf",
		"*.tfvars",
		"*.tfvars.json",
		"terraform/*",
		"infra/*",
		"infrastructure/*",
	}
}

// CanAnalyze determines if this plugin can analyze the given file
func (t *TerraformPlugin) CanAnalyze(file semantic.FileChange) bool {
	// Check extension
	for _, ext := range t.SupportedExtensions() {
		if strings.HasSuffix(strings.ToLower(file.Path), ext) {
			return true
		}
	}

	// Check content for Terraform syntax
	content := file.AfterContent
	if content == "" {
		content = file.BeforeContent
	}

	terraformKeywords := []string{
		"resource \"",
		"data \"",
		"variable \"",
		"output \"",
		"provider \"",
		"terraform {",
		"module \"",
	}

	for _, keyword := range terraformKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false
}

// AnalyzeFile analyzes a single Terraform file for semantic changes
func (t *TerraformPlugin) AnalyzeFile(ctx context.Context, file semantic.FileChange, analysisCtx semantic.AnalysisContext) (*semantic.SemanticChange, error) {
	switch file.ChangeType {
	case "added":
		return t.analyzeNewFile(file, analysisCtx)
	case "deleted":
		return t.analyzeDeletedFile(file, analysisCtx)
	case "modified":
		return t.analyzeModifiedFile(file, analysisCtx)
	default:
		return nil, fmt.Errorf("unknown change type: %s", file.ChangeType)
	}
}

// AnalyzeProject performs project-level analysis for Terraform
func (t *TerraformPlugin) AnalyzeProject(ctx context.Context, context semantic.AnalysisContext) (*semantic.SemanticChange, error) {
	// First, detect if this is actually a Terraform codebase
	if !t.isTerraformCodebase(context) {
		return nil, nil
	}

	// Look for Terraform project-level changes
	var terraformFiles []semantic.FileChange
	for _, file := range context.Files {
		if t.CanAnalyze(file) {
			terraformFiles = append(terraformFiles, file)
		}
	}

	if len(terraformFiles) == 0 {
		return nil, nil
	}

	// Use the sophisticated whole-changeset analysis
	return t.AnalyzeChangeset(terraformFiles)
}

// DefaultConfig returns the default configuration for the plugin
func (t *TerraformPlugin) DefaultConfig() map[string]string {
	return map[string]string{
		"detect_breaking_changes": "true",
		"analyze_security":        "true",
		"check_best_practices":    "true",
		"provider_sensitivity":    "high", // high, medium, low
	}
}

// ValidateConfig validates plugin configuration
func (t *TerraformPlugin) ValidateConfig(config map[string]string) error {
	validKeys := map[string]bool{
		"detect_breaking_changes": true,
		"analyze_security":        true,
		"check_best_practices":    true,
		"provider_sensitivity":    true,
	}

	for key := range config {
		if !validKeys[key] {
			return fmt.Errorf("unknown config key: %s", key)
		}
	}

	if sensitivity, ok := config["provider_sensitivity"]; ok {
		if sensitivity != "high" && sensitivity != "medium" && sensitivity != "low" {
			return fmt.Errorf("invalid provider_sensitivity: %s (must be high, medium, or low)", sensitivity)
		}
	}

	return nil
}

// analyzeNewFile analyzes a newly added Terraform file
func (t *TerraformPlugin) analyzeNewFile(file semantic.FileChange, _ semantic.AnalysisContext) (*semantic.SemanticChange, error) {
	content := file.AfterContent

	// Analyze what type of resources are being added
	resourceTypes := t.extractResourceTypes(content)

	scope := t.determineScope(file.Path, content)

	if len(resourceTypes) == 0 {
		return &semantic.SemanticChange{
			Type:        "feat",
			Scope:       scope,
			Description: fmt.Sprintf("add Terraform configuration file %s", t.getFileName(file.Path)),
			Intent:      "Infrastructure as Code setup",
			Impact:      "New infrastructure components defined",
			Files:       []string{file.Path},
			Confidence:  0.8,
			Reasoning:   "New Terraform file detected",
			Metadata:    map[string]string{"file_type": "terraform"},
		}, nil
	}

	// Analyze specific resource types
	change := t.analyzeResourceTypes(resourceTypes, file, scope)
	return change, nil
}

// analyzeDeletedFile analyzes a deleted Terraform file
func (t *TerraformPlugin) analyzeDeletedFile(file semantic.FileChange, _ semantic.AnalysisContext) (*semantic.SemanticChange, error) {
	content := file.BeforeContent
	resourceTypes := t.extractResourceTypes(content)
	scope := t.determineScope(file.Path, content)

	// Check if this is a breaking change
	breaking := t.isDeletionBreaking(resourceTypes)

	changeType := "refactor"
	if breaking {
		changeType = "feat" // Breaking changes are typically features
	}

	return &semantic.SemanticChange{
		Type:           changeType,
		Scope:          scope,
		Description:    fmt.Sprintf("remove Terraform configuration %s", t.getFileName(file.Path)),
		Intent:         "Infrastructure cleanup or refactoring",
		Impact:         "Infrastructure resources will be destroyed",
		BreakingChange: breaking,
		Files:          []string{file.Path},
		Confidence:     0.9,
		Reasoning:      fmt.Sprintf("Terraform file deleted with %d resource types", len(resourceTypes)),
		Metadata: map[string]string{
			"file_type":      "terraform",
			"resource_count": fmt.Sprintf("%d", len(resourceTypes)),
		},
	}, nil
}

// analyzeModifiedFile analyzes a modified Terraform file
func (t *TerraformPlugin) analyzeModifiedFile(file semantic.FileChange, _ semantic.AnalysisContext) (*semantic.SemanticChange, error) {
	beforeResources := t.extractResourceTypes(file.BeforeContent)
	afterResources := t.extractResourceTypes(file.AfterContent)

	added, removed, modified := t.compareResources(beforeResources, afterResources)

	scope := t.determineScope(file.Path, file.AfterContent)

	// Check if this file is a hotspot (modified repeatedly in recent commits)
	hotspots := t.detectHotspotFiles([]semantic.FileChange{file})
	isHotspot := hotspots[file.Path] > 0

	// Determine change type based on modifications
	changeType := "refactor" // default
	description := "update Terraform configuration"
	intent := "Infrastructure modification"
	breaking := false

	// Adjust for hotspot files
	if isHotspot {
		changeType = "fix"
		description = "stabilize Terraform configuration"
		intent = "Configuration stabilization and cleanup"
	}

	// Override with specific change patterns if not already a hotspot
	if !isHotspot {
		if len(added) > 0 && len(removed) == 0 {
			changeType = "feat"
			if len(added) == 1 {
				description = fmt.Sprintf("add %s resource", added[0])
			} else {
				description = fmt.Sprintf("add %d new resources", len(added))
			}
			intent = "Infrastructure expansion"
		} else if len(removed) > 0 && len(added) == 0 {
			changeType = "refactor"
			if len(removed) == 1 {
				description = fmt.Sprintf("remove %s resource", removed[0])
			} else {
				description = fmt.Sprintf("remove %d resources", len(removed))
			}
			intent = "Infrastructure cleanup"
			breaking = t.isRemovalBreaking(removed)
		} else if len(modified) > 0 {
			// Check if it's a fix or enhancement
			if t.isSecurityImprovement(file.DiffContent) {
				changeType = "fix"
				description = "improve infrastructure security configuration"
				intent = "Security hardening"
			} else if t.isPerformanceImprovement(file.DiffContent) {
				changeType = "perf"
				description = "optimize infrastructure performance"
				intent = "Performance optimization"
			} else if t.isBugFix(file.DiffContent) {
				changeType = "fix"
				description = "fix infrastructure configuration issues"
				intent = "Bug fix"
			} else {
				changeType = "refactor"
				description = "refactor infrastructure configuration"
				intent = "Configuration improvement"
			}

			// Check for breaking changes
			breaking = t.hasBreakingChanges(file.DiffContent)
		}
	} else {
		// For hotspot files, append the hotspot count to description
		hotspotCount := hotspots[file.Path]
		description = fmt.Sprintf("stabilize %s configuration (modified %d times recently)", t.getFileName(file.Path), hotspotCount)
	}

	confidence := t.calculateConfidence(added, removed, modified, file.DiffContent)

	// Prepare metadata
	metadata := map[string]string{
		"added_resources":    strings.Join(added, ","),
		"removed_resources":  strings.Join(removed, ","),
		"modified_resources": strings.Join(modified, ","),
	}

	// Add hotspot information to metadata
	if isHotspot {
		metadata["hotspot"] = "true"
		metadata["hotspot_count"] = fmt.Sprintf("%d", hotspots[file.Path])
		metadata["hotspot_reasoning"] = "File modified repeatedly in recent commits, indicating stabilization effort"
	}

	// Update reasoning to include hotspot information
	reasoning := t.generateReasoning(added, removed, modified)
	if isHotspot {
		reasoning = fmt.Sprintf("%s; Hotspot detected: modified %d times in last 5 commits", reasoning, hotspots[file.Path])
	}

	return &semantic.SemanticChange{
		Type:           changeType,
		Scope:          scope,
		Description:    description,
		Intent:         intent,
		Impact:         t.assessImpact(added, removed, modified),
		BreakingChange: breaking,
		Files:          []string{file.Path},
		Confidence:     confidence,
		Reasoning:      reasoning,
		Metadata:       metadata,
	}, nil
}

// extractResourceTypes extracts resource types from Terraform content
func (t *TerraformPlugin) extractResourceTypes(content string) []string {
	resourcePattern := regexp.MustCompile(`resource\s+"([^"]+)"\s+"([^"]+)"`)
	matches := resourcePattern.FindAllStringSubmatch(content, -1)

	var resources []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) >= 2 {
			resourceType := match[1]
			if !seen[resourceType] {
				resources = append(resources, resourceType)
				seen[resourceType] = true
			}
		}
	}

	return resources
}

// determineScope determines the scope based on file path and content
func (t *TerraformPlugin) determineScope(filePath, content string) string {
	// Path-based scoping
	pathParts := strings.Split(strings.ToLower(filePath), "/")

	for _, part := range pathParts {
		switch part {
		case "network", "networking", "vcn":
			return "network"
		case "security", "iam", "auth":
			return "security"
		case "storage", "objectstorage", "database", "db":
			return "storage"
		case "compute", "instances":
			return "compute"
		case "monitoring", "logs":
			return "monitoring"
		case "dns":
			return "dns"
		}
	}

	// Content-based scoping for OCI resources
	if strings.Contains(content, "oci_core_vcn") || strings.Contains(content, "oci_core_subnet") || strings.Contains(content, "oci_load_balancer") {
		return "network"
	}
	if strings.Contains(content, "oci_identity") || strings.Contains(content, "oci_core_security") {
		return "security"
	}
	if strings.Contains(content, "oci_objectstorage") || strings.Contains(content, "oci_database") || strings.Contains(content, "oci_mysql") {
		return "storage"
	}
	if strings.Contains(content, "oci_core_instance") || strings.Contains(content, "oci_containerengine") {
		return "compute"
	}

	return "infra"
}

// analyzeResourceTypes analyzes specific OCI resource types
func (t *TerraformPlugin) analyzeResourceTypes(resourceTypes []string, file semantic.FileChange, scope string) *semantic.SemanticChange {
	// Categorize resources by impact
	criticalResources := []string{"oci_core_vcn", "oci_database_autonomous_database", "oci_database_db_system", "oci_mysql_mysql_db_system"}
	securityResources := []string{"oci_identity_policy", "oci_core_security_list", "oci_identity_user", "oci_identity_group"}

	critical := false
	security := false

	for _, resource := range resourceTypes {
		for _, cr := range criticalResources {
			if resource == cr {
				critical = true
				break
			}
		}
		for _, sr := range securityResources {
			if resource == sr {
				security = true
				break
			}
		}
	}

	changeType := "feat"
	description := fmt.Sprintf("add %s infrastructure", scope)
	impact := "New infrastructure components"

	if critical {
		description = fmt.Sprintf("add critical %s infrastructure", scope)
		impact = "Critical infrastructure components added"
	} else if security {
		description = fmt.Sprintf("add %s security configuration", scope)
		impact = "Security infrastructure components added"
	}

	return &semantic.SemanticChange{
		Type:        changeType,
		Scope:       scope,
		Description: description,
		Intent:      "Infrastructure provisioning",
		Impact:      impact,
		Files:       []string{file.Path},
		Confidence:  0.85,
		Reasoning:   fmt.Sprintf("Added %d Terraform resources: %s", len(resourceTypes), strings.Join(resourceTypes, ", ")),
		Metadata: map[string]string{
			"resource_types": strings.Join(resourceTypes, ","),
			"critical":       fmt.Sprintf("%t", critical),
			"security":       fmt.Sprintf("%t", security),
		},
	}
}

// Helper methods for analysis
func (t *TerraformPlugin) compareResources(before, after []string) (added, removed, modified []string) {
	beforeMap := make(map[string]bool)
	afterMap := make(map[string]bool)

	for _, r := range before {
		beforeMap[r] = true
	}
	for _, r := range after {
		afterMap[r] = true
	}

	// Find added resources
	for _, r := range after {
		if !beforeMap[r] {
			added = append(added, r)
		}
	}

	// Find removed resources
	for _, r := range before {
		if !afterMap[r] {
			removed = append(removed, r)
		}
	}

	// Modified are resources present in both (we assume they're modified if file changed)
	for _, r := range after {
		if beforeMap[r] {
			modified = append(modified, r)
		}
	}

	return added, removed, modified
}

func (t *TerraformPlugin) isDeletionBreaking(resourceTypes []string) bool {
	breakingResources := []string{
		"oci_core_vcn", "oci_database_autonomous_database", "oci_database_db_system",
		"oci_objectstorage_bucket", "oci_containerengine_cluster", "oci_mysql_mysql_db_system",
	}

	for _, resource := range resourceTypes {
		for _, breaking := range breakingResources {
			if resource == breaking {
				return true
			}
		}
	}
	return false
}

func (t *TerraformPlugin) isRemovalBreaking(removed []string) bool {
	return t.isDeletionBreaking(removed)
}

func (t *TerraformPlugin) isSecurityImprovement(diff string) bool {
	securityPatterns := []string{
		"+.*encryption",
		"+.*security_list",
		"+.*iam_policy",
		"+.*identity_policy",
		"+.*https",
		"+.*443", // HTTPS port
		"+.*network_security_group",
		"-.*public_read",
		"-.*\"0.0.0.0/0\"",  // removing open CIDR access
		"+.*compartment_id", // proper compartment isolation
		"\\+.*min.*443",     // OCI security list HTTPS port
		"\\+.*max.*443",     // OCI security list HTTPS port
	}

	for _, pattern := range securityPatterns {
		matched, err := regexp.MatchString(pattern, diff)
		if err != nil {
			continue // Skip invalid regex patterns
		}
		if matched {
			return true
		}
	}
	return false
}

func (t *TerraformPlugin) isPerformanceImprovement(diff string) bool {
	perfPatterns := []string{
		"+.*shape.*\\.\\d+", // OCI compute shapes with more resources
		"+.*cpu_core_count.*[0-9]{2,}",
		"+.*memory_in_gbs.*[0-9]{3,}",
		"+.*is_auto_scaling_enabled.*true",
		"+.*backup_policy",
		"+.*high_availability",
	}

	for _, pattern := range perfPatterns {
		matched, err := regexp.MatchString(pattern, diff)
		if err != nil {
			continue // Skip invalid regex patterns
		}
		if matched {
			return true
		}
	}
	return false
}

func (t *TerraformPlugin) isBugFix(diff string) bool {
	bugFixPatterns := []string{
		"fix", "Fix", "bug", "Bug", "issue", "Issue",
		"-.*deprecated",
		"+.*latest",
		"correction", "correct",
	}

	for _, pattern := range bugFixPatterns {
		if strings.Contains(diff, pattern) {
			return true
		}
	}
	return false
}

func (t *TerraformPlugin) hasBreakingChanges(diff string) bool {
	breakingPatterns := []string{
		"-.*force_destroy.*false",
		"+.*force_destroy.*true",
		"-.*shape",          // changing compute shape
		"-.*compartment_id", // changing compartment
		"-.*vcn_id",
		"-.*subnet_id",
		"-.*availability_domain",
	}

	for _, pattern := range breakingPatterns {
		matched, err := regexp.MatchString(pattern, diff)
		if err != nil {
			continue // Skip invalid regex patterns
		}
		if matched {
			return true
		}
	}
	return false
}

func (t *TerraformPlugin) calculateConfidence(added, removed, _ []string, diff string) float64 {
	confidence := 0.7 // base confidence

	// Higher confidence for clear resource changes
	if len(added) > 0 || len(removed) > 0 {
		confidence += 0.2
	}

	// Higher confidence for well-structured diffs
	if strings.Contains(diff, "resource \"") {
		confidence += 0.1
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

func (t *TerraformPlugin) assessImpact(added, removed, modified []string) string {
	impacts := []string{}

	if len(added) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d resources will be created", len(added)))
	}
	if len(removed) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d resources will be destroyed", len(removed)))
	}
	if len(modified) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d resources will be modified", len(modified)))
	}

	if len(impacts) == 0 {
		return "Infrastructure configuration updated"
	}

	return strings.Join(impacts, "; ")
}

func (t *TerraformPlugin) generateReasoning(added, removed, modified []string) string {
	parts := []string{}

	if len(added) > 0 {
		parts = append(parts, fmt.Sprintf("Added resources: %s", strings.Join(added, ", ")))
	}
	if len(removed) > 0 {
		parts = append(parts, fmt.Sprintf("Removed resources: %s", strings.Join(removed, ", ")))
	}
	if len(modified) > 0 {
		parts = append(parts, fmt.Sprintf("Modified resources: %s", strings.Join(modified, ", ")))
	}

	if len(parts) == 0 {
		return "Terraform configuration changes detected"
	}

	return strings.Join(parts, "; ")
}

func (t *TerraformPlugin) getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// isTerraformCodebase detects if the current project is a Terraform codebase
func (t *TerraformPlugin) isTerraformCodebase(context semantic.AnalysisContext) bool {
	// Check for Terraform-specific indicators in the project
	terraformIndicators := []string{
		".terraform/",
		"terraform.tfstate",
		"terraform.tfstate.backup",
		"terraform.tfvars",
		"terraform.tfvars.json",
		".terraform.lock.hcl",
		"main.tf",
		"variables.tf",
		"outputs.tf",
		"provider.tf",
		"versions.tf",
	}

	// Check if any files in the changeset match Terraform patterns
	terraformFileCount := 0
	totalFiles := len(context.Files)

	for _, file := range context.Files {
		// Direct Terraform file check
		if t.CanAnalyze(file) {
			terraformFileCount++
			continue
		}

		// Check for Terraform-specific file names/patterns
		fileName := t.getFileName(file.Path)
		filePath := strings.ToLower(file.Path)

		for _, indicator := range terraformIndicators {
			if strings.Contains(filePath, indicator) || fileName == indicator {
				return true // Strong indicator - definitely a Terraform project
			}
		}

		// Check for Terraform directory patterns
		if strings.Contains(filePath, "/terraform/") ||
			strings.Contains(filePath, "/infra/") ||
			strings.Contains(filePath, "/infrastructure/") {
			terraformFileCount++
		}
	}

	// If more than 30% of files are Terraform-related, consider it a Terraform codebase
	if totalFiles > 0 && float64(terraformFileCount)/float64(totalFiles) > 0.3 {
		return true
	}

	// If we have any Terraform files at all, run the analysis
	return terraformFileCount > 0
}

// detectHotspotFiles checks if files have been modified repeatedly in recent commits
func (t *TerraformPlugin) detectHotspotFiles(files []semantic.FileChange) map[string]int {
	hotspots := make(map[string]int)

	for _, file := range files {
		// Sanitize and validate file path to prevent command injection
		cleanPath := filepath.Clean(file.Path)
		if strings.Contains(cleanPath, "..") || strings.Contains(cleanPath, ";") || 
		   strings.Contains(cleanPath, "|") || strings.Contains(cleanPath, "&") ||
		   strings.HasPrefix(cleanPath, "-") || len(cleanPath) == 0 {
			continue // Skip potentially malicious or invalid paths
		}
		
		// Use a safe, sanitized path for the git command
		// #nosec G204 - path is sanitized above
		cmd := exec.Command("git", "log", "-n", "5", "--name-only", "--pretty=", "--", cleanPath)
		output, err := cmd.Output()
		if err != nil {
			continue // Skip if git command fails
		}

		// Count occurrences of this file in recent commits
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		count := 0
		for _, line := range lines {
			if strings.TrimSpace(line) == cleanPath {
				count++
			}
		}

		if count > 1 { // File appears in multiple recent commits
			hotspots[file.Path] = count
		}
	}

	return hotspots
}


