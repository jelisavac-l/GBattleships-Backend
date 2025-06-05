<h1 align="center">⚓ GBattleships-Backend ⚓</h1>
<div align="center">
  <img src="https://img.shields.io/badge/Luka_Jelisavac-2022--0554-navy" alt="Luka Jelisavac Badge" />
  <img src="https://img.shields.io/badge/Vuk_Janjušević-2022--0xxx-firebrick" alt="Vuk Janjušević Badge" />
</div>
<br>
<div align="center">
  <img src="https://static.vecteezy.com/system/resources/previews/007/849/999/non_2x/abstract-background-of-blue-sea-and-summer-beach-for-banner-invitation-poster-or-website-design-vector.jpg" alt="Centered GIF" style="width: 100%; max-width: 700px;" />
</div>
<br>

## 🧾 Opis projekta

> Ovaj repozitorijum sadrži **isključivo serverski kod**, klijentski kod se može pogledati na sledećem [linku](https://www.youtube.com/watch?v=dQw4w9WgXcQ).

Ovaj repozitorijum sadrži backend deo projekta iz predmeta **Računarske mreže i telekomunikacije**. 
Projekat predstavlja implementaciju klasične igre *Potapanje brodova* (eng. Battleships), gde korisnici mogu da se povežu na server, odaberu protivnika, postave brodove i odigraju partiju.

Backend je razvijen u programskom jeziku **Go**, a komunikacija između servera i klijenata odvija se putem **WebSocket** protokola.

## ⚙️ Tehnologije

- [Go (Golang)](https://golang.org/) – za serversku logiku  
- [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) – za dvosmernu komunikaciju u realnom vremenu  
- JSON – za format razmene poruka  

## ▶️ Pokretanje servera

1. Klonirajte repozitorijum:
   ```bash
   git clone https://github.com/jelisavac-l/GBattleships-Backend.git
   cd GBattleships-Backend
   ```
2. Pokrenite server
   `go run cmd/server/main.go`

## 📚 Napomena

Ovaj projekat je urađen u sklopu praktičnog rada za predmet Računarske mreže i telekomunikacije na Fakultetu organizacionih nauka, Unverziteta u Beogradu.
> ⚠️ **_Commit log_ u ovom repozitorijumu predstavlja dokaz o autentičnosti i samostalnom razvoju projekta.**

Takođe, forma _commit_ poruka prati standarde definisane na [convetional commits](https://www.conventionalcommits.org/en/v1.0.0/).
