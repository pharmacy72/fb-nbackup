connect NBEXAMPLE;

CREATE TABLE TESTDATA (
    ID INTEGER PRIMARY KEY NOT NULL,
    NAME VARCHAR(200)
);

delete from TESTDATA;
-- test data
SET TERM ^ ;
execute block
as
declare i int = 0;
begin
  while (i < 10000) do
  begin
    insert into TESTDATA values (:i, 'DATA-'||:i);
    i = i + 1;
  end
end^
SET TERM ; ^

select count(1) from TESTDATA;

commit;