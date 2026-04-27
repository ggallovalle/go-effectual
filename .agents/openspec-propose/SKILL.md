---
name: openspec-propose
description: Propose a new change with all artifacts generated in one step. Use when the user wants to quickly describe what they want to build and get a complete proposal with design, specs, and tasks ready for implementation.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: openspec
  version: "2.0"
  generatedBy: "2.0"
---

Propose a new change - create the change and generate all artifacts in one step.

All artifacts are written in Typst format (.typ files).

Artifacts created:
- proposal.typ (what & why)
- design.typ (how)
- specs/*.typ (detailed requirements)
- tasks.typ (implementation steps)

When ready to implement, use openspec-apply-change skill.

---

**Input**: The user's request should include a change name (kebab-case) OR a description of what they want to build.

**Steps**

1. **If no clear input provided, ask what they want to build**

   Use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "What change do you want to work on? Describe what you want to build or fix."

   From their description, derive a kebab-case name (e.g., "add user authentication" → `add-user-auth`).

   **IMPORTANT**: Do NOT proceed without understanding what the user wants to build.

2. **Create the change directory**
   ```bash
   openspec new change "<name>"
   ```
   This creates a scaffolded change at `openspec/changes/<name>/`.

3. **Get the artifact build order**
   ```bash
   openspec status --change "<name>" --json
   ```
   Parse the JSON to get:
   - `applyRequires`: array of artifact IDs needed before implementation (e.g., `["tasks"]`)
   - `artifacts`: list of all artifacts with their status and dependencies

4. **Get instructions for each artifact**
   ```bash
   openspec instructions <artifact-id> --change "<name>" --json
   ```

   For each artifact, the instructions include:
   - `context`: Project background
   - `rules`: Artifact-specific rules
   - `template`: The structure to use
   - `instruction`: Schema-specific guidance
   - `outputPath`: Where to write the artifact
   - `dependencies`: Completed artifacts to read for context

5. **Read dependency artifacts for context**

   Before creating each artifact, read any completed dependency files.

6. **Create artifacts in Typst format**

   Use the template functions from `.agents/template/template.typ`:
   - `#import "template.typ": *`
   - Use `requirement()`, `scenario()`, `task()`, etc.

   **IMPORTANT**: All output paths use `.typ` extension, not `.md`.

   For specs, create `specs/<capability>/spec.typ` files.

7. **Track progress**

   Use the **TodoWrite tool** to track artifact creation.

   Loop: create artifact → re-check status → repeat until `applyRequires` artifacts are done.

8. **Final status**
   ```bash
   openspec status --change "<name>"
   ```

**Output**

After completing all artifacts:
```
## Change Created

**Name:** <change-name>
**Location:** openspec/changes/<name>/

### Artifacts Created

- proposal.typ
- design.typ
- specs/<cap1>/spec.typ
- specs/<cap2>/spec.typ
- tasks.typ

All artifacts in Typst format. Ready for implementation.
```

**Typst Format Guidelines**

All artifacts MUST use Typst syntax:

```typst
#import "../template.typ": *

= Proposal: <change-name>

== Why

<!-- content -->

== What Changes

#bullet-item([List item 1])
#bullet-item([List item 2])
#bullet-item([List item 3])

== Capabilities

=== New Capabilities

#bullet-item([`new-capability`: Description])

=== Modified Capabilities

#bullet-item([`existing-capability`: What changed])

== Impact

<!-- content -->

// Cross-references to other specs:
// @proposal.<change-name> - reference this proposal
// @specs.<capability>.<requirement-id> - reference a spec requirement
```

For requirements in specs:

```typst
#requirement("Requirement: <name>")[
  Description text with SHALL/MUST language.

  #scenario("Scenario name")[
    #when-then("condition", "outcome")
    #when-then("condition 2", "outcome 2")
  ]
]
```

For tasks:

```typst
#task-unchecked("1.1", [Task description])
#task-unchecked("1.2", [Task description])
#task-checked("2.1", [Completed task])
```

**Guardrails**
- Create ALL artifacts needed for implementation
- Always read dependency artifacts before creating new ones
- If context is critically unclear, ask the user
- Verify each artifact file exists after writing
- Output files use `.typ` extension, not `.md`
- Use cross-references to link artifacts together