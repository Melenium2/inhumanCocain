create table if not exists tickets 
(
    id bigserial not null PRIMARY KEY,
    title VARCHAR not null,
    description VARCHAR not null,
    section VARCHAR not null,
    from_user INT not null,
    helper INT,
    created_at INT,
    status varchar
);

create table if not EXISTS ticket_messages 
(
    id bigserial not null PRIMARY KEY,
    who int not null,
    ticket_id int not null,
    message_text VARCHAR not null,
    sended_at INT
);