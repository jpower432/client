Design: Collection Workflow
===
<!--toc-->
- [Design: Collection Workflow](#design-collection-workflow)
- [Collection Publishing](#collection-publishing)
- [Collection Pulling](#collection-pulling)
- [Anatomy of a Collection](#anatomy-of-a-collection)
- [Collection Lookup (Proposed)](#collection-lookup-proposed)
- [Event Engine Startup (Proposed)](#event-engine-startup-proposed)
- [Collection Processing (Proposed)](#collection-processing-proposed)

<!-- tocstop -->

# Collection Publishing

The workflow for collection publishing is very similar to the workflow used to build container images with most tooling options.
If the schema is published before the collection and associated to the collection, the specified attributes in the dataset configuration
will be validated against the schema during collection building.
This is demonstrated below by the diagram.

```mermaid
graph LR;
A[Build]-->B{Collection or Schema?};
B -- Schema --> C[Build Schema]
B -- Collection --> D[Build Collection]
C --> E[Push Schema]
E --> D
D --> F[Push Collection]
```

This would be the workflow if the ultimate goal was to publish a collection with or without a schema. Schema can be published
without a subsequent collection publish for later use as well.

```mermaid
graph LR;
A[Build]-->B{Collection or Schema?};
B -- Schema --> C[Build Schema]
B -- Collection --> D[Build Collection]
C --> E[Push]
D --> E
```

# Collection Pulling

Collections can be pulled as an entire OCI Artifact or filtered by an Attribute Query. The filtered OCI artifact is stored
in the build cache with the original manifest intact (sparse manifest) and the non-matching blobs (files) are not pulled into the cache.
All matching files are written to the cache and written to a user specified location.

The use of sparse manifest can pose a problem if re-tagging collections becomes part of the command line functionality in the future.
Some registries will reject manifests without all the blobs present. In this case, it may be of interest to reconstruct the manifest before pushing
and allow a flag to preserve the manifest, if desired.

# Anatomy of a Collection
> Note: The Collection and Schema representations depict current relationship made by the builder. The event engine is proposed.

```mermaid
graph TD;
    A[Collection Root] -->B[Collection Content];
    C[Schema Root] -->D[Schema Contents]
    A ---> C
    E[Event Engine Root] ---> F[Event Engine Logic]
    C ---> E
```
# Collection Lookup (Proposed)

```mermaid
sequenceDiagram
    participant Client
    participant Collection
    participant Schema
    participant Event Engine
    Client->>Collection: Where is your schema?
    Schema-->>Client: I'm here!
    Note right of Schema: Must validate<br>Collection!
    Client->>Schema: Where is the event engine?
    Schema-->>Client: "It's here!"
    Client->>Event Engine: "Pulling you now"
    Event Engine-->>Client: "On my way!"
```

# Event Engine Startup (Proposed)

```mermaid
sequenceDiagram
    participant User
    participant Client
    participant Event Engine
    participant Collection
    Note right of Client: "I know where the <br> Event Engine is!
    User->>Client: "Process the event with these arguments"
    Client->>Event Engine: "What do you need to start?"
    Event Engine-->>Client: "I need config file with extention=yaml"
    Client->>Collection: "Do you have content with a yaml extension?"
    Collection-->>Client: "Here is it!"
    Client-->>Event Engine: "Found it"
    Event Engine->>Event Engine: "Starting"
```


# Collection Processing (Proposed)

In order to ensure the Collection has all the required content 
for Event Engine processing, compatability between the Event Engine 
and Schema attributes is required. The Collection, upon building, is verified
against the schema to ensure it has the required keys and value types.

```mermaid
sequenceDiagram
    participant Event Engine
    participant Client
    participant Collection
     loop Execution
        Event Engine->>Event Engine: Execute entrypoint and proccess content
    end
    Event Engine->>Client: "Need content with attribute size=small"
    Client->>Collection: "Do you have anything with a small size?"
    Collection-->>Client: "Here you go!"
    Client-->>Event Engine: "Found it"
```



