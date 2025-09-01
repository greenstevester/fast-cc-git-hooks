= 🚀 Fast Conventional Commit Git Hooks

image:https://img.shields.io/github/v/release/greenstevester/fast-cc-git-hooks?style=for-the-badge&logo=github[GitHub release (latest SemVer),link=https://github.com/greenstevester/fast-cc-git-hooks/releases/latest]
image:https://img.shields.io/github/go-mod/go-version/greenstevester/fast-cc-git-hooks?style=for-the-badge&logo=go[Go Version,link=https://golang.org/]
image:https://img.shields.io/github/actions/workflow/status/greenstevester/fast-cc-git-hooks/release.yml?branch=main&style=for-the-badge&logo=github-actions[Build Status,link=https://github.com/greenstevester/fast-cc-git-hooks/actions]
image:https://img.shields.io/github/license/greenstevester/fast-cc-git-hooks?style=for-the-badge[License,link=https://github.com/greenstevester/fast-cc-git-hooks/blob/main/LICENSE]

image:https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge[Platform Support,link=https://github.com/greenstevester/fast-cc-git-hooks/releases/latest]
image:https://img.shields.io/badge/arch-amd64%20%7C%20arm64-lightgrey?style=for-the-badge[Architecture,link=https://github.com/greenstevester/fast-cc-git-hooks/releases/latest]
image:https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge[Conventional Commits,link=https://conventionalcommits.org]

== WHAT IS THIS?

*The fastest way for YOU, to get https://www.conventionalcommits.org/en/v1.0.0/[conventional commit] messages for ALL your git commits!* 

WHY CONVENTIONAL COMMITS?  WHY WOULD YOU CARE?
- clean, consistent, and easy to understand commit messages
- auto-generate CHANGELOG's
- infer semantic version bumps when required, just from your commit messages.
- Because when you use them, it damn makes your repo professional looking and as bonus, can:
- make your commit history data-mine-able, and more importantly, *easier to understand* for other developers.

It all stems from something simple, but usually this means more discipline AND (worst of all) strict adherence to a specific workflow. Most people automate this via their bash shell, but bash only goes so far - ergo this tool. 

Strict adherence to a specific workflow alone is enough to scare away most coders, who just want to build without the added "admin" overhead / accounting / audits while they're in flow.

This small set of performance-obsessed tools is for YOU and the wider community of coders to free the world of crappy commit messages and to prevent further pollution of potentially grep git repos.

This small tool-set is inspired by *_Boo_* ...after seeing lots of their git commits all with "feat: CGC-0000 added X" messages.

== WHY WOULD I USE THIS?

== It's ms fast "git add-on," that gives you ROI from your CPU, does the work (you simply don't have time to do):

* *`fcgh`* - installs a git-hook (locally or globally with local install always taking precedence), that checks your commit messages, ensuring they comply to conventional commit format.
* *`cc`* - generates conventional commit messages, based on your recent changes (and copies it to your clipboard)
* *`ccc`* - same as cc (above), except it ALSO does the git commit FOR YOU (3 c's for "Create Conventional Commit")

== More features?

* easily choose between system-wide or project-specific hook installation
* set a JIRA ticket to be auto included in your 10 next git commits, (although maybe you need to do a squash after all that)

🎯 Real-World Scenarios Handled:
- Team Developer: Can remove only local hooks while keeping global ones
- Personal Setup: Can remove global hooks without affecting project-specific local ones

== 🔄 Using fcgh + cc Together or separately - your choice!

The tools are *independent*, you can happily use one without the other but work beautifully together:

=== Usage Option 1: You write your own commit messages (with fcgh installed)

[source,bash]
----
# Write your commit message manually
git commit -m "feat: CGC-1245 Added login authentication"
# ↳ The fcgh hook, automatically validates your message ✅
----

=== Usage Option 2: Commit messages are generated based on your changes (using ccc) and validated by the fcgh hook

[source,bash]
----
# Let ccc generate the conventional commit message and execute it for you
ccc
# ↳ ccc analyzes your git changes and creates a conventional commit message
# ↳ When git commit runs, the git hook (installed via fcgh) validates the generated message ✅
# ↳ which means: compliant conventional commits, every time! 
----

=== You can mix and match:

* *Use fcgh alone*: Ensures your manual commit messages get automatic validation
* *Use cc alone*: Preview generated conventional commit messages, based on your changes, but you (ctrl-v), paste the whole git command, and you do the commit yourself
* *Use ccc alone*: Generated conventional commits, without validation (but why would you?)
* *Use ALL THREE*: 1) See generated conventional commits via cc, then 2) trigger the commit via cc PLUS 3) BONUS validation via fcgh being setup - talk about the perfect combo! 🎯

*Pro tip*: Start with `cc` to see what good commit messages look like, then graduate to writing your own!

== ⚡️ Installation

=== Prerequisites

* Git (obviously!)

=== Step 1: Install the tools

