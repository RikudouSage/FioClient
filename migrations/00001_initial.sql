-- +goose up
CREATE TABLE accounts (
    account_number           VARCHAR(34)    PRIMARY KEY NOT NULL,
    api_key                  VARCHAR(255)   NOT NULL UNIQUE,
    bank_code                VARCHAR(4)     NOT NULL,
    currency                 CHAR(3)        NOT NULL,
    iban                     VARCHAR(34)    NOT NULL,
    bic                      VARCHAR(11)    NOT NULL
);

CREATE TABLE transactions (
    id                       BIGINT         NOT NULL,
    account_number           VARCHAR(34)    NOT NULL,
    date                     DATETZ           NOT NULL,
    amount                   DECIMAL        NOT NULL,
    currency                 CHAR(3)        NOT NULL,
    counterparty_account     VARCHAR(34)    NOT NULL,
    counterparty_name        VARCHAR(255)   NOT NULL,
    counterparty_bank_code   VARCHAR(4)     NOT NULL,
    counterparty_bank_name   VARCHAR(255)   NOT NULL,
    constant_symbol          VARCHAR(10),
    variable_symbol          VARCHAR(10),
    specific_symbol          VARCHAR(10),
    user_identity            VARCHAR(255),
    transaction_type         VARCHAR(255)   NOT NULL,
    performed_by             VARCHAR(255),
    additional_info          TEXT,
    comment                  TEXT,
    bic                      VARCHAR(11),
    instruction_id           BIGINT,
    payer_reference          VARCHAR(255),
    PRIMARY KEY (account_number, id),
    FOREIGN KEY (account_number) REFERENCES accounts (account_number) ON DELETE CASCADE
);

CREATE INDEX transactions_account_date_idx ON transactions (account_number, date);

-- +goose down
DROP TABLE transactions;
DROP TABLE accounts;
