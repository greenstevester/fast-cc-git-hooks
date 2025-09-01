// Package plugins - Sophisticated whole-changeset analysis for Terraform
package plugins

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/semantic"
)

// TerraformChangesetAnalyzer provides whole-changeset analysis for Terraform
type TerraformChangesetAnalyzer struct {
	files         []semantic.FileChange
	addedFiles    []string
	modifiedFiles []string
	deletedFiles  []string
	allTerraform  bool
}

// AnalyzeChangeset performs sophisticated whole-changeset analysis for Terraform files
func (t *TerraformPlugin) AnalyzeChangeset(files []semantic.FileChange) (*semantic.SemanticChange, error) {
	analyzer := &TerraformChangesetAnalyzer{
		files: files,
	}

	// Categorize files
	analyzer.categorizeFiles()

	// Check if ALL files are Terraform-related
	analyzer.checkIfAllTerraform()

	// Perform whole-changeset analysis
	if change := analyzer.detectEnvironmentChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectModuleOnlyChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectVariableOnlyChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectProviderUpgrade(); change != nil {
		return change, nil
	}

	if change := analyzer.detectRefactoring(); change != nil {
		return change, nil
	}

	if change := analyzer.detectSecurityHardening(); change != nil {
		return change, nil
	}

	if change := analyzer.detectDataSourceOnlyChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectOutputOnlyChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectBackendConfigChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectStateManagementChanges(); change != nil {
		return change, nil
	}

	if change := analyzer.detectHotspotStabilization(); change != nil {
		return change, nil
	}

	// Analyze based on file patterns
	return analyzer.analyzeByFilePatterns(), nil
}

// categorizeFiles sorts files into added, modified, and deleted
func (a *TerraformChangesetAnalyzer) categorizeFiles() {
	for _, file := range a.files {
		switch file.ChangeType {
		case "added":
			a.addedFiles = append(a.addedFiles, file.Path)
		case "modified":
			a.modifiedFiles = append(a.modifiedFiles, file.Path)
		case "deleted":
			a.deletedFiles = append(a.deletedFiles, file.Path)
		}
	}
}

// checkIfAllTerraform checks if ALL files in changeset are Terraform-related
func (a *TerraformChangesetAnalyzer) checkIfAllTerraform() {
	a.allTerraform = true

	for _, file := range a.files {
		if !a.isTerraformFile(file.Path) && !a.isTerraformRelatedFile(file.Path) {
			a.allTerraform = false
			break
		}
	}
}

// isTerraformFile checks if file is a Terraform file
func (a *TerraformChangesetAnalyzer) isTerraformFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".tf" || ext == ".tfvars" || strings.HasSuffix(path, ".tfvars.json")
}

// isTerraformRelatedFile checks if file is Terraform-related (docs, configs, etc)
func (a *TerraformChangesetAnalyzer) isTerraformRelatedFile(path string) bool {
	name := filepath.Base(path)

	// Terraform-specific files
	terraformFiles := []string{
		".terraform.lock.hcl",
		"terraform.tfstate",
		"terraform.tfstate.backup",
		".terraformrc",
		"terragrunt.hcl",
	}

	for _, tf := range terraformFiles {
		if name == tf {
			return true
		}
	}

	// Check if in terraform directory
	if strings.Contains(path, "/terraform/") || strings.Contains(path, "/infra/") ||
		strings.Contains(path, "/infrastructure/") {
		return true
	}

	// Check for terraform documentation
	if strings.HasSuffix(path, ".md") && strings.Contains(strings.ToLower(path), "terraform") {
		return true
	}

	return false
}

