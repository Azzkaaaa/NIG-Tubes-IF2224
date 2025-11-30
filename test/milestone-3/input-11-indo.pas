program NestedFunctionExample;

fungsi OuterFunction(x: integer): integer;
  fungsi InnerFunction(y: integer): integer;
  mulai
    InnerFunction := y * 2;
  selesai;
mulai
  OuterFunction := InnerFunction(x) + 5;
selesai;

mulai
  write('Result: ', 'aa');
selesai.