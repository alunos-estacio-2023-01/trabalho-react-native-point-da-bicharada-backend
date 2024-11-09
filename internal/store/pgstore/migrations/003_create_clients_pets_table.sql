CREATE TABLE IF NOT EXISTS clientes_pets (
    "id"            BIGSERIAL       NOT NULL    PRIMARY KEY,
    "cliente_id"    BIGINT          NOT NULL,
    "nome"          TEXT            NOT NULL,
    "raca"          TEXT            NOT NULL,
    "especie"       TEXT            NOT NULL,

    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE IF EXISTS clientes_telefones;
