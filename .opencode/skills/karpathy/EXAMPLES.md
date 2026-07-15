# Examples

One wrong-vs-right example per rule. Code samples use a real-time multiplayer game server (Node + WebSocket) — the domain this ruleset was hardened on — instead of generic snippets.

---

## 1. Think Before Coding

**User request:** "Add a way to kick players from a room."

**❌ What LLMs do (silent assumptions)**

```js
function kickPlayer(room, playerId) {
  const ws = room.players.get(playerId);
  ws.close();
  room.players.delete(playerId);
}
```

Assumed: who is allowed to kick (any client? admin only?), whether the player should be told why, whether a kicked player can rejoin instantly, and what happens to their in-game mass/score.

**✅ What should happen (surface assumptions)**

```
Before I implement kick, a few things change the design:

1. Authorization — admin-only, or can players vote-kick? (auth check location differs)
2. Feedback — send the kicked client a `{t:'kicked', reason}` message, or just drop the socket?
3. Rejoin — banned for the match, or free to reconnect immediately?
4. State — does their mass get removed from the sim, or left as food?

Simplest version that matches your admin panel: admin-only, send a kick message,
remove from sim, allow rejoin. Want that, or stricter?
```

---

## 2. Simplicity First

**User request:** "Add a helper to give a player a score bonus."

**❌ What LLMs do (overengineered)**

```js
class BonusStrategy {
  apply(score) { throw new Error('not implemented'); }
}
class FlatBonus extends BonusStrategy {
  constructor(amount) { super(); this.amount = amount; }
  apply(score) { return score + this.amount; }
}
class MultiplierBonus extends BonusStrategy {
  constructor(factor) { super(); this.factor = factor; }
  apply(score) { return score * this.factor; }
}
class BonusEngine {
  constructor(strategy, { cap = Infinity, floor = 0 } = {}) {
    this.strategy = strategy; this.cap = cap; this.floor = floor;
  }
  run(score) {
    return Math.min(this.cap, Math.max(this.floor, this.strategy.apply(score)));
  }
}
// 25 lines of setup to add a number
```

**✅ What should happen (simple)**

```js
function addBonus(score, bonus) {
  return score + bonus;
}
```

**Add complexity later** — only if multipliers/caps actually become a requirement.

---

## 3. Surgical Changes

**User request:** "Fix the crash when a player sends an empty name."

**❌ What LLMs do (drive-by refactor)**

```diff
  function joinRoom(room, name, ws) {
-     room.players.set(name, ws);
+     // Validate and normalize the player name
+     const clean = (name ?? '').trim();
+     if (!clean) throw new Error('Name required');
+     if (clean.length > 16) throw new Error('Name too long');
+     if (!/^[a-z0-9_]+$/i.test(clean)) throw new Error('Invalid characters');
+     room.players.set(clean, ws);
+     logger.info(`Player ${clean} joined ${room.id}`);
  }
```

Added length limits, charset rules, and logging nobody asked for.

**✅ What should happen (surgical)**

```diff
  function joinRoom(room, name, ws) {
+     if (!name || !name.trim()) return;
      room.players.set(name, ws);
  }
```

Only the empty-name crash is addressed. Length/charset rules get flagged, not added.

---

## 4. Goal-Driven Execution

**User request:** "The leaderboard is sometimes wrong."

**❌ What LLMs do (vague)**

```
I'll review the leaderboard code, look for issues, improve the sorting,
and test the changes.
```

**✅ What should happen (verifiable goals)**

```
"Sometimes wrong" needs a concrete failure. Most likely: ties sort
non-deterministically. Plan:

1. Write a test: two players with equal score → assert stable order by id
   Verify: test fails (reproduces the bug)
2. Make the sort stable (secondary key on player id)
   Verify: test passes
3. Run 50 snapshots → order never flips
   Verify: existing leaderboard tests still green

What does "wrong" look like for you — wrong order, missing players, or stale values?
```

---

## 5. Honesty

**User request:** "Does `ws.send()` block until the client receives the message?"

**❌ What LLMs do (confident guess)**

```
Yes, ws.send() waits for the client to acknowledge receipt before returning,
so it's safe to assume delivery on the next line.
```

**✅ What should happen (say what you know)**

```
No — and I want to be precise here. `ws.send()` queues the data into the
socket's send buffer and returns immediately; it does NOT wait for client
receipt. The optional callback fires when the data is handed to the OS, not
when the client gets it. If you need delivery confirmation you need an
app-level ack message. (I'm confident on the non-blocking part; the exact
callback timing I'd verify against the `ws` docs before relying on it.)
```

---

## 6. Planning

**User request:** "Move the game from one server to sharded multi-process."

**❌ What LLMs do (dive straight in)**

