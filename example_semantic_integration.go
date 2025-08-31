// Example integration of semantic analysis with the cc command
package main

import (
	"fmt"
	"log"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/semantic"
	"github.com/greenstevester/fast-cc-git-hooks/pkg/semantic/plugins"
)

// Enhanced CC command with semantic analysis
func enhancedCCCommand(diff string, useSemanticAnalysis bool) {
	// Use ASCII heart for better terminal compatibility
	fmt.Println(">> Made with <3 for Boo")
	
	var commitMessage string
	var confidence float64
	
	if useSemanticAnalysis {
		// Initialize semantic analyzer
		analyzer := semantic.NewCCSemanticAnalyzer()
		
		// Register plugins
		terraformPlugin := plugins.NewTerraformPlugin()
		if err := analyzer.RegisterPlugins(terraformPlugin); err != nil {
			log.Printf("Failed to register Terraform plugin: %v", err)
			useSemanticAnalysis = false
		}
	}
	
	if useSemanticAnalysis {
		analyzer := semantic.NewCCSemanticAnalyzer()
		terraformPlugin := plugins.NewTerraformPlugin()
		analyzer.RegisterPlugins(terraformPlugin)
		
		// Analyze diff with semantic understanding
		semanticChange, err := analyzer.AnalyzeDiff(diff)
		if err != nil {
			log.Printf("Semantic analysis failed: %v", err)
			// Fall back to rule-based analysis
			commitMessage = generateRuleBasedCommitMessage(diff)
			confidence = 0.6
		} else if semanticChange != nil {
			// Use semantic analysis result
			commitMessage = formatSemanticCommitMessage(semanticChange)
			confidence = semanticChange.Confidence
			
			fmt.Println("🧠 Semantic Analysis Results:")
			fmt.Printf("   Type: %s\n", semanticChange.Type)
			fmt.Printf("   Scope: %s\n", semanticChange.Scope)
			fmt.Printf("   Intent: %s\n", semanticChange.Intent)
			fmt.Printf("   Impact: %s\n", semanticChange.Impact)
			fmt.Printf("   Breaking: %t\n", semanticChange.BreakingChange)
			fmt.Printf("   Confidence: %.1f%%\n", semanticChange.Confidence*100)
			fmt.Printf("   Reasoning: %s\n", semanticChange.Reasoning)
			fmt.Println()
		} else {
			// No semantic change detected, use rule-based
			commitMessage = generateRuleBasedCommitMessage(diff)
			confidence = 0.5
		}
	} else {
		// Use existing rule-based analysis
		commitMessage = generateRuleBasedCommitMessage(diff)
		confidence = 0.6
	}
	
	// Display results
	fmt.Println("─────────────────────────────────────────")
	fmt.Printf("\n>>> based on your changes, cc created the following git commit message for you (confidence: %.1f%%):\n", confidence*100)
	fmt.Println(commitMessage)
}

// formatSemanticCommitMessage formats a semantic change into a commit message
func formatSemanticCommitMessage(change *semantic.SemanticChange) string {
	subject := change.Type
	if change.Scope != "" {
		subject += fmt.Sprintf("(%s)", change.Scope)
	}
	subject += fmt.Sprintf(": %s", change.Description)
	
	// Add breaking change indicator
	if change.BreakingChange {
		if change.Scope != "" {
			subject = fmt.Sprintf("%s(%s)!: %s", change.Type, change.Scope, change.Description)
		} else {
			subject = fmt.Sprintf("%s!: %s", change.Type, change.Description)
		}
	}
	
	// Build body
	var body string
	if change.Intent != "" || change.Impact != "" {
		body = "\n"
		if change.Intent != "" {
			body += fmt.Sprintf("\n%s", change.Intent)
		}
		if change.Impact != "" {
			body += fmt.Sprintf("\nImpact: %s", change.Impact)
		}
	}
	
	// Add breaking change footer if needed
	var footer string
	if change.BreakingChange {
		footer = "\nBREAKING CHANGE: Infrastructure changes may affect existing resources"
	}
	
	return subject + body + footer
}

// generateRuleBasedCommitMessage - placeholder for existing rule-based logic
func generateRuleBasedCommitMessage(diff string) string {
	// This would be the existing cc command logic
	return "chore: update configuration files"
}

// Example usage scenarios
func demonstrateSemanticAnalysis() {
	fmt.Println("🧪 Semantic Analysis Plugin Architecture Demo\n")
	
	// Example 1: Terraform Infrastructure
	terraformDiff := `
diff --git a/infrastructure/vcn.tf b/infrastructure/vcn.tf
new file mode 100644
index 0000000..abc123
--- /dev/null
+++ b/infrastructure/vcn.tf
@@ -0,0 +1,18 @@
+resource "oci_core_vcn" "main" {
+  compartment_id = var.compartment_id
+  cidr_block     = "10.0.0.0/16"
+  dns_label      = "mainvcn"
+  
+  freeform_tags = {
+    Name = "main-vcn"
+    Environment = "production"
+  }
+}
+
+resource "oci_core_internet_gateway" "main" {
+  compartment_id = var.compartment_id
+  vcn_id         = oci_core_vcn.main.id
+  display_name   = "main-igw"
+}
`
	
	fmt.Println("Example 1: New OCI Infrastructure") 
	fmt.Println("Input diff: OCI VCN and internet gateway creation")
	enhancedCCCommand(terraformDiff, true)
	fmt.Println()
	
	// Example 2: Security Improvement
	securityDiff := `
diff --git a/security.tf b/security.tf
index def456..ghi789 100644
--- a/security.tf
+++ b/security.tf
@@ -5,7 +5,7 @@ resource "oci_core_security_list" "web" {
   ingress_security_rules {
     protocol = "6"
     source   = "0.0.0.0/0"
-    source   = "0.0.0.0/0"
+    source   = "10.0.0.0/8"
   }
 }
`
	
	fmt.Println("Example 2: Security Improvement")
	fmt.Println("Input diff: restricting security list access")
	enhancedCCCommand(securityDiff, true)
	fmt.Println()
	
	// Example 3: Non-Terraform file (fallback)
	goDiff := `
diff --git a/main.go b/main.go
index 123..456 100644
--- a/main.go
+++ b/main.go
@@ -10,6 +10,7 @@ func main() {
 	fmt.Println("Hello World")
+	fmt.Println("Added logging")
 }
`
	
	fmt.Println("Example 3: Non-OCI Change (Rule-based fallback)")
	fmt.Println("Input diff: Go code changes")
	enhancedCCCommand(goDiff, true)
}

func main() {
	demonstrateSemanticAnalysis()
}