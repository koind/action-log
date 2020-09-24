create table histories (
    id serial primary key,
    user_id bigint not null,
    project text not null,
    description text not null,
    created_at timestamp not null
);
create index user_id_index on histories (user_id);
create index project_index on histories (project);