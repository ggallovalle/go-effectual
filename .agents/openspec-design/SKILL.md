---
name: openspec-design
description: Create or update a design document for a change. Use when the user wants to work on the technical design after proposal is complete.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: openspec
  version: "2.0"
---

Create or update a design document in Typst format.

Design documents explain HOW to implement a change. They reference proposal (WHY) and specs (WHAT).

**Input**: Change name OR description of what to design.

**Steps**

1. **Resolve change name**
   - If provided, use it
   - If not, infer from conversation context
   - If ambiguous, use `openspec list --json` and let user select

2. **Read existing artifacts for context**
   ```bash
   openspec instructions design --change "<name>" --json
   ```
   This returns the template and dependencies.

3. **Read completed dependencies**
   - proposal.typ (if exists)
   - specs/*.typ (if exist)

4. **Create/update design.typ**

   Use Typst format with template imports.

   Structure:
   ```typst
   #import "../template.typ": *

   = Design: <change-name>

   == Context

   <!-- Background, current state, constraints -->

   == Goals / Non-Goals

   #decision("Goal 1", "What this achieves", "Why it matters")

   == Decisions

   #decision("Decision title")[
     Chosen: <option>
     Rationale: <why>
     Alternatives: <other options considered>
   ]

   == Risks / Trade-offs

   #note([Risk description → Mitigation])

   == Open Questions

   <!-- Outstanding items -->
   ```

5. **Write to**: `openspec/changes/<name>/design.typ`

**Output**

```
## Design Created

**Change:** <name>
**Artifact:** design.typ

Design document created with:
- Context section
- Goals / Non-Goals
- Decisions (with rationale)
- Risks / Trade-offs
- Open Questions

Cross-references: @proposal.<name>, @specs.<capability>
```

**Guardrails**
- Always read proposal and specs before creating design
- Use decision() function to highlight choices with rationale
- Use note() function for warnings and mitigations
- Cross-reference related artifacts
- Write to .typ extension, not .md