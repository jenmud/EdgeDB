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

- [ ] Move todos into github.
- [x] Add api and ui for nodes and edges.
- [x] Add a simple table displaying nodes and edges.
- [x] Tables should be filterable.
- [x] Swap /v1/ui to /ui/v1
- [x] This of a better way to represent the nodes and edges tables.
- [ ] Add in a detail page for a single node or edge.
- [ ] ID's in the nodes and edges table to be hyperlinks to the detail page of the node or edge.
- [x] Add in the force directed graph in to the UI. Look at using https://github.com/vasturiano/force-graph
- [ ] Make the store distrabuted using raft.
- [ ] Make the graph on the UI interactive. Example clicks etc..... This can be done when I have a better idea how I want to use it.
- [x] Change the upload to take a graph formatted upload file.
- [ ] I don't like how I am now passing the data to the graph UI. Think of something better. (works for now)
- [ ] There is no advantage to having two separate tables for nodes and edges. The data is too similar. Instead have a type column and also make the FTS figure out the type.


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
