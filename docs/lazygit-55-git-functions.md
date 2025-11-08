# Lazygit's 55 Git Command Files - Complete Breakdown

Complete reference of all files in `pkg/commands/git_commands/` from the lazygit repository.

---

## Core Infrastructure (4 files)

### 1. **common.go**
Common utilities and shared functionality across git commands.

### 2. **git_command_builder.go** ⭐ HIGH PRIORITY
**Key Exports:**
- Fluent builder pattern for constructing git commands
- `NewGitCmd()` - Creates new git command builder
- `Arg()`, `ArgIf()`, `ArgIfElse()` - Add arguments conditionally
- `Config()`, `ConfigIf()` - Add git config overrides
- `Dir()`, `DirIf()` - Set working directory
- `ToArgv()`, `ToString()` - Convert to executable format

**Value for fast-cc-git-hooks:** Replace all `exec.Command("git", ...)` calls

### 3. **git_command_builder_test.go**
Test coverage for command builder (reference for writing our tests).

### 4. **deps_test.go**
Dependency injection tests.

---

## Branch Operations (4 files)

### 5. **branch.go** ⭐ MEDIUM PRIORITY
**Key Exports (30+ functions):**

**Branch Creation:**
- `New()` - Create and checkout new branch
- `NewWithoutTracking()` - Create without upstream
- `NewWithoutCheckout()` - Create without checkout
- `CreateWithUpstream()` - Create with specified upstream

**Branch Info:**
- `CurrentBranchInfo()` - Get current branch name and detached HEAD status
- `CurrentBranchName()` - Get current branch name
- `PreviousRef()` - Get previously checked out branch
- `IsHeadDetached()` - Check if HEAD is detached

**Branch Management:**
- `LocalDelete()` - Delete local branches
- `Checkout()` - Checkout branch or commit
- `Rename()` - Rename branch

**Upstream Operations:**
- `SetCurrentBranchUpstream()` - Set upstream for current branch
- `SetUpstream()` - Set upstream for specified branch
- `UnsetUpstream()` - Remove upstream

**Differences:**
- `GetCurrentBranchUpstreamDifferenceCount()` - Count pushable/pullable commits
- `GetUpstreamDifferenceCount()` - Count differences for branch
- `GetCommitDifferences()` - Calculate push/pull counts

**Merge:**
- `Merge()` - Merge branches with strategies
- `CanDoFastForwardMerge()` - Check if FF merge possible
- `IsBranchMerged()` - Check if branch merged

**Graph:**
- `GetGraph()` - Get formatted log graph
- `GetGraphCmdObj()` - Create command for graph display
- `AllBranchesLogCmdObj()` - Log for all branches

**Value for fast-cc-git-hooks:** Branch detection for scope analysis, upstream info for status

### 6. **branch_loader.go**
Loads and parses branch information from git into structured data.

### 7. **branch_loader_test.go**
Tests for branch loading.

### 8. **branch_test.go**
Tests for branch operations.

---

## Commit Operations (7 files)

### 9. **commit.go** ⭐ HIGH PRIORITY
**Key Exports (25+ functions):**

**Author Management:**
- `ResetAuthor()` - Remove author info
- `SetAuthor()` - Set commit author
- `AddCoAuthor()` - Add co-author metadata
- `GetCommitAuthor()` - Retrieve author name/email

**Message Operations:**
- `GetCommitMessage()` - Get full commit message
- `GetCommitSubject()` - Get subject line only
- `GetCommitMessageFromHistory()` - Get historical messages
- `RewordLastCommit()` - Update latest commit message

**Commit Retrieval:**
- `GetHashesAndCommitMessagesFirstLine()` - Get hashes with subjects
- `GetCommitsOneline()` - Format commits in one line

**Commit Modifications:**
- `AmendHead()` - Amend HEAD commit
- `CreateFixupCommit()` - Create fixup commit
- `CreateAmendCommit()` - Create amend commit

**Utilities:**
- `ShowCmdObj()` - Get diff for commit
- `ShowFileContentCmdObj()` - Get file content at commit
- `ResetToCommit()` - Reset to specific commit

**Value for fast-cc-git-hooks:** Parse commit history for patterns, extract commit messages

### 10. **commit_file_loader.go**
Loads files changed in commits.

### 11. **commit_file_loader_test.go**
Tests for commit file loading.

