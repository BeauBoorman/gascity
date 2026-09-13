# SSE workflow consumer contract

| Field | Value |
|---|---|
| Bead | `ga-k0jndp` |
| Source | Architecture follow-up `ga-szsfhn`, item (c) |
| Scope | External-consumer confirmation only; no Gas City API change |
| Status | Resolved: no current affected consumer identified; risk is hypothetical |

## Question to settle

PR #6219 keeps ordered SSE event delivery responsive by omitting the optional
`workflow` projection while a subscriber is behind. The event itself still
arrives, but an omitted projection on `bead.created` or `bead.closed` does not
carry `requires_resync`; that flag is emitted only for `bead.updated`.

The remaining contract question is narrow: does any out-of-repo consumer apply
`workflow.changed_fields` incrementally, and, if so, does an event with no
`workflow` projection trigger a full workflow re-read?

## Evidence gathered

- The in-repo dashboard is ruled out by direct inspection of
  `internal/api/dashboardspa/web/frontend/src/hooks/useGcEvents.ts`. Its event
  handler validates only `type`, matches a configured prefix, and schedules a
  coalesced full REST refetch. It does not read `workflow`, `changed_fields`, or
  `requires_resync`.
- The only externally documented candidate, T3 Code, is not a Gas City SSE
  consumer in its public `main` branch at commit
  `20363c32c9bfdbf49c2716ef11d1f18483fcc01b`. Repository-wide code searches
  found no references to Gas City, `/v0/events/stream`, `changed_fields`, or
  `requires_resync`.
- T3 Code documents clients connecting to an environment over HTTP and typed
  WebSocket RPC. Its server registers the application transport at `/ws`.
- The integration direction in this repository is the reverse of the suspected
  consumer: `internal/runtime/t3bridge` is a Gas City runtime provider that
  connects to T3 Code's WebSocket API. It does not subscribe to Gas City's SSE
  API or inspect workflow projections.
- The mayor confirmed that no private or out-of-repo relay is known and found no
  `/v0/events/stream` consumer in any pack, script, or service in the city.
  Sibling repositories `tincan`, `tincan-iris`, `hold-court`, `MCDClient`, and
  `factory` have no consumer. The only sibling consumer is gascity-packs'
  slack-pack teardown subscriber, which decodes only event type and payload and
  deliberately omits workflow projection handling (mail `gm-wisp-84a1pn`).

## Resolution

No current consumer has been identified that applies
`workflow.changed_fields` incrementally. The omission behavior introduced by PR
#6219 therefore presents a hypothetical compatibility risk, not a demonstrated
stale-state defect. No implementation or architecture follow-up is warranted
without evidence of an affected consumer. A consumer outside this host may
exist, but the operator does not need to be awaited for this investigation.

## Outcome matrix

| Owner confirmation | Result for this bead | Follow-up |
|---|---|---|
| No current consumer exists | Record the review concern as hypothetical and close | None |
| Consumer re-reads when `workflow` is absent | Contract is safe and this bead closes | None |
| Consumer treats omission as no-op | Record the stale-state mechanism and close this investigation | Create one `needs-architecture` bead to decide resync semantics versus an explicitly accepted risk |
| Consumer ownership remains unknown | Keep this bead open as an external information dependency | Mayor identifies the owner; do not create speculative implementation work |

## Acceptance criteria

- The actual consumer, or the absence of one, is identified from owner-confirmed
  information.
- If a consumer exists, its behavior for an absent `workflow` projection on
  `bead.created` and `bead.closed` is recorded.
- No implementation bead is created unless a real consumer can become stale.
- Any required contract decision is routed to architecture; this investigation
  does not choose API semantics.

## References

- [Gas City PR #6219](https://github.com/gastownhall/gascity/pull/6219)
- [T3 Code remote architecture](https://github.com/pingdotgg/t3code/blob/20363c32c9bfdbf49c2716ef11d1f18483fcc01b/docs/internals/remote.md)
- [T3 Code WebSocket route](https://github.com/pingdotgg/t3code/blob/20363c32c9bfdbf49c2716ef11d1f18483fcc01b/apps/server/src/ws.ts)
