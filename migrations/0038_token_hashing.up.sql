TRUNCATE otp, session;
ALTER TABLE otp ALTER COLUMN token TYPE bytea USING token::bytea;
ALTER TABLE session ALTER COLUMN token TYPE bytea USING token::bytea;