*Option A: Download the Binary* 

Each release includes *both tools* - you get everything in one download!

** Windows:*
1. Go to https://github.com/greenstevester/fast-cc-git-hooks/releases/latest[Releases Page]
2. Download `fast-cc-git-hooks_windows_amd64.zip`
3. Extract and add `fcgh.exe` and `cc.exe` and `ccc.exe` to a directory on your PATH
NOTE: In windows, there's an option to change only your environment variables - use that as a preference, or just add both to your PATH.

*🐧 Linux:*

[source,bash]
----
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
----

*🍎 macOS:*

*Option A: Easy Installation Script (Recommended)*

[source,bash]
----
# One-line install - automatically detects your Mac type and handles PATH
curl -sSL https://raw.githubusercontent.com/greenstevester/fast-cc-git-hooks/refs/heads/main/install-macos-source.sh | bash
----

*Option B: Manual Installation*

[source,bash]
----
# Intel Macs
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_amd64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc ccc
sudo mv fcgh cc ccc /usr/local/bin/

# Apple Silicon (M1/M2/M3) - most common now
curl -L -o fcgh.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fcgh_darwin_arm64.tar.gz
tar -xzf fcgh.tar.gz
chmod +x fcgh cc ccc
sudo mv fcgh cc ccc /usr/local/bin/
----

*📱 For macOS Security (Both Options):*
If you see "fcgh cannot be opened" security warning:
1. Click "Done" in the dialog
2. Go to *System Preferences* → *Security &amp; Privacy* → *General* tab
3. Click "Allow Anyway" next to the fcgh message
4. Try running `fcgh` again and click "Open" when prompted

_Note: The installation script automatically handles the quarantine removal to minimize security warnings._

*Option C: Build from Source*

**⚠️ Requires Go 1.25 or later** - The project uses modern Go features and will not compile with older versions.

[source,bash]
----
# Clone the repository
git clone https://github.com/greenstevester/fast-cc-git-hooks.git
cd fast-cc-git-hooks

# Build all tools (requires Go 1.25+)
make build-all-tools

# Install to GOPATH/bin
make install-all
----

=== 🛠️ Step 2: Setup

==== Global Installation (Default - Recommended)

[source,bash]
----
fcgh setup      # For standard projects
fcgh setup-ent  # For enterprise projects with JIRA validation
cc
----

==== Local Installation (Current Repository Only)

[source,bash]
----
fcgh setup --local      # Install only for current repository
fcgh setup-ent --local  # Enterprise setup for current repository only
----

==== 🏆 Configuration Precedence (Important!)

When you have both local and global installations:

----
Local Repository Hook  ➤  ALWAYS WINS  ⚡
Global Git Hook        ➤  Ignored when local exists
----

=== Commands (fcgh)

[source,bash]
----
fcgh setup-ent # Set up everything (This creates a file called `fast-cc-config.yaml` in your home directory (`~/.fast-cc/`) that you can edit. )
fcgh remove    # Smart removal - prompts to choose local/global if both exist
fcgh remove --local   # Remove only from current repository
fcgh remove --global  # Remove only from global git configuration
----

*Testing fcgh:*

[source,bash]
----
fcgh validate "freak: my terrible message"  # Test if a message is bad
----

=== Commit Helper Commands (cc)

*Smart commit generation with semantic analysis:*

[source,bash]
----
cc                     # Generate commit message and copy to clipboard automatically
cc --no-copy           # Generate commit message without copying to clipboard
ccc                    # Generate commit message and commit
cc --verbose           # Show detailed analysis of your changes
cc --help              # Show all available options
----

*🧠 Semantic Analysis for Infrastructure Code:*
The `cc` command includes semantic analysis for Terraform files with Oracle OCI awareness. 
It understands the infrastructure changes and generates contextual commit messages. link:docs/semantic-analysis-examples.md[See examples →]

*📋 Clipboard Integration:*
The `cc` command automatically copies the generated git commit command to your clipboard by default. 
Perfect for quick copy-paste workflows - just run `cc` and press Ctrl+V in your terminal! Use `--no-copy` to disable this behavior.

*That's it!* Most people only ever need `fcgh setup-ent` and `ccc`.


== 🤔 Common Questions

*Q: Do I need all tools?*
A: No! They're completely independent. Install hooks for validation, ccc for generation, or both together.

*Q: What if I want to turn off validation temporarily?*
A: Just run `fcgh remove` and it's gone! Run `fcgh setup-ent` to turn it back on.

*Q: Will this mess up my code?*
A: Nope! The hooks only check commit messages, cc only generates them. Your code stays exactly the same.

*Q: What if cc generates a bad commit message?*
A: The hook will catch it! That's why they work so well together.

*Q: Can I use this on all my projects?*
A: Yes! When you run `fcgh setup-ent`, it works for ALL your Git projects. cc works in any git repo.

*Q: What if I don't like the cc generated message?*
A: Just use `cc` (without ccc) to preview first, then write your own and let the hook validate it!

