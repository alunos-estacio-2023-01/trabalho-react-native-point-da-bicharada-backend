-- name: CreateClient :one
INSERT INTO clientes ("cpf", "nome", "endereco", "email")
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateClientPhones :copyfrom
INSERT INTO clientes_telefones ("cliente_id", "telefone")
VALUES ($1, $2);

-- name: CreateClientPets :copyfrom
INSERT INTO clientes_pets ("cliente_id", "nome", "raca", "especie")
VALUES ($1, $2, $3, $4);

-- name: GetClient :one
SELECT * FROM clientes c
WHERE c.cpf = $1;

-- name: GetClients :many
SELECT * FROM clientes c;

-- name: GetClientPhones :many
SELECT telefone FROM clientes_telefones WHERE cliente_id = $1;

-- name: GetClientPets :many
SELECT * FROM clientes_pets WHERE cliente_id = $1;

-- name: DeleteClient :one
WITH deleted_rows AS (
    DELETE FROM clientes
    WHERE cpf = $1
    RETURNING 1
)
SELECT COUNT(*) AS count FROM deleted_rows;

-- name: UpdateCliente :one
UPDATE clientes
SET
    nome = COALESCE(sqlc.narg(nome), nome),
    email = COALESCE(sqlc.narg(email), email),
    endereco = COALESCE(sqlc.narg(endereco), endereco)
WHERE cpf = $1
RETURNING 1;

-- name: UpdateClientePhones :one
WITH id AS (
    SELECT id FROM clientes WHERE cpf = $1
),
deleted_rows AS (
    DELETE FROM clientes_telefones
    WHERE EXISTS (SELECT id FROM id)
        AND cliente_id = (SELECT id FROM id LIMIT 1)
    AND telefone <> ALL(ARRAY[sqlc.arg(telefone)]::text[])
    RETURNING 1
),
inserted_rows AS (
    INSERT INTO clientes_telefones (cliente_id, telefone)
    SELECT id, unnest(sqlc.arg(telefones)::text[])
    FROM id
    WHERE EXISTS (SELECT 1 FROM id)
    ON CONFLICT (cliente_id, telefone) DO NOTHING
)
SELECT EXISTS (SELECT 1 FROM id) AS success;

-- name: UpdateClientePets :one
WITH id AS (
    SELECT id FROM clientes WHERE cpf = $1
),
deleted_rows AS (
    DELETE FROM clientes_pets
    WHERE EXISTS (SELECT id FROM id)
      AND cliente_id = (SELECT id FROM id LIMIT 1)
    RETURNING 1
),
inserted_rows AS (
    INSERT INTO clientes_pets (cliente_id, nome, raca, especie)
    SELECT id, unnest(sqlc.arg(nomes)::text[]), unnest(sqlc.arg(racas)::text[]), unnest(sqlc.arg(especies)::text[])
    FROM id
    WHERE EXISTS (SELECT 1 FROM id)
)
SELECT EXISTS (SELECT 1 FROM id) AS success;
