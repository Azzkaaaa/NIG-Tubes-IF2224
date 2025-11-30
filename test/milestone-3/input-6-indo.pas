program ImplicitCastTest;
{ Test implicit casting from integer to real }
variabel
  x: integer;
  y, z: real;
mulai
  x := 5;
  y := 3.14;
  
  { Implicit cast: x (integer) will be cast to real }
  z := x + y;
  
  { Assignment with implicit cast }
  y := x;
  
  tulis('z = ', z);
  tulis('y = ', y)
selesai.
