![Build status](https://github.com/oam-dev/kubevela/workflows/E2E/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/oam-dev/kubevela)](https://goreportcard.com/report/github.com/oam-dev/kubevela)
![Docker Pulls](https://img.shields.io/docker/pulls/oamdev/vela-core)
[![codecov](https://codecov.io/gh/oam-dev/kubevela/branch/master/graph/badge.svg)](https://codecov.io/gh/oam-dev/kubevela)
[![LICENSE](https://img.shields.io/github/license/oam-dev/kubevela.svg?style=flat-square)](/LICENSE)
[![Releases](https://img.shields.io/github/release/oam-dev/kubevela/all.svg?style=flat-square)](https://github.com/oam-dev/kubevela/releases)
[![TODOs](https://img.shields.io/endpoint?url=https://api.tickgit.com/badge?repo=github.com/oam-dev/kubevela)](https://www.tickgit.com/browse?repo=github.com/oam-dev/kubevela)
[![Twitter](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Ftwitter.com%2Foam_dev)](https://twitter.com/oam_dev)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubevela)](https://artifacthub.io/packages/search?repo=kubevela)

![alt](docs/resources/KubeVela-03.png)

*Make shipping applications more enjoyable.*

# KubeVela

KubeVela is the platform engine to create *PaaS-like* experience on Kubernetes, in a scalable approach.

## Community

- Slack:  [CNCF Slack](https://slack.cncf.io/) #kubevela channel
- Gitter: [Discussion](https://gitter.im/oam-dev/community)
- Bi-weekly Community Call: [Meeting Notes](https://docs.google.com/document/d/1nqdFEyULekyksFHtFvgvFAYE-0AMHKoS3RMnaKsarjs)

> NOTE: KubeVela is still iterating quickly. It's currently under pre-beta release.

## What problems does it solve?

Building **developer-centric platforms** with Kubernetes requires higher level primitives which is out-of-scope of Kubernetes itself. Hence, we platform teams build abstractions.

However, great in flexibility and extensibility, the existing solutions such as IaC (Infrastructure-as-Code) and client-side templating tools all lead to ***Configuration Drift*** (i.e. the generated instances are not in line with the expected configuration) which is a nightmare in production.

KubeVela allows platform teams to create developer-centric abstractions with IaC but maintain them with the battle tested [Kubernetes Control Loop](https://kubernetes.io/docs/concepts/architecture/controller/). Think about a plug-in that turns your Kubernetes cluster into a *"PaaS"* via abstractions designed by yourself.

## Getting Started

- [Installation](https://kubevela.io/docs/install)
- [Quick start](https://kubevela.io/docs/quick-start)
- [How it works](https://kubevela.io/docs/concepts)

## Features

- **Robust, repeatable and extensible approach to create and maintain abstractions** - design your abstractions with [CUE](https://cuelang.org/) or [Helm](https://helm.sh), ship them to end users by `kubectl apply -f`, automatically generating GUI forms, upgrade your abstractions at runtime, and let Kubernetes controller guarantee determinism of the abstractions, no configuration drift.
- **Generic progressive rollout framework** - built-in rollout framework and strategies to upgrade your microservice regardless of its workload type (e.g. stateless, stateful, or even custom operators etc), seamless integration with observability systems.
- **Multi-enviroment app delievry model (WIP)** - built-in model to deliver or rollout your apps across multiple enviroments and/or clusters, seamless integration with Service Mesh for traffic management. 
- **Simple and Kubernetes native** - KubeVela is just a simple custom controller, all its app delivery abstractions and features are defined as [Kubernetes Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) so they naturally work with any CI/CD or GitOps tools.

## Documentation

Visit the [KubeVela documentation site](https://kubevela.io/) to find *Installation Instruction*, *Platform Builder Guide* and *Developer Experience Guide*.

## Talks and Conferences

| Engagement | Link        |
|:-----------|:------------|
| 🎤  Talks | - [KubeVela - The Modern App Delivery System in Alibaba](https://docs.google.com/presentation/d/1CWCLcsKpDQB3bBDTfdv2BZ8ilGGJv2E8L-iOA5HMrV0/edit?usp=sharing) |
| 🌎 KubeCon | - [ [NA 2020] Standardizing Cloud Native Application Delivery Across Different Clouds](https://www.youtube.com/watch?v=0yhVuBIbHcI) <br> - [ [EU 2021] Zero Pain Microservice Development and Deployment with Dapr and KubeVela](https://sched.co/iE4S) |
| 📺 Conferences | - [Dapr, Rudr, OAM: Mark Russinovich presents next gen app development & deployment](https://www.youtube.com/watch?v=eJCu6a-x9uo) <br> - [Mark Russinovich presents "The Future of Cloud Native Applications with OAM and Dapr"](https://myignite.techcommunity.microsoft.com/sessions/82059)|

## Contributing
Check out [CONTRIBUTING](./CONTRIBUTING.md) to see how to develop with KubeVela.

## Code of Conduct
KubeVela adopts [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).
