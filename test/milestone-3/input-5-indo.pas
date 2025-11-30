program ArrayTest;
konstanta
  kons = 67;
tipe
  angka = integer;
  mobil = rekaman
    tahun : angka;
    merek : string;
  selesai;
variabel
  arr: larik[1..5] dari integer;
  i: integer;
  mobil1: mobil;
mulai
  mobil1.tahun := 2004;
  mobil1.merek := 'honda';
  untuk i := 1 ke 5 lakukan
    arr[i] := i * 2;
selesai.
