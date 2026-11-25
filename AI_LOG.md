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
