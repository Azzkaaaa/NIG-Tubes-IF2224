program StaticRangeTest;
{ Test static evaluation for array range operator }
konstanta
  MIN_INDEX = 1;
  MAX_INDEX = 10;
  ARRAY_SIZE = MAX_INDEX - MIN_INDEX + 1;

variabel
  numbers: larik[MIN_INDEX..MAX_INDEX] dari integer;
  values: larik[1..ARRAY_SIZE] dari real;
  flags: larik[0..4] dari boolean;
  i: integer;

mulai
  { Initialize arrays using constant range }
  untuk i := MIN_INDEX ke MAX_INDEX lakukan
    numbers[i] := i * 2;
  
  untuk i := 1 ke ARRAY_SIZE lakukan
    values[i] := i * 1.5;
  
  untuk i := 0 ke 4 lakukan
    flags[i] := (i mod 2 = 0);
  
  { Display some values }
  tulis('numbers[5] = ', numbers[5]);
  tulis('values[3] = ', values[3]);
  tulis('flags[2] = ', flags[2])
selesai.
