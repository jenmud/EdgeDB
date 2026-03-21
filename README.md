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


# Tasks

## Todo

- [ ] Move todos into github.
- [ ] Add in a detail page for a single node or edge.
- [ ] ID's in the nodes and edges table to be hyperlinks to the detail page of the node or edge.
- [ ] Make the store distrabuted using raft.
- [ ] Make the graph on the UI interactive. Example clicks etc..... This can be done when I have a better idea how I want to use it.
- [ ] Graph is using the limit which means for large graphs the search limits to the default 1000 and that means only nodes returned. We need nodes and edges.
      One thing that we can do is for a node, it will automatically include the edges that nodes belongs to, and for edges it will fetch the nodes that are not included in the search result set.
- [ ] Make graph be able to run in WASM - this would be a nice to have if possible.
- [ ] UI to add and delete nodes and edges
- [ ] Implement some sort of query language
- [ ] Add in pagination

## Done

- [x] Add api and ui for nodes and edges.
- [x] Add a simple table displaying nodes and edges.
- [x] Tables should be filterable.
- [x] Swap /v1/ui to /ui/v1
- [x] This of a better way to represent the nodes and edges tables.
- [x] Add in the force directed graph in to the UI. Look at using https://github.com/vasturiano/force-graph
- [x] Change the upload to take a graph formatted upload file.
- [x] There is no advantage to having two separate tables for nodes and edges. The data is too similar. Instead have a type column and also make the FTS figure out the type.
- [x] I don't like how I am now passing the data to the graph UI. Think of something better. (works for now)