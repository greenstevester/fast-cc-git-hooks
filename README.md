# 🚀 Fast Conventional Commit Git Hooks

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/greenstevester/fast-cc-git-hooks?style=for-the-badge&logo=github)](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/greenstevester/fast-cc-git-hooks?style=for-the-badge&logo=go)](https://golang.org/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/greenstevester/fast-cc-git-hooks/release.yml?branch=main&style=for-the-badge&logo=github-actions)](https://github.com/greenstevester/fast-cc-git-hooks/actions)
[![License](https://img.shields.io/github/license/greenstevester/fast-cc-git-hooks?style=for-the-badge)](https://github.com/greenstevester/fast-cc-git-hooks/blob/main/LICENSE)

[![Platform Support](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
[![Architecture](https://img.shields.io/badge/arch-amd64%20%7C%20arm64-lightgrey?style=for-the-badge)](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## WHAT IS THIS?

**The fastest way for YOU, to get [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) messages for ALL your git commits!** 

WHY CONVENTIONAL COMMITS? 
- 
- Because when you use them, it damn makes your repo professional looking and as bonus, can:
- auto generate CHANGELOG's
- infer semantic version bumps when required, just from your commit messages.
- make your commit history data-mine-able, and more importantly, **easier to understand** for other developers.

It all stems from something simple, but usually this means more discipline AND (worst of all) strict adherence to a specific workflow. Most people automate this via their bash shell, but bash only goes so far - ergo this tool. 

Strict adherence to a specific workflow alone, is enough to scare away most coders, who just want to build without the added "admin" overhead / accounting / audits while they're in flow.

Its for YOU and the wider community of coders having an *absent of crappy commit messages* polluting an otherwise great repo.

Inspired by **_Boo_**, after seeing lots of great work being committed all with "feat: CGC-0000 added X" messages.

## WHY WOULD I USE THIS?
 
## It's ms fast "git add-on," that gives you ROI from your CPU, does the work (you simply don't have time to do):
- **`fcgh`** - installs a git-hook (locally or globally with local install always taking precedence), that checks your commit messages, ensuring they comply to conventional commit format.
- **`cc`** - preview generated conventional commit messages, based on your changes
- **`ccc`** - generate conventional commit messages AND commit (3 c's for "Create Conventional Commit") (most popular)

## FEATURES PLEASE?
 - Flexibility: easily choose between system-wide or project-specific hook installation
 - Team collaboration: Local installation perfect for team projects with specific requirements
 - Personal workflow: Global installation ideal for individual developers wanting consistent validation
 
🎯 Real-World Scenarios Handled:
- Team Developer: Can remove only local hooks while keeping global ones
- Personal Setup: Can remove global hooks without affecting project-specific local ones


## 🔄 Using Tools Together or separately - your choice!

The tools are **independent**, you can happily use one without the other but work beautifully together:

### Option 1: You write your own commit messages (with fcgh installed)
```bash
# Write your commit message manually
git commit -m "feat: CGC-1245 Added login authentication"
# ↳ The fcgh hook, automatically validates your message ✅
```

### Option 2: Commit messages are generated based on your changes (using ccc) and validated by the fcgh hook
```bash
# Let ccc generate the conventional commit message and execute it for you
ccc
# ↳ ccc analyzes your git changes and creates a conventional commit message
# ↳ When git commit runs, the git hook (installed via fcgh) validates the generated message ✅
# ↳ which means: compliant conventional commits, every time! 
```

### You can mix and match:
- **Use fcgh alone**: Ensures your manual commit messages get automatic validation
- **Use cc alone**: Preview generated conventional commit messages, based on your changes, but you do the commit yourself
- **Use ccc alone**: Generated compliant conventional commits, without validation (but why would you?)
- **Use both**: Generated conventional commits PLUS BONUS validation, talk about the perfect combo! 🎯

**Pro tip**: Start with `ccc` to see what good commit messages look like, then graduate to writing your own!


## ⚡️ Quick Setup / a.k.a. "Take my money! (except its free)"

### Step 1: Install the tools

**Option A: Download the Binary** 

Each release includes **both tools** - you get everything in one download!

*** Windows:**
1. Go to [Releases Page](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
2. Download `fast-cc-git-hooks_windows_amd64.zip`
3. Extract and add both `fcgh.exe` and `cc.exe` to a directory on your PATH
NOTE: In windows, there's an option to change only your environment variables - use that as a preference, or just add both to your PATH.

**🐧 Linux:**
```bash
# Most common (AMD64)
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_linux_amd64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc
sudo mv fcgh cc /usr/local/bin/

# ARM64 (Raspberry Pi, etc.)
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_linux_arm64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc
sudo mv fcgh cc /usr/local/bin/
```

**🍎 macOS:**

**Option A: Easy Installation Script (Recommended)**
```bash
# One-line install - automatically detects your Mac type and handles PATH
curl -sSL https://raw.githubusercontent.com/greenstevester/fast-cc-git-hooks/refs/heads/main/install-macos-source.sh | bash
```

**Option B: Manual Installation**
```bash
# Intel Macs
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_amd64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc
sudo mv fcgh cc /usr/local/bin/

# Apple Silicon (M1/M2/M3) - most common now
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_arm64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc
sudo mv fcgh cc /usr/local/bin/
```

**📱 For macOS Security (Both Options):**
If you see "fcgh cannot be opened" security warning:
1. Click "Done" in the dialog
2. Go to **System Preferences** → **Security & Privacy** → **General** tab
3. Click "Allow Anyway" next to the fcgh message
4. Try running `fcgh` again and click "Open" when prompted

*Note: The installation script automatically handles the quarantine removal to minimize security warnings.*


### Step 2: Verify tools work on your machine (a.k.a. "Is it installed ...smoke test?")

```bash
fcgh version  # shows version info
cc --verbose           # runs cc in verbose mode
ccc --verbose           # runs cc in verbose mode
```

If you see version/help info for both, you're good to go! If not, make sure all tool are on your PATH.

### 🛠️ Step 3: Setup

#### Global Installation (Default - Recommended)
```bash
fcgh setup      # For standard projects
fcgh setup-ent  # For enterprise projects with JIRA validation
```

#### Local Installation (Current Repository Only)
```bash
fcgh setup --local      # Install only for current repository
fcgh setup-ent --local  # Enterprise setup for current repository only
```

#### What the Setup Commands Do

- **Automatically**: 
  - Creates/checks configuration: Ensures a config file exists (creates default if needed)
  - Installs hooks: Then installs the git hooks
- **Guided experience**: Shows exactly what it's doing with clear output
- **Global vs Local Installation**:
  - **Global (default)**: Installs hooks for ALL git repositories on your system
  - **Local (--local)**: Installs hooks only for the current repository
  - **⚡ Precedence Rule**: When both exist, **local always takes precedence** over global configuration

#### 🏆 Configuration Precedence (Important!)

When you have both local and global installations:

```
Local Repository Hook  ➤  ALWAYS WINS  ⚡
Global Git Hook        ➤  Ignored when local exists
```
## Every time you make a commit, compliant conventional commit messages will be your default 

## ✨ How to write good commit messages

Instead of writing messy commits like:
```bash
git commit -m "fixed stuff"
git commit -m "update"
git commit -m "CGC-0000 blah blah blah"
```

Write clear commits like:
```bash
git commit -m "feat: add login button"
git commit -m "fix: CGC-4561 Resolved login bug"  
git commit -m "docs: update README"
```

The format is simple: `type: what you did`

**Common types:**
- `feat` - when you add something new
- `fix` - when you fix a bug
- `docs` - when you update documentation
- `test` - when you add tests
- `chore` - when you do maintenance stuff


### Commands (fcgh)

```bash
fcgh setup-ent # Set up everything (start here!)
fcgh remove    # Smart removal - prompts to choose local/global if both exist
fcgh remove --local   # Remove only from current repository
fcgh remove --global  # Remove only from global git configuration
```

**Test things:**
```bash
fcgh validate "freak: my terrible message"  # Test if a message is bad
```

### Commit Helper Commands (cc)

**Smart commit generation with semantic analysis:**

```bash
cc                     # Generate commit message and copy to clipboard automatically
cc --no-copy           # Generate commit message without copying to clipboard
ccc                    # Generate perfect commit message and commit
cc --verbose           # Show detailed analysis of your changes
cc --help              # Show all available options
```

**🧠 Semantic Analysis for Infrastructure Code:**
The `cc` command now includes intelligent semantic analysis for Terraform files with Oracle OCI awareness. It understands the actual impact of your infrastructure changes and generates contextual commit messages. [See examples →](docs/semantic-analysis-examples.md)

**📋 Clipboard Integration:**
The `cc` command automatically copies the generated git commit command to your clipboard by default. Perfect for quick copy-paste workflows - just run `cc` and press Ctrl+V in your terminal! Use `--no-copy` to disable this behavior.

**That's it!** Most people only ever need `fcgh setup-ent` and `ccc`.

## 🤔 Common Questions

**Q: Do I need all tools?**
A: No! They're completely independent. Install hooks for validation, ccc for generation, or both together.

**Q: What if I want to turn off validation temporarily?**
A: Just run `fcgh remove` and it's gone! Run `fcgh setup-ent` to turn it back on.

**Q: Will this mess up my code?**
A: Nope! The hooks only check commit messages, cc only generates them. Your code stays exactly the same.

**Q: What if cc generates a bad commit message?**
A: The hook will catch it! That's why they work so well together.

**Q: Can I use this on all my projects?**
A: Yes! When you run `fcgh setup-ent`, it works for ALL your Git projects. cc works in any git repo.

**Q: What if I don't like the cc generated message?**
A: Just use `cc` (without ccc) to preview first, then write your own and let the hook validate it!

## 🎯 Examples of Good vs Bad Commits

❌ **Bad commits** (these will be rejected):
```bash
git commit -m "fix"
git commit -m "updated stuff"  
git commit -m "asdf"
git commit -m "WIP"
```

✅ **Good Enterprise commits** (these will work):
```bash
git commit -m "feat: CGC-5641 Added user login page"
git commit -m "fix: CAHC-4132 Resolve password validation bug"
git commit -m "docs: LIND-4777 installation instructions"
```

## ⚙️ Want to Customize? (Optional)

**Most people don't need to do this!** The tool works great out of the box.

But if you want to customize the rules, run:
```bash
fcgh init
```

This creates a file called `fast-cc-config.yaml` in your home directory (`~/.fast-cc-git-hooks/`) that you can edit. It has comments explaining everything!

## 🚨 Troubleshooting

**Problem: The tool says my commit message is bad, but I think it's fine!**

Run this to test your message:
```bash
fcgh validate "your message here"
```

It will tell you exactly what's wrong and how to fix it.

**Problem: I want to turn it off for just one commit**

You can't bypass it easily (that's the point!), but you can run:
```bash
fcgh remove
git commit -m "your message"  
fcgh setup
```

**Problem: It's not working at all**

Try setting it up again:
```bash
fcgh remove
fcgh setup-ent
```

## 🏗️ Advanced Examples

**Want to include ticket numbers?** (like JIRA tickets)
```bash
git commit -m "feat: CGC-1234 add user login"
git commit -m "fix: PROJ-456 resolve password bug"
```

**Want to use scopes?** (optional grouping)
```bash
git commit -m "feat(auth): add login form"
git commit -m "fix(api): resolve timeout issue"
git commit -m "docs(readme): update setup instructions"  
```

The part in `()` is called a "scope" - it's like a category for your change.

## 👨‍💻 For Developers

**Want to help make this tool better?**

```bash
git clone https://github.com/greenstevester/fast-cc-git-hooks.git
cd fast-cc-git-hooks
make build
make test
```

All commits to this project must follow conventional format too! 😄

## 🔄 Changelog Generation Tools

Once you're using conventional commits, you can automatically generate changelogs and manage versioning! Here are the most popular tools:

### 🚀 **Fully Automated Release Tools** (Recommended)

- **[semantic-release](https://github.com/semantic-release/semantic-release)** - The gold standard for automated releases
  - Fully automates: version bumping, changelog generation, git tagging, and publishing
  - Works with GitHub, GitLab, npm, and more
  - Perfect for CI/CD pipelines

- **[commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version)** - Drop-in replacement for `npm version`
  - Handles version bumping, tagging, and CHANGELOG generation
  - Great for manual releases with automation

### 📊 **Changelog-Focused Tools**

- **[conventional-changelog](https://github.com/conventional-changelog/conventional-changelog)** - The original changelog generator
  - Generate changelog from git metadata
  - Multiple presets (Angular, Atom, etc.)
  - Highly customizable

- **[git-cliff](https://git-cliff.org/)** - Modern Rust-based changelog generator
  - Highly customizable templates
  - Fast and reliable
  - Great for complex projects

### 🏢 **Enterprise & Monorepo Tools**

- **[cocogitto](https://github.com/oknozor/cocogitto)** - Complete conventional commits toolkit
  - Version bumping, changelog generation, and commit linting
  - Great for complex workflows

- **[versio](https://github.com/chaaz/versio)** - Monorepo-compatible versioning
  - Handles project dependencies
  - Generates tags and changelogs

### ⚙️ **Language-Specific Tools**

- **Go**: [chglog](https://github.com/goreleaser/chglog) - Template-based changelog generation
- **Python**: [python-semantic-release](https://github.com/relekang/python-semantic-release)
- **PHP**: [php-conventional-changelog](https://github.com/marcocesarato/php-conventional-changelog)
- **Java**: [git-changelog-lib](https://github.com/tomasbjerre/git-changelog-lib)

### 🎯 **Quick Start Recommendations**

1. **For automated CI/CD**: Use `semantic-release`
2. **For manual releases**: Use `commit-and-tag-version`
3. **For just changelogs**: Use `conventional-changelog` or `git-cliff`
4. **For monorepos**: Use `versio` or `cocogitto`

All these tools work perfectly with the conventional commit messages that `fcgh` enforces! 🎉

## 📚 Documentation

- [Semantic Analysis Examples](docs/semantic-analysis-examples.md) - See how the intelligent commit message generation works with real-world infrastructure code examples
- [CLAUDE.md](CLAUDE.md) - Project development guide for Claude Code

## 📝 License

MIT License - do whatever you want with this code!

---

**That's everything!** Remember: just run `fcgh setup-ent` and start writing better commits! 🚀