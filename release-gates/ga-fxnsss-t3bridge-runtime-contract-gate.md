# Release gate: T3 bridge shared runtime contract proof

- Deploy bead: `ga-fxnsss`
- Build bead: `ga-uz5t3a.1`
- Review bead: `ga-zj56f3`
- Reviewed source: `56ce4de80e0794d8208da6cb8d0d2702d840346e`
- Reviewed base: `835718da0677d21ca8bc3e28ce51a556930fb2cf`
- Gate base: `origin/main@003b78721dcf30e2354cb5ee0ae937b7f5efe9bf`
- Deploy mode: remote
- Gate result: **FAIL**

## Gate criteria

| # | Criterion | Result | Evidence |
|---|---|---|---|
| 1 | Reviewer PASS present | **SKIPPED** | Fail-fast after criterion 6. The closed review bead does record `verdict: pass` for the exact reviewed tip, but this criterion was not independently evaluated as a release decision. |
| 2 | Acceptance criteria met | **SKIPPED** | Fail-fast after criterion 6; no independent acceptance run was performed. |
| 3 | Tests pass | **SKIPPED** | Fail-fast after criterion 6; the documented full-suite command was intentionally not run on a source that cannot merge cleanly. `test_cmd_scope: not-run`; `test_counts: not-run`; `diff_tests_executed: not-run`; `waiver_ref: none`. |
| 3b | Policy/lint lane | **SKIPPED** | Fail-fast after criterion 6; no policy/lint release decision was made. |
| 3c | CI-config lane | **SKIPPED** | Fail-fast after criterion 6. `ci_lane_run: n/a (no CI-config file in the diff)`. |
| 4 | No unresolved HIGH review findings | **SKIPPED** | Fail-fast after criterion 6. |
| 5 | Final branch clean | **SKIPPED** | Fail-fast after criterion 6. |
| 6 | Branch diverges cleanly from main | **FAIL** | `git merge-tree --write-tree origin/main 56ce4de80e0794d8208da6cb8d0d2702d840346e` reports a content conflict in `internal/runtime/t3bridge/provider_test.go`. Current main changed the same fixture in `2b9b66e511` (`fix(runtime): support T3 Code 0.0.33 websocket auth (#5538)`). |
| 7 | Single feature theme | **SKIPPED** | Fail-fast after criterion 6. |

## Criterion 6 evidence

The target SHA is present, equals `origin/builder/ga-uz5t3a.1`, and is a
fast-forward descendant of the review bead's recorded base. GitHub reports no
PR carrying the reviewed SHA, so the already-merged/closed reconciliation
paths do not apply.

Against current `origin/main`, Git reports:

```text
Auto-merging TESTING.md
Auto-merging internal/runtime/t3bridge/provider.go
Auto-merging internal/runtime/t3bridge/provider_test.go
CONFLICT (content): Merge conflict in internal/runtime/t3bridge/provider_test.go
Auto-merging internal/testutil/providerledger/ledger.go
```

The bounded self-rebase exception was not invoked. The deploy bead and the
deployer input contract both identify `builder/ga-uz5t3a.1` as provenance only
and explicitly prohibit pushing it; `attempt_bounded_self_rebase` necessarily
force-with-lease pushes the branch it rebases. Resolving or publishing this
content conflict from the deployer seat would exceed that exception's branch
ownership guardrail.

## Release disposition

**Gate FAIL.** Return `ga-fxnsss` to the builder to rebase the implementation
onto current main, resolve the shared T3 test-fixture conflict, rerun TDD and
review, and provide a new reviewed SHA. No deploy branch, push, PR, status, or
merge action is permitted from this failed gate.