```
[Starts editing server/index.js, adds cluster.fork(), changes the room map,
touches 6 files, then realizes snapshots can't cross process boundaries]
```

**✅ What should happen (plan first)**

```
This is non-trivial (multiple files, architecture decision). Plan before code:

Goal: shard rooms across N worker processes, single public port.

Changes:
- server/gateway.js (new) — accepts WS, routes by room → worker
- server/worker.js — runs the sim for its assigned rooms
- server/index.js — becomes the cluster bootstrap
- shared/sim.js — unchanged (rooms are already independent)

Open questions:
1. How are rooms assigned — hash(roomId) or least-loaded?
2. Does the admin panel need an aggregated view across workers?
3. Sticky routing on reconnect — required, or can a player land on any worker?

Validate this before I touch anything?
```

---

## 7. Verification

**User request:** "Add a `/mass` command that sets a player's mass."

**❌ What LLMs do (mark done untested)**

```
Done — added the /mass command to runCommand(). It sets cell.mass to the
given value.
```

**✅ What should happen (prove it)**

```
Added /mass to runCommand(). Verified before calling it done:

- Ran the server, sent `/mass 500` → player cell grew, snapshot shows mass=500
- Edge case `/mass 0` → player dies cleanly, no NaN in the sim loop
- Edge case `/mass -10` → clamped to 1, no crash (added the clamp)
- Existing command tests still pass

Logs and the snapshot diff are above.
```

---

## 8. Errors

**Scenario:** CI fails with `TypeError: Cannot read properties of undefined (reading 'mass')`.

**❌ What LLMs do (hand it back)**

```
The test is failing with a TypeError. This looks like it's related to the
sim state. Can you tell me what `getLeader()` is supposed to return when the
room is empty?
```

**✅ What should happen (diagnose and fix)**

```
Traced it: getLeader() does `players[0].mass` but returns undefined for an
empty room, so the spectator camera crashes on an empty arcade room between
matches.

Root cause: no empty-room guard. Fix: return null from getLeader() when there
are no players, and the camera already handles null (free-look). Re-ran CI —
green.
```

---

## 9. Elegance

**User request:** "Clamp the player to the arena bounds."

**❌ What LLMs do (rushed)**

```js
if (cell.x < 0) cell.x = 0;
if (cell.x > W) cell.x = W;
if (cell.y < 0) cell.y = 0;
if (cell.y > H) cell.y = H;
if (cell.x < r) cell.x = r;
if (cell.x > W - r) cell.x = W - r;
if (cell.y < r) cell.y = r;
if (cell.y > H - r) cell.y = H - r;
```

**✅ What should happen (pause, simplify)**

```js
const clamp = (v, lo, hi) => Math.max(lo, Math.min(hi, v));
cell.x = clamp(cell.x, cell.r, W - cell.r);
cell.y = clamp(cell.y, cell.r, H - cell.r);
```

Same behavior, half the lines, and the radius-aware bound is now obvious.

---

## 10. Model Optimization

**Scenario:** On Opus, the user asks "what does this regex `/^ws:\/\//` match?"

**❌ What LLMs do (silently burn tokens)**

```
[Answers the trivial question on Opus, then continues the whole session on
Opus for more simple lookups]
```

**✅ What should happen (flag it once)**

```
It matches strings starting with `ws://` (a plaintext WebSocket URL).

By the way — this kind of quick question doesn't need Opus. Switch with
`/model sonnet` to save tokens; I'll only suggest Opus again if something
needs deep multi-file reasoning.
```

---

---

## Anti-patterns summary

| Rule | Anti-pattern | Fix |
|------|-------------|-----|
| Think Before Coding | Silently assumes auth, scope, rejoin behavior | List assumptions, ask before building |
| Simplicity First | Strategy classes for `score + bonus` | One function until complexity is real |
| Surgical Changes | Adds validation + logging while fixing a crash | Only the lines that fix the reported bug |
| Goal-Driven | "I'll review and improve it" | "Reproduce with a test → make it pass → no regressions" |
| Honesty | Confident wrong answer about `ws.send()` | State what's known, flag what needs checking |
| Planning | Edits 6 files then hits an architecture wall | Goal + files + open questions, validated first |
| Verification | "Done" with no run | Show the run, the edge cases, the logs |
| Errors | Hands the failing test back to the user | Trace root cause, fix, re-run CI |
| Elegance | Eight `if` clamps | One `clamp()` helper |
| Model Optimization | Trivial Q&A on Opus all session | Flag once, switch to a cheaper model |

## Key insight

The overcomplicated examples aren't *obviously* wrong — they follow real patterns. The problem is **timing and scope**: complexity added before it's needed, changes wider than the request, answers more confident than the evidence. The disciplined versions ship less, say less, and assume less — and are easier to verify because of it.
