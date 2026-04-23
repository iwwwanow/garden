# Planning: Phase 0 — Zig + raylib + Sciter

## Context

Building a 3D fantasy farm sandbox game (lo-poly, cel-shading aesthetic).  
Stack: Zig + raylib (3D rendering) + Sciter (HTML/CSS UI overlay).  
Target platforms: Desktop (native) + Web (WASM via Emscripten).  
Backend: Go (separate, not in scope for this phase).

## Goal of this phase

Prove the core tech stack works end-to-end:

- Zig can call raylib via C bindings
- 3D scene renders in a window
- Sciter renders HTML/CSS panel on top of the 3D scene
- Input (mouse/keyboard) flows correctly between raylib and Sciter
- Sciter can call Zig functions (UI event → game action)
- Same Zig + raylib code compiles to WASM for browser

## Constraints

- Zig version: 0.13.0 (stable, not nightly) — stored in `.zigversion`
- raylib version: 5.x (latest stable)
- Sciter: use C API only (no C++ wrappers) — called via `@cImport` in Zig
- No game logic in this phase — only tech validation
- No Go backend in this phase — all local, no networking
- Keep it simple: one window, one cube, one HTML panel

## Project structure

```
/client-3d
  /src
    main.zig          # entry point
    renderer.zig      # raylib draw calls
    ui.zig            # Sciter integration
    input.zig         # keyboard/mouse handling
  /ui
    hud.html          # sample HTML panel
    hud.css           # styles
  /libs
    raylib/           # raylib C source or prebuilt
    sciter/           # Sciter SDK (headers + dylib)
  build.zig           # build configuration
  .zigversion         # "0.13.0"
```

## Tasks

### Task 1: Project scaffold

**Goal:** compiling Zig project with raylib linked  
**Steps:**

1. Create `build.zig` that links raylib as C library via `@cImport`
2. Write minimal `main.zig` that opens a 640x480 window and closes on ESC
3. Confirm: `zig build run` opens and closes window with no errors

**Success criteria:** window opens, no crash, clean exit

---

### Task 2: 3D scene

**Goal:** render a lo-poly cube in 3D with camera  
**Steps:**

1. Initialize raylib 3D camera (perspective, orbiting)
2. Draw a `DrawCube` at origin
3. Draw a ground plane (`DrawPlane`)
4. Add basic ambient + directional light
5. Camera orbits with mouse drag (raylib `UpdateCamera`)

**Success criteria:** visible 3D cube on a ground plane, camera moves with mouse

---

### Task 3: Custom GLSL cel-shading shader

**Goal:** apply cel-shading to the cube  
**Steps:**

1. Write `cel.vert` and `cel.frag` GLSL shaders
2. Load shader via `raylib.LoadShader`
3. Apply to cube material
4. Tune: 3 light bands (dark / mid / bright), black outline via backface expansion

**Success criteria:** cube has visible cartoon shading with outline

---

### Task 4: Sciter integration

**Goal:** render an HTML panel on top of the 3D scene  
**Steps:**

1. Add Sciter SDK headers to `/libs/sciter/`
2. In `build.zig`: link Sciter dylib, add include path
3. In `ui.zig`: initialize Sciter, load `ui/hud.html` into an offscreen surface
4. Each frame: call `SciterRender` → get texture → `DrawTexture` via raylib
5. Position panel in top-left corner (300x200px)

**Sciter C API calls needed:**

```c
SciterSetOption(NULL, SCITER_SET_SCRIPT_RUNTIME_FEATURES, ...);
SciterCreateWindow(...)
SciterLoadFile(hwnd, L"ui/hud.html")
SciterRender(hwnd, &bitmap)
```

**Success criteria:** HTML panel visible over 3D scene, styles applied correctly

---

### Task 5: Input routing

**Goal:** mouse and keyboard events reach the correct handler  
**Steps:**

1. Each frame: check if mouse cursor is over Sciter panel rect
2. If yes → forward `WM_MOUSEMOVE`, `WM_LBUTTONDOWN` to Sciter via `SciterProcND`
3. If no → raylib handles mouse (camera orbit)
4. Keyboard: ESC always goes to raylib (close window), other keys routed by focus state

**Success criteria:** clicking HTML button works, camera still orbits when clicking outside panel

---

### Task 6: Sciter → Zig callback

**Goal:** HTML button click triggers action in Zig game code  
**Steps:**

1. Register a native function handler via `SciterSetCallback`
2. In `hud.html`: add button `<button id="test">Hello from HTML</button>`
3. In JS inside HTML: on click → call `sciter.call("onButtonClick")`
4. In `ui.zig`: implement `onButtonClick` — print to stdout or change cube color

**Success criteria:** clicking HTML button changes cube color in 3D scene

---

### Task 7: WASM build

**Goal:** same Zig + raylib code runs in browser  
**Steps:**

1. Install Emscripten SDK (emcc), add to PATH
2. Add WASM target to `build.zig`:
   ```zig
   const target = b.resolveTargetQuery(.{
       .cpu_arch = .wasm32,
       .os_tag = .emscripten,
   });
   ```
3. Compile: `zig build -Dtarget=wasm32-emscripten`
4. Serve `game.html` via local HTTP server
5. Note: Sciter is NOT available in WASM — UI panel is skipped for this build
6. Confirm: 3D cube with cel-shading renders in browser canvas

**Success criteria:** cube visible in Chrome/Firefox, no console errors

---

## Notes for agent

- Sciter license: free for non-commercial. For commercial use ~$800 one-time. Plan accordingly.
- Sciter in WASM: not supported. Web build uses DOM overlay instead of Sciter.
- raylib-zig bindings: use `@cImport(@cInclude("raylib.h"))` directly — do not use third-party Zig wrappers, they may be outdated for Zig 0.13.0.
- If Sciter integration proves too complex, fallback option is `dear imgui` via `rlimgui` — but HTML/CSS kit will not transfer.
- All hardcoded values (window size, panel position, cube color) should be constants at top of file — they will be extracted to world config in a later phase.
- Do not add game logic, networking, or asset loading in this phase.

## Definition of done

- [ ] `zig build run` → window with cel-shaded cube + HTML panel overlay
- [ ] Mouse correctly routes between Sciter and raylib camera
- [ ] HTML button click changes something in the 3D scene
- [ ] `zig build -Dtarget=wasm32-emscripten` → cube renders in browser
- [ ] No memory leaks (test with `-Doptimize=ReleaseSafe`)
- [ ] Code is split into `renderer.zig`, `ui.zig`, `input.zig` — no spaghetti in `main.zig`
