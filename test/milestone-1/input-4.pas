program ProcFuncTest;
const
  PI = 3.14;
type
  letter = char;
var
  c: letter;
  r: real;

procedure ShowChar(ch: char); (* Tampilin karakter*)
begin
  write('Character: ', ch);
end;

function Area(radius: real): real; {* Hitung luas lingkaran *}
begin
  Area := PI * radius * radius;
end;

begin
  c := 'A';
  r := 25e1;
  ShowChar(c);
  write('Area = ', Area(r));
end.

