package conventionalcommit

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    *Commit
		wantErr bool
	}{
		{
			name:    "simple commit",
			message: "feat: add new feature",
			want: &Commit{
				Type:        "feat",
				Description: "add new feature",
				Raw:         "feat: add new feature",
			},
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
		name   string
		commit *Commit
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
		name   string
		commit *Commit
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
		_, _ = parser.Parse(message)
	}
}

func BenchmarkParser_ParseSimple(b *testing.B) {
	parser := DefaultParser()
	message := "feat: add new feature"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(message)
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