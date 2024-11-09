CREATE TABLE IF NOT EXISTS clientes (
    "id"        BIGSERIAL       NOT NULL    PRIMARY KEY,
    "cpf"       VARCHAR(14)     NOT NULL    UNIQUE,
    "nome"      TEXT            NOT NULL,
    "endereco"  TEXT            NOT NULL,
    "email"     TEXT            NOT NULL                    DEFAULT ''
);

---- create above / drop below ----

DROP TABLE IF EXISTS clientes;
