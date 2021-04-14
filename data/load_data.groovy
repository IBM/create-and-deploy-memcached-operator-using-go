graph = TinkerGraph.open()
// graph.createIndex('airport', Vertex.class) //1

g = traversal().withEmbedded(graph)

g.addV('airport').property('object_type', 'airport').property(
            'id', '9600276f-608f-4325-a037-f185848f2e28'
        ).property('name', 'Los Angeles International Airport').property(
            'is_hub', true
        ).property(
            'is_destination', true
        ).property(
            'type', 'large_airport'
        ).property(
            'country', 'United States'
        ).property(
            'city', 'Los Angeles'
        ).property(
            'latitude', 33.94250107
        ).property(
            'longitude', -118.4079971
        ).property(
            'gps_code', 'KLAX'
        ).property(
            'iata_code', 'LAX'
        ).next()
g.addV('airport').property('object_type', 'airport').property(
            'id', 'ebc645cd-ea42-40dc-b940-69456b64d2dd'
        ).property('name', 'John F. Kennedy International Airport').property(
            'is_hub', true
        ).property(
            'is_destination', true
        ).property(
            'type', 'large_airport'
        ).property(
            'country', 'United States'
        ).property(
            'city', 'New York'
        ).property(
            'latitude', 40.63980103
        ).property(
            'longitude', -73.77890015
        ).property(
            'gps_code', 'KJFK'
        ).property(
            'iata_code', 'JFK'
        ).next()
g.addV('airport').property('object_type', 'airport').property(
            'id', '30a2b6e8-fcc2-4e59-a61d-3aff713b23b0'
        ).property('name', 'Dallas Fort Worth International Airport').property(
            'is_hub', true
        ).property(
            'is_destination', false
        ).property(
            'type', 'large_airport'
        ).property(
            'country', 'United States'
        ).property(
            'city', 'Dallas-Fort Worth'
        ).property(
            'latitude', 32.896801
        ).property(
            'longitude', -97.038002
        ).property(
            'gps_code', 'KDFW'
        ).property(
            'iata_code', 'DFW'
        ).next()
g.addV('airport').property('object_type', 'airport').property(
            'id', '6275b58b-55a8-4c3f-93fa-372529ee0b2f'
        ).property('name', 'Lester B. Pearson International Airport').property(
            'is_hub', true
        ).property(
            'is_destination', true
        ).property(
            'type', 'large_airport'
        ).property(
            'country', 'Canada'
        ).property(
            'city', 'Toronto'
        ).property(
            'latitude', 43.6772003174
        ).property(
            'longitude', -79.63059997559999
        ).property(
            'gps_code', 'YYZ'
        ).property(
            'iata_code', 'YYZ'
        ).next()
g.addV('flight').property('object_type', 'flight').property(
            'id', 'e7c3d85d-c523-4634-93ef-a84f55aeb1e5'
        ).property(
            'source_airport_id', '9600276f-608f-4325-a037-f185848f2e28'
        ).property(
            'destination_airport_id', 'ebc645cd-ea42-40dc-b940-69456b64d2dd'
        ).property(
            'flight_time', 345
        ).property(
            'flight_duration', 324.384941546607
        ).property(
            'cost', 584.833849847819
        ).property(
            'airlines', 'MilkyWay Airlines'
        ).next()
g.addV('flight').property('object_type', 'flight').property(
            'id', 'fa3448ff-f157-4690-9180-0e06700ac909'
        ).property(
            'source_airport_id', '9600276f-608f-4325-a037-f185848f2e28'
        ).property(
            'destination_airport_id', 'ebc645cd-ea42-40dc-b940-69456b64d2dd'
        ).property(
            'flight_time', 945
        ).property(
            'flight_duration', 324.384941546607
        ).property(
            'cost', 483.08914733391794
        ).property(
            'airlines', 'Phoenix Airlines'
        ).next()
g.addV('flight').property('object_type', 'flight').property(
            'id', 'd76b8e14-3519-42eb-84ff-9a7406b43234'
        ).property(
            'source_airport_id', '9600276f-608f-4325-a037-f185848f2e28'
        ).property(
            'destination_airport_id', '30a2b6e8-fcc2-4e59-a61d-3aff713b23b0'
        ).property(
            'flight_time', 615
        ).property(
            'flight_duration', 176.90181296430302
        ).property(
            'cost', 345.5986288722164
        ).property(
            'airlines', 'Spartan Airlines'
        ).next()
g.addV('flight').property('object_type', 'flight').property(
            'id', 'cd29ac89-1b5e-4f8c-8e7c-d53404e6b092'
        ).property(
            'source_airport_id', '30a2b6e8-fcc2-4e59-a61d-3aff713b23b0'
        ).property(
            'destination_airport_id', 'ebc645cd-ea42-40dc-b940-69456b64d2dd'
        ).property(
            'flight_time', 1080
        ).property(
            'flight_duration', 195.54786163329285
        ).property(
            'cost', 368.5532617453382
        ).property(
            'airlines', 'Spartan Airlines'
        ).next()
g.addV('flight').property('object_type', 'flight').property(
            'id', '0656f352-8903-406f-a55c-e83e70028302'
        ).property(
            'source_airport_id', 'ebc645cd-ea42-40dc-b940-69456b64d2dd'
        ).property(
            'destination_airport_id', '6275b58b-55a8-4c3f-93fa-372529ee0b2f'
        ).property(
            'flight_time', 1350
        ).property(
            'flight_duration', 73.59971345361168
        ).property(
            'cost', 221.63444620384936
        ).property(
            'airlines', 'Spartan Airlines'
        ).next()

lax = g.V().has('iata_code', 'LAX').next()
jfk = g.V().has('iata_code', 'JFK').next()
dfw = g.V().has('iata_code', 'DFW').next()
yyz = g.V().has('iata_code', 'YYZ').next()
f1 = g.V().has('id', 'e7c3d85d-c523-4634-93ef-a84f55aeb1e5').next()
f2 = g.V().has('id', 'fa3448ff-f157-4690-9180-0e06700ac909').next()
f3 = g.V().has('id', 'd76b8e14-3519-42eb-84ff-9a7406b43234').next()
f4 = g.V().has('id', 'cd29ac89-1b5e-4f8c-8e7c-d53404e6b092').next()
f5 = g.V().has('id', '0656f352-8903-406f-a55c-e83e70028302').next()

g.addE('departing').from(f1).to(lax).iterate()
g.addE('arriving').from(f1).to(jfk).iterate()
g.addE('departing').from(f2).to(lax).iterate()
g.addE('arriving').from(f2).to(jfk).iterate()
g.addE('departing').from(f3).to(lax).iterate()
g.addE('arriving').from(f3).to(dfw).iterate()
g.addE('departing').from(f4).to(dfw).iterate()
g.addE('arriving').from(f4).to(jfk).iterate()
g.addE('departing').from(f5).to(jfk).iterate()
g.addE('arriving').from(f5).to(yyz).iterate()