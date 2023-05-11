---
type: project
title: Telegram integration
number:
state: open
visibility: public
owner: lifeosm
creator: kamilsk
url:
short_description: Read-only personal-data ingestion from Telegram into the Sparkle index — dialogs, messages, topics, contacts, multi-account.
created_at:
updated_at:
milestones:
spec:
tags:
  - type/project
  - topic/telegram
---
# Telegram integration

Bring Telegram into the `indexit` ingest pipeline as a first-class, read-only
source for the Sparkle index. The `indexit telegram` subcommand streams the
user's own Telegram data (dialogs, message history, forum topics, contacts, and
Story viewers) as stable JSONL that a downstream Sparkle indexer can consume
without re-fetching.

Design and scope live in the spec: [[Telegram fetcher, PoC implementation plan]].

## Milestones

| Milestone                               | State | Progress                                           |
| --------------------------------------- | ----- | -------------------------------------------------- |
| [[Telegram fetcher PoC implementation]] | open  | Phase 1 functionally complete; 7 / 16 tasks closed |

Later milestones (not yet opened) will graduate the PoC into the wider pipeline —
e.g. durable peer/message store, and wiring the JSONL stream into the Sparkle
indexer itself.

## Tracking

- Spec / plan: [[Telegram fetcher, PoC implementation plan]]
- Spec-compliance audits: `.github/reports/`
- Code-review iterations: `.github/reviews/`
