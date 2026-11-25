# XP radionica

Funkcionalnost: uvoz i izvoz artikala

## Priča 1: Uvoz artikala iz CSV-a u bazu

Procjena: 3 boda

Kao administrator želim uvesti artikle iz CSV datoteke, kako bi se odmah spremili u katalog i ne bih ih unosio jedan po jedan.

Kriteriji prihvatljivosti:
- Zadano: CSV s 8 artikala s novim šiframa. Kada: odaberem datoteku i dodirnem „Uvezi”. Tada: svih 8 artikala spremljeno je u bazu i vidim ih u katalogu.
- Zadano: CSV od 8 artikala u kojem redak 4 ima negativnu cijenu. Kada: uvezem datoteku. Tada: spremljeno je 7 artikala, a vidim poruku „redak 4: neispravna cijena”.

## Priča 2: Slanje artikala e-poštom

Procjena: 3 boda

Kao administrator želim poslati popis artikala e-poštom, kako bih imao katalog izvan aplikacije i mogao ga podijeliti.

Kriteriji prihvatljivosti:
- Zadano: u katalogu je 12 artikala. Kada: dodirnem „Pošalji e-poštom”. Tada: otvara se dijeljenje s datotekom artikli.csv od 13 redaka (zaglavlje + 12 artikala).
- Zadano: otvoreno dijeljenje. Kada: odaberem aplikaciju za e-poštu. Tada: artikli.csv je privitak nove poruke.

## Priča 3: Pregled prije uvoza

Procjena: 2 boda

Kao administrator želim vidjeti što će se uvesti prije spremanja, kako ne bih greškom uvezao krivu datoteku.

Kriteriji prihvatljivosti:
- Zadano: odabran CSV s 8 artikala. Kada: dodirnem „Pregledaj”. Tada: vidim popis 8 artikala i ništa još nije spremljeno.
- Zadano: otvoren pregled. Kada: dodirnem „Odustani”. Tada: u bazi nema nijednog novog artikla.
