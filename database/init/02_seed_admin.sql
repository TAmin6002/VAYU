-- Default admin account. Login: username "admin" / password "12345678".
-- The password is stored as a bcrypt hash, never in plain text.
-- CHANGE THIS PASSWORD before deploying anywhere public.
INSERT INTO admins (full_name, email, password_hash)
VALUES ('admin', 'admin@gmail.com', '$2a$10$5uolVTgH5Yd0huI1Gz5gBeOdj.qObBt.HavCCozED1DpsWn3HSuFy')
ON CONFLICT (email) DO NOTHING;
