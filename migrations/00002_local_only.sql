-- +goose up
alter table transactions add column local_only boolean not null default false;

-- +goose down
alter table transactions drop column local_only;
