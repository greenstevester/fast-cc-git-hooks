# ⏱️ 10 Second History

**The lightning-fast story of why these tools exist.**

## 🎯 The Problem

**Bad commit messages everywhere.**

```bash
git log --oneline
a1b2c3d WIP
e4f5g6h fix stuff  
h7i8j9k updated files
k0l1m2n asdf
```

Sound familiar? 😅

## 💡 The Inspiration

Inspired by **_Boo_** ...after seeing lots and lots of git commit messages, all having something to the tune of "feat: CGC-0000 bla bla bla...".

```bash
git log --oneline
a1b2c3d CGC-0000 bla bla this...
e4f5g6h CGC-0000 bla bla that...
h7i8j9k CGC-0000 fixed...
k0l1m2n CGC-0000 updated...
```

## 🚀 The Solution

**"What if everyone could write commits with almost little or no effort?"**

**Three simple tools:**
1. **`fcgh`** - Validates your commit messages (no more "WIP" commits!)
2. **`ccg`** - Generates perfect messages based on your changes
3. **`ccdo`** - Does everything automatically (add + generate + commit)

## 🎪 The Magic

**From amateur hour to professional in 30 seconds:**

```bash
# Before (amateur hour)
git add .
git commit -m "fix stuff"

# After (professional)
git add .
ccdo  # → "fix(auth): resolve JWT token expiration handling"
```

## 🌍 The Vision

**Free the world from crappy commit messages.**

Because every developer deserves:
- 📈 **Auto-generated changelogs** from commit history
- 🏷️ **Semantic versioning** based on commit types
- 🔍 **Searchable git history** that actually makes sense
- 👥 **Team collaboration** without confusion
- 🤖 **CI/CD integration** that works seamlessly

It also depends really, if you care about:
- clean, consistent, and easy-to-understand commit messages (sure)
- auto-generate CHANGELOG's (maybe)
- infer semantic version bumps when required, just from your commit messages. (yes, take that boilerplate out of your life)

It all stems from something simple, but usually this means more discipline AND (worst of all) strict adherence to a specific workflow. Most people automate this via their bash shell, but bash only goes so far - ergo this tool.

And on that note, strict adherence to a specific workflow is enough to scare away most coders, who just want to build without the added "admin" overhead / accounting / audits while they're in flow.

So, this small set of performance-obsessed tools is for YOU (and the wider community of coders), to free the world of *crappy commit messages* AND to prevent further pollution of potentially GREAT git repos.

## 🎉 The Result

**Perfect commits. Every time. Without thinking about it.**

---

*Now go make Boo proud with your beautiful commit messages! 🚀*