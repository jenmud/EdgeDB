import requests
import json
import re
import logging
from collections import defaultdict


# Define the API URL for uploading the graph
GRAPH_API_URL = "http://localhost:8080/api/v1/graph"

# Define the URL for downloading the Bible text
BIBLE_API_URL = "https://bible-api.com/data/kjv"

_id = 1

# nextID returns the next ID.
def nextID() -> int:
    global _id
    i = _id
    _id += 1
    return i


# A dictionary to store the graph structure
graph = {
    "nodes": [],
    "edges": []
}

# Helper function to add a node to the graph
def add_node(label: str, properties: dict = {}) -> dict:
    if properties is None:
        properties = {}

    node = {
        "id": nextID(),
        "label": label,
        "properties": properties
    }

    graph['nodes'].append(node)
    return node

# Helper function to add an edge to the graph
def add_edge(from_id: int, label: str, to_id: int, weight: int = 1, properties: dict = {}) -> dict:
    if properties is None:
        properties = {}

    edge ={
        "id": nextID(),
        "from_id": from_id,
        "label": label,
        "to_id": to_id,
        "weight": weight,
        "properties": properties
    }

    graph['edges'].append(edge)
    return edge


# add in the root node
root_node = add_node("bible", {})

def fetch(url: str):
    try:
        logging.info("fetching book %s", url)
        resp = requests.get(url)
        data = resp.json()
    except Exception as e:
        logging.exception("error fetching books")
        return

    for book in data.get("books", []):
        chapter_url = book.get("url", "")

        book_node = add_node(
            label="book",
            properties={"url": chapter_url, "name": book.get("name", "")},
        )

        add_edge(
            from_id=root_node["id"],
            label="HAS_BOOK",
            to_id=book_node["id"],
            weight=1,
            properties={"url": url},
        )

        try:
            logging.info("fetching chapter %s", chapter_url)
            resp = requests.get(chapter_url)
            chap_data = resp.json()
        except Exception as e:
            logging.exception("error fetching chapter")
            continue

        for c in chap_data.get("chapters", []):
            verse_url = c.get("url", "")

            chap_node = add_node(
                label="chapter",
                properties={"url": verse_url, "chapter": int(c.get("chapter", 0))},
            )

            add_edge(
                from_id=book_node["id"],
                label="HAS_CHAPTER",
                to_id=chap_node["id"],
                weight=1,
                properties={"url": chapter_url},
            )

            try:
                logging.info("fetching verse %s", verse_url)
                resp = requests.get(chapter_url)
                verse_data = resp.json()
            except Exception as e:
                logging.exception("error fetching verse")
                continue

            for v in verse_data.get("verses", []):
                verse_num = int(v.get("verse", 0))
                text = v.get("text", "")

                verse_node = add_node(
                    label="verse",
                    properties={"verse": verse_num, "content": text},
                )

                add_edge(
                    from_id=chap_node["id"],
                    label="HAS_VERSE",
                    to_id=verse_node["id"],
                    weight=1,
                    properties={"url": chap_url},
                )


# Function to upload the graph to the server
def upload_graph(url: str):
    headers = {"Content-Type": "application/json"}
    response = requests.put(url, json=graph, headers=headers)
    if response.status_code == 200:
        print("Graph uploaded successfully.")
    else:
        print(f"Failed to upload graph: {response.status_code}, {response.text}")


# Main function to run the script
def main():
    print("Generating Bible graph...")
    fetch(BIBLE_API_URL)
    print("Uploading graph to server...")
    upload_graph(GRAPH_API_URL)


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    main()
