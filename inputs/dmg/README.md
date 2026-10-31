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

## Animated grid cells (for gbart grid)

Per visual: `inputs/dmg/<v>/grid.json` (`name`, `fps`, `params` with
`values`, `label`, `labels`; 72 combos each, so 7 visuals fit 504 of MBC5's 511 banks) and
`inputs/dmg/<v>/cells/cell-NNNN/frame-KK.png`: 24 frames per cell, 320x288,
NNNN = combo index with the last param varying fastest, frame 0 follows
frame 23 seamlessly (checked on every cell). Cells are gitignored; render
them with

    inputs/dmg/cells.sh [visual...]   # cmd/gridcells over <v>-play.json + <v>/grid.json

Tiles-mode quantization per cell (gbart build, % pixels changed, mean /
worst): blobs 0.6/1.5, planet 0.4/0.9, tunnel 1.7/3.5, topo 2.2/9.4,
ripples 4.9/11.7, voronoi 5.2/12.4, truchet 5.4/8.6. Detail density drives
it (cells, rings, zoom x terraces), so those knobs are capped; line width
barely matters.

`variant` (labels A/B/C, planet A-D) is the random layout, added to the seed by the
runner; every other knob keeps the layout and changes one aspect.

## Blog gallery set (36 combos)

`blog/<v>/grid.json` is the set behind the gbart gallery ROM on mamonas.dev
(personal-blog `content/posts/gbemu-player-demo/gallery.gb.gz`): each
visual's 72-combo grid with one value dropped from two knobs, so 36 combos
each and 253 banks for all seven (4 MiB, 1.32 MB gzip -9) instead of 505
(8 MiB, 2.68 MB). The densest values go first, which also lowers the worst
tiles-mode quantization (voronoi 13.0% -> 8.3%, topo 9.9% -> 6.2%). Bars are
relabelled to the new length.

| visual  | dropped              |
|---------|----------------------|
| topo    | terraces 14, variant C |
| ripples | sources 4, variant C |
| blobs   | blobs 4, variant C   |
| truchet | tiles 12, variant C  |
| planet  | variant C and D      |
| tunnel  | pattern spokes, segments 10 |
| voronoi | cells many, variant C |

    SET=blog inputs/dmg/cells.sh       # -> inputs/dmg/blog/<v>/cells
    gbart gallery -spec inputs/dmg/blog/gallery.json -out out/gallery.gb   # in gbemu

The cells render from the same `<v>-play.json`, so they match the 72-combo
cells they were cut from, and the ROM rebuilds byte for byte.

## Play mode (always moving, big knobs)

`<name>-play.json` is a 24-frame, 2 s seamless loop (12 fps) that keeps
playing while idle; `<name>/grid.json` has 4 knobs of 3 values each
(81 combos; numbers below were measured on an earlier 3-4 value set) whose steps change the look a lot. Every combo is a loop;
a knob step should land on the same frame index of the new combo's loop.
`variant` is the random layout (added to the seed by the runner), the one
knob meant to jump. Needs gbart grid to play loops per cell (not yet).

Measured on DMG previews (% of pixels changed; frame = idle motion per
frame, mean/max; steps = mean change per knob step):

| engine  | frame      | knob steps                                              |
|---------|------------|---------------------------------------------------------|
| topo    | 7.3 / 16.2 | zoom 64, terraces 74, style 52, variant 67              |
| ripples | 12.8 / 23.3| sources 58, rings 75, edge 43, variant 76               |
| blobs   | 10.4 / 23.5| blobs 38, size 48, halo 30, variant 45                  |
| truchet | 6.9 / 9.4  | tiles 61, weight 38, order 40, variant 51               |
| planet  | 3.3 / 5.3  | sea 8, land 10, clouds 12, variant 16 (planet is ~half the screen) |
| tunnel  | 17.1 / 32.4| pattern 46, segments 46, bend 49, fog 64                |
| voronoi | 10.9 / 24.4| cells 57, shading 31, lead 26, variant 58               |
