import requests
import json
import re
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
rootNode = add_node("bible", {})

def fetch_books():
    resp = requests.get(BIBLE_API_URL)
    data = resp.json()

    for book in data.get("books", []):
        n = add_node(
            label="book",
            properties={"url": book.get("url", ""), "name": book.get("name", "")},
        )

        add_edge(
            from_id=rootNode["id"],
            label="HAS_BOOK",
            to_id=n["id"],
            weight=1,
            properties={"url": BIBLE_API_URL},
        )


# Function to upload the graph to the server
def upload_graph():
    headers = {"Content-Type": "application/json"}
    response = requests.put(GRAPH_API_URL, json=graph, headers=headers)
    if response.status_code == 200:
        print("Graph uploaded successfully.")
    else:
        print(f"Failed to upload graph: {response.status_code}, {response.text}")


# Main function to run the script
def main():
    print("Generating Bible graph...")
    fetch_books()
    print("Uploading graph to server...")
    upload_graph()


if __name__ == "__main__":
    main()
