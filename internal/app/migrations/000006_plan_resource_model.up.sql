ALTER TABLE plan.resource_packages
  ADD COLUMN IF NOT EXISTS resource_model TEXT;

CREATE INDEX IF NOT EXISTS idx_plan_resource_packages_resource_model
  ON plan.resource_packages (resource_model);
