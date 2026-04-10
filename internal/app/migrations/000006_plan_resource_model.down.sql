DROP INDEX IF EXISTS idx_plan_resource_packages_resource_model;

ALTER TABLE plan.resource_packages
  DROP COLUMN IF EXISTS resource_model;
