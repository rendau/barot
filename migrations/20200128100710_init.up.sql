create table banner
(
    id      bigint not null
        constraint banner_pk primary key,
    slot_id bigint not null,
    note    text   not null default ''
);
create index banner_slot_id_idx
    on banner (slot_id);

create table stat
(
    banner_id   bigint not null
        constraint stat_fk_banner_id references banner on update cascade on delete cascade,
    slot_id     bigint not null,
    usr_type_id bigint not null,
    show_cnt    bigint not null default 0,
    click_cnt   bigint not null default 0
);
create index stat_banner_id_idx
    on stat (banner_id);
create index stat_slot_id_idx
    on stat (slot_id);
create index stat_usr_type_id_idx
    on stat (usr_type_id);
