# 🔄 Changelog Generation Tools

Once you're using conventional commits with **fcgh**, you can automatically generate changelogs and manage versioning! Here are the most popular tools that work perfectly with conventional commit messages.

## 🚀 Fully Automated Release Tools (Recommended)

### **[semantic-release](https://github.com/semantic-release/semantic-release)**
**The gold standard for automated releases**
- Fully automates: version bumping, changelog generation, git tagging, and publishing
- Works with GitHub, GitLab, npm, and more
- Perfect for CI/CD pipelines
- Zero manual intervention needed

**Quick Setup:**
```bash
npm install --save-dev semantic-release
# Add to package.json scripts:
# "release": "semantic-release"
```

### **[commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version)**
**Drop-in replacement for `npm version`**
- Handles version bumping, tagging, and CHANGELOG generation
- Great for manual releases with automation
- More control over release timing than semantic-release

**Quick Setup:**
```bash
npm install --save-dev commit-and-tag-version
# Run: npx commit-and-tag-version
```

## 📊 Changelog-Focused Tools

### **[conventional-changelog](https://github.com/conventional-changelog/conventional-changelog)**
**The original changelog generator**
- Generate changelog from git metadata
- Multiple presets (Angular, Atom, etc.)
- Highly customizable templates
- Just generates changelogs (no version bumping)

**Quick Setup:**
```bash
npm install --save-dev conventional-changelog-cli
npx conventional-changelog -p angular -i CHANGELOG.md -s
```

### **[git-cliff](https://git-cliff.org/)**
**Modern Rust-based changelog generator**
- Highly customizable templates
- Fast and reliable performance
- Great for complex projects
- Template-based configuration

**Quick Setup:**
```bash
# Install via cargo or download binary
cargo install git-cliff
git cliff -o CHANGELOG.md
```

## 🏢 Enterprise & Monorepo Tools

### **[cocogitto](https://github.com/oknozor/cocogitto)**
**Complete conventional commits toolkit**
- Version bumping, changelog generation, and commit linting
- Great for complex workflows
- Built-in conventional commit validation

**Quick Setup:**
```bash
cargo install cocogitto
cog init  # Initialize configuration
cog changelog
```

### **[versio](https://github.com/chaaz/versio)**
**Monorepo-compatible versioning**
- Handles project dependencies
- Generates tags and changelogs for multiple packages
- Perfect for complex monorepos

## ⚙️ Language-Specific Tools

### **Go**
- **[chglog](https://github.com/goreleaser/chglog)** - Template-based changelog generation

### **Python**
- **[python-semantic-release](https://github.com/relekang/python-semantic-release)** - Python version of semantic-release

### **PHP**
- **[php-conventional-changelog](https://github.com/marcocesarato/php-conventional-changelog)** - PHP changelog generator

### **Java**
- **[git-changelog-lib](https://github.com/tomasbjerre/git-changelog-lib)** - Java library for changelog generation

## 🎯 Quick Start Recommendations

| Use Case | Recommended Tool | Why |
|----------|------------------|-----|
| **Automated CI/CD releases** | `semantic-release` | Zero manual intervention, works with all major platforms |
| **Manual releases with automation** | `commit-and-tag-version` | More control, less complexity than semantic-release |
| **Just need changelogs** | `conventional-changelog` or `git-cliff` | Focused tools, no version management |
| **Monorepos** | `versio` or `cocogitto` | Handle multiple packages and dependencies |
| **Go projects** | `chglog` + `goreleaser` | Integrates well with Go toolchain |
| **Python projects** | `python-semantic-release` | Native Python tooling |

## 🔗 Integration Examples

### GitHub Actions + semantic-release
```yaml
name: Release
on:
  push:
    branches: [main]
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: npm ci
      - run: npm run test
      - run: npx semantic-release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Manual Release Workflow
```bash
# 1. Make conventional commits using fcgh
git add .
ccdo  # Generates conventional commit

# 2. When ready to release
npx commit-and-tag-version

# 3. Push release
git push --follow-tags origin main
```

## 💡 Pro Tips

1. **Start Simple**: Begin with `conventional-changelog` to see your commit history in action
2. **Automate Gradually**: Move to `semantic-release` once you're comfortable with conventional commits
3. **Test First**: Use `--dry-run` flags to preview what tools will generate
4. **Customize Templates**: Most tools support custom templates to match your project style
5. **Combine Tools**: Use `fcgh` for commit validation + your chosen changelog tool

## 📚 Further Reading

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

---

**Remember**: All these tools work perfectly with the conventional commit messages that `fcgh` enforces! The better your commit messages, the better your generated changelogs will be. 🎉