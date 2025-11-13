# raumzeitalpaka

<img src="app/ports/www/static/timekeeper.png" width="161">

----

raumzeitalpaka (formerly timekeeper) is an open source system for managing event schedules. 
It is created and optimized for [Jugend hackt](https://jugendhackt.org/) events. 

It provides basic mechanisms for the creation and management of time schedules, locations and rooms for events.
It can also export schedules as Calendar (ical), [VOC Schedule](https://github.com/voc/schedule/blob/master/validator/json/schema.json) (Info Beamer) and Markdown tables.

<div>
<img src="app/ports/www/static/pixelhack/flag_pride.svg" width="69">
<img src="app/ports/www/static/pixelhack/resistor_pride.svg" width="45">
<img src="app/ports/www/static/pixelhack/flag_nonbinary.svg" width="69">
<img src="app/ports/www/static/pixelhack/resitor_nonbinary.svg" width="45">
<img src="app/ports/www/static/pixelhack/flag_trans.svg" width="69">
<img src="app/ports/www/static/pixelhack/resistor_trans.svg" width="45">
<img src="app/ports/www/static/pixelhack/reisealpaka.svg" width="50">
<img src="app/ports/www/static/pixelhack/resistor_trans.svg" width="45">
<img src="app/ports/www/static/pixelhack/flag_trans.svg" width="69">
<img src="app/ports/www/static/pixelhack/resitor_nonbinary.svg" width="45">
<img src="app/ports/www/static/pixelhack/flag_nonbinary.svg" width="69">
<img src="app/ports/www/static/pixelhack/resistor_pride.svg" width="45">
<img src="app/ports/www/static/pixelhack/flag_pride.svg" width="69">
</div>


# Get Started

## Prerequisites
- Postgresql

## Container
```shell
podman run ghcr.io/m4schini/raumzeitalpaka:latest \
  -p 8080:80 \
  --env TIMEKEEPER_BASE_URL= \
  --env DATABASE_CONNECTIONSTRING= \
  --env JWT_SECRET= 

```

## Helm
```shell
helm install raumzeitalpaka ./helm
```

## Oldschool

# Config

| Environment Variable         | Required | Default       | Description                                                             |
|------------------------------|----------|---------------|-------------------------------------------------------------------------|
| DATABASE_CONNECTIONSTRING    | REQUIRED |               |                                                                         |
| JWT_SECRET                   | REQUIRED |               |                                                                         |
| TIMEKEEPER_TIMEZONE          |          | Europe/Berlin |                                                                         |
| TIMEKEEPER_TELEMETRY_ENABLED |          | false         | Enables prometheus metrics endpoint at /metrics and debug level logging |
| PORT                         |          | 80            |                                                                         |
| TIMEKEEPER_BASE_URL          |          |               |                                                                         |
| TIMEKEEPER_ADMIN_PASSWORD    |          |               |                                                                         |

# Attributions
## Bundled pixelart: PixelHack
CC-BY-SA 4.0 Jugend Hackt (Hanno Sternberg)

[License](http://creativecommons.org/licenses/by-sa/4.0/)

## Bundled font: Atkinson Hyperlegible Next
Copyright 2020, Braille Institute of America, Inc. (https://www.brailleinstitute.org/), with Reserved Font Names: “ATKINSON” and “HYPERLEGIBLE”.
This Font Software is licensed under the SIL Open Font License, Version 1.1. This license is copied below, and is also available with a FAQ at: https://openfontlicense.org

[License](https://codeberg.org/aur0ra/timekeeper/blob/8a312ed579c1be3845f79f92294fee7c61a771cc/ports/www/static/font/Atkinson-Hyperlegible-SIL-OPEN-FONT-LICENSE-Version%201.1-v2%20ACC.pdf)