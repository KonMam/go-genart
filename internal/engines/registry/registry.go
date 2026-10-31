// Package registry maps engine names, as configs give them, to engines.
package registry

import (
	"genart/internal/core"
	"genart/internal/engines/blackhole"
	"genart/internal/engines/blobs"
	"genart/internal/engines/cells"
	"genart/internal/engines/contourlines"
	"genart/internal/engines/flow"
	"genart/internal/engines/flowfield"
	"genart/internal/engines/perlinpearls"
	"genart/internal/engines/planet"
	"genart/internal/engines/ripples"
	"genart/internal/engines/strata"
	"genart/internal/engines/swirl"
	"genart/internal/engines/topo"
	"genart/internal/engines/truchet"
	"genart/internal/engines/tunnel"
	"genart/internal/engines/voronoi"
)

// Engines returns every engine by name.
func Engines() map[string]core.Engine {
	return map[string]core.Engine{
		"flowfield":    flowfield.Engine{},
		"contourlines": contourlines.Engine{},
		"blackhole":    blackhole.Engine{},
		"perlinpearls": perlinpearls.Engine{},
		"swirl":        swirl.Engine{},
		"flow":         flow.Engine{},
		"strata":       strata.Engine{},
		"cells":        cells.Engine{},
		"blobs":        blobs.Engine{},
		"truchet":      truchet.Engine{},
		"planet":       planet.Engine{},
		"tunnel":       tunnel.Engine{},
		"voronoi":      voronoi.Engine{},
		"topo":         topo.Engine{},
		"ripples":      ripples.Engine{},
	}
}
