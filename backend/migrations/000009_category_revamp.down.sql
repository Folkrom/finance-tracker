-- Remove partial unique indexes
DROP INDEX IF EXISTS idx_categories_global_unique;
DROP INDEX IF EXISTS idx_categories_user_unique;

-- Delete all global categories
DELETE FROM categories WHERE user_id IS NULL;

-- Restore NOT NULL on user_id
ALTER TABLE categories ALTER COLUMN user_id SET NOT NULL;

-- Restore original unique constraint
ALTER TABLE categories ADD CONSTRAINT categories_user_id_name_domain_key UNIQUE (user_id, name, domain);

-- Drop is_system column
ALTER TABLE categories DROP COLUMN is_system;
