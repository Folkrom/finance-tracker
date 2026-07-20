-- Backfill the catalogue-default "Cash" payment method for existing users.
-- New users get it seeded by the profile middleware on first request.
INSERT INTO payment_methods (id, user_id, name, type, created_at, updated_at)
SELECT gen_random_uuid(), p.user_id, 'Cash', 'cash', now(), now()
FROM profiles p
WHERE NOT EXISTS (
  SELECT 1 FROM payment_methods pm
  WHERE pm.user_id = p.user_id AND pm.type = 'cash'
);
