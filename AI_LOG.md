# AI_LOG

AI agent: Claude Code (Claude Opus 5.5)

## 25. 11. 2026., planiranje izdanja

Upit: prijedlozi korisničkih priča za uvoz i izvoz artikala.

AI je predložio priče, kriterije prihvatljivosti (Zadano/Kada/Tada) i procjene u bodovima te slanje e-poštom preko Androidovog dijeljenja umjesto s backenda. Tim je odabrao tri priče (uvoz u bazu, slanje e-poštom, pregled prije uvoza) i prilagodio ih.

#Radionica - Prvi user story
Gdje je AI pomogao:
- Razrada priča i kriterija prihvatljivosti u obliku Zadano / Kada / Tada.
- Parser: obrada zareza i točke u cijeni, cijene u centima, BOM na početku datoteke, zaglavlje s točnim nazivima stupaca.
- Spremanje u jednoj transakciji: ili se spreme svi artikli ili nijedan. ON CONFLICT (code) ažurira postojeće artikle, a nepostojeća kategorija se stvara.
- Objašnjenje grešaka (npr. 405 u testu znači da POST ruta nije registrirana u Handler-u koji test koristi).
## 25. 11. 2026., testovi za uvoz artikala

Upit: testovi u Go-u za priču „Uvoz artikala iz CSV-a u bazu” prema njezinim kriterijima prihvatljivosti, prije pisanja koda.

AI je predložio testove za čitanje i provjeru CSV-a (importer_test.go), spremanje artikala uz lažnu bazu (service_test.go) i API za uvoz (handler_test.go). Prva verzija nije odgovarala modelu projekta (cijena kao float umjesto u centima, bez kategorije i količine), pa je tim AI-u dao pristup repozitoriju i testovi su prilagođeni. Tim je potvrdio da testovi padaju (crvena faza) i commitao ih prije implementacije. AI-jev prijedlog implementacije nije korišten; kod piše drugi član tima prema testovima.
