create table if not EXISTS opaque_store (
    id bigserial not NULL PRIMARY KEY,
    user_id int not null,
    jwt VARCHAR not null,
    opaque VARCHAR not null,
    created_at TIMESTAMP not null
)