-- Drop triggers
DROP TRIGGER IF EXISTS trigger_items_updated_at ON items;
DROP TRIGGER IF EXISTS trigger_items_category_updated_at ON items_category;
DROP TRIGGER IF EXISTS trigger_update_category_quantity ON items;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_category_quantity();

-- Drop index
DROP INDEX IF EXISTS idx_items_category_id;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS items_category;