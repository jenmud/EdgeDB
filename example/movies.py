import sys
import json
import requests
import time

TMDB_API_KEY = "0c8aad208a02d7f489c3420d04d02c44"

UPLOAD_URL = "http://localhost:8080/api/v1/upload"

session = requests.Session()

node_id = 1
edge_id = 1

nodes = []
edges = []

person_map = {}
genre_map = {}


def new_node(label, props):
    global node_id

    node = {
        "id": node_id,
        "label": label,
        "Properties": props
    }

    node_id += 1
    nodes.append(node)

    return node


def new_edge(src, label, dst, props=None):
    global edge_id

    edge = {
        "id": edge_id,
        "from": src,
        "label": label,
        "to": dst,
        "Properties": props or {}
    }

    edge_id += 1
    edges.append(edge)


def get_or_create_person(person):

    pid = person["id"]

    if pid in person_map:
        return person_map[pid]

    node = new_node(
        "actor",
        {
            "name": person["name"],
            "tmdb_id": pid,
            "known_for_department": person.get("known_for_department")
        }
    )

    person_map[pid] = node["id"]
    return node["id"]


def get_or_create_director(person):

    pid = person["id"]

    if pid in person_map:
        return person_map[pid]

    node = new_node(
        "director",
        {
            "name": person["name"],
            "tmdb_id": pid
        }
    )

    person_map[pid] = node["id"]

    return node["id"]


def get_or_create_genre(genre):

    gid = genre["id"]

    if gid in genre_map:
        return genre_map[gid]

    node = new_node(
        "genre",
        {
            "name": genre["name"],
            "tmdb_id": gid
        }
    )

    genre_map[gid] = node["id"]

    return node["id"]


def fetch_movies(page):

    url = "https://api.themoviedb.org/3/movie/popular"

    r = session.get(
        url,
        params={
            "api_key": TMDB_API_KEY,
            "page": page
        }
    )

    return r.json()["results"]


def fetch_movie_details(movie_id):

    url = f"https://api.themoviedb.org/3/movie/{movie_id}"

    r = session.get(
        url,
        params={"api_key": TMDB_API_KEY}
    )

    return r.json()


def fetch_credits(movie_id):

    url = f"https://api.themoviedb.org/3/movie/{movie_id}/credits"

    r = session.get(
        url,
        params={"api_key": TMDB_API_KEY}
    )

    return r.json()


def import_movie(movie):

    details = fetch_movie_details(movie["id"])
    credits = fetch_credits(movie["id"])

    movie_node = new_node(
        "movie",
        {
            "title": details["title"],
            "release_date": details.get("release_date"),
            "rating": details.get("vote_average"),
            "vote_count": details.get("vote_count"),
            "tmdb_id": details["id"],
            "overview": details.get("overview")
        }
    )

    movie_id = movie_node["id"]

    # genres
    for g in details.get("genres", []):

        gid = get_or_create_genre(g)

        new_edge(
            movie_id,
            "HAS_GENRE",
            gid
        )

    # actors
    for cast in credits.get("cast", [])[:10]:

        actor_id = get_or_create_person(cast)

        new_edge(
            actor_id,
            "ACTED_IN",
            movie_id,
            {
                "character": cast.get("character"),
                "billing_order": cast.get("order")
            }
        )

    # director
    for crew in credits.get("crew", []):

        if crew["job"] == "Director":

            director_id = get_or_create_director(crew)

            new_edge(
                director_id,
                "DIRECTED",
                movie_id
            )


def upload_graph(data):

    print("Uploading graph...")

    r = session.put(UPLOAD_URL, json=data)

    if r.status_code not in (200, 201):
        print("Upload failed:", r.text)
    else:
        print("Upload successful")


def main(pages=50):

    for page in range(1, pages + 1):

        print("Fetching page", page)

        movies = fetch_movies(page)

        for movie in movies:

            try:
                import_movie(movie)
                #time.sleep(0.2)
            except Exception as e:
                print("error:", e)

    graph = {
        "nodes": nodes,
        "edges": edges
    }

    print("Nodes:", len(nodes))
    print("Edges:", len(edges))

    with open("movies.json", "w") as f:
        json.dump(graph, f)

    upload_graph(graph)


main(int(sys.argv[1]))