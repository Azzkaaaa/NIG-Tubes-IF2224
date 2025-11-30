program RecordTypeTest;
{ Test record type analysis }
konstanta
  MAX_STUDENTS = 100;

tipe
  Mahasiswa = rekaman
    nama: char;
    umur: integer;
    ipk: real
  selesai;
  
  Kelas = rekaman
    kode: char;
    jumlah: integer
  selesai;

variabel
  mhs1, mhs2: Mahasiswa;
  kls: Kelas;

mulai
  { Initialize mahasiswa 1 }
  mhs1.nama := 'A';
  mhs1.umur := 20;
  mhs1.ipk := 3.75;
  
  { Initialize mahasiswa 2 }
  mhs2.nama := 'B';
  mhs2.umur := 21;
  mhs2.ipk := 3.50;
  
  { Initialize kelas }
  kls.kode := 'X';
  kls.jumlah := 30;
  
  tulis('Mahasiswa 1: ', mhs1.nama, ' umur=', mhs1.umur, ' IPK=', mhs1.ipk);
  tulis('Mahasiswa 2: ', mhs2.nama, ' umur=', mhs2.umur, ' IPK=', mhs2.ipk);
  tulis('Kelas: ', kls.kode, ' jumlah=', kls.jumlah)
selesai.
