async function MakeGraph(target, apiQuery) {
    const res = await fetch(apiQuery);
    
    if (!res.ok) {
        console.error("fetch failed", res.status, await res.text());
        return;
    }

    const data = await res.json();

    ForceGraph()(document.getElementById(target))
        .graphData(data)
        .nodeLabel('label')
        .nodeAutoColorBy('label')
        .linkDirectionalArrowLength(6)
        .linkDirectionalArrowRelPos(1);
}