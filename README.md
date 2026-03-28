# EdgeDB

EdgeDB is a datastore with graph type characteristics.

This is a experimental project mainly for learning.

## Environment variables

See `.env` file for the default environment variables supported.


## Examples

See in the `example` folder for pre-populated examples.

```bash
$ EDGEDB_STORE_DSN="file:./example/movies.sqlite?_fk=1" make run-reload
```

Once running, navigate to http://localhost:8080/v1/ui


![image](https://github.com/jenmud/EdgeDB/blob/fill-missing-edge-nodes-when-fetching-edges/simple-graph-example.png)
