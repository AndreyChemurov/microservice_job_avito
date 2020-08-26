-- name: create-user-table
CREATE TABLE IF NOT EXISTS user_job (
    user_id BIGSERIAL PRIMARY KEY NOT NULL
);

-- name: create-balance-table
CREATE TABLE IF NOT EXISTS balance_job (
    balance_id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER REFERENCES user_job(user_id),
    count INTEGER NOT NULL
);

--name: get-user-balance
SELECT count FROM balance_job;

--name: remittance-from
UPDATE balance_job SET count = count - ? WHERE user_id = ? AND count > 0;

--name: remittance-to
UPDATE balance_job SET count = count + ? WHERE user_id = ?;