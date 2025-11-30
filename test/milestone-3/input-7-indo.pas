program StaticRangeTest;
{ Test static evaluation for array range operator }
konstanta
  MIN_INDEX = 1;
  MAX_INDEX = 10;

variabel
  numbers: larik[MIN_INDEX..MAX_INDEX] dari integer;
  values: larik[1..(MAX_INDEX - MIN_INDEX + 1)] dari real;
  flags: larik[0..4] dari boolean;
  i: integer;

mulai
  { Initialize arrays using constant range }
  untuk i := MIN_INDEX ke MAX_INDEX lakukan
    numbers[i] := i * 2;
  
  untuk i := 1 ke (MAX_INDEX - MIN_INDEX + 1) lakukan
    values[i] := i * 1.5;
  
  untuk i := 0 ke 4 lakukan
    flags[i] := (i mod 2 = 0);
selesai.
