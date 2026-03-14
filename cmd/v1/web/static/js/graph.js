function MakeGraph(target, data) {

    // lets add some particle speeds to the edges
    // FIXME: this should be read from a property on the edge/link
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
    .linkSource('from')
    .linkTarget('to')
    .linkAutoColorBy(d => d.from.label)
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
}