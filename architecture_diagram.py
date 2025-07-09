#!/usr/bin/env python3
from diagrams import Diagram, Cluster, Edge, Node
from diagrams.programming.language import Go
from diagrams.onprem.inmemory import Redis
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.client import User, Client
from diagrams.onprem.network import Internet
from diagrams.programming.framework import React
from diagrams.generic.compute import Rack
from diagrams.oci.connectivity import CDN

# Set graph attributes
graph_attr = {
    "fontsize": "30",
    "bgcolor": "white",
    "rankdir": "TB",
    "pad": "2.0",
    "splines": "ortho",
    "nodesep": "0.8",
    "ranksep": "1.0"
}

# Create a diagram with a custom filename and title
with Diagram("CrawlerX Architecture", show=True, filename="crawlerx_architecture", 
             outformat="png", graph_attr=graph_attr):
    
    # External clients/users
    with Cluster("Frontend"):
        frontend = React("CrawlerX UI")
        
    # Create a cluster for the API layer
    with Cluster("API Layer"):
        api = Go("API Server")
        websocket = Internet("WebSocket")
    
    # Create a cluster for the job queue
    with Cluster("Job Queue"):
        redis = Redis("Redis")
    
    # Create a cluster for the workers
    with Cluster("Worker Pool"):
        workers = [Go("Worker 1"), Go("Worker 2"), Go("Worker 3")]
    
    # Create a cluster for the crawler
    with Cluster("Crawler Engine"):
        crawler = Rack("Crawler")
    
    # Create a cluster for the database
    with Cluster("Database"):
        db = PostgreSQL("PostgreSQL")
    
    # Create a cluster for the target websites
    with Cluster("Target Websites"):
        websites = CDN("Web")
    
    # Define the connections with labels and styles
    frontend >> Edge(color="black", style="bold") >> api
    frontend >> Edge(color="blue", style="dashed", label="real-time updates") >> websocket
    
    api >> Edge(color="red", label="enqueue jobs") >> redis
    
    for i, worker in enumerate(workers):
        redis >> Edge(color="red", label=f"job {i+1}") >> worker
        worker >> Edge(color="green") >> crawler
    
    crawler >> Edge(color="orange", style="dashed") >> websites
    crawler >> Edge(color="black", label="store results") >> db
    crawler >> Edge(color="blue", style="dashed", label="progress") >> websocket
    
    api >> Edge(color="black", label="query") >> db
