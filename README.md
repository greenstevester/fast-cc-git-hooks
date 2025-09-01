# ⚡ Fast Conventional Commits

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/greenstevester/fast-cc-git-hooks?style=for-the-badge&logo=github)](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
[![Platform Support](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

> **Write [conventional](https://www.conventionalcommits.org/en/v1.0.0/) git commits automatically.** Never think about allowed commit formats again, this toolset does it for you.

```bash
# Instead of this:
git commit -m "fix stuff"
git commit -m "WIP"
git commit -m "updated files"

# Get this automatically:
git commit -m "feat(auth): add JWT token validation"
git commit -m "fix(api): resolve timeout in user login"
git commit -m "docs(readme): update installation instructions"
```

## WHAT IS THIS?

Lightening-fast "git add-ons," that do the repetitive work:

* **`fcgh`** - installs a git-hook (locally or globally), to check your commit messages, ensuring they comply to conventional commit format.
* **`ccg`** - generates conventional commit messages, based on your recent changes (and copies it to your clipboard)
* **`ccdo`** - same as ccg (above), except it ALSO does the git commit FOR YOU


## 🚀 Quick Start

**Step 1:** Install

**For Linux/Windows:** Download from [releases](https://github.com/greenstevester/fast-cc-git-hooks/releases) and extract to your PATH

```bash
# macOS 
curl -sSL https://raw.githubusercontent.com/greenstevester/fast-cc-git-hooks/main/install-macos.sh | bash

# Linux/Windows - Manual install (choose your platform below)
```

**Step 2:** Set up once, globally (for all repos - this is the default)
```bash
fcgh setup-ent
```
Or Locally (from within the current repo)

```bash
fcgh setup-ent --local
```

**Step 3:** set a JIRA ticket reference (when you start working on a ticket)
```bash
ccg set-jira CGC-1234
```

**Step 4.a:** use `ccg` to preview the commit message (copied to the clipboard)

NOTE: You don't need to run git add . manually - both ccg and ccdo handle this automatically as part of their workflow.

```bash
ccg  # Generates commit message with ticket reference (if you set one) and copies to the clipboard
```
OR 

**Step 4.b:** Generate + Commit automatically (fastest)
```bash
ccdo  # Generates commit message with ticket reference (if you set one) + commits automatically
```

🎉 **Done!** Every commit is now formatted.


## 📖 To recap, choose your workflow

### **🤖 Fully Automated** (fastest)
Let the tools do everything for you:
```bash
ccdo  # Analyzes changes, creates perfect commit message, commits it
```
- **No heavy lifting for message formulation required**
- **Always follows conventions** 
- **Perfect for daily development**

### **👀 Preview First** 
See the generated message before committing:
```bash
ccg   # Shows suggested commit message + copies to clipboard
# Review the message, then paste with Ctrl+V:
git commit -m "feat(api): add user authentication endpoint"
```
- **Review before committing**
- **Learn good commit patterns**
- **Full control over final message**

### **✍️ Manual + Validation**
Write your own messages with automatic validation:
```bash
git commit -m "feat: add cool feature"
# ✅ Commit message is valid (fcgh hook is invoked and validates automatically)

git commit -m "fix stuff" 
# ❌ Commit message validation failed: invalid format
```
- **Write your own messages**
- **Automatic format checking** 
- **Learn by doing**

<details>

**Windows:**
1. Download `fast-cc-git-hooks_windows_amd64.zip` from [releases](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
2. Extract `fcgh.exe`, `ccg.exe`, and `ccdo.exe` 
3. Add to your PATH

**Linux:**
```bash
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_linux_amd64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh ccg ccdo
sudo mv fcgh ccg ccdo /usr/local/bin/
```

**macOS:**
```bash
# Intel Macs
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_amd64.tar.gz
tar -xzf fcgh.tar.gz && chmod +x fcgh ccg ccdo && sudo mv fcgh ccg ccdo /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_arm64.tar.gz
tar -xzf fcgh.tar.gz && chmod +x fcgh ccg ccdo && sudo mv fcgh ccg ccdo /usr/local/bin/
```
</details>

<details>
<summary><strong>🔧 Build from Source</strong></summary>

Requires Go 1.25+
```bash
git clone https://github.com/greenstevester/fast-cc-git-hooks.git
cd fast-cc-git-hooks
make build-all-tools
make install-all
```
</details>

## ⚙️ Configuration

**Basic Setup (Works for 90% of projects):**
```bash
fcgh setup-ent  # Sets up everything with enterprise features
```

**Project-Only Setup:**
```bash
fcgh setup-ent --local  # Only for current repository
```

**Custom Configuration:**
```bash
fcgh init  # Creates ~/.fast-cc/fast-cc-config.yaml for customization
```

<details>
<summary><strong>🏢 Enterprise Features</strong></summary>

- **JIRA Integration**: Auto-include ticket numbers
- **Team Scopes**: Predefined scopes (api, web, cli, db, etc.)
- **Custom Rules**: Company-specific validation
- **Advanced Patterns**: Complex commit requirements

```bash
ccg set-jira PROJ-1234    # Auto-include JIRA ticket
ccg jira-status           # Check current ticket
```
</details>

## 🛠️ Commands Reference

| Command | Purpose | Example |
|---------|---------|---------|
| `fcgh setup-ent` | One-time setup with validation | Sets up git hooks |
| `ccdo` | Generate + commit automatically | `ccdo` |
| `ccg` | Generate message (preview only) | `ccg --verbose` |
| `fcgh validate` | Test a commit message | `fcgh validate "feat: add feature"` |
| `fcgh status` | Show hook and JIRA status | Shows installation status and current ticket |

## ❓ Common Questions

<details>
<summary><strong>Q: Do I need all three tools?</strong></summary>

No! They work independently:
- **Just hooks**: `fcgh setup-ent` validates your manual messages
- **Just generation**: `ccg` generates messages without validation  
- **Full automation**: All three together for zero-thought commits
</details>

<details>
<summary><strong>Q: What if I don't like the generated message?</strong></summary>

Use `ccg` (without `ccdo`) to preview first. Copy the generated command and modify it before running.
</details>

<details>
<summary><strong>Q: Does this change my code?</strong></summary>

No! It only affects commit messages. Your code stays exactly the same.
</details>

<details>
<summary><strong>Q: Can I turn it off temporarily?</strong></summary>

```bash
fcgh remove      # Remove hooks
fcgh setup-ent   # Add them back
```
</details>

## 🎯 Examples

**✅ Good commits (auto-generated):**
```
feat(auth): add JWT token validation
fix(api): resolve timeout in user login  
docs(readme): update installation instructions
test(user): add integration tests for signup flow
```

**❌ Bad commits (blocked by validation):**
```
fix stuff
WIP
updated files
asdf
```

## 🔗 Why Conventional Commits?

- **📈 Auto-generate changelogs** from commit history
- **🏷️ Semantic version bumps** based on commit types  
- **🔍 Searchable history** with consistent formatting
- **👥 Team collaboration** with clear change communication
- **🤖 CI/CD integration** for automated workflows

<details>
<summary><strong>📚 Advanced Topics</strong></summary>

### Semantic Analysis
The tools include intelligent analysis for infrastructure code, particularly Terraform with Oracle OCI awareness.

### JIRA Integration
```bash
ccg set-jira PROJ-1234     # Set ticket for next 10 commits
ccg clear-jira            # Remove ticket
ccg jira-history          # View ticket history
```

### Custom Scopes
Edit `~/.fast-cc/fast-cc-config.yaml` to add project-specific scopes:
```yaml
scopes:
  - api
  - web  
  - cli
  - database
  - auth
  - docs
```

### Multiple Install Types
- **Global**: Works for all Git repositories on your machine
- **Local**: Works only for current repository
- **Local always wins** when both are installed

### Changelog Generation Tools
Once using conventional commits, you can automate your entire release process:
- **[See our complete guide →](docs/ChangelogGenerationTools.md)** - Detailed comparison of all tools
- **Quick picks**: [semantic-release](https://github.com/semantic-release/semantic-release) (automated), [conventional-changelog](https://github.com/conventional-changelog/conventional-changelog) (manual)
</details>

## 🤝 Contributing

Found a bug? Want a feature? [Open an issue](https://github.com/greenstevester/fast-cc-git-hooks/issues) or submit a PR.

## 📄 License

MIT License - do whatever you want with this code!

---

**TL;DR: `fcgh setup-ent` + `ccdo` = perfect commits forever** 🚀

---

📚 **Learn More:**
- **[Why these tools exist →](docs/10SecondHistory.md)** - The inspiration behind fcgh
- **[Changelog automation →](docs/ChangelogGenerationTools.md)** - Next-level release management