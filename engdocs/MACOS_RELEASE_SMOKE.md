# macOS Release Smoke — validating *released* gc artifacts on darwin

This is the operator-facing companion to `.github/workflows/macos-release-smoke.yml`
(and the existing `homebrew-tap-smoke.yml`). The tap smoke proves the **tap**
formula installs; this procedure proves the artifacts users actually get —
**Homebrew core** installs and **goreleaser darwin binaries** — survive
day-one platform and packaging regressions before users report them.

It exists because one didn't: v1.4.1 shipped with an embedded bd older than
the bd formula it depends on, and every native store path (`gc mail count`,
nudge delivery, session wake) failed against a current store while
`gc version` — the only thing the tap smoke checked — passed.
(gastownhall/gascity#6393, gastownhall/gascity#6416.)

## When to run it

- Before tagging any release (release captain, ~10 minutes)
- On any Homebrew core formula bump that changes the `beads` dependency range
- Whenever a user reports a darwin-only regression — this is the reproducer harness

## The procedure

Run on a clean or brew-clean macOS machine (arm64 or amd64):

### 1. Homebrew core artifact

```sh
brew uninstall --force gascity beads 2>/dev/null || true
brew install gascity                      # homebrew-core formula, bottled
```

### 2. Binary sanity (the tap smoke's existing floor)

```sh
gc version                                # prints a version, exit 0
```

### 3. Store interop sanity (the floor the tap smoke missed)

```sh
mkdir -p /tmp/gc-smoke && cd /tmp/gc-smoke
bd init                                   # creates a store at the CLI's schema version
gc bd list                                # gc's embedded bd must read the CLI's store schema
gc mail count                             # native store path must not error
```

If `gc bd` works but `gc mail count` errors with a schema-version complaint,
the release embeds an older bd than the formula's — **do not ship**. This is
exactly the v1.4.1 failure mode: bottled gc + current bd CLI + migrated store
= broken native paths with a green `gc version`.

### 4. goreleaser artifact (one arch per release)

```sh
curl -fsSL -o gc.tar.gz \
  "https://github.com/gastownhall/gascity/releases/download/v<VERSION>/gc_<VERSION>_darwin_arm64.tar.gz"
tar xzf gc.tar.gz && ./gc version
./gc bd list                              # same store-interop check against a real store
```

### 5. Version-contains check

```sh
git tag --contains "$(go list -m -f '{{.Version}}' github.com/steveyegge/beads 2>/dev/null)" || \
  echo "WARNING: embedded bd commit is not in any tag — the interop check above is load-bearing"
```

## CI

`.github/workflows/macos-release-smoke.yml` runs steps 1–3 on
`macos-latest` daily and on demand. It fails on any schema-version mismatch
between the bottled gc and the current bottled `beads` — the class of
regression that shipped in v1.4.1.

## Failure triage

| Symptom | Likely cause | Reference |
| --- | --- | --- |
| `gc bd` works, `gc mail count` errors on schema version | embedded bd older than formula bd | #6393 |
| `gc hook --claim` exceeds its window on a shared dolt server | work-query probe starvation, not packaging | #6416 |
| install succeeds, `gc version` fails | bottle/link issue in the formula | file with `brew doctor` output |
