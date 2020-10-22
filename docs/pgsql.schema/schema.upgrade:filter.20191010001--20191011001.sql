BEGIN;
  ALTER TABLE filter.confidentiality_mapping ADD COLUMN isID uuid NULL;
  ALTER TABLE filter.confidentiality_mapping ADD CONSTRAINT __fk_flcm_isID FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping ADD COLUMN componentID uuid NULL;
  ALTER TABLE filter.confidentiality_mapping ADD CONSTRAINT __fk_flcm_compID FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping ADD COLUMN subgroupID uuid NULL;
  ALTER TABLE filter.confidentiality_mapping ADD CONSTRAINT __fk_flcm_groupID FOREIGN KEY ( subgroupID ) REFERENCES ix.logical_component_subgroup ( subgroupID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping ALTER CONSTRAINT __fk_flcm_confID DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping ALTER CONSTRAINT __fk_flcm_tlsID DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping ALTER CONSTRAINT __fk_flcm_prodID DEFERRABLE;
  ALTER TABLE filter.confidentiality_mapping DROP CONSTRAINT __flbm_uq_object;
  ALTER TABLE filter.confidentiality_mapping ADD CONSTRAINT __flbm_uq_object CHECK (
                 ((tlsID IS NOT NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NOT NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NOT NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NOT NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NOT NULL)));
  ALTER TABLE filter.integrity_mapping ALTER COLUMN tlsID DROP NOT NULL;
  ALTER TABLE filter.integrity_mapping ALTER CONSTRAINT __fk_flim_tlsID DEFERRABLE;
  ALTER TABLE filter.integrity_mapping DROP CONSTRAINT __fk_flim_confID;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __fk_flim_intID FOREIGN KEY ( integrityID ) REFERENCES filter.integrity ( integrityID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.integrity_mapping ADD COLUMN productID uuid NULL;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __fk_flim_prodID FOREIGN KEY ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.integrity_mapping ADD COLUMN isID uuid NULL;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __fk_flim_isID FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.integrity_mapping ADD COLUMN componentID uuid NULL;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __fk_flim_compID FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.integrity_mapping ADD COLUMN subgroupID uuid NULL;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __fk_flim_grpID FOREIGN KEY ( subgroupID ) REFERENCES ix.logical_component_subgroup ( subgroupID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.integrity_mapping ADD CONSTRAINT __flim_uq_object CHECK (
                 ((tlsID IS NOT NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NOT NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NOT NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NOT NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NOT NULL)));
  ALTER TABLE filter.availability_mapping ALTER COLUMN tlsID DROP NOT NULL;
  ALTER TABLE filter.availability_mapping ALTER CONSTRAINT __fk_flam_tlsID DEFERRABLE;
  ALTER TABLE filter.availability_mapping DROP CONSTRAINT __fk_flam_confID;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __fk_flam_availID FOREIGN KEY ( availabilityID ) REFERENCES filter.availability ( availabilityID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.availability_mapping ADD COLUMN productID uuid NULL;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __fk_flam_prodID FOREIGN KEY ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.availability_mapping ADD COLUMN isID uuid NULL;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __fk_flam_isID FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.availability_mapping ADD COLUMN componentID uuid NULL;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __fk_flam_compID FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.availability_mapping ADD COLUMN subgroupID uuid NULL;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __fk_flam_grpID FOREIGN KEY ( subgroupID ) REFERENCES ix.logical_component_subgroup ( subgroupID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.availability_mapping ADD CONSTRAINT __flam_uq_object CHECK (
                 ((tlsID IS NOT NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NOT NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NOT NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NOT NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NOT NULL)));
  ALTER TABLE filter.authenticity_mapping ALTER COLUMN tlsID DROP NOT NULL;
  ALTER TABLE filter.authenticity_mapping ALTER CONSTRAINT __fk_flum_tlsID DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping DROP CONSTRAINT __fk_flum_confID;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __fk_flum_authID FOREIGN KEY ( authenticityID ) REFERENCES filter.authenticity ( authenticityID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping ADD COLUMN productID uuid NULL;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __fk_flum_prodID FOREIGN KEY ( productID ) REFERENCES ix.product ( productID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping ADD COLUMN isID uuid NULL;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __fk_flum_isID FOREIGN KEY ( isID ) REFERENCES ix.information_system ( isID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping ADD COLUMN componentID uuid NULL;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __fk_flum_compID FOREIGN KEY ( componentID ) REFERENCES ix.functional_component ( componentID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping ADD COLUMN subgroupID uuid NULL;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __fk_flum_grpID FOREIGN KEY ( subgroupID ) REFERENCES ix.logical_component_subgroup ( subgroupID ) ON DELETE RESTRICT DEFERRABLE;
  ALTER TABLE filter.authenticity_mapping ADD CONSTRAINT __flum_uq_object CHECK (
                 ((tlsID IS NOT NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NOT NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NOT NULL) AND (componentID IS NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NOT NULL) AND (subgroupID IS NULL))
              OR ((tlsID IS NULL) AND (productID IS NULL) AND (isID IS NULL) AND (componentID IS NULL) AND (subgroupID IS NOT NULL)));
  INSERT INTO public.schema_versions ( schema, version, description )
              VALUES ( 'filter', 20191011001, 'CIAA mapping update');
COMMIT;
