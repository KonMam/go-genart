#!/bin/sh
# Renders the animated grid cells gbemu's gbart grid reads, for every visual
# (or the ones named): inputs/dmg/<v>/cells/cell-NNNN/frame-KK.png from
# inputs/dmg/<v>-play.json and inputs/dmg/<v>/grid.json.
set -e
cd "$(dirname "$0")/../.."
go build -o bin/gridcells ./cmd/gridcells
for v in ${@:-topo ripples blobs truchet planet tunnel voronoi}; do
	rm -rf "inputs/dmg/$v/cells"
	bin/gridcells -config "inputs/dmg/$v-play.json" -grid "inputs/dmg/$v/grid.json" -out "inputs/dmg/$v/cells"
done
