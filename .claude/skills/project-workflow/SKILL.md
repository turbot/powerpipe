---
name: project-workflow
description: Plan and organize complex multi-task projects. Use for work requiring 3+ related tasks, parallel execution, or careful coordination. Creates .claude/wip/ project structure.
---

# Project Workflow Skill

This skill provides patterns for planning, organizing, and executing complex projects with multiple related tasks. Whether working sequentially or with parallel agents, this skill helps structure larger efforts effectively.

## Purpose

When working on projects with multiple related tasks (e.g., fixing several bugs, implementing a feature set, refactoring a subsystem), this skill helps:
- Plan and break down complex work into manageable tasks
- Organize work for sequential or parallel execution
- Maintain clear task boundaries and dependencies
- Track overall progress across tasks and agents
- Prevent merge conflicts using git worktrees (for parallel work)

## Core Workflow

### 1. Create Work Directory

```bash
mkdir -p .claude/wip/<topic-name>
```

Examples:
- `.claude/wip/bug-fixes-wave-1/`
- `.claude/wip/refactor-auth-layer/`
- `.claude/wip/oci-policy-additions/`

### 2. Create Project Plan and Task Files

Create individual task files and a coordination plan (whether working alone or with multiple agents):

```bash
# In .claude/wip/<topic>/
task-1-description.md      # First task
task-2-description.md      # Second task
task-3-description.md      # Third task
plan.md                    # Overall project plan and dependencies
status.md                  # Progress tracking (optional)
```

#### Task File Format

Each `task-N-*.md` file should contain:

```markdown
# Task N: [Brief Title]

## Objective
Clear, specific goal for this task

## Context
- Why this task exists
- How it relates to other tasks
- Key constraints or requirements

## Dependencies
- Files/modules this task will modify
- Other tasks that must complete first (if any)
- Shared resources to coordinate

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Tests pass

## Notes
Any additional context or considerations
```

#### Plan File Format

The `plan.md` should contain:

```markdown
# Project Plan: [Topic Name]

## Overview
High-level description of the project and its goals

## Task Breakdown
1. **Task 1**: Brief description (Status: Pending)
2. **Task 2**: Brief description (Status: Pending)
3. **Task 3**: Brief description (Status: Pending)

## Execution Strategy
- Sequential: Tasks run one after another
- Parallel: Multiple agents work simultaneously on independent tasks

## Dependencies
- Task 2 depends on Task 1 completing
- Task 3 can run in parallel with Task 1
- Task 4 requires Tasks 2 and 3 to be complete

## Shared Resources
- Files that multiple tasks might touch
- Coordination strategy for conflicts (if working in parallel)

## Integration Plan
How to merge work back together (especially for parallel work)

## Notes
Additional coordination details
```

### 3. Git Worktree Setup (Recommended for Parallel Work)

To prevent merge conflicts when agents work in parallel:

```bash
# Create worktree for each agent
git worktree add .claude/wip/<topic>/worktree-agent-1 -b wip/<topic>/task-1
git worktree add .claude/wip/<topic>/worktree-agent-2 -b wip/<topic>/task-2
git worktree add .claude/wip/<topic>/worktree-agent-3 -b wip/<topic>/task-3
```

Then agents work in their respective worktrees:
```bash
cd .claude/wip/<topic>/worktree-agent-1
# Agent 1 does work here

cd .claude/wip/<topic>/worktree-agent-2
# Agent 2 does work here
```

When done:
```bash
# Merge branches back to main
git checkout main
git merge wip/<topic>/task-1
git merge wip/<topic>/task-2

# Clean up worktrees
git worktree remove .claude/wip/<topic>/worktree-agent-1
git worktree remove .claude/wip/<topic>/worktree-agent-2
```

### 4. Executing Tasks

**Sequential execution** (work through tasks one at a time):
```
"Read plan in .claude/wip/bug-fixes/ and execute task-1"
[Wait for completion]
"Read plan in .claude/wip/bug-fixes/ and execute task-2"
[Wait for completion]
"Read plan in .claude/wip/bug-fixes/ and execute task-3"
```

**Parallel execution** (launch multiple agents simultaneously):
```
Launch 3 agents in parallel:
- Agent 1: "Read .claude/wip/bug-fixes/task-1-*.md and execute"
- Agent 2: "Read .claude/wip/bug-fixes/task-2-*.md and execute"
- Agent 3: "Read .claude/wip/bug-fixes/task-3-*.md and execute"
```

**Mixed approach** (some sequential, some parallel):
```
"Read and execute task-1"  # Must complete first
[Wait for completion]

Launch tasks 2 and 3 in parallel:  # Both depend on task 1
- Agent 1: "Read and execute task-2"
- Agent 2: "Read and execute task-3"
```

## When to Use This Pattern

**Use when:**
- Working on complex projects with 3+ related tasks
- Need to track progress across multiple work items
- Want to break down large efforts into manageable pieces
- Working with clear task boundaries and dependencies
- Want option to parallelize independent work for speed

**Don't use when:**
- Single straightforward task
- Tightly coupled changes requiring iterative exploration
- Unclear requirements (explore first, then create plan)
- Very small changes (overkill for simple edits)

## Tips for Effective Project Planning

1. **Clear boundaries**: Ensure tasks are well-defined with specific goals
2. **Document dependencies**: Be explicit about task ordering and relationships
3. **Track progress**: Use `status.md` or update task files as work completes
4. **Atomic tasks**: Each task should be independently testable
5. **Plan integration**: Know how pieces fit together before starting
6. **Choose execution mode**: Sequential for dependent work, parallel for independent tasks
7. **Use worktrees for parallel work**: Prevents git conflicts when multiple agents work simultaneously

## Example: Bug Fix Wave

```bash
# Setup
mkdir -p .claude/wip/bug-fixes-wave-1

# Create tasks (coordinator)
echo "Fix null pointer in auth.ts" > .claude/wip/bug-fixes-wave-1/task-1-auth-null.md
echo "Fix race condition in sync.ts" > .claude/wip/bug-fixes-wave-1/task-2-sync-race.md
echo "Fix memory leak in cache.ts" > .claude/wip/bug-fixes-wave-1/task-3-cache-leak.md

# Execute tasks (choose sequential or parallel)
# Sequential:
# "Execute .claude/wip/bug-fixes-wave-1/task-1-auth-null.md"

# Parallel:
# Launch 3 agents in parallel with task-1, task-2, task-3
```

## Cleanup

After all tasks complete and merge:

```bash
# Archive the work plan
mv .claude/wip/<topic> .claude/wip/archive/<topic>-$(date +%Y%m%d)

# Or just delete if not needed
rm -rf .claude/wip/<topic>
```
