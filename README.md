# malcover

malcover generates MAL cover CSS in the `#more1234` format.

Two paths are provided to support both anime and manga lists.

- `/<username>/anime.css`
- `/<username>/manga.css`

Each support minified CSS via a URL query string, like:

- `/<username>/anime.css?minify=true`
- `/<username>/manga.css?minify=true`

The server is rate limited by default to 5 requests at a time,
with a timeout of 5 seconds.