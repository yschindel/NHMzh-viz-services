# IFC Service

This service exists because:

We want to have a performant 3D viewing experience in the final dashboard. This means we need to convert the IFC file to something more compressed and faster to load in a webGL viewer.

Code available at: [https://github.com/yschindel/NHMzh-viz-services/tree/main/viz_ifc](https://github.com/yschindel/NHMzh-viz-services/tree/main/viz_ifc)

Responsibilites:

- Kafka consumer for notifications when a file is uploaded
- Convert the IFC file to fragments
- Send fragments file to another services via a POST request.
- Extract element properties from the IFC file.
- Send element data to another service via a POST request

## Kafka

This services is part of the core NHM infrastructure and is therefore connected to Kafka.
It consumes a Kafka topic, which is specified in the `KAFKA_IFC_TOPIC` environment variable.

## MinIO

The service will download a file from MinIO.

## Processing IFC Files

### Geometry

The IFC files are converted to fragments using [@ThatOpen/](https://www.npmjs.com/org/thatopen)
This enables their use in the web viewer also provided by ThatOpen.

### Properties

The element properties are extracted using 'web-ifc'.
Properties can be specified by setting the `IFC_PROPERTIES_TO_INCLUDE` environment variable:

```bash
IFC_PROPERTIES_TO_INCLUDE=Reference,LoadBearing,FireRating,Name,category,level,ebkph
```
