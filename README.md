# 🛑 longblock (debug module)

A Cosmos SDK **debug** module that enables pausing and resuming block processing in a running chain — useful for debugging, testing, or controlled maintenance.

[![Tests Status](https://github.com/hmoragrega/longblock/actions/workflows/test.yml/badge.svg)](https://github.com/hmoragrega/longblock/actions)
[![Security](https://github.com/hmoragrega/longblock/actions/workflows/gosec.yml/badge.svg)](https://github.com/hmoragrega/longblock/actions)
[![Analysis](https://github.com/hmoragrega/longblock/actions/workflows/codeql.yml/badge.svg)](https://github.com/hmoragrega/longblock/actions)
[![Lint](https://github.com/hmoragrega/longblock/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/hmoragrega/longblock/actions)

---

## ✅ Features

- **Pause / Resume** via `BeginBlock` hook
- Controlled via CLI flags:
    - `--debug.pause-allowed`: enables pause/resume behavior
    - `--debug.pause-on-each-block`: pause after every block
    - `--debug.pause-skip N`: number of initial blocks to skip before pausing

---

## 🏗️ Installation

```go
import (
  debug "github.com/hmoragrega/longblock/debug"
  // ...
)

// In your app constructor
appDebug := debug.NewAppModuleFromOpts(appOpts)

app.mm = module.NewManager(
  // other modules...
  appDebug,
)

app.mm.SetOrderBeginBlockers(
  debug.ModuleName,
  // other modules...
)
```

## ⚙️ Configuration
Add the module into your app's module manager and set it to run in the `BeginBlock` phase.

If you use cobra for CLI, register the flags in your root command:

```go
bebug.AddModuleInitFlags(cobraCli)
```

Compile with `--tags debug` to enable the debug module.

Then when launching your node configure the debug behaviour using CLI flags:

### 🚩 Flags
* `debug.pause-allowed`: enables the pause/resume functionality
* `debug.pause-on-each-block`: pauses after every block
* `debug.pause-skip N`: skips the first N blocks before pausing

```bash
longblocked start \
--debug.pause-allowed=true \ 
--debug.pause-on-each-block=true \
--debug.pause-skip=5
```

**NOTE:** Without `--debug.pause-allowed`, the module does nothing.

## 🔁 Runtime Behavior
* `BeginBlock` hook calls `PauseService.HoldIfPaused(ctx)` to pause if conditions are met.
* The pauser package implements the pause logic; a no-op version is used when disabled.


## 🧩 gRPC Query Support
To pause/resume you have to use the `PauseService` methods:

* `hmoragrega.longblock.debug.v1.Pause`: pauses block processing
* `hmoragrega.longblock.debug.v1.Resume`: resumes block processing

(As defined in types/query.proto and registered via RegisterServices)

## 🚀 Use-Cases
* Local debugging: step through block execution manually
* Testnets or forked chains: safely pause during setup
* Maintenance: freeze the chain during off-chain operations, then resume


## 💡 Tips
* Register this module first in BeginBlockers for immediate effect
* Add CLI or gRPC commands for Pause() and Resume()
* In tests, use `NewNoOpPauser` to avoid halting logic

## 🎯 Summary
longblock/debug offers a simple yet powerful mechanism to pause block processing, 
enhancing visibility and control for developers and operators.