# chacha20

We forked https://github.com/convto/ChaCha20 for the ChaCha64 for Firebird and made some modifications.

See https://firebirdsql.org/docs/drivers/java/6.0.x/release_notes.html#chacha64 .

- ChaCha: 96 bit nonce and 32 bit counter
- ChaCha64: 64 bit nonce and 64 bit counter

Implemented following [Bernstein, D., "ChaCha, a variant of Salsa20", January 2008](http://cr.yp.to/chacha/chacha-20080128.pdf) and passed through the [RFC8439 Appendix A.1](https://datatracker.ietf.org/doc/html/rfc8439#appendix-A.1) test suite.

# blog (ja-JP)
https://convto.hatenablog.com/entry/2024/02/26/121013
