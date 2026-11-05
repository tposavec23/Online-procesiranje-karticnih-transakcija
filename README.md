# 7pay – Online procesiranje kartičnih transakcija

Sustav za iniciranje i procesiranje online kartičnih transakcija. Projekt se razvija s tvrtkom 7pay u sklopu kolegija Analiza i razvoj programa na Fakultetu organizacije i informatike, ak. god. 2026./2027.

Sustav se sastoji od Android aplikacije i backend servisa u Go-u.

## Struktura repozitorija

```
android/                   Android aplikacija (Kotlin)
services/
  account/                 Account management service
  catalogue/               Catalogue management service
  transaction/             Transaction management service
  bank-communication/      Bank communication service
  bank-simulator/          Bank simulator service
```

## Servisi

| Servis | Zadatak |
|---|---|
| account | Korisnički računi, prijava i uloge (trgovac, admin) |
| catalogue | Artikli: dodavanje i brisanje te dohvat |
| transaction | Zaprimanje transakcija, prepoznavanje banke i routanje, spremanje i dohvat transakcija |
| bank-communication | Komunikacija s bankom: Banka 1 preko ISO8583 socketa, Banka 2 preko HTTP API-ja |
| bank-simulator | Simulator banke koji odobrava ili odbija transakcije |

Testni redak za probni pull request.
