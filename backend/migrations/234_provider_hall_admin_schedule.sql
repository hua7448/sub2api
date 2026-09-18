-- Preserve the scheduling behavior of existing installations only.
ALTER TABLE provider_hall_config ADD COLUMN auto_schedule_enabled boolean NOT NULL DEFAULT false;
UPDATE provider_hall_config SET auto_schedule_enabled = tasks_enabled;
ALTER TABLE provider_hall_targets ADD COLUMN auto_schedule_enabled boolean NOT NULL DEFAULT false;
UPDATE provider_hall_targets SET auto_schedule_enabled = true;
