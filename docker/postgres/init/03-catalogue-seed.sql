INSERT INTO catalogue.kategorija (naziv, opis) VALUES
('Torte', 'Kriške torti'),
('Kolači', 'Pojedinačni kolači i deserti');

INSERT INTO catalogue.artikl (kategorija_id, sifra, naziv, opis, cijena, kolicina_na_zalihi) VALUES
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Torte'),  'S-001', 'Čokoladna torta', 'Kriška, tamna čokolada i vrhnje', 4.50, 8),
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Torte'),  'S-002', 'Sacher kocka', 'Marelica i čokoladna glazura', 3.80, 12),
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Kolači'), 'S-003', 'Krempita', 'Domaća, s vanilija kremom', 3.20, 0),
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Kolači'), 'S-004', 'Palačinka s Nutellom', '2 komada, posipana lješnjacima', 3.00, 2),
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Kolači'), 'S-005', 'Tiramisu', 'Mascarpone, kava i piškote', 3.50, 6),
((SELECT kategorija_id FROM catalogue.kategorija WHERE naziv = 'Torte'),  'S-006', 'Cheesecake od jagoda', 'Kriška, preljev od svježih jagoda', 4.20, 1);