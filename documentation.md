---
layout: page
title: Documentation
description: Find Gemara guides, schema and SDK reference, and explanations of the GRC engineering model.
permalink: /docs/
---

The Gemara documentation is organized by what you need right now: doing a task, looking something up, or understanding how the model fits together.

{% for section in site.data.sidebar %}
## {{ section.title }}

{{ section.description }}

{% for item in section.children -%}
- [{{ item.title }}]({{ item.url | relative_url }})
{% endfor %}
{% endfor %}
