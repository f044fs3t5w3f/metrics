CREATE TABLE metric(
    id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    name varchar(255),
    "type" varchar(7), -- Можно потом сделать оптимальнее
    "value" double precision,
    delta BIGINT,
    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX type_name_1762799044404_index ON public.metric USING btree (type, name);