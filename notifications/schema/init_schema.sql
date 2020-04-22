create table if not EXISTS notifications (
    id bigserial not null PRIMARY KEY,
    message VARCHAR not null,
    created_at INTEGER not null,
    noti_status VARCHAR not NULL,
    for_user INTEGER not NULL,
    checked boolean not null
)