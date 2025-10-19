program LogicalTest;
var
  x, y: integer;
  flag: boolean;
begin
  x := 3;
  y := 7;
  flag := (x < y) and not (x = 0) or (y >= 10);
  if flag then
    write('Condition true')
  else
    write('Condition false');
end.
