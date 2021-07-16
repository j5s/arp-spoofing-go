create database arp;
use arp;

create table hosts (
    id int unsigned not null primary key auto_increment,
    ip varchar(100) not null,
    mac varchar(100) not null,
    mac_info varchar(100) not null
);