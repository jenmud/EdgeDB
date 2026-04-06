// FetchGraph will fetch graph data from the given URL.
async function FetchGraph(graphURL) {
    const res = await fetch(graphURL);
    const data = await res.json();
    return data;
}


// FillGraph takes graph data and returns the populated Graph.
function FillGraph(target, data) {

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
        .nodeLabel('label')
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
        .linkLabel("label")
        .onNodeDragEnd(node => {
              node.fx = node.x;
              node.fy = node.y;
        })
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
        })
        .linkCanvasObjectMode(() => 'after')
        .linkCanvasObject((link, ctx) => {
          const MAX_FONT_SIZE = 4;
          const LABEL_NODE_MARGIN = Graph.nodeRelSize() * 1.5;

          const start = link.source;
          const end = link.target;

          // ignore unbound links
          if (typeof start !== 'object' || typeof end !== 'object') return;

          // midpoint
          const textPos = {
            x: start.x + (end.x - start.x) / 2,
            y: start.y + (end.y - start.y) / 2
          };

          const dx = end.x - start.x;
          const dy = end.y - start.y;

          const linkLength = Math.sqrt(dx * dx + dy * dy);
          const maxTextLength = linkLength - LABEL_NODE_MARGIN * 2;

          // angle
          let textAngle = Math.atan2(dy, dx);

          // keep text upright
          if (textAngle > Math.PI / 2) textAngle -= Math.PI;
          if (textAngle < -Math.PI / 2) textAngle += Math.PI;

          const label = `[${link.label ?? ''}]`;
          if (!label.trim()) return;

          // estimate font size
          ctx.font = '1px Sans-Serif';
          const textWidth = ctx.measureText(label).width;
          const fontSize = Math.min(MAX_FONT_SIZE, maxTextLength / textWidth);

          if (fontSize <= 0) return;

          ctx.font = `${fontSize}px Sans-Serif`;
          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';

          ctx.save();
          ctx.translate(textPos.x, textPos.y);
          ctx.rotate(textAngle);

          // --- CUT OUT THE LINE UNDER THE TEXT ---
          ctx.globalCompositeOperation = 'destination-out';

          // add a bit of padding so no line peeks through
          ctx.lineWidth = fontSize * 0.4;
          ctx.strokeText(label, 0, 0);
          ctx.fillText(label, 0, 0);

          // --- DRAW TEXT NORMALLY ---
          ctx.globalCompositeOperation = 'source-over';
          ctx.fillStyle = 'darkgrey';
          ctx.fillText(label, 0, 0);

          ctx.restore();
        });
        
        Graph.onNodeClick(async (node, event) => {
            const data = await FetchGraph(`http://localhost:7331/api/v1/graph/nodes/${node.id}`);
            const current = Graph.graphData();

            data.edges.forEach((i) => {
                i.value = 1;
            });

            // Ensure no duplicate nodes by filtering out nodes that already exist in the graph
            const existingNodeIds = new Set(current.nodes.map((n) => n.id));
            const newNodes = data.nodes.filter((n) => !existingNodeIds.has(n.id));
            
            // Ensure no duplicate nodes by filtering out nodes that already exist in the graph
            const existingEdgeIds = new Set(current.links.map((e) => e.id));
            const newEdges = data.edges.filter((e) => !existingEdgeIds.has(e.id));

            Graph.graphData({
                nodes: [...current.nodes, ...newNodes],
                links: [...current.links, ...newEdges]
            });
        });

        // Graph.d3Force('center', null);
        // Graph.onEngineStop(() => Graph.zoomToFit(400)); // fit to canvas when engine stops
}
