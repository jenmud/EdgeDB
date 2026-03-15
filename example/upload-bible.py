import json
import requests

KJV_JSON_URL = "https://raw.githubusercontent.com/thiagobodruk/bible/refs/heads/master/json/en_kjv.json"
API_URL = "http://localhost:8080/api/v1/graph"  # replace with your API endpoint


def download_kjv_json(url):
    response = requests.get(url)
    if response.status_code == 200:
        text = response.content.decode("utf-8-sig")  # remove BOM
        return json.loads(text)
    else:
        raise Exception(f"Failed to download KJV JSON: {response.status_code}")


def create_graph(bible_data):
    nodes = []
    edges = []
    node_id = 1
    edge_id = 1
    book_nodes = {}

    for book in bible_data:
        book_name = book["name"]
        chapters = book["chapters"]

        # Book node
        book_node_id = node_id
        nodes.append({
            "id": book_node_id,
            "label": book_name,
            "properties": {"type": "book"}
        })
        node_id += 1
        book_nodes[book_name] = book_node_id

        previous_chapter_node = None

        for chapter_idx, verses in enumerate(chapters):
            chapter_num = chapter_idx + 1
            chapter_node_id = node_id
            nodes.append({
                "id": chapter_node_id,
                "label": f"{book_name} {chapter_num}",
                "properties": {"type": "chapter", "chapter": chapter_num}
            })
            node_id += 1

            # Edge: book -> chapter
            edges.append({
                "id": edge_id,
                "from_id": book_node_id,
                "to_id": chapter_node_id,
                "label": "has_chapter",
                "weight": 1,
                "properties": {}
            })
            edge_id += 1

            # Optional: link previous chapter -> current chapter
            if previous_chapter_node:
                edges.append({
                    "id": edge_id,
                    "from_id": previous_chapter_node,
                    "to_id": chapter_node_id,
                    "label": "next_chapter",
                    "weight": 1,
                    "properties": {}
                })
                edge_id += 1
            previous_chapter_node = chapter_node_id

            previous_verse_node = None
            for verse_idx, verse_text in enumerate(verses):
                verse_num = verse_idx + 1
                verse_node_id = node_id
                nodes.append({
                    "id": verse_node_id,
                    "label": f"{book_name} {chapter_num}:{verse_num}",
                    "properties": {
                        "type": "verse",
                        "text": verse_text,
                        "word_count": len(verse_text.split())
                    }
                })
                node_id += 1

                # Edge: chapter -> verse
                edges.append({
                    "id": edge_id,
                    "from_id": chapter_node_id,
                    "to_id": verse_node_id,
                    "label": "has_verse",
                    "weight": 1,
                    "properties": {}
                })
                edge_id += 1

                # Edge: previous verse -> current verse
                if previous_verse_node:
                    edges.append({
                        "id": edge_id,
                        "from_id": previous_verse_node,
                        "to_id": verse_node_id,
                        "label": "next_verse",
                        "weight": 1,
                        "properties": {}
                    })
                    edge_id += 1
                previous_verse_node = verse_node_id

    return {"nodes": nodes, "edges": edges}


def upload_graph(graph_data, api_url):
    with open("graph-nodes.json", "w") as w:
        w.write(json.dumps(graph_data.get("nodes")))

    with open("graph-edges.json", "w") as w:
        w.write(json.dumps(graph_data.get("edges")))

    headers = {"Content-Type": "application/json"}
    response = requests.put(api_url, headers=headers, data=json.dumps(graph_data))
    if response.status_code == 200:
        print("Graph uploaded successfully!")
        with open("graph-uploaded", "w") as w:
            w.write(response.text)
    else:
        print(f"Failed to upload graph: {response.status_code} {response.text}")


if __name__ == "__main__":
    print("Downloading KJV JSON...")
    kjv_data = download_kjv_json(KJV_JSON_URL)
    print("Creating graph...")
    graph_data = create_graph(kjv_data)
    print(f"Graph has {len(graph_data['nodes'])} nodes and {len(graph_data['edges'])} edges.")
    print("Uploading graph...")
    upload_graph(graph_data, API_URL)