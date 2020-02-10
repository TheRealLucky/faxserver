create table v_fax_users
(
    fax_user_uuid uuid not null
        constraint v_fax_users_pkey
            primary key,
    domain_uuid   uuid,
    fax_uuid      uuid,
    user_uuid     uuid
);

create table v_domains
(
    domain_uuid        uuid not null
        constraint v_domains_pkey
            primary key,
    domain_parent_uuid uuid,
    domain_name        text,
    domain_enabled     text,
    domain_description text
);

insert into v_fax_users values ('30d59b83-3afc-4d32-a417-0bec1eaa2f12', '28dc4965-8d0b-484d-bf8a-49986c53ef4e', '4d902414-82f5-427c-93df-bb3cb494756a', '3685ab1d-2e05-44fe-8c31-87805021e189');

insert into v_domains values ('28dc4965-8d0b-484d-bf8a-49986c53ef4e', 'd8731274-2e01-4cbd-9b98-90b10c14f948', 'domain1', 'true', 'domain1 description');