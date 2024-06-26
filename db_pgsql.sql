-- Active: 1713865706485@@127.0.0.1@5432@postgres@public
-- Generated by the database client.
CREATE TABLE geofences(
                          id SERIAL NOT NULL,
                          name varchar(255),
                          status smallint,
                          created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                          updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                          city_name varchar(255),
                          city_coords USER-DEFINED,
                          boundary USER-DEFINED,
                          manager_ids integer[],
                          alert_area USER-DEFINED,
                          description text,
                          alert_distans integer,
                          PRIMARY KEY(id)
);
CREATE INDEX idx_geofences_boundary ON geofences USING gist ("boundary");