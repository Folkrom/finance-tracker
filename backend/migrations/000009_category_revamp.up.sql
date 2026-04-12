-- Step 1: Add is_system column
ALTER TABLE categories ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT false;

-- Step 2: Make user_id nullable
ALTER TABLE categories ALTER COLUMN user_id DROP NOT NULL;

-- Step 3: Drop old unique constraint
ALTER TABLE categories DROP CONSTRAINT categories_user_id_name_domain_key;

-- Step 4: Insert global categories (user_id = NULL)
-- Income
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Salary', 'income', 0, false),
  ('Bonus', 'income', 1, false),
  ('Freelance', 'income', 2, false),
  ('Dividends', 'income', 3, false),
  ('Interest', 'income', 4, false),
  ('Side Hustle', 'income', 5, false),
  ('Other', 'income', 6, true);

-- Expense
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Home Expenses', 'expense', 0, false),
  ('Eating Out', 'expense', 1, false),
  ('Self Care', 'expense', 2, false),
  ('Coffee/Drink', 'expense', 3, false),
  ('Entertainment', 'expense', 4, false),
  ('Transportation', 'expense', 5, false),
  ('Groceries', 'expense', 6, false),
  ('Utilities', 'expense', 7, false),
  ('Clothes', 'expense', 8, false),
  ('Card Payments', 'expense', 9, false),
  ('Savings/Investment', 'expense', 10, false),
  ('Taxes', 'expense', 11, false),
  ('Knowledge', 'expense', 12, false),
  ('Tech', 'expense', 13, false),
  ('Other', 'expense', 14, true);

-- Wishlist
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Electronics', 'wishlist', 0, false),
  ('Clothing', 'wishlist', 1, false),
  ('Home & Kitchen', 'wishlist', 2, false),
  ('Books & Media', 'wishlist', 3, false),
  ('Sports & Outdoors', 'wishlist', 4, false),
  ('Other', 'wishlist', 5, true);

-- Step 5: Reassign FK references from user categories to matching globals, then delete user dupes
DO $$
DECLARE
  user_cat RECORD;
  global_id UUID;
BEGIN
  FOR user_cat IN
    SELECT uc.id AS user_cat_id, uc.name, uc.domain, uc.user_id
    FROM categories uc
    WHERE uc.user_id IS NOT NULL
      AND EXISTS (
        SELECT 1 FROM categories gc
        WHERE gc.user_id IS NULL
          AND gc.name = uc.name
          AND gc.domain = uc.domain
      )
  LOOP
    SELECT gc.id INTO global_id
    FROM categories gc
    WHERE gc.user_id IS NULL
      AND gc.name = user_cat.name
      AND gc.domain = user_cat.domain;

    -- Update all FK references
    UPDATE incomes SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE expenses SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE debts SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE budgets SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE wishlist_items SET category_id = global_id WHERE category_id = user_cat.user_cat_id;

    -- Delete the user's duplicate
    DELETE FROM categories WHERE id = user_cat.user_cat_id;
  END LOOP;
END $$;

-- Step 6: Create partial unique indexes
CREATE UNIQUE INDEX idx_categories_global_unique
  ON categories (name, domain) WHERE user_id IS NULL;
CREATE UNIQUE INDEX idx_categories_user_unique
  ON categories (user_id, name, domain) WHERE user_id IS NOT NULL;
