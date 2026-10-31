#!/bin/sh
# Renders the animated grid cells gbemu's gbart grid reads, for every visual
# (or the ones named): inputs/dmg/<v>/cells/cell-NNNN/frame-KK.png from
# inputs/dmg/<v>-play.json and inputs/dmg/<v>/grid.json.
#
# SET=blog renders the blog's 36-combo set instead, from the same play
# configs into inputs/dmg/blog/<v>/cells.
set -e
cd "$(dirname "$0")/../.."
dir=inputs/dmg${SET:+/$SET}
go build -o bin/gridcells ./cmd/gridcells
for v in ${@:-topo ripples blobs truchet planet tunnel voronoi}; do
	rm -rf "$dir/$v/cells"
	bin/gridcells -config "inputs/dmg/$v-play.json" -grid "$dir/$v/grid.json" -out "$dir/$v/cells"
done
