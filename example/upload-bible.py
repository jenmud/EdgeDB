import requests
import json
import re
import logging
import time
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


# Add this helper function for retry logic
def fetch_with_retry(url: str, max_retries: int = 5, initial_wait: int = 1):
    retries = 0
    wait_time = initial_wait  # start with 1 second

    while retries < max_retries:
        try:
            logging.info("Fetching URL: %s", url)
            resp = requests.get(url)
            
            # If status is 200, return the response data
            if resp.status_code == 200:
                return resp.json()

            # If status is 429, handle rate limiting
            elif resp.status_code == 429:
                retries += 1
                logging.warning("Rate limit exceeded. Retrying in %d seconds...", wait_time)
                time.sleep(wait_time)  # wait for a while before retrying
                wait_time *= 2  # exponential backoff
            else:
                # For other errors, log and break
                logging.error("Failed to fetch data, status code: %d", resp.status_code)
                break
        
        except requests.RequestException as e:
            logging.exception("Error fetching data")
            retries += 1
            time.sleep(wait_time)  # wait before retrying
            wait_time *= 2  # exponential backoff

    return None  # return None if all retries fail


def fetch(url: str):
    logging.info("fetching book %s", url)
    data = fetch_with_retry(url)

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

        logging.info("fetching chapter %s", chapter_url)
        chap_data = fetch_with_retry(chapter_url)

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

            logging.info("fetching verse %s", verse_url)
            verse_data = fetch_with_retry(verse_url)

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
                    properties={"url": chapter_url},
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