### 12. **commit_loader.go** ⭐ MEDIUM PRIORITY
Loads commit history with parsing and filtering.

### 13. **commit_loader_test.go**
Tests for commit loading.

### 14. **commit_loading_shared.go**
Shared utilities for commit loading operations.

### 15. **commit_test.go**
Tests for commit operations.

---

## Configuration (1 file)

### 16. **config.go** ⭐ MEDIUM PRIORITY
**Key Exports:**

**GPG Configuration:**
- `NeedsGpgSubprocess()` - Check if GPG subprocess needed
- `GetGpgTagSign()` - Get GPG tag signing preference

**Editor & Remote:**
- `GetCoreEditor()` - Get configured editor
- `GetRemoteURL()` - Get remote repository URL

**Repository Behavior:**
- `GetShowUntrackedFiles()` - Get untracked files setting
- `GetPushToCurrent()` - Check push-to-current setting
- `GetRebaseUpdateRefs()` - Get rebase update refs preference
- `GetMergeFF()` - Get merge fast-forward config

**Advanced:**
- `Branches()` - Get branches from git config
- `GetGitFlowPrefixes()` - Get gitflow prefixes
- `GetCoreCommentChar()` - Get comment character (default '#')
- `DropConfigCache()` - Clear config cache

**Value for fast-cc-git-hooks:** Detect user preferences, repository settings

---

## Diff Operations (1 file)

### 17. **diff.go** ⭐ HIGH PRIORITY
Diff generation with various formats, word diff, numstat parsing.

**Value for fast-cc-git-hooks:** Replace all manual diff parsing in `advanced_git_analyzer.go`

---

## File Operations (4 files)

### 18. **file.go** ⭐ MEDIUM PRIORITY
File-level operations and metadata.

### 19. **file_loader.go** ⭐ MEDIUM PRIORITY
Loads file status, staged/unstaged detection, change type parsing.

**Value for fast-cc-git-hooks:** Better file analysis for commit generation

### 20. **file_loader_test.go**
Tests for file loading.

### 21. **file_test.go**
Tests for file operations.

---

## Rebase Operations (3 files)

### 22. **rebase.go** 🔵 LOW PRIORITY (Advanced Feature)
**Key Exports (40+ functions):**

**Reword Operations:**
- `RewordCommit()` - Modify commit message via rebase
- `RewordCommitInEditor()` - Reword using editor

**Author Modifications:**
- `ResetCommitAuthor()` - Remove author info
- `SetCommitAuthor()` - Set new author
- `AddCommitCoAuthor()` - Add co-author

**Commit Movement:**
- `MoveCommitsDown()` - Move commits to later positions
- `MoveCommitsUp()` - Move commits to earlier positions

**Interactive Rebase:**
- `InteractiveRebase()` - Apply todo action to commits
- `EditRebase()` - Insert breakpoint in rebase
- `EditRebaseFromBaseCommit()` - Rebase from base commit
- `PrepareInteractiveRebaseCommand()` - Build rebase command
- `GitRebaseEditTodo()` - Execute rebase --edit-todo

**Fixup and Squash:**
- `AmendTo()` - Create and move fixup commit
- `MoveFixupCommitDown()` - Relocate fixup commit
- `SquashAllAboveFixupCommits()` - Auto-squash fixups

**Todo Management:**
- `EditRebaseTodo()` - Modify rebase todo actions
- `DeleteUpdateRefTodos()` - Remove todos
- `MoveTodosDown()` / `MoveTodosUp()` - Reorder todos

**Rebase Control:**
- `BeginInteractiveRebaseForCommit()` - Start rebase for single commit
- `BeginInteractiveRebaseForCommitRange()` - Start rebase for range
- `RebaseBranch()` - Rebase onto branch
- `ContinueRebase()` - Continue after conflict
- `AbortRebase()` - Cancel rebase

**Cherry-pick:**
- `CherryPickCommits()` - Apply commits as cherry-picks
- `DropMergeCommit()` - Remove merge commit

**Value for fast-cc-git-hooks:** Not needed for current use case (hook validation)

### 23. **rebase_test.go**
Tests for rebase operations.

### 24. **deps_test.go**
Already listed above.

---

## Reflog Operations (2 files)

### 25. **reflog_commit_loader.go** 🔵 LOW PRIORITY
Loads commits from git reflog for undo/redo functionality.

### 26. **reflog_commit_loader_test.go**
Tests for reflog loading.

