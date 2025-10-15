## Cara run
1. pada root repository ketik:
windows:
```bash
    go build -o bin/pslex.exe ./cmd/pslex
``` 

linux:
```bash
    go build -o bin/pslex ./cmd/pslex
``` 

2. Lalu ketik:
```bash
    ./bin/pslex.exe --rules src/rules/tokenizer.json --input test/milestone-1/(nama file pascal.pas)
```

contoh:
./bin/pslex.exe --rules src/rules/tokenizer.json --input test/milestone-1/test1.pas