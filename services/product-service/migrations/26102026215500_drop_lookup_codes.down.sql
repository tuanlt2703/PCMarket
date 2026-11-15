DO $$
DECLARE
    table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'component_types', 'manufacturers', 'cpu_sockets', 'memory_types',
        'pcie_generations', 'motherboard_form_factors', 'display_connector_types',
        'power_connector_types', 'electrical_connector_types', 'product_colors',
        'lighting_types', 'cpu_architectures', 'cpu_series', 'cpu_instruction_sets',
        'gpu_chipsets', 'gpu_architectures', 'vram_types', 'cooling_types',
        'graphics_api_versions', 'upscaling_technologies', 'motherboard_chipsets',
        'memory_module_types', 'ram_heat_spreader_types', 'storage_device_types',
        'storage_form_factors', 'storage_interfaces', 'storage_protocols', 'nand_types',
        'encryption_standards', 'psu_efficiency_certifications', 'psu_standards',
        'psu_form_factors', 'psu_modular_types', 'psu_protection_types', 'case_types',
        'case_panel_materials', 'drive_bay_types', 'fan_sizes', 'radiator_sizes',
        'io_port_types', 'fan_bearing_types', 'wireless_standards', 'bluetooth_versions',
        'audio_codecs', 'usb_port_types', 'ram_profile_types'
    ] LOOP
        EXECUTE format('ALTER TABLE %I ADD COLUMN code VARCHAR(50)', table_name);
        EXECUTE format('UPDATE %I SET code = id::name', table_name);
        EXECUTE format('ALTER TABLE %I ALTER COLUMN code SET NOT NULL, ADD UNIQUE (code)', table_name);
    END LOOP;
END $$;