// detectEnvironmentChanges detects environment-specific changes (dev/staging/prod)
func (a *TerraformChangesetAnalyzer) detectEnvironmentChanges() *semantic.SemanticChange {
	envPatterns := map[string]string{
		"dev":     "development",
		"staging": "staging",
		"stage":   "staging",
		"prod":    "production",
		"prd":     "production",
	}

	envChanges := make(map[string][]string)

	for _, file := range a.files {
		for pattern, env := range envPatterns {
			if strings.Contains(file.Path, "/"+pattern+"/") ||
				strings.Contains(file.Path, "-"+pattern+".") ||
				strings.Contains(file.Path, "_"+pattern+".") ||
				strings.Contains(file.Path, "."+pattern+".") {
				envChanges[env] = append(envChanges[env], file.Path)
			}
		}
	}

	// If all changes are in one environment
	if len(envChanges) == 1 && a.allTerraform {
		for env, files := range envChanges {
			scope := fmt.Sprintf("infra-%s", env)

			// Determine change type based on operations
			changeType := a.determineChangeTypeFromOperations()

			return &semantic.SemanticChange{
				Type:        changeType,
				Scope:       scope,
				Description: fmt.Sprintf("update %s environment infrastructure", env),
				Intent:      fmt.Sprintf("%s environment infrastructure changes", strings.Title(env)),
				Impact:      fmt.Sprintf("Changes isolated to %s environment", env),
				Files:       files,
				Confidence:  0.95,
				Reasoning:   fmt.Sprintf("All %d changes are in %s environment", len(files), env),
				Metadata: map[string]string{
					"environment":   env,
					"all_terraform": "true",
				},
			}
		}
	}

	return nil
}

// detectModuleOnlyChanges detects if only module definitions changed
func (a *TerraformChangesetAnalyzer) detectModuleOnlyChanges() *semantic.SemanticChange {
	allModules := true
	moduleFiles := []string{}

	for _, file := range a.files {
		if a.isTerraformFile(file.Path) {
			baseName := filepath.Base(file.Path)
			if strings.HasPrefix(baseName, "module") ||
				strings.Contains(file.Path, "/modules/") {
				moduleFiles = append(moduleFiles, file.Path)
			} else if !strings.Contains(file.AfterContent+file.BeforeContent, "module ") {
				// File doesn't contain module definitions
				allModules = false
				break
			}
		}
	}

	if allModules && len(moduleFiles) > 0 && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "refactor",
			Scope:       "infra-modules",
			Description: "refactor Terraform module structure",
			Intent:      "Module organization and structure improvements",
			Impact:      "Infrastructure module organization changed",
			Files:       moduleFiles,
			Confidence:  0.9,
			Reasoning:   fmt.Sprintf("All changes (%d files) are module-related", len(moduleFiles)),
			Metadata: map[string]string{
				"change_type": "module_refactoring",
			},
		}
	}

	return nil
}

// detectVariableOnlyChanges detects if only variables/tfvars changed
func (a *TerraformChangesetAnalyzer) detectVariableOnlyChanges() *semantic.SemanticChange {
	allVariables := true
	varFiles := []string{}

	for _, file := range a.files {
		if strings.HasSuffix(file.Path, ".tfvars") ||
			strings.HasSuffix(file.Path, ".tfvars.json") ||
			filepath.Base(file.Path) == "variables.tf" ||
			filepath.Base(file.Path) == "vars.tf" {
			varFiles = append(varFiles, file.Path)
		} else if a.isTerraformFile(file.Path) {
			// Check if file only contains variable definitions
			content := file.AfterContent
			if content == "" {
				content = file.BeforeContent
			}

			if !a.isOnlyVariableDefinitions(content) {
				allVariables = false
				break
			}
		}
	}

	if allVariables && len(varFiles) > 0 && a.allTerraform {
		// Check if it's adding new config or just updating values
		changeType := "chore"
		description := "update infrastructure configuration values"

		if len(a.addedFiles) > 0 {
			changeType = "feat"
			description = "add new infrastructure configuration variables"
		}

		return &semantic.SemanticChange{
			Type:        changeType,
			Scope:       "infra-config",
			Description: description,
			Intent:      "Configuration management",
			Impact:      "Infrastructure parameters updated",
			Files:       varFiles,
			Confidence:  0.85,
			Reasoning:   fmt.Sprintf("All changes (%d files) are variable/configuration related", len(varFiles)),
			Metadata: map[string]string{
				"change_type": "variable_changes",
			},
		}
	}

	return nil
}

