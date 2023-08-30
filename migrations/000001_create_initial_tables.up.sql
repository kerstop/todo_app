BEGIN;

CREATE TABLE users 
(
    id serial,
    username varchar(255) NOT NULL,
    passwd_hash varchar(255) NOT NULL,
    salt varchar(255) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_id_key UNIQUE (id),
    CONSTRAINT users_username_key UNIQUE (username)
);

CREATE TABLE todo_entries
(
    id serial NOT NULL,
    descript text NOT NULL,
    complete boolean NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT todo_entries_pkey PRIMARY KEY (id),
    CONSTRAINT todo_entries_id_key UNIQUE (id),
    CONSTRAINT todo_entries_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES users (id)
);

COMMIT;