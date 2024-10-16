CREATE TABLE users(
    id UUID NOT NULL PRIMARY KEY,
    email varchar(63) NOT NULL UNIQUE,
    username varchar(63) NOT NULL UNIQUE,
    password_hash varchar(63) NOT NULL,
    confirmed_email boolean NOT NULL DEFAULT false
);

CREATE TABLE tokens(
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    token varchar(63) NOT NULL,
    exp_date Timestamp NOT NULL
);

CREATE TABLE users_info(
    user_id UUID REFERENCES users(id) PRIMARY KEY,
    first_name varchar(32) NOT NULL,
    second_name varchar(32),
    status varchar(32),
    birthday Timestamp,
    gender varchar(16),
    country varchar(32),
    city varchar(32)
);

CREATE TABLE friends_invites(
    from_user_id UUID REFERENCES users(id),
    to_user_id UUID REFERENCES users(id),
    confirmed boolean NOT NULL DEFAULT 'f',
    PRIMARY KEY(from_user_id, to_user_id),
    CHECK (from_user_id != to_user_id)
);

create extension pg_trgm;

create or replace function fullname(first_name varchar, second_name varchar)
returns text
language plpgsql
immutable
as $$
begin
  return regexp_replace(
    lower(
      trim(
        coalesce(first_name, '') || ' ' ||
        coalesce(second_name, '')
      )
    ),
    'ё',
    'е',
    'g'
   );
exception
  when others then raise exception '%', sqlerrm;
end;
$$;