// detectProviderUpgrade detects provider version upgrades
func (a *TerraformChangesetAnalyzer) detectProviderUpgrade() *semantic.SemanticChange {
	providerFiles := []string{}
	hasVersionChange := false

	for _, file := range a.files {
		if filepath.Base(file.Path) == ".terraform.lock.hcl" ||
			filepath.Base(file.Path) == "versions.tf" ||
			filepath.Base(file.Path) == "provider.tf" ||
			filepath.Base(file.Path) == "providers.tf" {
			providerFiles = append(providerFiles, file.Path)

			// Check for version changes
			if strings.Contains(file.DiffContent, "version") ||
				strings.Contains(file.DiffContent, "constraints") {
				hasVersionChange = true
			}
		}
	}

	if hasVersionChange && len(providerFiles) == len(a.files) && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "build",
			Scope:       "deps",
			Description: "update Terraform provider versions",
			Intent:      "Provider dependency management",
			Impact:      "Infrastructure provider versions updated",
			Files:       providerFiles,
			Confidence:  0.95,
			Reasoning:   "Only provider version files changed",
			Metadata: map[string]string{
				"change_type": "provider_upgrade",
			},
		}
	}

	return nil
}

// detectRefactoring detects pure refactoring (renames, moves)
func (a *TerraformChangesetAnalyzer) detectRefactoring() *semantic.SemanticChange {
	// Check if files are mostly renames/moves
	renamedCount := 0

	for range a.files {
		// Simple heuristic: if deleted and added files have similar names
		for _, deleted := range a.deletedFiles {
			for _, added := range a.addedFiles {
				if a.isSimilarFileName(deleted, added) {
					renamedCount++
				}
			}
		}
	}

	// If most changes are renames/moves
	if renamedCount > len(a.files)/2 && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "refactor",
			Scope:       "infra",
			Description: "reorganize Terraform configuration structure",
			Intent:      "Infrastructure code organization",
			Impact:      "No functional changes, structure improved",
			Files:       a.getAllFiles(),
			Confidence:  0.8,
			Reasoning:   fmt.Sprintf("Detected %d file renames/moves out of %d changes", renamedCount, len(a.files)),
			Metadata: map[string]string{
				"change_type": "restructuring",
			},
		}
	}

	return nil
}

// detectSecurityHardening detects security-focused changes
func (a *TerraformChangesetAnalyzer) detectSecurityHardening() *semantic.SemanticChange {
	securityFiles := []string{}
	securityKeywords := []string{
		"security", "encryption", "tls", "ssl", "iam", "policy",
		"firewall", "nacl", "security_group", "identity", "auth",
		"certificate", "key", "secret", "vault", "kms",
	}

	allSecurity := true

	for _, file := range a.files {
		isSecurity := false

		// Check filename
		lowerPath := strings.ToLower(file.Path)
		for _, keyword := range securityKeywords {
			if strings.Contains(lowerPath, keyword) {
				isSecurity = true
				securityFiles = append(securityFiles, file.Path)
				break
			}
		}

		// Check content if not identified by filename
		if !isSecurity && a.isTerraformFile(file.Path) {
			content := strings.ToLower(file.AfterContent + file.BeforeContent)
			securityResourceCount := 0

			for _, keyword := range securityKeywords {
				if strings.Contains(content, keyword) {
					securityResourceCount++
				}
			}

			// If significant security content
			if securityResourceCount >= 3 {
				isSecurity = true
				securityFiles = append(securityFiles, file.Path)
			} else {
				allSecurity = false
			}
		}
	}

	if allSecurity && len(securityFiles) > 0 && a.allTerraform {
		changeType := "fix"
		description := "harden infrastructure security configuration"

		if len(a.addedFiles) > len(a.modifiedFiles) {
			changeType = "feat"
			description = "add infrastructure security controls"
		}

		return &semantic.SemanticChange{
			Type:        changeType,
			Scope:       "security",
			Description: description,
			Intent:      "Security hardening and compliance",
			Impact:      "Infrastructure security posture improved",
			Files:       securityFiles,
			Confidence:  0.9,
			Reasoning:   fmt.Sprintf("All %d changes are security-related", len(securityFiles)),
			Metadata: map[string]string{
				"change_type": "security_hardening",
			},
		}
	}

	return nil
}

