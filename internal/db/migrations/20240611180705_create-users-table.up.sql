CREATE TABLE users (
   id serial primary key,
   name character varying(100) not null,
   email character varying(100) not null,
   status character varying(30) default NULL,
   username character varying(100) not NULL,
   password character varying(100) not NULL
);