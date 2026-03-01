import requests
import json
import time

KJV_URL = "https://raw.githubusercontent.com/thiagobodruk/bible/master/json/en_kjv.json"
API_URL = "http://localhost:8080/api/v1/upload"

BATCH_SIZE = 500


def download_bible():
    print("Downloading KJV Bible...")
    response = requests.get(KJV_URL, timeout=30)
    response.raise_for_status()
    return json.loads(response.content.decode("utf-8-sig"))


def generate_graph(bible_json):
    node_id = 0

    for book_order, book in enumerate(bible_json, start=1):
        book_id = node_id
        testament = "Old" if book_order <= 39 else "New"

        # Book node
        yield {
            "type": "node",
            "data": {
                "id": book_id,
                "label": "letter",
                "properties": {
                    "name": book["name"],
                    "testament": testament,
                    "order": book_order,
                    "chapter_count": len(book["chapters"]),
                    "type": "book",
                },
            },
        }
        node_id += 1

        for chapter_index, chapter in enumerate(book["chapters"], start=1):
            chapter_id = node_id

            # Chapter node
            yield {
                "type": "node",
                "data": {
                    "id": chapter_id,
                    "label": "chapter",
                    "properties": {
                        "book": book["name"],
                        "chapter": chapter_index,
                        "testament": testament,
                        "bookOrder": book_order,
                        "verse_count": len(chapter),
                        "reference": f"{book['name']} {chapter_index}",
                        "type": "chapter",
                    },
                },
            }
            node_id += 1

            # Edge: Book → Chapter
            yield {
                "type": "edge",
                "data": {
                    "from_id": book_id,
                    "label": "contains",
                    "to_id": chapter_id,
                    "weight": 1,
                    "properties": {
                        "relation": "book_to_chapter"
                    },
                },
            }

            for verse_index, verse_text in enumerate(chapter, start=1):
                verse_id = node_id

                # Verse node
                yield {
                    "type": "node",
                    "data": {
                        "id": verse_id,
                        "label": "verse",
                        "properties": {
                            "book": book["name"],
                            "chapter": chapter_index,
                            "verse": verse_index,
                            "testament": testament,
                            "bookOrder": book_order,
                            "reference": f"{book['name']} {chapter_index}:{verse_index}",
                            "content": verse_text,
                            "word_count": len(verse_text.split()),
                            "char_count": len(verse_text),
                            "type": "verse",
                        },
                    },
                }
                node_id += 1

                # Edge: Chapter → Verse
                yield {
                    "type": "edge",
                    "data": {
                        "from_id": chapter_id,
                        "label": "contains",
                        "to_id": verse_id,
                        "weight": 1,
                        "properties": {
                            "relation": "chapter_to_verse"
                        },
                    },
                }

                # Optional: Book → Verse (fast traversal)
                yield {
                    "type": "edge",
                    "data": {
                        "from_id": book_id,
                        "label": "contains",
                        "to_id": verse_id,
                        "weight": 1,
                        "properties": {
                            "relation": "book_to_verse"
                        },
                    },
                }


def upload_batch(nodes, edges):
    payload = {
        "nodes": nodes,
        "edges": edges,
    }

    response = requests.put(
        API_URL,
        headers={"Content-Type": "application/json"},
        data=json.dumps(payload),
        timeout=60,
    )

    if response.status_code >= 400:
        print("Upload failed:", response.status_code)
        print(response.text)
        raise Exception("Upload failed")

    print(f"Uploaded {len(nodes)} nodes and {len(edges)} edges")


def upload_in_batches(generator):
    nodes_batch = []
    edges_batch = []

    for item in generator:
        if item["type"] == "node":
            nodes_batch.append(item["data"])
        else:
            edges_batch.append(item["data"])

        if len(nodes_batch) >= BATCH_SIZE:
            upload_batch(nodes_batch, edges_batch)
            nodes_batch.clear()
            edges_batch.clear()
            time.sleep(0.2)

    if nodes_batch:
        upload_batch(nodes_batch, edges_batch)


def main():
    bible_json = download_bible()
    graph_generator = generate_graph(bible_json)
    upload_in_batches(graph_generator)
    print("Done uploading full Bible graph.")


if __name__ == "__main__":
    main()