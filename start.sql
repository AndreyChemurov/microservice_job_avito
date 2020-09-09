-- name: create-user-table
CREATE TABLE IF NOT EXISTS "user_job" (
    "user_id" TEXT NOT NULL PRIMARY KEY UNIQUE
);

-- name: create-balance-table
CREATE TABLE IF NOT EXISTS "balance_job" (
    "balance_id" BIGSERIAL PRIMARY KEY NOT NULL,
    "user_id" TEXT REFERENCES user_job(user_id),
    "amount" NUMERIC(1000, 2) CHECK ("amount" >= 0.0) NOT NULL
);

--name: get-user-balance
SELECT amount FROM balance_job WHERE user_id = $1

--name: remittance-from
UPDATE balance_job SET amount = amount - $1 WHERE user_id = $2

--name: remittance-to
UPDATE balance_job SET amount = amount + $1 WHERE user_id = $2

--name: drop-user-table
DROP TABLE IF EXISTS user_job

--name: drop-balance-table
DROP TABLE IF EXISTS balance_job

--name: create-user
INSERT INTO user_job VALUES ($1)

--name: create-balance
INSERT INTO balance_job VALUES (DEFAULT, $1, 0)

--name: check-user-exists
SELECT user_id FROM user_job WHERE user_id = $1

--name: drop-balance
DELETE FROM balalce_job WHERE user_id = $1

--name: drop-user
DELETE FROM user_job WHERE user_id = $1