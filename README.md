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


# Todo

- [x] Add api and ui for nodes and edges.
- [x] Add a simple table displaying nodes and edges.
- [x] Tables should be filterable.
- [ ] Swap /v1/ui to /ui/v1
- [ ] This of a better way to represent the nodes and edges tables.
- [ ] Add in a detail page for a single node or edge.
- [ ] ID's in the nodes and edges table to be hyperlinks to the detail page of the node or edge.
- [ ] Add in the force directed graph in to the UI. Look at using https://github.com/vasturiano/force-graph
- [ ] Make the store distrabuted using raft.


```js
// Possible example
<script src="https://unpkg.com/force-graph"></script>

<div id="graph"></div>

<script>
fetch("/graph.json")
.then(res => res.json())
.then(data => {

  const Graph = ForceGraph()
    (document.getElementById('graph'))
      .graphData(data)
      .nodeLabel('name')
      .nodeAutoColorBy('label')
      .linkDirectionalArrowLength(6)
      .linkDirectionalArrowRelPos(1)

})
</script>
```
