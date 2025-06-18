# YouTube Playlist Downloader

Este é um utilitário em Go para baixar vídeos de uma playlist do YouTube e convertê-los para arquivos MP3.

## Como buildar

Para compilar o programa, execute o seguinte comando no terminal:

```sh
go build -o youtube_downloader.exe main.go
```

## Exemplo de uso

```sh
.\youtube_downloader.exe "PLAYLIST_LINK" "OUTPUT_PATH"
```

### Notas adicionais:

- Certifique-se de ter o Go instalado em sua máquina antes de tentar compilar o programa.
- Substitua `PLAYLIST_LINK` pelo URL da playlist do YouTube.
- Substitua `OUTPUT_PATH` pelo diretório onde deseja salvar os arquivos MP3.

Este utilitário pode ser muito útil para criar uma biblioteca de músicas offline a partir de playlists do YouTube.


## Usando Docker(linux)

Faça a construção da imagem docker

```sh
sudo docker build -t youtube_downloader -f Dockerfile.multi .
```

Crie o executável diretamente no seu diretório

```sh
sudo docker run --rm -v $(pwd):/mnt youtube_downloader sh -c "cp /youtube_downloader /mnt"
```

Execute para baixar

```sh
./youtube_downloader "https://youtube.com/playlist?list=PLop9WaDyI8nDyfJgiCxyWIPzrDxxu-Hhu&si=4roqhi0cyvme_7bS" "/home/petini96/musics"
```

