INSERT INTO timekeeper.locations (name, file) VALUES ('Betahaus | Schanze', 'betahaus.png');

INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Hackspace', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Makespace', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Ruhe Räume', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Infodesk', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Cafe', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Tresen', 0, 0, 0, 0);
INSERT INTO timekeeper.rooms (location, name, location_x, location_y, location_w, location_h)
VALUES (1, 'Void', 0, 0, 0, 0);

INSERT INTO timekeeper.events (name, start) VALUES ('Jugend hackt Hamburg 2025', date(now()));

-- Friday
INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Check in Hostel', '', 0, '15:00:00', 4, 'Organizer');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Infodesk besetzt', '', 0, '16:00:00', 4, 'Organizer');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Ankunft', '', 0, '16:30:00', 4, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Mentor*innen Begrüßung', 'Erstes Meeting', 0, '16:30:00', 7, 'Mentor');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Eröffnung', '', 0, '17:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Input Talks', 'Findet parallel in verschiedenen räumen statt', 0, '17:30:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Abendessen', '', 0, '18:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Brainstorming', 'Findet parallel in verschiedenen räumen statt', 0, '19:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Vorstellung der Ideen', '', 0, '21:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Tagesabschluss', 'Zulip, Hardware + Projektidee Gallery\nGruppenfoto\nSpiel & Spaß, andere Lustige Abendgestaltung\nTagesabschluss, Gruppenfoto und ins Bett', 0, '22:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Aufbruch ins Hostel', '', 0, '22:30:00', 4, 'Participant');

-- Saturday
INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Frühstück', '', 1, '08:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Tageseröffnung', '', 1, '09:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Gruppenfindung', '', 1, '09:30:00', 1, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Arbeiten in gruppen', '', 1, '10:00:00', 1, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Mittagessen', '', 1, '13:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Abendessen', '', 1, '18:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Präsentationstraining', '', 1, '19:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Gemeinsamer Tagesabschluss', '', 1, '22:00:00', 5, 'Participant');

--Sunday
INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Frühstück', '', 2, '08:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Tageseröffnung', 'Evaluation', 2, '09:00:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Mittagssnack', '', 2, '12:30:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Abschlusspräsentation', '', 2, '13:30:00', 5, 'Participant');

INSERT INTO timekeeper.timeslots (event, title, note, day, start, room, role)
VALUES (1, 'Ende der Veranstaltung', '', 2, '15:00:00', 5, 'Participant');
