CREATE TYPE unit_type AS ENUM ('unit','kilos','grams','liters', 'box', 'size');
CREATE TYPE product_status AS ENUM ('available','pending','inactive');

create table if not exists products
(
    id              UUID PRIMARY KEY,
    name            varchar (256) NOT NULL,
    description     varchar (256) NULL,
    unit_type       unit_type NOT NULL,
    unit            varchar (50) NOT NULL,
    brand           varchar (50) NOT NULL,
    color           varchar (50) NOT NULL,
    style           varchar (50) NOT NULL,
    status          product_status NOT NULL,
    audit_user      varchar (50) NOT NULL,
    creation_date   timestamp NOT NULL DEFAULT now(),
    update_date     timestamp NOT NULL DEFAULT now()
);
