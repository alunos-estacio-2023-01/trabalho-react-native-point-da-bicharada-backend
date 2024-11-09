CREATE TABLE IF NOT EXISTS clientes_telefones (
    "cliente_id"    BIGINT          NOT NULL,
    "telefone"      VARCHAR(12)     NOT NULL    UNIQUE,

    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    PRIMARY KEY (cliente_id, telefone)
);

---- create above / drop below ----

DROP TABLE IF EXISTS clientes_telefones;
