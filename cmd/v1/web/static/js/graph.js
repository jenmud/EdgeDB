function FetchGraph(target, graphURL) {
    fetch(graphURL).then(res => res.json()).then(data => {
        // lets add some particle speeds to the edges
        data.edges.forEach((i) => {
            i.value = 1;
        });

        const Graph = new ForceGraph(document.getElementById(target))
            .graphData({
                nodes: data.nodes,
                links: data.edges // pass edges as links
            })
            .nodeId('id')
            .nodeLabel('snippet')
            .nodeAutoColorBy('label')
            .cooldownTicks(100)
            //.maxZoom(10)    // closer zoom limit
            //.minZoom(3)   // zoomed out limit
            .linkSource('from_id')
            .linkTarget('to_id')
            .linkAutoColorBy(d => d.from_id.label)
            .linkCurvature('curvature')
            .linkDirectionalParticles("value")
            .linkDirectionalParticleSpeed(d => d.value * 0.001)
            .nodeCanvasObject((node, ctx, globalScale) => {
                const r = 4 + (node.val || 1); // approximate default scaling

                ctx.beginPath();
                ctx.arc(node.x, node.y, r, 0, 2 * Math.PI);
                ctx.fillStyle = node.color;
                ctx.fill();

                if (r * globalScale < 10) return;

                const fontSize = 12 / globalScale;
                ctx.font = `${fontSize}px Sans-Serif`;
                ctx.textAlign = "center";
                ctx.textBaseline = "middle";
                ctx.fillStyle = "white";

                ctx.fillText(node.id, node.x, node.y);
            });

            Graph.d3Force('center', null);
            Graph.onEngineStop(() => Graph.zoomToFit(400)); // fit to canvas when engine stops
    })
}