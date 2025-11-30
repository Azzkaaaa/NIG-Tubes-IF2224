program TestFeatures;
variabel
    x : integer;
    y : real;

prosedur swap(variabel a : integer; variabel b : integer);
variabel
    temp : integer;
mulai
    temp := a;
    a := b;
    b := temp
selesai;

fungsi tambah(a : integer; b : integer) : integer;
mulai
    tambah := a + b
selesai;

mulai
    x := 5;
    y := 3.14;
    y := x + y;
    x := tambah(3, 7)
selesai.
