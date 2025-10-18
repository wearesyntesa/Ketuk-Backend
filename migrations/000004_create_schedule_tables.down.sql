-- Drop triggers
DROP TRIGGER IF EXISTS trigger_schedule_from_ticket_updated_at ON schedule_ticket;
DROP TRIGGER IF EXISTS trigger_semester_unblocking_updated_at ON unblocking;
DROP TRIGGER IF EXISTS trigger_schedule_reguler_updated_at ON schedule_reguler;

-- Drop function
DROP FUNCTION IF EXISTS update_schedule_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_schedule_from_ticket_ticket_id;
DROP INDEX IF EXISTS idx_schedule_from_ticket_user_id;
DROP INDEX IF EXISTS idx_schedule_from_ticket_start_date;
DROP INDEX IF EXISTS idx_schedule_from_ticket_end_date;

DROP INDEX IF EXISTS idx_semester_unblocking_user_id;
DROP INDEX IF EXISTS idx_semester_unblocking_tahun;
DROP INDEX IF EXISTS idx_semester_unblocking_semester;

DROP INDEX IF EXISTS idx_schedule_reguler_user_id;
DROP INDEX IF EXISTS idx_schedule_reguler_start_date;
DROP INDEX IF EXISTS idx_schedule_reguler_end_date;

-- Drop tables
DROP TABLE IF EXISTS schedule_ticket;
DROP TABLE IF EXISTS unblocking;
DROP TABLE IF EXISTS schedule_reguler;