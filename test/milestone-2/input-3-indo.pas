program LoopRangeTest;
variabel
  i: integer;
  arr: larik [1..5] dari integer;
mulai
  untuk i := 1 ke 5 lakukan
        arr[i] := i * 2;
  untuk i := 5 turun_ke 1 lakukan
        write(arr[i]);
  i := 0;
  selama i < 5 lakukan
  mulai
    i := i + 1
  selesai;
selesai.
