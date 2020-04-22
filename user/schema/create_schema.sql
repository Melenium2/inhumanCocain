create table if not EXISTS users (
    id bigserial not null PRIMARY KEY,
    username VARCHAR not null,
    email VARCHAR not null UNIQUE,
    encrypted_password varchar not null,
    created_at INTEGER,
    token VARCHAR not null,
    contacts VARCHAR,
    role VARCHAR not null,
    is_active boolean
)