-- +goose Up

create table if not exists users (
  id bigserial primary key,
  username text unique not null,
  created_at timestamp with time zone not null,
  invite_id bigint
);

create table if not exists credentials (
  id bigserial primary key,
  user_id bigint not null,
  login text not null,
  password_hash text not null,
  created_at timestamp with time zone not null,
  obsoleted_at timestamp with time zone
);

create table if not exists invites (
  id bigserial primary key,
  code text unique not null,
  user_id bigint not null,
  created_at timestamp with time zone not null,
  used_at timestamp with time zone
);

create table if not exists posts (
  id bigserial primary key,
  user_id bigint not null,
  body text unique not null,
  created_at timestamp with time zone not null,
  expires_at timestamp with time zone not null
);