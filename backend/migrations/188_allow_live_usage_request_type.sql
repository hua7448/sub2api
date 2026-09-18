-- Existing rows already satisfy the previous, narrower 0..4 constraint.
-- NOT VALID keeps this metadata-only on large production tables while still
-- enforcing the expanded constraint for every new or updated row.
ALTER TABLE usage_logs
    DROP CONSTRAINT IF EXISTS usage_logs_request_type_check;

ALTER TABLE usage_logs
    ADD CONSTRAINT usage_logs_request_type_check
    CHECK (request_type >= 0 AND request_type <= 5) NOT VALID;
