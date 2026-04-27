---
name: openspec-sync-specs
description: Sync delta specs from a change to main specs. Use when the user wants to reconcile specs from an archived change into the main specs directory.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: openspec
  version: "2.0"
---

Sync delta specs from a change to the main specs.

**Input**: Change name (required). If omitted, check if it can be inferred from conversation context.

**Steps**

1. **Resolve the change name**

   If a name is provided, use it. Otherwise:
   - Infer from conversation context if the user mentioned a change
   - Run `openspec list --json` to get available changes
   - Use the **AskUserQuestion tool** to let the user select

   Always announce: "Using change: <name>"

2. **Read delta specs**

   Read all delta spec files from `openspec/changes/<name>/specs/**/*.typ`

   For each delta spec, parse the Typst content to identify operation types:
   - Look for `## ADDED Requirements` → add operations
   - Look for `## MODIFIED Requirements` → modify operations
   - Look for `## REMOVED Requirements` → remove operations
   - Look for `## RENAMED Requirements` → rename operations

3. **Read main specs**

   For each capability in delta specs:
   - Check if main spec exists at `openspec/specs/<capability>/spec.typ`
   - If exists, read and parse it
   - If not exists, treat as new capability

4. **Reconcile at requirement level**

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
   - Parse the FROM:/TO: format
   - If FROM requirement exists → rename to TO

   **For new capability:**
   - Create new main spec file at `openspec/specs/<capability>/spec.typ`

5. **Write main specs**

   Write reconciled spec back to `openspec/specs/<capability>/spec.typ`

   Use Typst format with proper structure.

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

**Guardrails**
- Always prompt for change selection if not provided
- Preserve requirement content exactly as specified in delta
- Idempotent: running multiple times produces same result
- Show clear summary of what was applied per capability
- If no delta specs exist for a change, report and exit
- Use Typst format for all output files (.typ extension)