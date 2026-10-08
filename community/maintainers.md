---
layout: page
title: Maintainers
nav-title: Maintainers
---

Gemara is maintained by a group of contributors across the community. These
folks review contributions, guide the project's direction, and are happy to
help you get started.

{% for maintainer in site.data.maintainers.maintainers %}
- {{ maintainer.name }}, {{ maintainer.organization }} (@{{ maintainer.github }})
{% endfor %}

Want to get involved? See the [Contributing guide](/community/contributing)
or come find us on [Slack](https://openssf.slack.com/archives/C09A9PP765Q).
