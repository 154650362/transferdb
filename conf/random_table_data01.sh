#! /bin/bash
sqlplus test/test@helowin >>random1.log << EOF
alter session set nls_date_format = 'yyyy-mm-dd hh24:mi:ss';
create table test.marvin3(
ID number primary key,
INC_DATETIME date,
RANDOM_ID number,
RANDOM_STRING varchar2(1000)
);
create index idx_marvin3_RANDOM_ID on test.marvin3(RANDOM_ID);

insert into test.marvin3
  (ID, INC_DATETIME,RANDOM_ID,RANDOM_STRING)
  select rownum as id,
         to_date(sysdate + rownum / 24 / 3600, 'yyyy-mm-dd hh24:mi:ss') as inc_datetime,
         trunc(dbms_random.value(0, 100)) as random_id,
         dbms_random.string('x', 20) random_string
   from xmltable('1 to 100');
EOF