// detectDataSourceOnlyChanges detects if only data sources changed
func (a *TerraformChangesetAnalyzer) detectDataSourceOnlyChanges() *semantic.SemanticChange {
	allDataSources := true
	dataFiles := []string{}

	for _, file := range a.files {
		if a.isTerraformFile(file.Path) {
			content := file.AfterContent
			if content == "" {
				content = file.BeforeContent
			}

			// Check if file only contains data sources
			if strings.Contains(content, "data \"") && !strings.Contains(content, "resource \"") {
				dataFiles = append(dataFiles, file.Path)
			} else if strings.Contains(content, "resource \"") {
				allDataSources = false
				break
			}
		}
	}

	if allDataSources && len(dataFiles) > 0 && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "chore",
			Scope:       "infra",
			Description: "update infrastructure data source configurations",
			Intent:      "Data source management",
			Impact:      "No infrastructure changes, only data lookups modified",
			Files:       dataFiles,
			Confidence:  0.85,
			Reasoning:   fmt.Sprintf("All %d changes are data source related", len(dataFiles)),
			Metadata: map[string]string{
				"change_type": "data_sources",
			},
		}
	}

	return nil
}

// detectOutputOnlyChanges detects if only outputs changed
func (a *TerraformChangesetAnalyzer) detectOutputOnlyChanges() *semantic.SemanticChange {
	allOutputs := true
	outputFiles := []string{}

	for _, file := range a.files {
		baseName := filepath.Base(file.Path)
		if baseName == "outputs.tf" || baseName == "output.tf" {
			outputFiles = append(outputFiles, file.Path)
		} else if a.isTerraformFile(file.Path) {
			content := file.AfterContent
			if content == "" {
				content = file.BeforeContent
			}

			// Check if file only contains outputs
			if !a.isOnlyOutputDefinitions(content) {
				allOutputs = false
				break
			}
		}
	}

	if allOutputs && len(outputFiles) > 0 && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "chore",
			Scope:       "infra",
			Description: "update infrastructure output values",
			Intent:      "Output configuration",
			Impact:      "Infrastructure outputs modified for downstream consumption",
			Files:       outputFiles,
			Confidence:  0.85,
			Reasoning:   fmt.Sprintf("All %d changes are output related", len(outputFiles)),
			Metadata: map[string]string{
				"change_type": "outputs",
			},
		}
	}

	return nil
}

// detectBackendConfigChanges detects backend configuration changes
func (a *TerraformChangesetAnalyzer) detectBackendConfigChanges() *semantic.SemanticChange {
	backendFiles := []string{}

	for _, file := range a.files {
		if filepath.Base(file.Path) == "backend.tf" ||
			filepath.Base(file.Path) == "backend-config.tf" ||
			strings.Contains(file.DiffContent, "backend \"") {
			backendFiles = append(backendFiles, file.Path)
		}
	}

	if len(backendFiles) == len(a.files) && a.allTerraform {
		return &semantic.SemanticChange{
			Type:        "chore",
			Scope:       "infra",
			Description: "update Terraform backend configuration",
			Intent:      "State management configuration",
			Impact:      "Terraform state storage configuration modified",
			Files:       backendFiles,
			Confidence:  0.9,
			Reasoning:   "Only backend configuration files changed",
			Metadata: map[string]string{
				"change_type": "backend_config",
			},
		}
	}

	return nil
}

// detectStateManagementChanges detects state management related changes
func (a *TerraformChangesetAnalyzer) detectStateManagementChanges() *semantic.SemanticChange {
	stateFiles := []string{}

	for _, file := range a.files {
		baseName := filepath.Base(file.Path)
		if strings.Contains(baseName, "terraform.tfstate") ||
			strings.Contains(file.DiffContent, "terraform import") ||
			strings.Contains(file.DiffContent, "moved {") {
			stateFiles = append(stateFiles, file.Path)
		}
	}

	if len(stateFiles) > 0 && len(stateFiles) == len(a.files) {
		return &semantic.SemanticChange{
			Type:        "chore",
			Scope:       "infra",
			Description: "manage Terraform state",
			Intent:      "State management and resource import",
			Impact:      "Terraform state adjusted without infrastructure changes",
			Files:       stateFiles,
			Confidence:  0.85,
			Reasoning:   "State management operations detected",
			Metadata: map[string]string{
				"change_type": "state_management",
			},
		}
	}

	return nil
}

