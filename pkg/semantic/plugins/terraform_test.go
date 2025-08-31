package plugins

import (
	"context"
	"testing"

	"github.com/greenstevester/fast-cc-git-hooks/pkg/semantic"
)

func TestTerraformPlugin(t *testing.T) {
	plugin := NewTerraformPlugin()

	t.Run("plugin metadata", func(t *testing.T) {
		if plugin.Name() != "terraform" {
			t.Errorf("expected name 'terraform', got %s", plugin.Name())
		}

		if plugin.Version() == "" {
			t.Error("version should not be empty")
		}

		extensions := plugin.SupportedExtensions()
		if len(extensions) == 0 {
			t.Error("should support at least one extension")
		}

		found := false
		for _, ext := range extensions {
			if ext == ".tf" {
				found = true
				break
			}
		}
		if !found {
			t.Error("should support .tf extension")
		}
	})

	t.Run("can analyze terraform files", func(t *testing.T) {
		tests := []struct {
			name     string
			file     semantic.FileChange
			expected bool
		}{
			{
				name: "terraform file by extension",
				file: semantic.FileChange{
					Path:         "main.tf",
					AfterContent: `resource "aws_instance" "example" {}`,
				},
				expected: true,
			},
			{
				name: "terraform file by content",
				file: semantic.FileChange{
					Path:         "infrastructure",
					AfterContent: `resource "aws_vpc" "main" { cidr_block = "10.0.0.0/16" }`,
				},
				expected: true,
			},
			{
				name: "non-terraform file",
				file: semantic.FileChange{
					Path:         "main.go",
					AfterContent: `package main`,
				},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := plugin.CanAnalyze(tt.file)
				if result != tt.expected {
					t.Errorf("CanAnalyze() = %v, expected %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("analyze new terraform file", func(t *testing.T) {
		file := semantic.FileChange{
			Path:       "vpc.tf",
			ChangeType: "added",
			AfterContent: `
resource "oci_core_vcn" "main" {
  cidr_block     = "10.0.0.0/16"
  compartment_id = var.compartment_id
  
  freeform_tags = {
    Name = "main-vcn"
  }
}

resource "oci_core_subnet" "public" {
  vcn_id         = oci_core_vcn.main.id
  cidr_block     = "10.0.1.0/24"
  compartment_id = var.compartment_id
}
`,
		}

		ctx := semantic.AnalysisContext{}
		change, err := plugin.AnalyzeFile(context.Background(), file, ctx)

		if err != nil {
			t.Fatalf("AnalyzeFile() error = %v", err)
		}

		if change == nil {
			t.Fatal("expected semantic change, got nil")
		}

		if change.Type != "feat" {
			t.Errorf("expected type 'feat', got %s", change.Type)
		}

		if change.Scope != "network" {
			t.Errorf("expected scope 'network', got %s", change.Scope)
		}

		if !containsString(change.Files, "vpc.tf") {
			t.Errorf("expected files to contain 'vpc.tf', got %v", change.Files)
		}

		if change.Confidence < 0.8 {
			t.Errorf("expected high confidence, got %f", change.Confidence)
		}
	})

	t.Run("analyze modified terraform file", func(t *testing.T) {
		file := semantic.FileChange{
			Path:       "security.tf",
			ChangeType: "modified",
			BeforeContent: `
resource "oci_core_security_list" "web" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.main.id
  display_name   = "web-security-list"
  
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    
    tcp_options {
      min = 80
      max = 80
    }
  }
}
`,
			AfterContent: `
resource "oci_core_security_list" "web" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.main.id
  display_name   = "web-security-list"
  
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    
    tcp_options {
      min = 443
      max = 443
    }
  }
  
  ingress_security_rules {
    protocol = "6"
    source   = "10.0.0.0/8"
    
    tcp_options {
      min = 80
      max = 80
    }
  }
}
`,
			DiffContent: `
-      min = 80
-      max = 80
+      min = 443
+      max = 443
+  }
+  
+  ingress_security_rules {
+    protocol = "6"
+    source   = "10.0.0.0/8"
+    
+    tcp_options {
+      min = 80
+      max = 80
`,
		}

		ctx := semantic.AnalysisContext{}
		change, err := plugin.AnalyzeFile(context.Background(), file, ctx)

		if err != nil {
			t.Fatalf("AnalyzeFile() error = %v", err)
		}

		if change == nil {
			t.Fatal("expected semantic change, got nil")
		}

		// This should be detected as a security improvement
		if change.Type != "fix" {
			t.Errorf("expected type 'fix' for security improvement, got %s", change.Type)
		}

		if change.Scope != "network" {
			t.Errorf("expected scope 'network', got %s", change.Scope)
		}
	})

	t.Run("analyze deleted terraform file", func(t *testing.T) {
		file := semantic.FileChange{
			Path:       "database.tf",
			ChangeType: "deleted",
			BeforeContent: `
resource "oci_database_autonomous_database" "main" {
  compartment_id           = var.compartment_id
  cpu_core_count          = 1
  data_storage_size_in_tbs = 1
  db_name                 = "maindb"
  admin_password          = var.admin_password
}
`,
		}

		ctx := semantic.AnalysisContext{}
		change, err := plugin.AnalyzeFile(context.Background(), file, ctx)

		if err != nil {
			t.Fatalf("AnalyzeFile() error = %v", err)
		}

		if change == nil {
			t.Fatal("expected semantic change, got nil")
		}

		// Deleting RDS should be a breaking change
		if !change.BreakingChange {
			t.Error("expected breaking change for RDS deletion")
		}

		if change.Type != "feat" {
			t.Errorf("expected type 'feat' for breaking change, got %s", change.Type)
		}
	})

	t.Run("extract resource types", func(t *testing.T) {
		content := `
resource "oci_core_vcn" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "oci_core_subnet" "public" {
  vcn_id = oci_core_vcn.main.id
}

resource "oci_core_vcn" "secondary" {
  cidr_block = "172.16.0.0/16"
}
`
		plugin := &TerraformPlugin{}
		resources := plugin.extractResourceTypes(content)

		expected := []string{"oci_core_vcn", "oci_core_subnet"}
		if len(resources) != len(expected) {
			t.Errorf("expected %d resources, got %d", len(expected), len(resources))
		}

		for _, expectedResource := range expected {
			if !containsString(resources, expectedResource) {
				t.Errorf("expected resource %s not found in %v", expectedResource, resources)
			}
		}
	})

	t.Run("determine scope from path and content", func(t *testing.T) {
		plugin := &TerraformPlugin{}

		tests := []struct {
			path     string
			content  string
			expected string
		}{
			{
				path:     "network/vpc.tf",
				content:  "",
				expected: "network",
			},
			{
				path:     "modules/security/main.tf",
				content:  "",
				expected: "security",
			},
			{
				path:     "main.tf",
				content:  `resource "oci_objectstorage_bucket" "main" {}`,
				expected: "storage",
			},
			{
				path:     "compute.tf",
				content:  `resource "oci_core_instance" "web" {}`,
				expected: "compute",
			},
			{
				path:     "unknown.tf",
				content:  `resource "oci_unknown" "test" {}`,
				expected: "infra",
			},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				result := plugin.determineScope(tt.path, tt.content)
				if result != tt.expected {
					t.Errorf("determineScope(%s, content) = %s, expected %s", tt.path, result, tt.expected)
				}
			})
		}
	})

	t.Run("config validation", func(t *testing.T) {
		plugin := &TerraformPlugin{}

		// Valid config
		validConfig := map[string]string{
			"detect_breaking_changes": "true",
			"analyze_security":        "false",
			"provider_sensitivity":    "medium",
		}

		if err := plugin.ValidateConfig(validConfig); err != nil {
			t.Errorf("ValidateConfig() with valid config returned error: %v", err)
		}

		// Invalid config key
		invalidConfig := map[string]string{
			"unknown_key": "value",
		}

		if err := plugin.ValidateConfig(invalidConfig); err == nil {
			t.Error("ValidateConfig() with invalid config should return error")
		}

		// Invalid provider sensitivity
		invalidSensitivity := map[string]string{
			"provider_sensitivity": "invalid",
		}

		if err := plugin.ValidateConfig(invalidSensitivity); err == nil {
			t.Error("ValidateConfig() with invalid provider_sensitivity should return error")
		}
	})
}

// Helper function to check if slice contains string
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
