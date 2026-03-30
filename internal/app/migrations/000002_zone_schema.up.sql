CREATE SCHEMA IF NOT EXISTS zone;

CREATE TABLE IF NOT EXISTS zone.zones (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS zone.zone_objects (
  object_type TEXT NOT NULL,
  object_id TEXT NOT NULL,
  zone_id UUID NOT NULL REFERENCES zone.zones(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (object_type, object_id)
);

CREATE INDEX IF NOT EXISTS idx_zone_zone_objects_zone_id
  ON zone.zone_objects (zone_id);