---

## Remote Operations (2 files)

### 27. **remote.go** 🔵 LOW PRIORITY
Remote repository management (add, remove, rename remotes).

### 28. **remote_loader.go**
Loads and parses remote repository information.

---

## Repository Path Management (2 files)

### 29. **repo_paths.go** ⭐ MEDIUM PRIORITY
Manages repository paths (.git directory, worktrees, etc.).

**Value for fast-cc-git-hooks:** Proper .git directory detection for hooks

### 30. **repo_paths_test.go**
Tests for repository path utilities.

---

## Stash Operations (4 files)

### 31. **stash.go** 🔵 LOW PRIORITY
Stash creation, application, deletion.

### 32. **stash_loader.go**
Loads stash list.

### 33. **stash_loader_test.go**
Tests for stash loading.

### 34. **stash_test.go**
Tests for stash operations.

---

## Status Operations (1 file)

### 35. **status.go** ⭐ HIGH PRIORITY
Repository status checking (clean, dirty, conflicts, etc.).

**Value for fast-cc-git-hooks:** Pre-commit validation

---

## Sync Operations (3 files)

### 36. **sync.go** 🔵 LOW PRIORITY
**Key Exports:**

**Push Operations:**
- `PushCmdObj()` - Build push command with options
- `Push()` - Execute push with force/force-with-lease support

**Fetch Operations:**
- `FetchCmdObj()` - Create fetch command
- `Fetch()` - Interactive fetch with credentials
- `FetchBackgroundCmdObj()` - Background fetch command
- `FetchBackground()` - Non-interactive background fetch
- `FetchRemote()` - Fetch specific remote

**Pull Operations:**
- `Pull()` - Pull with fast-forward-only support
- `FastForward()` - Update branch to match remote

**Value for fast-cc-git-hooks:** Not needed (no push/pull in hooks)

### 37. **sync_test.go**
Tests for sync operations.

### 38. (see above - sync_test.go included here)

---

## Tag Operations (3 files)

### 39. **tag.go** 🔵 LOW PRIORITY
Tag creation, deletion, pushing.

### 40. **tag_loader.go**
Loads tag list with metadata.

### 41. **tag_loader_test.go**
Tests for tag loading.

### 42. **tag_test.go**
Tests for tag operations.

---

## Version Detection (2 files)

### 43. **version.go** ⭐ MEDIUM PRIORITY
Git version detection and feature compatibility checking.

**Value for fast-cc-git-hooks:** Ensure git version compatibility

### 44. **version_test.go**
Tests for version detection.

---

## Working Tree Operations (4 files)

### 45. **working_tree.go** ⭐ HIGH PRIORITY
**Key Exports (40+ functions):**

**File Staging:**
- `StageFile()` - Stage single file
- `StageFiles()` - Stage multiple files
- `StageAll()` - Stage all files

**Unstaging:**
- `UnstageAll()` - Remove all from staging
- `UnStageFile()` - Unstage single file
- `UnstageTrackedFiles()` - Unstage tracked files
- `UnstageUntrackedFiles()` - Remove untracked from index

**Discard Changes:**
- `DiscardAllFileChanges()` - Discard all changes to file
- `DiscardAllDirChanges()` - Discard directory changes
- `DiscardUnstagedDirChanges()` - Discard unstaged in directory
- `DiscardUnstagedFileChanges()` - Discard unstaged changes
- `DiscardAnyUnstagedFileChanges()` - Discard all unstaged globally

**File Tracking:**
- `Ignore()` - Add to .gitignore
- `Exclude()` - Add to .git/info/exclude
- `RemoveUntrackedDirFiles()` - Remove untracked from directory
- `RemoveTrackedFiles()` - Delete tracked files
- `RemoveConflictedFile()` - Remove conflicted file
- `RemoveUntrackedFiles()` - Execute git clean

**Diff:**
- `WorktreeFileDiff()` - Get file diff
- `WorktreeFileDiffCmdObj()` - Create diff command
- `ShowFileDiff()` - Show diff between commits
- `ShowFileDiffCmdObj()` - Create show-diff command

**Checkout & Reset:**
- `CheckoutFile()` - Checkout file from commit
- `ResetHard()` - Hard reset to reference
- `ResetSoft()` - Soft reset to reference
- `ResetMixed()` - Mixed reset to reference
- `ResetAndClean()` - Reset and remove untracked

