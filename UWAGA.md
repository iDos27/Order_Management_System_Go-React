# PROBLEM Z SERWISEM POWIADOMIEŃ

Po skończeniu pisania aplikacji zdałem sobie sprawę z tego, że serwis powiadomień desktopowych musi być odpalony jako binarka na komputerach użytkowników.
Wrzucenie serwisu do Docker Compose albo Podman Play Kube spowoduje, że tylko serwer będzie z tego korzystał a nie użytkownicy.
Konteneryzacja serwisu powoduje, że powiadomienia będą uruchamiane TYLKO w kontenerze a nie na pulpicie użytkownika co też przeczy idei tego serwisu.

**OPCJE**:
* 1. Zostawić tak jak jest. Skonteneryzować serwis autoryzacji, zamówień i raportów a powiadomienia zostawić w formie binarki
* 2. Zamienić serwis powiadomień na wysyłanie wiadomości email albo jako wiadomość na mattermost/slack/discord