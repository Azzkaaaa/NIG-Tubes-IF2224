program ProcFuncTest;
const
  PI = 3.14;
type
  letter = char;
var
  c: letter;
  r: real;

procedure ShowChar(ch: char);
begin
  write('Character: ', ch);
end;

function Area(radius: real): real;
begin
  Area := PI * radius * radius;
end;

begin
  c := 'A';
  r := 2.5;
  ShowChar(c);
  write('Area = ', Area(r));
end.

