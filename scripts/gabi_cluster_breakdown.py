#!/usr/bin/env python3

import json
import os
import sys
import requests
from uuid import UUID

GABI_URL = os.getenv("GABI_URL", "")
GABI_TOKEN = os.getenv("GABI_TOKEN", "")

HEADERS = {"Authorization": f"Bearer {GABI_TOKEN}"}


def is_valid_uuid(uuid_to_test, version=4):
    try:
        uuid_obj = UUID(uuid_to_test, version=version)
    except ValueError:
        return False
    return str(uuid_obj) == uuid_to_test


def query(query):
    data = {"query": query}
    r = requests.get(GABI_URL, headers=HEADERS, json=data)
    if r.status_code == 200:
        return r.json()["result"]
    else:
        print(f"Query failed: {query}, HTTP code: {r.status_code}", file=sys.stderr)
        sys.exit(3)


def extract_image_shas(cluster_workload):
    images = {}
    for sha in cluster_workload["images"]:
        images[sha] = {"matched": None}
    for ns in cluster_workload["namespaces"].values():
        for shape in ns["shapes"]:
            for container in shape["containers"] or []:
                sha = container["imageID"]
                images[sha] = {"matched": None}
            for container in shape["initContainers"] or []:
                sha = container["imageID"]
                images[sha] = {"matched": None}
    return images


def match_image(row, sha_type, images):
    images[row[2]] = {"matched": sha_type,
                      "pyxis_id": row[1],
                      "name": f"{row[5]}/{row[6]}",
                      "cves": [x for x in json.loads(row[7]) if x]}


def match_images(cluster_image_rows, images):
    if len(cluster_image_rows) <= 1:
        return
    cluster_image_rows = cluster_image_rows[1:]
    for row in cluster_image_rows:
        if row[2] in images:
            match_image(row, "manifest_schema2_digest", images)
            continue
        elif row[3] in images:
            match_image(row, "manifest_list_digest", images)
            continue
        elif row[4] in images:
            match_image(row, "docker_image_digest", images)
            continue
        print(f"ERR: Image id not found: {row[0]}")


def main():
    if not GABI_URL or not GABI_TOKEN:
        print("GABI_URL or GABI_TOKEN env variable not defined!", file=sys.stderr)
        sys.exit(1)
    cluster_uuid = sys.argv[1]
    if len(sys.argv) != 2 or not is_valid_uuid(cluster_uuid):
        print(f"Usage: {sys.argv[0]} <cluster_id>", file=sys.stderr)
        sys.exit(2)

    print(f"Gabi URL: {GABI_URL}")
    print(f"Gabi token: ***")
    print(f"Cluster UUID: {cluster_uuid}")
    print("")

    cluster_rows = query(f"SELECT id, workload FROM cluster WHERE uuid = '{cluster_uuid}';")
    if len(cluster_rows) <= 1:
        print(f"Cluster not found: {cluster_uuid}", file=sys.stderr)
        sys.exit(4)

    cluster_id = cluster_rows[1][0]
    cluster_workload = json.loads(cluster_rows[1][1])

    images = extract_image_shas(cluster_workload)

    print("=========================")
    print("workload_info.json checks")
    print("=========================")
    print(f"Expected images: {cluster_workload['imageCount']}")
    print(f"Images found: {len(images)}")
    print(f"")

    expected_matches_cnt = query(f"SELECT COUNT(*) FROM cluster_image WHERE cluster_id = {cluster_id};")[1][0]
    cluster_image_rows = query(f"""
        SELECT i.id,
               i.pyxis_id,
               i.manifest_schema2_digest,
               i.manifest_list_digest,
               i.docker_image_digest,
               r.registry,
               r.repository,
               json_agg(cve.name)
        FROM cluster_image ci
        JOIN image i ON ci.image_id = i.id
        JOIN arch a on i.arch_id = a.id
        JOIN repository_image ri on i.id = ri.image_id
        JOIN repository r on ri.repository_id = r.id
        LEFT JOIN image_cve ic on i.id = ic.image_id
        LEFT JOIN cve on ic.cve_id = cve.id
        WHERE ci.cluster_id = {cluster_id}
          AND a.name = 'amd64'
        GROUP BY i.id,
                 i.pyxis_id,
                 i.manifest_schema2_digest,
                 i.manifest_list_digest,
                 i.docker_image_digest,
                 r.registry,
                 r.repository;
    """)
    match_images(cluster_image_rows, images)

    print("================")
    print("DB tables checks")
    print("================")
    print(f"Expected cluster_image rows: {expected_matches_cnt}")
    print(f"cluster_image rows validated: {len([x for x in images.values() if x['matched']])}")
    print(f"")
    print("======")
    print("Report")
    print("======")
    for image in sorted(images):
        details = images[image]
        print(f"{image} {details}")


if __name__ == "__main__":
    main()
