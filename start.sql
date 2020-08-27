-- name: create-user-table
CREATE TABLE IF NOT EXISTS user_job (
    user_id BIGSERIAL PRIMARY KEY NOT NULL
);

-- name: create-balance-table
CREATE TABLE IF NOT EXISTS balance_job (
    balance_id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER REFERENCES user_job(user_id),
    amount NUMERIC NOT NULL
);

--name: get-user-balance
SELECT count FROM balance_job WHERE user_id = ?

--name: remittance-from
UPDATE balance_job SET count = count - ? WHERE user_id = ? AND count > 0

--name: remittance-to
UPDATE balance_job SET count = count + ? WHERE user_id = ?

--name: drop-balance-table
DROP TABLE IF EXISTS balance_job

--name: create-user
INSERT INTO user_job VALUES (DEFAULT)

--name: check-users-exist
SELECT count(*) FROM (SELECT * FROM user_job LIMIT 1) AS t