CREATE TABLE IF NOT EXISTS sended_order
(
    result      text not null,
    order_id    text PRIMARY KEY,
    side        text         not null,
    status      text         not null,
    symbol      text         not null,
    quantity    real         not null,
    price       real         not null,
    type        text         not null,
    server_time TIMESTAMP    not null
);

CREATE TABLE IF NOT EXISTS client_users
(
    user_id bigint not null PRIMARY KEY
);
