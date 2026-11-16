DROP INDEX idx_jobs_status_created
ON jobs;


ALTER TABLE jobs
DROP COLUMN locked_by,
DROP COLUMN locked_at,
DROP COLUMN failed_at;
