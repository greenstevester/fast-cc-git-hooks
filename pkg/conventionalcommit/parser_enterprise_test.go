package conventionalcommit

import (
	"reflect"
	"testing"
)

func TestParser_ParseEnterpriseFormats(t *testing.T) {
	tests := []struct {
		want    *Commit
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "CGC format - standard",
			message: "feat(db): CGC-1425 Added new database",
			want: &Commit{
				Type:        "feat",
				Scope:       "db",
				Description: "CGC-1425 Added new database",
				Raw:         "feat(db): CGC-1425 Added new database",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-1425", Raw: "CGC-1425"}},
			},
		},
		{
			name:    "CGC format - different ticket number",
			message: "fix(api): CGC-99 Fixed authentication issue",
			want: &Commit{
				Type:        "fix",
				Scope:       "api",
				Description: "CGC-99 Fixed authentication issue",
				Raw:         "fix(api): CGC-99 Fixed authentication issue",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-99", Raw: "CGC-99"}},
			},
		},
		{
			name:    "CGC format - five digit ticket",
			message: "docs(web): CGC-12345 Updated API documentation",
			want: &Commit{
				Type:        "docs",
				Scope:       "web",
				Description: "CGC-12345 Updated API documentation",
				Raw:         "docs(web): CGC-12345 Updated API documentation",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-12345", Raw: "CGC-12345"}},
			},
		},
		{
			name:    "CGC format with body",
			message: "feat(core): CGC-567 Implemented new validation\n\nThis adds comprehensive input validation.",
			want: &Commit{
				Type:        "feat",
				Scope:       "core",
				Description: "CGC-567 Implemented new validation",
				Body:        "This adds comprehensive input validation.",
				Raw:         "feat(core): CGC-567 Implemented new validation\n\nThis adds comprehensive input validation.",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-567", Raw: "CGC-567"}},
			},
		},
		{
			name:    "CGC format with breaking change",
			message: "feat(db)!: CGC-890 Breaking database schema change",
			want: &Commit{
				Type:        "feat",
				Scope:       "db",
				Description: "CGC-890 Breaking database schema change",
				Breaking:    true,
				Raw:         "feat(db)!: CGC-890 Breaking database schema change",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-890", Raw: "CGC-890"}},
			},
		},
		{
			name:    "CGC format with footer",
			message: "fix(auth): CGC-2001 Fixed token expiration\n\nFixes: CGC-2001\nReviewed-by: John Doe",
			want: &Commit{
				Type:        "fix",
				Scope:       "auth",
				Description: "CGC-2001 Fixed token expiration",
				Body:        "Fixes: CGC-2001\nReviewed-by: John Doe",
				Footer:      "",
				Raw:         "fix(auth): CGC-2001 Fixed token expiration\n\nFixes: CGC-2001\nReviewed-by: John Doe",
				TicketRefs:  []TicketRef{{Type: "JIRA", ID: "CGC-2001", Raw: "CGC-2001"}},
			},
		},
		{
			name:    "Multiple enterprise tickets",
			message: "feat(api): CGC-100 PROJ-200 Added multi-tenant support",
			want: &Commit{
				Type:        "feat",
				Scope:       "api",
				Description: "CGC-100 PROJ-200 Added multi-tenant support",
				Raw:         "feat(api): CGC-100 PROJ-200 Added multi-tenant support",
				TicketRefs: []TicketRef{
					{Type: "JIRA", ID: "CGC-100", Raw: "CGC-100"},
					{Type: "JIRA", ID: "PROJ-200", Raw: "PROJ-200"},
				},
			},
		},
		{
			name:    "Enterprise ticket with GitHub reference",
			message: "fix(cli): CAVH-3334 Fixed CLI parsing (#456)",
			want: &Commit{
				Type:        "fix",
				Scope:       "cli",
				Description: "CAVH-3334 Fixed CLI parsing (#456)",
				Raw:         "fix(cli): CAVH-3334 Fixed CLI parsing (#456)",
				TicketRefs: []TicketRef{
					{Type: "GITHUB", ID: "456", Raw: "#456"},
					{Type: "JIRA", ID: "CAVH-3334", Raw: "CAVH-3334"},
				},
			},
		},
	}

	parser := DefaultParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParser_CGCFormatValidation(t *testing.T) {
	parser := DefaultParser()

	tests := []struct {
		check   func(*testing.T, *Commit)
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "Valid CGC format with uppercase description",
			message: "feat(db): CGC-1425 Added new database connection pooling",
			check: func(t *testing.T, c *Commit) {
				if c.Type != "feat" {
					t.Errorf("Expected type 'feat', got %s", c.Type)
				}
				if c.Scope != "db" {
					t.Errorf("Expected scope 'db', got %s", c.Scope)
				}
				if !c.HasJIRATicket() {
					t.Error("Expected JIRA ticket, but none found")
				}
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-1425" {
					t.Errorf("Expected JIRA ticket CGC-1425, got %+v", jiraTickets)
				}
				expectedDesc := "CGC-1425 Added new database connection pooling"
				if c.Description != expectedDesc {
					t.Errorf("Expected description '%s', got '%s'", expectedDesc, c.Description)
				}
			},
		},
		{
			name:    "CGC format with single digit ticket",
			message: "fix(core): CGC-1 Fixed critical startup bug",
			check: func(t *testing.T, c *Commit) {
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-1" {
					t.Errorf("Expected JIRA ticket CGC-1, got %+v", jiraTickets)
				}
			},
		},
		{
			name:    "CGC format with max digits (5)",
			message: "refactor(api): CGC-99999 Restructured API endpoints",
			check: func(t *testing.T, c *Commit) {
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-99999" {
					t.Errorf("Expected JIRA ticket CGC-99999, got %+v", jiraTickets)
				}
			},
		},
		{
			name:    "CGC ticket in body",
			message: "feat(auth): Implemented OAuth2\n\nThis implements OAuth2 as specified in CGC-777",
			check: func(t *testing.T, c *Commit) {
				if c.Body != "This implements OAuth2 as specified in CGC-777" {
					t.Errorf("Unexpected body: %s", c.Body)
				}
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-777" {
					t.Errorf("Expected JIRA ticket CGC-777 from body, got %+v", jiraTickets)
				}
			},
		},
		{
			name:    "Multiple CGC tickets",
			message: "feat(db): CGC-100 CGC-101 Implemented sharding and replication",
			check: func(t *testing.T, c *Commit) {
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 2 {
					t.Errorf("Expected 2 JIRA tickets, got %d: %+v", len(jiraTickets), jiraTickets)
				}
				expectedIDs := map[string]bool{"CGC-100": false, "CGC-101": false}
				for _, ticket := range jiraTickets {
					if _, ok := expectedIDs[ticket.ID]; ok {
						expectedIDs[ticket.ID] = true
					}
				}
				for id, found := range expectedIDs {
					if !found {
						t.Errorf("Missing expected JIRA ticket: %s", id)
					}
				}
			},
		},
		{
			name:    "CGC format with special characters in description",
			message: "fix(api): CGC-2024 Fixed API response for /users/:id endpoint",
			check: func(t *testing.T, c *Commit) {
				expectedDesc := "CGC-2024 Fixed API response for /users/:id endpoint"
				if c.Description != expectedDesc {
					t.Errorf("Expected description '%s', got '%s'", expectedDesc, c.Description)
				}
			},
		},
		{
			name:    "CGC ticket with other enterprise tickets",
			message: "feat(integration): CGC-500 INTG-600 SAP-700 Added enterprise connectors",
			check: func(t *testing.T, c *Commit) {
				if len(c.TicketRefs) != 3 {
					t.Errorf("Expected 3 ticket references, got %d: %+v", len(c.TicketRefs), c.TicketRefs)
				}
				// Check that CGC ticket is recognized
				foundCGC := false
				for _, ref := range c.TicketRefs {
					if ref.ID == "CGC-500" {
						foundCGC = true
						break
					}
				}
				if !foundCGC {
					t.Error("CGC-500 ticket not found in references")
				}
			},
		},
		{
			name: "CGC format in complex commit",
			message: `feat(db)!: CGC-3000 Major database refactoring

BREAKING CHANGE: This completely changes the database schema.

The new schema provides:
- Better performance
- Improved scalability
- Reduced storage requirements

Implements: CGC-3000, CGC-3001, CGC-3002
Fixes: #123
Co-authored-by: Jane Doe`,
			check: func(t *testing.T, c *Commit) {
				if !c.Breaking {
					t.Error("Expected breaking change")
				}
				// Should find 4 CGC tickets (one in header, three in footer) and 1 GitHub issue
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) < 3 {
					t.Errorf("Expected at least 3 unique JIRA tickets, got %d: %+v", len(jiraTickets), jiraTickets)
				}
				// Check for specific CGC tickets
				expectedCGCs := map[string]bool{
					"CGC-3000": false,
					"CGC-3001": false,
					"CGC-3002": false,
				}
				for _, ticket := range jiraTickets {
					if _, ok := expectedCGCs[ticket.ID]; ok {
						expectedCGCs[ticket.ID] = true
					}
				}
				for id, found := range expectedCGCs {
					if !found {
						t.Errorf("Missing expected CGC ticket: %s", id)
					}
				}
			},
		},
		{
			name:    "CGC with empty scope",
			message: "feat(): CGC-4444 Global feature implementation",
			check: func(t *testing.T, c *Commit) {
				if c.Scope != "" {
					t.Errorf("Expected empty scope, got '%s'", c.Scope)
				}
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-4444" {
					t.Errorf("Expected JIRA ticket CGC-4444, got %+v", jiraTickets)
				}
			},
		},
		{
			name:    "CGC format without scope",
			message: "feat: CGC-5555 Added new feature",
			check: func(t *testing.T, c *Commit) {
				if c.Scope != "" {
					t.Errorf("Expected no scope, got '%s'", c.Scope)
				}
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "CGC-5555" {
					t.Errorf("Expected JIRA ticket CGC-5555, got %+v", jiraTickets)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit, err := parser.Parse(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if commit != nil && tt.check != nil {
				tt.check(t, commit)
			}
		})
	}
}

