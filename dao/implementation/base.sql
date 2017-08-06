drop database parser;
create database parser;
use parser;

create table content (
url VARCHAR(255) not null,
context LONGTEXT not null,
primary key (url))