== 🎯 Examples of Good vs Bad Commits

❌ *Bad commits* (these will be rejected):

[source,bash]
----
git commit -m "fix"
git commit -m "updated stuff"  
git commit -m "asdf"
git commit -m "WIP"
----

✅ *Good Enterprise commits* (these will work):

[source,bash]
----
git commit -m "feat: CGC-5641 Added user login page"
git commit -m "fix: CAHC-4132 Resolve password validation bug"
git commit -m "docs: LIND-4777 installation instructions"
----

== ⚙️ Want to Customize? (Optional)

*Most people don't need to do this!* The tool works great out of the box.

But if you want to customize the rules, run:

[source,bash]
----
fcgh init
----

This creates a file called `fast-cc-config.yaml` in your home directory (`~/.fast-cc/`) that you can edit. 

== 🚨 Troubleshooting

*Problem: The tool says my commit message is bad, but I think it's fine!*

Run this to test your message:

[source,bash]
----
fcgh validate "your message here"
----

It will tell you exactly what's wrong and how to fix it.

*Problem: I want to turn it off for just one commit*

You can't bypass it easily (that's the point!), but you can run:

[source,bash]
----
fcgh remove
git commit -m "your message"  
fcgh setup
----

*Problem: It's not working at all*

Try setting it up again:

[source,bash]
----
fcgh remove
fcgh setup-ent
----

== 🏗️ Advanced Examples

*Want to include ticket numbers?* (like JIRA tickets)

[source,bash]
----
git commit -m "feat: CGC-1234 add user login"
git commit -m "fix: PROJ-456 resolve password bug"
----

*Want to use scopes?* (optional grouping)

[source,bash]
----
git commit -m "feat(auth): add login form"
git commit -m "fix(api): resolve timeout issue"
git commit -m "docs(readme): update setup instructions"  
----

The part in `()` is called a "scope" - it's like a category for your change.

== 👨‍💻 For Developers

*Want to help make this tool better?*

[source,bash]
----
git clone https://github.com/greenstevester/fast-cc-git-hooks.git
cd fast-cc-git-hooks
make build
make test
----

All commits to this project must follow conventional format too! 😄

== 🔄 Changelog Generation Tools

Once you're using conventional commits, you can automatically generate changelogs and manage versioning! Here are the most popular tools:

=== 🚀 *Fully Automated Release Tools* (Recommended)

* *https://github.com/semantic-release/semantic-release[semantic-release]* - The gold standard for automated releases
* Fully automates: version bumping, changelog generation, git tagging, and publishing
* Works with GitHub, GitLab, npm, and more
* Perfect for CI/CD pipelines
* *https://github.com/absolute-version/commit-and-tag-version[commit-and-tag-version]* - Drop-in replacement for `npm version`

* Handles version bumping, tagging, and CHANGELOG generation
* Great for manual releases with automation

=== 📊 *Changelog-Focused Tools*

* *https://github.com/conventional-changelog/conventional-changelog[conventional-changelog]* - The original changelog generator
* Generate changelog from git metadata
* Multiple presets (Angular, Atom, etc.)
* Highly customizable
* *https://git-cliff.org/[git-cliff]* - Modern Rust-based changelog generator

* Highly customizable templates
* Fast and reliable
* Great for complex projects

=== 🏢 *Enterprise &amp; Monorepo Tools*

* *https://github.com/oknozor/cocogitto[cocogitto]* - Complete conventional commits toolkit
* Version bumping, changelog generation, and commit linting
* Great for complex workflows
* *https://github.com/chaaz/versio[versio]* - Monorepo-compatible versioning

* Handles project dependencies
* Generates tags and changelogs

=== ⚙️ *Language-Specific Tools*

* *Go*: https://github.com/goreleaser/chglog[chglog] - Template-based changelog generation
* *Python*: https://github.com/relekang/python-semantic-release[python-semantic-release]
* *PHP*: https://github.com/marcocesarato/php-conventional-changelog[php-conventional-changelog]
* *Java*: https://github.com/tomasbjerre/git-changelog-lib[git-changelog-lib]

=== 🎯 *Quick Start Recommendations*

. *For automated CI/CD*: Use `semantic-release`
. *For manual releases*: Use `commit-and-tag-version`
. *For just changelogs*: Use `conventional-changelog` or `git-cliff`
. *For monorepos*: Use `versio` or `cocogitto`

All these tools work perfectly with the conventional commit messages that `fcgh` enforces! 🎉

== 📚 Documentation

* link:docs/semantic-analysis-examples.md[Semantic Analysis Examples] - See how the intelligent commit message generation works with real-world infrastructure code examples
* link:CLAUDE.md[CLAUDE.md] - Project development guide for Claude Code

== 📝 License

MIT License - do whatever you want with this code!

'''

*That's everything!* Remember: just run `fcgh setup-ent` and start writing better commits! 🚀