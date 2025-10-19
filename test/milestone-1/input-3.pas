program LoopRangeTest;
var
  i: integer;
  arr: array [1..5] of integer;
begin
  for i := 1 to 5 do
    arr[i] := i * 2;

  for i := 5 downto 1 do
    write(arr[i]);

  i := 0;
  while i < 5 do
  begin
    i := i + 1;
  end;
end.