// detectHotspotStabilization detects changes focused on stabilizing frequently modified files
func (a *TerraformChangesetAnalyzer) detectHotspotStabilization() *semantic.SemanticChange {
	// Only apply to modified files (not new/deleted files)
	if len(a.modifiedFiles) == 0 || len(a.addedFiles) > 0 || len(a.deletedFiles) > 0 {
		return nil
	}

	// Use the terraform plugin's hotspot detection
	plugin := &TerraformPlugin{}
	hotspots := plugin.detectHotspotFiles(a.files)

	// Check if majority of files are hotspots
	hotspotCount := len(hotspots)
	totalModified := len(a.modifiedFiles)

	if hotspotCount > 0 && float64(hotspotCount)/float64(totalModified) >= 0.5 {
		// Build list of hotspot files with their counts
		hotspotDetails := []string{}
		for filePath, count := range hotspots {
			hotspotDetails = append(hotspotDetails, fmt.Sprintf("%s (%d times)", filepath.Base(filePath), count))
		}

		return &semantic.SemanticChange{
			Type:        "fix",
			Scope:       "infra",
			Description: fmt.Sprintf("stabilize Terraform configuration hotspots (%d files)", hotspotCount),
			Intent:      "Configuration stabilization and cleanup",
			Impact:      fmt.Sprintf("Stabilizes %d frequently modified files", hotspotCount),
			Files:       a.modifiedFiles,
			Confidence:  0.9,
			Reasoning:   fmt.Sprintf("Hotspot files detected: %s", strings.Join(hotspotDetails, ", ")),
			Metadata: map[string]string{
				"change_type":   "hotspot_stabilization",
				"hotspot_count": fmt.Sprintf("%d", hotspotCount),
				"hotspot_files": strings.Join(hotspotDetails, ", "),
			},
		}
	}

	return nil
}

// analyzeByFilePatterns performs analysis based on file patterns
func (a *TerraformChangesetAnalyzer) analyzeByFilePatterns() *semantic.SemanticChange {
	if !a.allTerraform {
		return nil
	}

	// Determine primary change type
	changeType := a.determineChangeTypeFromOperations()

	// Build scope from common directory patterns
	scope := a.determineScopeFromPaths()

	// Generate description based on operations
	description := a.generateDescription()

	return &semantic.SemanticChange{
		Type:        changeType,
		Scope:       scope,
		Description: description,
		Intent:      "Infrastructure configuration changes",
		Impact:      a.assessOverallImpact(),
		Files:       a.getAllFiles(),
		Confidence:  0.7,
		Reasoning:   fmt.Sprintf("Terraform changeset with %d files", len(a.files)),
		Metadata: map[string]string{
			"added_files":    fmt.Sprintf("%d", len(a.addedFiles)),
			"modified_files": fmt.Sprintf("%d", len(a.modifiedFiles)),
			"deleted_files":  fmt.Sprintf("%d", len(a.deletedFiles)),
		},
	}
}

// Helper methods

func (a *TerraformChangesetAnalyzer) determineChangeTypeFromOperations() string {
	if len(a.deletedFiles) > len(a.addedFiles)+len(a.modifiedFiles) {
		return "chore" // Cleanup
	}
	if len(a.addedFiles) > len(a.modifiedFiles) {
		return "feat" // New features
	}
	if len(a.modifiedFiles) > len(a.addedFiles) {
		// Check if mostly deletions in modifications
		deletionHeavy := a.checkIfDeletionHeavy()
		if deletionHeavy {
			return "refactor"
		}
		return "fix" // Assuming fixes if modifying existing
	}
	return "refactor" // Default for mixed changes
}

func (a *TerraformChangesetAnalyzer) checkIfDeletionHeavy() bool {
	deletionCount := 0
	additionCount := 0

	for _, file := range a.files {
		if file.ChangeType == "modified" {
			// Simple heuristic: count + and - lines
			deletionCount += strings.Count(file.DiffContent, "\n-")
			additionCount += strings.Count(file.DiffContent, "\n+")
		}
	}

	return deletionCount > additionCount*2
}

