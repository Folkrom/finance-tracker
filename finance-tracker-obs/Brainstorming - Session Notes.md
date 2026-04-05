# Brainstorming - Session Notes

> **Date:** 2026-03-30 (started), 2026-04-03 (design review completed)
> **Session:** Initial brainstorming, approach selection & design review

---

## What We Did

1. **Explored Notion** — fetched and analyzed all database schemas from "Personal Finance Tracker - 2026" and "Plantilla"
2. **Clarifying questions** — resolved 9 key design decisions (see [[Finance Tracker - Design Spec#Key Design Decisions]])
3. **Evaluated 3 approaches** — selected Approach C (Go API + Next.js SSR)
4. **Design review** — walked through all 6 sections, resolved open items, verified against Notion charts and Wishlist

---

## Questions Asked & Answers

| # | Question | Answer |
|---|---|---|
| 1 | Keep Cuentas por Cobrar & PPR Teresa? | Drop both |
| 2 | Shared or independent categories for Expenses/Debt? | Single shared list |
| 3 | Year-scoped or continuous database? | Year-scoped workspaces, cross-year queryable |
| 4 | UI language? | i18n (EN + ES with toggle) |
| 5 | Payment methods: fixed or user-defined? | User-defined + typed |
| 6 | "Daily cash expenses" meaning? | All expenses by day (cash = cash-flow) |
| 7 | Data migration? | Later concern |
| 8 | Card usage calculation? | Auto from Debt + manual override |
| 9 | Auth & multi-user? | Multi-tenant, goal to commercialize |

---

## Open Items Resolved (Session 2)

| Item | Resolution |
|---|---|
| "Debt" category vs module naming | Renamed category to "Card Payments" |
| "Alcohol/Drugs" vs "Alcohol" | Split into "Alcohol" and "Drugs" as separate categories |
| Card health thresholds | Green 0-20%, Yellow 21-30%, Orange 31-70%, Red 71-100%+ |
| Budget recurrence | Recurring defaults + per-month overrides |
| Dashboard charts | 6 charts total: 3 from Notion (line, 2 donuts) + 3 new (donut, 2 bar charts) |
| Wishlist details | Verified against Notion: file upload + URL fallback for images, up to 5 buy links, kanban statuses, 3 views (Gallery/Table/Board), separate category list |
| Currency | Configurable per-field (USD/MXN for now, extensible) |
| Default language | English, with i18n for Spanish |

---

## Approaches Evaluated

### A: Monolith API + Next.js SPA
- Simple but no SSR, monolith gets heavy over time

### B: Next.js Full-Stack + Go Microservices
- Fast to prototype but splits logic across runtimes, more TS than Go

### C: Go API + Next.js SSR (Selected)
- Clean separation, plays to Go strengths, SSR for dashboards, scales naturally
- Two services to deploy (trivial with Docker Compose)

---

## Next Steps

- [ ] User reviews written spec in [[Finance Tracker - Design Spec]]
- [ ] Invoke writing-plans skill to create implementation plan
