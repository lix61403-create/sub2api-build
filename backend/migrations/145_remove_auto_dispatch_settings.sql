-- Remove the retired auto-dispatch settings key.
-- This only clears the standalone auto-dispatch configuration row and leaves
-- normal account scheduling, stream-timeout, and other system settings intact.
DELETE FROM settings
WHERE key = 'auto_dispatch_settings';
