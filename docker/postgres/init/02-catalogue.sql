CREATE TABLE catalogue.kategorija (
    kategorija_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    naziv varchar(100) NOT NULL,
    opis text,
    datum_kreiranja timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalogue.artikl (
    artikl_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kategorija_id integer NOT NULL REFERENCES catalogue.kategorija (kategorija_id),
    sifra varchar(50) NOT NULL UNIQUE,
    naziv varchar(150) NOT NULL,
    opis text,
    cijena numeric(10,2) NOT NULL CHECK (cijena >= 0),
    kolicina_na_zalihi integer NOT NULL DEFAULT 0 CHECK (kolicina_na_zalihi >= 0),
    datum_kreiranja timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    datum_zadnje_izmjene timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);