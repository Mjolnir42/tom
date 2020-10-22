BEGIN;
  ALTER VIEW view.dictionary RENAME TO dictionary_details;
  ALTER FUNCTION view.dictionary_at RENAME TO dictionary_details_at;
  ALTER VIEW view.information_system RENAME TO information_system_details;
  ALTER FUNCTION view.information_system_at RENAME TO information_system_details_at;
  ALTER VIEW view.logical_component_subgroup RENAME TO logical_component_subgroup_details;
  ALTER FUNCTION view.logical_component_subgroup_at RENAME TO logical_component_subgroup_details_at;
  ALTER VIEW view.information_system_component RENAME TO functional_component_details;
  ALTER FUNCTION view.information_system_component_at RENAME TO functional_component_details_at;
  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'view', 20191011001, 'rename views so they show up in \d');
COMMIT;
