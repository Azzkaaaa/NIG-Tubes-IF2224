program ArrayTest;
variabel
  arr: larik[1..5] dari integer;
  i: integer;
mulai
  untuk i := 1 ke 5 lakukan
    arr[i] := i * 2;
  tulis(arr[3])
selesai.
