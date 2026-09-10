-- +goose up
create table accounts (
    account_number           varchar(34)    primary key not null,
    api_key                  varchar(255)   not null,
    bank_code                varchar(4)     not null,
    currency                 char(3)        not null,
    iban                     varchar(34)    not null,
    bic                      varchar(11)    not null
);

create table transactions (
    id                       bigint         not null,
    account_number           varchar(34)    not null,
    date                     datetz         not null,
    amount                   decimal        not null,
    currency                 char(3)        not null,
    counterparty_account     varchar(34)    not null,
    counterparty_name        varchar(255)   not null,
    counterparty_bank_code   varchar(4)     not null,
    counterparty_bank_name   varchar(255)   not null,
    constant_symbol          varchar(10),
    variable_symbol          varchar(10),
    specific_symbol          varchar(10),
    user_identity            varchar(255),
    transaction_type         varchar(255)   not null,
    performed_by             varchar(255),
    additional_info          text,
    comment                  text,
    bic                      varchar(11),
    instruction_id           bigint,
    payer_reference          varchar(255),
    primary key (account_number, id),
    foreign key (account_number) references accounts (account_number) on delete cascade
);

create index transactions_account_date_idx on transactions (account_number, date);
create unique index accounts_api_key_unique_idx on accounts(api_key) where api_key != '';

-- +goose down
drop table transactions;
drop table accounts;