func (a *TerraformChangesetAnalyzer) determineScopeFromPaths() string {
	// Look for common directory patterns
	commonScopes := map[string][]string{
		"network":    {"network", "networking", "vpc", "subnet", "firewall"},
		"compute":    {"compute", "instances", "vm", "container", "kubernetes", "k8s"},
		"storage":    {"storage", "database", "rds", "s3", "bucket"},
		"security":   {"security", "iam", "identity", "policy", "rbac"},
		"monitoring": {"monitoring", "logging", "metrics", "observability"},
		"ci":         {"pipeline", "github", "gitlab", "jenkins"},
	}

	scopeCounts := make(map[string]int)

	for _, file := range a.files {
		lowerPath := strings.ToLower(file.Path)
		for scope, patterns := range commonScopes {
			for _, pattern := range patterns {
				if strings.Contains(lowerPath, pattern) {
					scopeCounts[scope]++
					break
				}
			}
		}
	}

	// Find most common scope
	maxCount := 0
	selectedScope := "infra"

	for scope, count := range scopeCounts {
		if count > maxCount {
			maxCount = count
			selectedScope = scope
		}
	}

	return selectedScope
}

func (a *TerraformChangesetAnalyzer) generateDescription() string {
	if len(a.addedFiles) > 0 && len(a.deletedFiles) == 0 && len(a.modifiedFiles) == 0 {
		if len(a.addedFiles) == 1 {
			return fmt.Sprintf("add %s", filepath.Base(a.addedFiles[0]))
		}
		return fmt.Sprintf("add %d Terraform configurations", len(a.addedFiles))
	}

	if len(a.deletedFiles) > 0 && len(a.addedFiles) == 0 && len(a.modifiedFiles) == 0 {
		if len(a.deletedFiles) == 1 {
			return fmt.Sprintf("remove %s", filepath.Base(a.deletedFiles[0]))
		}
		return fmt.Sprintf("remove %d Terraform configurations", len(a.deletedFiles))
	}

	if len(a.modifiedFiles) > 0 && len(a.addedFiles) == 0 && len(a.deletedFiles) == 0 {
		if len(a.modifiedFiles) == 1 {
			return fmt.Sprintf("update %s", filepath.Base(a.modifiedFiles[0]))
		}
		return fmt.Sprintf("update %d Terraform configurations", len(a.modifiedFiles))
	}

	// Mixed changes
	return "update infrastructure configuration"
}

func (a *TerraformChangesetAnalyzer) assessOverallImpact() string {
	impacts := []string{}

	if len(a.addedFiles) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d new configurations", len(a.addedFiles)))
	}
	if len(a.modifiedFiles) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d modified configurations", len(a.modifiedFiles)))
	}
	if len(a.deletedFiles) > 0 {
		impacts = append(impacts, fmt.Sprintf("%d removed configurations", len(a.deletedFiles)))
	}

	if len(impacts) == 0 {
		return "Infrastructure configuration updated"
	}

	return strings.Join(impacts, ", ")
}

func (a *TerraformChangesetAnalyzer) getAllFiles() []string {
	files := []string{}
	for _, file := range a.files {
		files = append(files, file.Path)
	}
	return files
}

func (a *TerraformChangesetAnalyzer) isSimilarFileName(file1, file2 string) bool {
	base1 := filepath.Base(file1)
	base2 := filepath.Base(file2)

	// Simple similarity check
	return strings.Contains(base1, base2) || strings.Contains(base2, base1)
}

func (a *TerraformChangesetAnalyzer) isOnlyVariableDefinitions(content string) bool {
	// Check if content only has variable blocks
	lines := strings.Split(content, "\n")
	nonVariableLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if !strings.Contains(line, "variable") && !strings.Contains(line, "{") && !strings.Contains(line, "}") &&
			!strings.Contains(line, "description") && !strings.Contains(line, "type") && !strings.Contains(line, "default") {
			nonVariableLines++
		}
	}

	return nonVariableLines < 5 // Allow some non-variable lines
}

func (a *TerraformChangesetAnalyzer) isOnlyOutputDefinitions(content string) bool {
	// Check if content only has output blocks
	lines := strings.Split(content, "\n")
	nonOutputLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if !strings.Contains(line, "output") && !strings.Contains(line, "{") && !strings.Contains(line, "}") &&
			!strings.Contains(line, "value") && !strings.Contains(line, "description") {
			nonOutputLines++
		}
	}

	return nonOutputLines < 5 // Allow some non-output lines
}
