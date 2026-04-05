import requests
import json
import re
from collections import defaultdict

# Define the API URL for uploading the graph
GRAPH_API_URL = "http://localhost:8080/api/v1/graph"

# Define the URL for downloading the Bible text
BIBLE_API_URL = "https://bible-api.com/data/kjv"

nodeID = 1
edgeID = 1

# nextNodeID returns the next node ID.
def nextNodeID() -> int:
    global nodeID
    n = nodeID
    nodeID += 1
    return n

# nextEdgeID returns the next edge ID.
def nextEdgeID() -> int:
    global edgeID
    e = edgeID
    edgeID += 1
    return e


# A dictionary to store the graph structure
graph = {
    "nodes": [],
    "edges": []
}

# Helper function to add a node to the graph
def add_node(node_id: int, label: str, properties: dict = {}) -> dict:
    if not node_id or node_id <= 0:
        node_id = nextNodeID()

    if properties is None:
        properties = {}

    node = {
        "id": node_id,
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
        "id": nextEdgeID(),
        "from_id": from_id,
        "label": label,
        "to_id": to_id,
        "weight": weight,
        "properties": properties
    }

    graph['edges'].append(edge)
    return edge


# add in the root node
rootNode = add_node(0, "bible", {})

def fetch_books():
    resp = requests.get(BIBLE_API_URL)
    data = resp.json()

    for book in data.get("books", []):
        n = add_node(
            node_id=0,
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
