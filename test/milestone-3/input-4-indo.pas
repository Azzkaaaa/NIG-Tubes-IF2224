program FunctionTest;

fungsi tambah(a: integer; b: integer): integer;
mulai
  tambah := a + b
selesai;

prosedur cetak(nilai: integer);
mulai
  tulis('Nilai: ', nilai)
selesai;

variabel
  x, y, hasil: integer;

mulai
  x := 10;
  y := 20;
  hasil := tambah(x, y);
  cetak(hasil)
selesai.
