# Server and file encode for Alice's chess

## Target of this project

* Create working server for Alice's chess
* Add archive system to store played games
* Add terminal way to see already played games

## Way to store games
### Encode
#### Piece encode
Pieces are incoded on frontend part as 16 bit integer. Then they will be send on backend to resend to opponent and save in archive. First 6 bit of integer is position on boarde next 3 bit is a number of piece, next 3 bit is a piece type next 2 bits is a team and side (of board left/right) last 2 bits is bit for castling and mate.
### File encode
I create first version of my file type. It can encode 8bit or 16 bit data into 16 and 8 bit uint. For all file you only need 7 bytes for file description.First 4 bytes is a string "pice" for checking files, next 3 bytes : version of encode, data length, word size. All other data will be read with size of word.
