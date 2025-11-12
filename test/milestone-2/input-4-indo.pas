program ProcFuncTest;
konstanta
  PI = 3.14;
tipe
  letter = char;
variabel
  c: letter;
  r: real;

prosedur ShowChar(ch: char); (* Tampilin karakter*)
mulai
  write('Character: ', ch);
selesai;

fungsi Area(radius: real): real; {* Hitung luas lingkaran *}
mulai
  Area := PI * radius * radius;
selesai;

mulai
  c := 'A';
  r := 25e1;
  ShowChar(c);
  write('Area = ', Area(r));
selesai.

