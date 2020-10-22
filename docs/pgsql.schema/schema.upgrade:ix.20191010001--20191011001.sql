BEGIN;
  ALTER TABLE ix.information_system_component RENAME TO functional_component;
  ALTER INDEX __pk_ixisc RENAME TO __pk_ixfc;
  ALTER TABLE ix.functional_component RENAME CONSTRAINT __fk_ixisc_dictID TO __fk_ixfc_dictID;
  ALTER TABLE ix.functional_component RENAME CONSTRAINT __ixisc_fk_origin TO __ixfc_fk_origin;

  ALTER TABLE ix.information_system_component_attribute RENAME TO functional_component_attr_values;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __fk_ixisca_compID TO __fk_ixfcav_compID;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __fk_ixisca_attrID TO __fk_ixfcav_attrID;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __fk_ixisca_dictID TO __fk_ixfcav_dictID;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __fk_ixisca_uq_dct TO __fk_ixfcav_uq_dct;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __fk_ixisca_uq_att TO __fk_ixfcav_uq_att;
  ALTER TABLE ix.functional_component_attr_values RENAME CONSTRAINT __ixisca_temporal TO __ixfcav_temporal;

  ALTER TABLE ix.information_system_component_unique RENAME TO functional_component_attr_uniq_values;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __fk_ixiscq_compID TO __fk_ixfcqv_compID;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __fk_ixiscq_attrID TO __fk_ixfcqv_attrID;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __fk_ixiscq_dictID TO __fk_ixfcqv_dictID;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __fk_ixiscq_uq_dct TO __fk_ixfcqv_uq_dct;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __fk_ixiscq_uq_att TO __fk_ixfcqv_uq_att;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __ixiscq_temporal TO __ixfcqv_temporal;
  ALTER TABLE ix.functional_component_attr_uniq_values RENAME CONSTRAINT __ixiscq_temp_uniq TO __ixfcqv_temp_uniq;

  ALTER TABLE ix.mapping_information_system_component RENAME TO mapping_functional_component;
  ALTER TABLE ix.mapping_functional_component RENAME CONSTRAINT __fk_ixmisc_compID TO __fk_ixmfc_compID;
  ALTER TABLE ix.mapping_functional_component RENAME CONSTRAINT __fk_ixmisc_grpID TO __fk_ixmfc_grpID;
  ALTER TABLE ix.mapping_functional_component RENAME CONSTRAINT __ixmisc_temporal TO __ixmfc_temporal;
  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'ix', 20191011001, 'rename information_system_component to functional_component');
COMMIT;
