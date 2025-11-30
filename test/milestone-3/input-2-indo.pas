program LogicalTest;
variabel
  x, y: integer;
  flag: boolean;
mulai
  x := 3;
  y := 7;
  flag := (x < y) dan tidak (x = 0) atau (y >= 10);
  jika flag maka
    write('Condition true', 'ad')
  selain_itu
    write('Condition false', 'ad');
selesai.
