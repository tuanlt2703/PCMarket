package lookup

import (
	"sort"
	"strings"
)

// tables lists every lookup table exposed by the Lookup API. component_types and
// manufacturers have the same shape but belong to the Catalog module.
var tables = []string{
	// Shared across component types.
	"cpu_sockets",
	"memory_types",
	"pcie_generations",
	"motherboard_form_factors",
	"display_connector_types",
	"power_connector_types",
	"electrical_connector_types",
	"product_colors",
	"lighting_types",
	"io_port_types",
	"usb_port_types",
	"fan_sizes",
	"radiator_sizes",
	"fan_bearing_types",
	"wireless_standards",
	"bluetooth_versions",
	"encryption_standards",
	"cooling_types",
	// CPU
	"cpu_architectures",
	"cpu_series",
	"cpu_instruction_sets",
	// GPU
	"gpu_chipsets",
	"gpu_architectures",
	"vram_types",
	"graphics_api_versions",
	"upscaling_technologies",
	// Motherboard
	"motherboard_chipsets",
	"audio_codecs",
	// RAM
	"memory_module_types",
	"ram_heat_spreader_types",
	"ram_profile_types",
	// Storage
	"storage_device_types",
	"storage_form_factors",
	"storage_interfaces",
	"storage_protocols",
	"nand_types",
	// PSU
	"psu_efficiency_certifications",
	"psu_standards",
	"psu_form_factors",
	"psu_modular_types",
	"psu_protection_types",
	// Case
	"case_types",
	"case_panel_materials",
	"drive_bay_types",
}

type Registry struct {
	bySlug map[string]Table
	slugs  []string
}

// NewRegistry exposes each table under its kebab-case slug, e.g. cpu_sockets -> cpu-sockets.
func NewRegistry() *Registry {
	r := &Registry{bySlug: make(map[string]Table, len(tables))}
	for _, name := range tables {
		slug := strings.ReplaceAll(name, "_", "-")
		r.bySlug[slug] = Table{Name: name, Label: strings.ReplaceAll(name, "_", " ")}
		r.slugs = append(r.slugs, slug)
	}
	sort.Strings(r.slugs)
	return r
}

func (r *Registry) Get(slug string) (Table, bool) {
	t, ok := r.bySlug[slug]
	return t, ok
}

func (r *Registry) Slugs() []string { return r.slugs }
