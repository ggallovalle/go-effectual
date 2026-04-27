---
name: openspec-sync-specs
description: Sync delta specs from a change to main specs. Use when the user wants to reconcile specs from an archived change into the main specs directory.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: openspec
  version: "1.0"
  generatedBy: "1.3.1"
---

Sync delta specs from a change to the main specs.

**Input**: Change name (required). If omitted, check if it can be inferred from conversation context. If vague or ambiguous you MUST prompt for available changes.

**Steps**

1. **Resolve the change name**

   If a name is provided, use it. Otherwise:
   - Infer from conversation context if the user mentioned a change
   - Run `openspec list --json` to get available changes
   - Filter to changes that have delta specs at `openspec/changes/<name>/specs/`
   - Use the **AskUserQuestion tool** to let the user select

   Always announce: "Using change: <name>"

2. **Read delta specs**

   Read all delta spec files from `openspec/changes/<name>/specs/**/spec.md`

   For each delta spec, parse section headers to identify operation types:
   - `## ADDED Requirements` → add operations
   - `## MODIFIED Requirements` → modify operations
   - `## REMOVED Requirements` → remove operations
   - `## RENAMED Requirements` → rename operations

3. **Read main specs**

   For each capability in delta specs:
   - Check if main spec exists at `openspec/specs/<capability>/spec.md`
   - If exists, read and parse it
   - If not exists, treat as new capability

4. **Reconcile**

   Process each delta spec:

   **For ADDED requirements:**
   - If requirement doesn't exist in main spec → add it
   - If requirement exists → update to match delta version

   **For MODIFIED requirements:**
   - If requirement exists in main spec → replace with delta version
   - If requirement doesn't exist → add it (treated as ADDED)

   **For REMOVED requirements:**
   - If requirement exists in main spec → remove it

   **For RENAMED requirements:**
   - Parse FROM:/TO: format
   - If FROM requirement exists → rename to TO

   **For new capability:**
   - Create new main spec file at `openspec/specs/<capability>/spec.md`

5. **Write main specs**

   Write reconciled spec back to `openspec/specs/<capability>/spec.md`

6. **Verify idempotency**

   Re-read the just-written main spec and delta, verify they match.
   If mismatch, retry reconciliation once.

**Output On Success**

```
## Sync Complete

**Change:** <name>
**Capabilities synced:** <list>

<capability1>:
- <N> requirements added
- <N> requirements modified
- <N> requirements removed
- <N> requirements renamed
<capability2>: ...
```

**Output When Already Synced**

```
## Specs Already In Sync

**Change:** <name>
**Specs:** No changes needed - main specs match delta specs
```

**Output On Error**

```
## Sync Failed

**Error:** <error message>

Retry or contact support.
```

**Guardrails**
- Always prompt for change selection if not provided or ambiguous
- Preserve requirement content exactly as specified in delta
- Idempotent: running multiple times produces same result
- Show clear summary of what was applied per capability
- If no delta specs exist for a change, report and exit
