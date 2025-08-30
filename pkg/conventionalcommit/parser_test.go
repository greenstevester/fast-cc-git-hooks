package conventionalcommit

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		want    *Commit
		name    string
		message string
		wantErr bool
	}{
		{
			want: &Commit{
				Type:        "feat",
				Scope:       "",
				Description: "add new feature",
				Body:        "",
				Footer:      "",
				Raw:         "feat: add new feature",
				TicketRefs:  nil,
				Breaking:    false,
			},
			name:    "simple commit",
			message: "feat: add new feature",
		},
		{
			name:    "commit with scope",
			message: "fix(api): resolve authentication issue",
			want: &Commit{
				Type:        "fix",
				Scope:       "api",
				Description: "resolve authentication issue",
				Raw:         "fix(api): resolve authentication issue",
			},
		},
		{
			name:    "breaking change with exclamation",
			message: "feat!: breaking API change",
			want: &Commit{
				Type:        "feat",
				Breaking:    true,
				Description: "breaking API change",
				Raw:         "feat!: breaking API change",
			},
		},
		{
			name:    "commit with scope and breaking change",
			message: "refactor(core)!: reorganize module structure",
			want: &Commit{
				Type:        "refactor",
				Scope:       "core",
				Breaking:    true,
				Description: "reorganize module structure",
				Raw:         "refactor(core)!: reorganize module structure",
			},
		},
		{
			name:    "commit with body",
			message: "feat: add feature\n\nThis is the body of the commit message.",
			want: &Commit{
				Type:        "feat",
				Description: "add feature",
				Body:        "This is the body of the commit message.",
				Raw:         "feat: add feature\n\nThis is the body of the commit message.",
			},
		},
		{
			name:    "commit with footer",
			message: "fix: bug fix\n\nSome body text\n\nFixes: #123",
			want: &Commit{
				Type:        "fix",
				Description: "bug fix",
				Body:        "Some body text",
				Footer:      "Fixes: #123",
				Raw:         "fix: bug fix\n\nSome body text\n\nFixes: #123",
				TicketRefs:  []TicketRef{{Type: "GITHUB", ID: "123", Raw: "#123"}},
			},
		},
		{
			name:    "breaking change in footer",
			message: "feat: new feature\n\nBREAKING CHANGE: This breaks the API",
			want: &Commit{
				Type:        "feat",
				Description: "new feature",
				Breaking:    true,
				Footer:      "BREAKING CHANGE: This breaks the API",
				Raw:         "feat: new feature\n\nBREAKING CHANGE: This breaks the API",
			},
		},
		{
			name:    "empty scope allowed",
			message: "docs(): update README",
			want: &Commit{
				Type:        "docs",
				Scope:       "",
				Description: "update README",
				Raw:         "docs(): update README",
			},
		},
		{
			name:    "multiple footers",
			message: "feat: feature\n\nBody\n\nSigned-off-by: John Doe\nCo-authored-by: Jane Doe",
			want: &Commit{
				Type:        "feat",
				Description: "feature",
				Body:        "Body",
				Footer:      "Signed-off-by: John Doe\nCo-authored-by: Jane Doe",
				Raw:         "feat: feature\n\nBody\n\nSigned-off-by: John Doe\nCo-authored-by: Jane Doe",
			},
		},
		{
			name:    "empty message",
			message: "",
			wantErr: true,
		},
		{
			name:    "invalid format in strict mode",
			message: "This is not a conventional commit",
			wantErr: true,
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

func TestCommit_Format(t *testing.T) {
	tests := []struct {
		commit *Commit
		name   string
		want   string
	}{
		{
			name: "simple commit",
			commit: &Commit{
				Type:        "feat",
				Description: "add new feature",
			},
			want: "feat: add new feature",
		},
		{
			name: "commit with scope",
			commit: &Commit{
				Type:        "fix",
				Scope:       "api",
				Description: "fix bug",
			},
			want: "fix(api): fix bug",
		},
		{
			name: "breaking change",
			commit: &Commit{
				Type:        "feat",
				Breaking:    true,
				Description: "breaking change",
			},
			want: "feat!: breaking change",
		},
		{
			name: "full commit",
			commit: &Commit{
				Type:        "feat",
				Scope:       "core",
				Breaking:    true,
				Description: "major update",
				Body:        "This is the body",
				Footer:      "BREAKING CHANGE: API changed",
			},
			want: "feat(core)!: major update\n\nThis is the body\n\nBREAKING CHANGE: API changed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.Format(); got != tt.want {
				t.Errorf("Commit.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_Header(t *testing.T) {
	tests := []struct {
		commit *Commit
		name   string
		want   string
	}{
		{
			name: "simple header",
			commit: &Commit{
				Type:        "feat",
				Description: "new feature",
			},
			want: "feat: new feature",
		},
		{
			name: "header with scope and breaking",
			commit: &Commit{
				Type:        "fix",
				Scope:       "api",
				Breaking:    true,
				Description: "fix critical bug",
			},
			want: "fix(api)!: fix critical bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.Header(); got != tt.want {
				t.Errorf("Commit.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	parser := DefaultParser()
	message := "feat(scope): add new feature\n\nThis is the body of the commit.\nIt has multiple lines.\n\nBREAKING CHANGE: This is a breaking change\nFixes: #123\nSigned-off-by: John Doe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parser.Parse(message); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParser_ParseSimple(b *testing.B) {
	parser := DefaultParser()
	message := "feat: add new feature"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parser.Parse(message); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCommit_Format(b *testing.B) {
	commit := &Commit{
		Type:        "feat",
		Scope:       "core",
		Breaking:    true,
		Description: "major update",
		Body:        "This is a long body with multiple lines\nand more content here",
		Footer:      "BREAKING CHANGE: API changed\nFixes: #123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = commit.Format()
	}
}

func TestParseTicketRefs(t *testing.T) {
	tests := []struct {
		expected []TicketRef
		name     string
		message  string
	}{
		{
			name:    "JIRA ticket in header",
			message: "feat: add user auth PROJ-123",
			expected: []TicketRef{
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
			},
		},
		{
			name:    "multiple JIRA tickets",
			message: "feat: implement auth PROJ-123 DEFG-456",
			expected: []TicketRef{
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
				{Type: "JIRA", ID: "DEFG-456", Raw: "DEFG-456"},
			},
		},
		{
			name:    "GitHub issue",
			message: "fix: resolve bug #123",
			expected: []TicketRef{
				{Type: "GITHUB", ID: "123", Raw: "#123"},
			},
		},
		{
			name:    "GitHub issue with GH prefix",
			message: "fix: resolve bug GH-456",
			expected: []TicketRef{
				{Type: "GITHUB", ID: "456", Raw: "GH-456"},
			},
		},
		{
			name:    "Linear ticket (treated as JIRA by default)",
			message: "feat: add feature ABC-123",
			expected: []TicketRef{
				{Type: "JIRA", ID: "ABC-123", Raw: "ABC-123"},
			},
		},
		{
			name:    "Generic bracketed ticket",
			message: "feat: implement [PROJ-789]",
			expected: []TicketRef{
				{Type: "GENERIC", ID: "PROJ-789", Raw: "[PROJ-789]"},
			},
		},
		{
			name:    "mixed ticket types",
			message: "feat: implement auth PROJ-123 #456 [ABC-789]",
			expected: []TicketRef{
				{Type: "GITHUB", ID: "456", Raw: "#456"},
				{Type: "GENERIC", ID: "ABC-789", Raw: "[ABC-789]"},
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
			},
		},
		{
			name: "tickets in footer",
			message: `feat: add authentication

Implements OAuth2 login flow with JWT tokens.

Fixes PROJ-123
Closes #456`,
			expected: []TicketRef{
				{Type: "GITHUB", ID: "456", Raw: "#456"},
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
			},
		},
		{
			name:     "no tickets",
			message:  "feat: add new feature without references",
			expected: []TicketRef{},
		},
		{
			name:    "duplicate tickets",
			message: "feat: implement PROJ-123 and fix PROJ-123",
			expected: []TicketRef{
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := parseTicketRefs(tt.message)

			if len(refs) != len(tt.expected) {
				t.Fatalf("parseTicketRefs() returned %d refs, expected %d", len(refs), len(tt.expected))
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

func TestCommit_HasTicketRefs(t *testing.T) {
	tests := []struct {
		commit   *Commit
		name     string
		expected bool
	}{
		{
			commit: &Commit{
				TicketRefs: []TicketRef{
					{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
				},
			},
			name:     "has tickets",
			expected: true,
		},
		{
			name:     "no tickets",
			commit:   &Commit{TicketRefs: []TicketRef{}},
			expected: false,
		},
		{
			name:     "nil tickets",
			commit:   &Commit{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.HasTicketRefs(); got != tt.expected {
				t.Errorf("Commit.HasTicketRefs() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCommit_HasJIRATicket(t *testing.T) {
	tests := []struct {
		commit   *Commit
		name     string
		expected bool
	}{
		{
			name: "has JIRA ticket",
			commit: &Commit{
				TicketRefs: []TicketRef{
					{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
					{Type: "GITHUB", ID: "456", Raw: "#456"},
				},
			},
			expected: true,
		},
		{
			name: "no JIRA tickets",
			commit: &Commit{
				TicketRefs: []TicketRef{
					{Type: "GITHUB", ID: "456", Raw: "#456"},
				},
			},
			expected: false,
		},
		{
			name:     "no tickets at all",
			commit:   &Commit{TicketRefs: []TicketRef{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.HasJIRATicket(); got != tt.expected {
				t.Errorf("Commit.HasJIRATicket() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCommit_GetJIRATickets(t *testing.T) {
	tests := []struct {
		expected []TicketRef
		commit   *Commit
		name     string
	}{
		{
			name: "multiple JIRA tickets",
			commit: &Commit{
				TicketRefs: []TicketRef{
					{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
					{Type: "GITHUB", ID: "456", Raw: "#456"},
					{Type: "JIRA", ID: "ABC-789", Raw: "ABC-789"},
				},
			},
			expected: []TicketRef{
				{Type: "JIRA", ID: "PROJ-123", Raw: "PROJ-123"},
				{Type: "JIRA", ID: "ABC-789", Raw: "ABC-789"},
			},
		},
		{
			name: "no JIRA tickets",
			commit: &Commit{
				TicketRefs: []TicketRef{
					{Type: "GITHUB", ID: "456", Raw: "#456"},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.commit.GetJIRATickets()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Commit.GetJIRATickets() = %+v, expected %+v", got, tt.expected)
			}
		})
	}
}

func TestParser_ParseWithTickets(t *testing.T) {
	parser := DefaultParser()

	tests := []struct {
		check   func(*testing.T, *Commit)
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "conventional commit with JIRA ticket",
			message: "feat(auth): implement OAuth2 PROJ-123",
			wantErr: false,
			check: func(t *testing.T, c *Commit) {
				if c.Type != "feat" {
					t.Errorf("Expected type 'feat', got %s", c.Type)
				}
				if c.Scope != "auth" {
					t.Errorf("Expected scope 'auth', got %s", c.Scope)
				}
				if !c.HasJIRATicket() {
					t.Error("Expected JIRA ticket, but none found")
				}
				jiraTickets := c.GetJIRATickets()
				if len(jiraTickets) != 1 || jiraTickets[0].ID != "PROJ-123" {
					t.Errorf("Expected JIRA ticket PROJ-123, got %+v", jiraTickets)
				}
			},
		},
		{
			name: "commit with multiple ticket types",
			message: `fix: resolve authentication bug PROJ-456

This fixes the OAuth2 token validation issue.
Also addresses GitHub issue #789.

Closes: PROJ-456
Fixes: #789`,
			wantErr: false,
			check: func(t *testing.T, c *Commit) {
				if len(c.TicketRefs) != 2 { // PROJ-456 and #789 (duplicates removed)
					t.Errorf("Expected 2 unique ticket references, got %d: %+v", len(c.TicketRefs), c.TicketRefs)
				}
				if !c.HasJIRATicket() {
					t.Error("Expected JIRA ticket, but none found")
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
