program ReferenceParameterTest;
{ Test reference parameters using variabel keyword }

prosedur swap(variabel a: integer; variabel b: integer);
variabel
  temp: integer;
mulai
  temp := a;
  a := b;
  b := temp
selesai;

prosedur increment(variabel x: integer);
mulai
  x := x + 1
selesai;

variabel
  num1, num2: integer;

mulai
  num1 := 10;
  num2 := 20;
  
  tulis('Sebelum swap: num1=', num1, ' num2=', num2);
  swap(num1, num2);
  tulis('Setelah swap: num1=', num1, ' num2=', num2);
  
  increment(num1);
  tulis('Setelah increment: num1=', num1)
selesai.
