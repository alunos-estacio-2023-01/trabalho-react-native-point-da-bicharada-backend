CREATE UNIQUE INDEX unique_non_empty_clients_email ON clientes (email)
WHERE email <> '';

---- create above / drop below ----

DROP INDEX IF EXISTS unique_non_empty_clients_email;