**Merge & Conflicts:**
- `OpenMergeToolCmdObj()` - Open merge tool
- `BeforeAndAfterFileForRename()` - Get file states for rename
- `ShowFileAtStage()` - Show file at merge stage
- `ObjectIDAtStage()` - Get object ID at merge stage
- `MergeFileForFiles()` - Merge from three paths
- `MergeFileForObjectIDs()` - Merge using object IDs

**Value for fast-cc-git-hooks:** Validate working tree state, detect staged files

### 46. **working_tree_test.go**
Tests for working tree operations.

### 47. (see sync_test.go entry above - numbering adjustment)

---

## Worktree Operations (3 files)

### 48. **worktree.go** 🔵 LOW PRIORITY
Git worktree management (add, remove, switch worktrees).

### 49. **worktree_loader.go**
Loads worktree list.

### 50. **worktree_loader_test.go**
Tests for worktree loading.

---

## Git Flow (2 files)

### 51. **flow.go** 🔵 LOW PRIORITY
Git-flow workflow support (feature/hotfix/release branches).

### 52. **flow_test.go**
Tests for git-flow operations.

---

## Advanced Features

### 53. **bisect.go** 🔵 LOW PRIORITY
Git bisect operations for finding problematic commits.

### 54. **bisect_info.go**
Bisect state information.

### 55. **blame.go** 🔵 LOW PRIORITY
Git blame functionality.

---

## Custom Commands

### (Correction - custom.go is #17 above)
### **custom.go** 🔵 LOW PRIORITY
Custom user-defined git command execution.

---

## Patch Operations

### (Correction - patch.go covered above)
### **patch.go** 🔵 LOW PRIORITY
Patch creation and manipulation.

---

## Submodule Operations

### **submodule.go** 🔵 LOW PRIORITY
Git submodule management.

---

## Main Branches

### **main_branches.go** ⭐ MEDIUM PRIORITY
Main branch detection (main, master, develop).

**Value for fast-cc-git-hooks:** Scope detection, branch naming

---

## Summary by Priority for Fast-CC-Git-Hooks

### 🔴 HIGH PRIORITY (Must Extract - Core Value)
1. **git_command_builder.go** - Foundation for all git commands
2. **commit.go** - Commit message parsing and history analysis
3. **diff.go** - Diff operations (replaces 9 manual calls)
4. **status.go** - Repository status validation
5. **working_tree.go** - Working tree and staging operations

### 🟡 MEDIUM PRIORITY (Should Extract - Good Value)
6. **branch.go** - Branch info for scope detection
7. **commit_loader.go** - Structured commit history
8. **file_loader.go** - File status and change detection
9. **config.go** - Git configuration access
10. **repo_paths.go** - Repository path utilities
11. **version.go** - Git version compatibility
12. **main_branches.go** - Main branch detection

### 🔵 LOW PRIORITY (Nice to Have - Future Features)
13. **rebase.go** - Advanced interactive features (not needed for hooks)
14. **sync.go** - Push/pull (not needed for hooks)
15. **stash.go** - Stash operations
16. **tag.go** - Tag management
17. **remote.go** - Remote management
18. **worktree.go** - Worktree support
19. **flow.go** - Git-flow support
20. **bisect.go** - Bisect operations
21. **blame.go** - Blame functionality
22. **patch.go** - Patch manipulation
23. **submodule.go** - Submodule operations
24. **reflog_commit_loader.go** - Reflog/undo features
25. **custom.go** - Custom commands

---

## Extraction Recommendation

**Phase 1 (Weeks 1-2):** Extract HIGH priority files
- Start with `git_command_builder.go` as foundation
- Add `status.go` and `working_tree.go` for validation
- Add `diff.go` to replace manual parsing

**Phase 2 (Weeks 3-4):** Extract MEDIUM priority files
- Add `commit.go` and `commit_loader.go` for better commit analysis
- Add `file_loader.go` for improved file detection
- Add `config.go`, `repo_paths.go`, `version.go` for infrastructure

**Phase 3 (Optional):** LOW priority as needed
- Only extract if specific features are requested
- Most are UI/interactive features not needed for git hooks

**Total Core Extraction:** ~12-17 files from 55 total
**Estimated Benefit:** 80% of value from 20-30% of code

This targeted extraction provides maximum ROI while minimizing complexity and maintenance burden.
