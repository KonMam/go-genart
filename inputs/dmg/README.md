# DMG set

Seven engines made for the Game Boy via gbemu's `gbart`: bold shapes that
survive 160x144 in 4 shades, and parameters that move the picture a little
per step. Every parameter in the grids below is continuous, each engine
draws the same random numbers whatever the parameters, and nothing is
simulated, so a step never re-lays out the picture.

Per engine: `<name>.json` (still), `<name>-grid.json` (gbart grid spec,
2 params x 12 steps = 144 cells, so seven fit one MBC5 ROM) and
`<name>-loop.json` (seamless loop: `fixed_seed`, `loop`, PNG frames).

    gbart grid -spec inputs/dmg/topo-grid.json -keep out/cells -out out/topo.gb -capture
    go run ./cmd/genart -config inputs/dmg/topo-loop.json   # -> outputs/dmg/topo-loop/frame-NNNN.png
    gbart capture -in outputs/dmg/topo-loop -mode tiles -fps 12 -out out/topo.gb

HUD: every axis shows a plain name for what it does (`label`) and a
12-step bar, `#-----------` to `############` (`labels`), so all seven read
on one scale; step 1 is the start, right always means more. The table
below uses the engine parameter names.

Smoothness: mean (max) % of screen pixels whose shade changes per step,
after gbart's conversion (shared levels, Bayer), over all 132 neighbouring
pairs per axis; loops: per frame, including the seam back to frame 0.

| engine  | axis 1             | axis 2               | loop (frames)         |
|---------|--------------------|----------------------|-----------------------|
| topo    | phase 6.3 (7.9)    | offset 9.2 (10.3)    | phase 8.9, seam 10.3 (120) |
| ripples | phase 5.2 (5.6)    | drift 6.4 (7.2)      | phase x2 19.2, seam 18.6 (60) |
| blobs   | phase 5.2 (8.0)    | radius 4.6 (5.9)     | phase 10.6, seam 10.5 (120) |
| truchet | phase 4.7 (6.8)    | lw 6.4 (7.1)         | phase 5.5, seam 5.0 (120) |
| planet  | spin 3.6 (7.3)     | light 4.9 (7.7)      | spin 1.6, seam 1.7 (120) |
| tunnel  | depth 6.3 (10.1)   | turn 6.4 (6.8)       | depth 4.1, seam 5.0 (48) |
| voronoi | phase 6.3 (9.2)    | wander 6.7 (7.8)     | phase 12.1, seam 11.5 (120) |

Most steps move every edge on screen by about a pixel; fewer % needs
narrower ranges (less change across the axis) or more steps.
