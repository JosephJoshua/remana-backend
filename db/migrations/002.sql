-- +migrate Up
CREATE TABLE login_codes (
  login_code_id UUID NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (user_id),
  login_code TEXT NOT NULL
);

-- +migrate Down
DROP TABLE login_codes;
