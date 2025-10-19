# timekeeper

<img src="https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png" width="100">

----

timekeeper is an open source system for managing event schedules. 
It is created and optimized for [Jugend hackt](https://jugendhackt.org/) events. 

It provides basic mechanisms for the creation and management of time schedules, locations and rooms for events.
It can also export schedules as Calendar (ical), [VOC Schedule](https://github.com/voc/schedule/blob/master/validator/json/schema.json) (Info Beamer) and Markdown tables.

----

## To start using timekeeper

To use Kubernetes code as a library in other applications, see the [list of published components](https://git.k8s.io/kubernetes/staging/README.md).
Use of the `k8s.io/kubernetes` module or `k8s.io/kubernetes/...` packages as libraries is not supported.

## To start developing K8s

The [community repository] hosts all information about
building Kubernetes from source, how to contribute code
and documentation, who to contact about what, etc.

If you want to build Kubernetes right away there are two options:

##### You have a working [Go environment].

```
git clone https://github.com/kubernetes/kubernetes
cd kubernetes
make
```

##### You have a working [Docker environment].

```
git clone https://github.com/kubernetes/kubernetes
cd kubernetes
make quick-release
```

For the full story, head over to the [developer's documentation].