func TestParser_EnterpriseTicketPatterns(t *testing.T) {
	tests := []struct {
		expected []TicketRef
		name     string
		message  string
	}{
		{
			name:    "CGC ticket pattern",
			message: "feat: implement auth CGC-1425",
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-1425", Raw: "CGC-1425"},
			},
		},
		{
			name:    "Multiple CGC tickets",
			message: "feat: implement CGC-100 CGC-200 CGC-300",
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-100", Raw: "CGC-100"},
				{Type: "JIRA", ID: "CGC-200", Raw: "CGC-200"},
				{Type: "JIRA", ID: "CGC-300", Raw: "CGC-300"},
			},
		},
		{
			name:    "Mixed enterprise tickets",
			message: "feat: integrate systems CGC-1000 SAP-2000 INTG-3000",
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-1000", Raw: "CGC-1000"},
				{Type: "JIRA", ID: "SAP-2000", Raw: "SAP-2000"},
				{Type: "JIRA", ID: "INTG-3000", Raw: "INTG-3000"},
			},
		},
		{
			name:    "CGC with other formats",
			message: "fix: resolve CGC-999 #123 [TASK-456]",
			expected: []TicketRef{
				{Type: "GITHUB", ID: "123", Raw: "#123"},
				{Type: "GENERIC", ID: "TASK-456", Raw: "[TASK-456]"},
				{Type: "JIRA", ID: "CGC-999", Raw: "CGC-999"},
			},
		},
		{
			name:    "CGC in various positions",
			message: "CGC-111 feat: start CGC-222 middle CGC-333 end",
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-111", Raw: "CGC-111"},
				{Type: "JIRA", ID: "CGC-222", Raw: "CGC-222"},
				{Type: "JIRA", ID: "CGC-333", Raw: "CGC-333"},
			},
		},
		{
			name:    "CGC with different case (should match uppercase only)",
			message: "feat: implement cgc-123 CGC-456 Cgc-789",
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-456", Raw: "CGC-456"},
			},
		},
		{
			name: "CGC in complex message",
			message: `feat(api): CGC-1001 Implement new API endpoints

This commit adds the following endpoints as specified in CGC-1001:
- GET /api/v2/users (CGC-1002)
- POST /api/v2/users (CGC-1003)
- DELETE /api/v2/users/:id (CGC-1004)

Related to: CGC-1000
Fixes: CGC-1001, CGC-1002, CGC-1003, CGC-1004`,
			expected: []TicketRef{
				{Type: "JIRA", ID: "CGC-1001", Raw: "CGC-1001"},
				{Type: "JIRA", ID: "CGC-1002", Raw: "CGC-1002"},
				{Type: "JIRA", ID: "CGC-1003", Raw: "CGC-1003"},
				{Type: "JIRA", ID: "CGC-1004", Raw: "CGC-1004"},
				{Type: "JIRA", ID: "CGC-1000", Raw: "CGC-1000"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := parseTicketRefs(tt.message)

			if len(refs) != len(tt.expected) {
				t.Fatalf("parseTicketRefs() returned %d refs, expected %d\nGot: %+v\nExpected: %+v",
					len(refs), len(tt.expected), refs, tt.expected)
			}

			for i, ref := range refs {
				expected := tt.expected[i]
				if ref.Type != expected.Type || ref.ID != expected.ID || ref.Raw != expected.Raw {
					t.Errorf("parseTicketRefs()[%d] = %+v, expected %+v", i, ref, expected)
				}
			}
		})
	}
}

func BenchmarkParser_ParseCGCFormat(b *testing.B) {
	parser := DefaultParser()
	message := "feat(db): CGC-1425 Added new database connection pooling with improved performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parser.Parse(message); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParser_ParseComplexCGC(b *testing.B) {
	parser := DefaultParser()
	message := `feat(api)!: CGC-2000 Major API refactoring

BREAKING CHANGE: Complete API restructuring for v2.

This commit implements:
- New REST endpoints (CGC-2001)
- GraphQL support (CGC-2002)
- WebSocket connections (CGC-2003)

Implements: CGC-2000, CGC-2001, CGC-2002, CGC-2003
Fixes: #456, #457
Related: SAP-1000, INTG-2000`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parser.Parse(message); err != nil {
			b.Fatal(err)
		}
	}
}