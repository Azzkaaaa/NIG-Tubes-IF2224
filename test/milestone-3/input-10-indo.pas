program ComprehensiveTest;
{ Comprehensive test combining multiple features }
konstanta
  PI = 3.14159;
  MAX = 5;

tipe
  Point = rekaman
    x: real;
    y: real
  selesai;

variabel
  p1, p2: Point;
  distances: larik[1..MAX] dari real;
  i: integer;
  total: real;

prosedur initPoint(variabel pt: Point; xVal: integer; yVal: integer);
mulai
  { Implicit cast: integer to real }
  pt.x := xVal;
  pt.y := yVal
selesai;

fungsi distance(p: Point): real;
mulai
  distance := p.x * p.x + p.y * p.y
selesai;

mulai
  { Initialize points using reference parameters }
  initPoint(p1, 3, 4);
  initPoint(p2, 5, 12);
  
  { Calculate distances }
  distances[1] := distance(p1);
  distances[2] := distance(p2);
  
  { Sum with implicit casting }
  total := 0;
  untuk i := 1 ke 2 lakukan
    total := total + distances[i];
  
  tulis('Point 1: (', p1.x, ', ', p1.y, ')');
  tulis('Point 2: (', p2.x, ', ', p2.y, ')');
  tulis('Total distance squared: ', total)
selesai.
