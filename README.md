# 🚀 fast-cc-hooks

**The fastest (and laziest) way for YOU, to write [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) messages for ALL your git commits!** 

Why Use Conventional Commits? Because when you use them, there's tooling (listed later in this page) that can:
- auto-generate CHANGELOG's AND even determine semantic version bumps 
- make it dead simple for YOU and the community of coders to contribute to projects *without lazy commit messages*.
 
## Value Prop: This **performance focused** "git companion" toolkit, provides two (2) complementary tools that do the heavy lifting:
- **`cc`** - preview generated conventional commit messages, based on your changes
- **`ccc`** - generate conventional commit messages AND commit (3 c's for "Create Conventional Commit") (most popular)
- **`fast-cc-hooks`** - installer of a git-hook, that checks your commit messages, ensuring they comply to conventional commit format.

Use them together for the perfect commit workflow! 🎯 

## 🔄 Perfect Workflow: Using ALL Tools Together

The tools are **completely independent** but work beautifully together:

### Option 1: You write your own commit messages (with fast-cc-hooks installed)
```bash
# Write your commit message manually
git commit -m "feat: CGC-1245 Added login authentication"
# ↳ The fast-cc-hooks hook, automatically validates your message ✅
```

### Option 2: Commit messages are generated based on your changes (using ccc) and validated by the fast-cc-hooks hook
```bash
# Let ccc generate the conventional commit message and execute it for you
ccc
# ↳ ccc analyzes your git changes and creates a conventional commit message
# ↳ When git commit runs, the git hook (installed via fast-cc-hooks) validates the generated message ✅
# ↳ which means: compliant conventional commits, every time! 
```

### You can mix and match:
- **Use fast-cc-hooks alone**: Ensures your manual commit messages get automatic validation
- **Use cc alone**: Preview generated conventional commit messages, based on your changes, but you do the commit yourself
- **Use ccc alone**: Generated compliant conventional commits, without validation (but why would you?)
- **Use both**: Generated conventional commits PLUS BONUS validation, talk about the perfect combo! 🎯

**Pro tip**: Start with `ccc` to see what good commit messages look like, then graduate to writing your own!

## ⚡️ Quick Setup (In 3 steps!)

### Step 1: Install the tools

**Option A: Download the Binary** 

Each release includes **both tools** - you get everything in one download!

*** Windows:**
1. Go to [Releases Page](https://github.com/greenstevester/fast-cc-git-hooks/releases/latest)
2. Download `fast-cc-hooks_windows_amd64.zip`
3. Extract and add both `fast-cc-hooks.exe` and `cc.exe` to a directory on your PATH
NOTE: In windows, there's an option to change only your environment variables - use that as a preference, or just add both to your PATH.

**🐧 Linux:**
```bash
# Most common (AMD64)
curl -L -o fast-cc-hooks.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_linux_amd64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks cc
sudo mv fast-cc-hooks cc /usr/local/bin/

# ARM64 (Raspberry Pi, etc.)
curl -L -o fast-cc-hooks.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_linux_arm64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks cc
sudo mv fast-cc-hooks cc /usr/local/bin/
```

**🍎 macOS:**
```bash
# Intel Macs
curl -L -o fast-cc-hooks.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_darwin_amd64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks cc
sudo mv fast-cc-hooks cc /usr/local/bin/

# Apple Silicon (M1/M2/M3) - most common now
curl -L -o fast-cc-hooks.tar.gz https://github.com/greenstevester/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_darwin_arm64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks cc
sudo mv fast-cc-hooks cc /usr/local/bin/
```

**Option C: Using Homebrew** (macOS/Linux - coming soon!)
```bash
# Will be available after first release
brew install greenstevester/tap/fast-cc-hooks
```

### Step 2: Verify all tools work

```bash
fast-cc-hooks version  # shows version info
cc --verbose           # runs cc in verbose mode
ccc --verbose           # runs cc in verbose mode
```

If you see version/help info for both, you're good to go! If not, make sure `/usr/local/bin` is in your PATH.

### Step 3: Set it up (takes 5 seconds!)

```bash
fast-cc-hooks setup OR fast-cc-hooks setup-ent
```

**That's it!** 🎉 Now every time you make a commit, it will automatically check that your message is compliant conventional commit message.


## ✨ How to write good commit messages

Instead of writing messy commits like:
```bash
git commit -m "fixed stuff"
git commit -m "update"
git commit -m "asdfasdf"
```

Write clear commits like:
```bash
git commit -m "feat: add login button"
git commit -m "fix: resolve login bug"  
git commit -m "docs: update README"
```

The format is simple: `type: what you did`

**Common types:**
- `feat` - when you add something new
- `fix` - when you fix a bug
- `docs` - when you update documentation
- `test` - when you add tests
- `chore` - when you do maintenance stuff

## 🛠️ Simple Commands

### Git Hook Commands (fast-cc-hooks)

**The only commands you need:**

```bash
fast-cc-hooks setup     # Set up everything (start here!)
fast-cc-hooks remove    # Remove everything if you want to stop using it
```

**Test things:**
```bash
fast-cc-hooks validate "feat: my commit message"  # Test if a message is good
```

### Commit Helper Commands (cc)

**Smart commit generation:**

```bash
cc                      # Preview generated commit message
ccc                    # Generate perfect commit message and commit
cc --verbose           # Show detailed analysis of your changes
cc --help             # Show all available options
```

**That's it!** Most people only ever need `fast-cc-hooks setup` and `ccc`.

## 🤔 Common Questions

**Q: Do I need both tools?**
A: No! They're completely independent. Use hooks for validation, cc for generation, or both together.

**Q: What if I want to turn off validation temporarily?**
A: Just run `fast-cc-hooks remove` and it's gone! Run `fast-cc-hooks setup` to turn it back on.

**Q: Will this mess up my code?**
A: Nope! The hooks only check commit messages, cc only generates them. Your code stays exactly the same.

**Q: What if cc generates a bad commit message?**
A: The hook will catch it! That's why they work so well together. Plus, cc learns from conventional commit standards.

**Q: Can I use this on all my projects?**
A: Yes! When you run `fast-cc-hooks setup`, it works for ALL your Git projects. cc works in any git repo.

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

✅ **Good commits** (these will work):
```bash
git commit -m "feat: add user login page"
git commit -m "fix: resolve password validation bug"
git commit -m "docs: add installation instructions"
git commit -m "test: add login form tests"
```

## ⚙️ Want to Customize? (Optional)

**Most people don't need to do this!** The tool works great out of the box.

But if you want to customize the rules, run:
```bash
fast-cc-hooks init
```

This creates a file called `default-config.yaml` in your home directory (`~/.fast-cc-git-hooks/`) that you can edit. It has comments explaining everything!

## 🚨 Troubleshooting

**Problem: The tool says my commit message is bad, but I think it's fine!**

Run this to test your message:
```bash
fast-cc-hooks validate "your message here"
```

It will tell you exactly what's wrong and how to fix it.

**Problem: I want to turn it off for just one commit**

You can't bypass it easily (that's the point!), but you can run:
```bash
fast-cc-hooks remove
git commit -m "your message"  
fast-cc-hooks setup
```

**Problem: It's not working at all**

Try setting it up again:
```bash
fast-cc-hooks remove
fast-cc-hooks setup
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

## 📝 License

MIT License - do whatever you want with this code!

---

**That's everything!** Remember: just run `fast-cc-hooks setup` and start writing better commits! 🚀