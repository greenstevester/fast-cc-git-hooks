# 🚀 fast-cc-hooks

**The easiest way to write better commit messages!** 

This tool automatically checks your Git commit messages to make sure they're clear and helpful. It's so easy, anyone can use it! 

## ⚡️ Super Quick Setup (3 steps!)

### Step 1: Install the tool

**Option A: Download Binary** (easiest! 🎯)

**🐧 Linux:**
```bash
# Most common (AMD64)
curl -L -o fast-cc-hooks.tar.gz https://github.com/stevengreensill/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_linux_amd64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks
sudo mv fast-cc-hooks /usr/local/bin/

# ARM64 (Raspberry Pi, etc.)
curl -L -o fast-cc-hooks.tar.gz https://github.com/stevengreensill/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_linux_arm64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks
sudo mv fast-cc-hooks /usr/local/bin/
```

**🍎 macOS:**
```bash
# Intel Macs
curl -L -o fast-cc-hooks.tar.gz https://github.com/stevengreensill/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_darwin_amd64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks
sudo mv fast-cc-hooks /usr/local/bin/

# Apple Silicon (M1/M2/M3) - most common now
curl -L -o fast-cc-hooks.tar.gz https://github.com/stevengreensill/fast-cc-git-hooks/releases/latest/download/fast-cc-hooks_darwin_arm64.tar.gz
tar -xzf fast-cc-hooks.tar.gz
chmod +x fast-cc-hooks
sudo mv fast-cc-hooks /usr/local/bin/
```

**🪟 Windows:**
1. Go to [Releases Page](https://github.com/stevengreensill/fast-cc-git-hooks/releases/latest)
2. Download `fast-cc-hooks_windows_amd64.zip`
3. Extract and add `fast-cc-hooks.exe` to your PATH

**Option B: Using Go**
```bash
go install github.com/stevengreensill/fast-cc-git-hooks/cmd/fast-cc-hooks@latest
```

**Option C: Using Homebrew** (macOS/Linux - coming soon!)
```bash
# Will be available after first release
brew install greenstevester/tap/fast-cc-hooks
```

### Step 2: Verify it works

```bash
fast-cc-hooks --version
```

If you see version info, you're good to go! If not, make sure `/usr/local/bin` is in your PATH.

### Step 3: Set it up (takes 5 seconds!)

```bash
fast-cc-hooks setup
```

**That's it!** 🎉 Now every time you make a commit, it will automatically check that your message is good!

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

**The only commands you need:**

```bash
fast-cc-hooks setup     # Set up everything (start here!)
fast-cc-hooks remove    # Remove everything if you want to stop using it
```

**Test things:**
```bash
fast-cc-hooks validate "feat: my commit message"  # Test if a message is good
```

**That's it!** Most people only ever need `setup` and `remove`.

## 🤔 Common Questions

**Q: What if I want to turn it off temporarily?**
A: Just run `fast-cc-hooks remove` and it's gone! Run `fast-cc-hooks setup` to turn it back on.

**Q: Will this mess up my code?**
A: Nope! It only checks your commit messages. Your code stays exactly the same.

**Q: What if I make a mistake in my commit message?**
A: The tool will tell you what's wrong and you can try again. It's very friendly!

**Q: Can I use this on all my projects?**
A: Yes! When you run `setup`, it works for ALL your Git projects automatically.

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

This creates a file called `.fast-cc-hooks.yaml` that you can edit. It has comments explaining everything!

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
git clone https://github.com/stevengreensill/fast-cc-git-hooks.git
cd fast-cc-git-hooks
make build
make test
```

All commits to this project must follow conventional format too! 😄

## 📝 License

MIT License - do whatever you want with this code!

---

**That's everything!** Remember: just run `fast-cc-hooks setup` and start writing better commits! 🚀