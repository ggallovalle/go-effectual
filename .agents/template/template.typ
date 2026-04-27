// OpenSpec Template
// Base template for all OpenSpec artifacts in Typst

// ============================================================================
// TYPOGRAPHY
// ============================================================================

#let body-font = "New Computer Modern"
#let heading-font = "New Computer Modern"

#show heading.where(level: 1): set text(20pt, weight: "bold", font: heading-font)
#show heading.where(level: 2): set text(16pt, weight: "bold", font: heading-font)
#show heading.where(level: 3): set text(13pt, weight: "semibold", font: heading-font)
#show heading.where(level: 4): set text(11pt, weight: "semibold", font: heading-font)

// ============================================================================
// COLORS
// ============================================================================

#let color-primary = rgb("#2563EB")
#let color-secondary = rgb("#64748B")
#let color-accent = rgb("#7C3AED")
#let color-success = rgb("#059669")
#let color-warning = rgb("#D97706")
#let color-danger = rgb("#DC2626")

#let color-req-bg = rgb("#F0F9FF")
#let color-req-border = rgb("#2563EB")
#let color-scenario-bg = rgb("#F8FAFC")
#let color-scenario-border = rgb("#94A3B8")
#let color-note-bg = rgb("#FEF9C3")
#let color-note-border = rgb("#CA8A04")

// ============================================================================
// REQUIREMENT BOX
// ============================================================================

#let requirement(name, body) = block(
  width: 100%,
  fill: color-req-bg,
  stroke: (left: 3pt + color-req-border),
  inset: (left: 12pt, top: 8pt, bottom: 8pt, right: 12pt),
  radius: (left: 4pt),
)[
  #text(13pt, weight: "bold", fill: color-req-border)[#name]
  #linebreak()
  #body
]

// ============================================================================
// SCENARIO BLOCK
// ============================================================================

#let scenario(name, content) = block(
  width: 100%,
  fill: color-scenario-bg,
  stroke: (left: 2pt + color-scenario-border),
  inset: (left: 10pt, top: 6pt, bottom: 6pt, right: 10pt),
  radius: (left: 3pt),
)[
  #text(11pt, weight: "bold", fill: color-scenario-border)[Scenario: #name]
  #linebreak()
  #content
]

// ============================================================================
// WHEN/THEN LINE
// ============================================================================

#let when-then(condition, outcome) = {
  [WHEN ] + text(weight: "bold")[#condition]
  [ THEN ] + text(weight: "bold")[#outcome]
}

// ============================================================================
// TASK CHECKBOX
// ============================================================================

#let task(num, description, completed: false) = {
  let checkbox = if completed { "✓" } else { "○" }
  let checkbox-color = if completed { color-success } else { color-secondary }
  text(checkbox-color)[#checkbox] + " " + text(weight: "bold")[#num] + " " + description
}

// ============================================================================
// INCOMPLETE TASK CHECKBOX
// ============================================================================

#let task-unchecked(num, description) = task(num, description, completed: false)

// ============================================================================
// COMPLETED TASK CHECKBOX
// ============================================================================

#let task-checked(num, description) = task(num, description, completed: true)

// ============================================================================
// CROSS-REFERENCE
// ============================================================================

#let ref(target, label) = {
  let styled-target = text(color-primary, style: "italic")[#target]
  link(target)[#styled-target]
}

// ============================================================================
// SECTION LABELS (for cross-referencing)
// ============================================================================

#let section-label(name) = {
  hide(label(name))
  text(weight: "bold", color: color-secondary)[#name]
}

// ============================================================================
// NOTE/WARNING BOX
// ============================================================================

#let note(content) = block(
  width: 100%,
  fill: color-note-bg,
  stroke: (left: 2pt + color-note-border),
  inset: (left: 10pt, top: 6pt, bottom: 6pt, right: 10pt),
  radius: (left: 3pt),
)[
  #text(11pt, weight: "bold", fill: color-note-border)[NOTE]
  #linebreak()
  #content
]

// ============================================================================
// DECISION BLOCK
// ============================================================================

#let decision(title, chosen, rationale, alternatives: none) = block(
  width: 100%,
  fill: rgb("#F5F5F5"),
  stroke: (left: 3pt + color-accent),
  inset: (left: 12pt, top: 8pt, bottom: 8pt, right: 12pt),
  radius: (left: 4pt),
)[
  #text(13pt, weight: "bold", fill: color-accent)[Decision: #title]
  #linebreak()
  #text(weight: "bold")[Chosen: ] + chosen
  #linebreak()
  #text(weight: "bold")[Rationale: ] + rationale
  #if alternatives != none [
    #linebreak()
    #text(weight: "bold")[Alternatives: ] + alternatives
  ]
]

// ============================================================================
// CODE BLOCK
// ============================================================================

#let code-block(content, lang: none) = block(
  width: 100%,
  fill: rgb("#1E1E1E"),
  inset: (top: 8pt, bottom: 8pt, left: 12pt, right: 12pt),
  radius: 4pt,
)[
  #text(10pt, font: "JetBrains Mono", fill: rgb("#D4D4D4"))[#content]
]

// ============================================================================
// EXPORT ALL
// ============================================================================

#let exports = (
  requirement: requirement,
  scenario: scenario,
  when-then: when-then,
  task: task,
  task-unchecked: task-unchecked,
  task-checked: task-checked,
  ref: ref,
  section-label: section-label,
  note: note,
  decision: decision,
  code-block: code-block,
)