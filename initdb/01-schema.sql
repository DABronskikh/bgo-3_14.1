CREATE TABLE clients
(
    id        BIGSERIAL PRIMARY KEY,
    login     TEXT      NOT NULL,
    password  TEXT      NOT NULL,
    full_name TEXT      NOT NULL,
    passport  TEXT      NOT NULL,
    birthday  DATE      NOT NULL,
    status    TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE cards
(
    id       BIGSERIAL PRIMARY KEY,
    number   TEXT      NOT NULL,
    balance  BIGINT    NOT NULL DEFAULT 0,
    issuer   TEXT      NOT NULL CHECK (issuer IN ('Visa', 'MasterCard', 'MIR')),
    holder   TEXT      NOT NULL,
    owner_id BIGINT    NOT NULL REFERENCES clients,
    status   TEXT      NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE supplier_icons
(
    id  BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL
);

CREATE TABLE mcc
(
    id          BIGSERIAL PRIMARY KEY,
    mcc         TEXT NOT NULL,
    description TEXT NOT NULL
);


CREATE TABLE transactions
(
    id               BIGSERIAL PRIMARY KEY,
    card_id          BIGINT    NOT NULL REFERENCES cards,
    amount           BIGINT    NOT NULL,
    created          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status           TEXT      NOT NULL DEFAULT 'Обрабатывается' CHECK (status IN ('Обрабатывается', 'Исполнена', 'Отклонена')),
    mcc_id           BIGINT    NOT NULL REFERENCES mcc,
    description      TEXT      NOT NULL,
    supplier_icon_id BIGINT    NOT NULL REFERENCES supplier_icons
);
