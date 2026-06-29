create database if not exists test;

use test;

create table if not exists `user` (
	id bigint primary key auto_increment,
	username varchar(64) not null unique,
	password varchar(255) not null,
	created_at datetime not null
);

create table if not exists issues (
	id bigint primary key auto_increment,
	user_id bigint not null,
	title varchar(255) not null,
	content text not null,
	status varchar(32) not null,
	priority int not null,
	created_at datetime not null,
	updated_at datetime not null,
	index idx_issues_user_id (user_id)
);
