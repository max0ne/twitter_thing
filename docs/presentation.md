<!-- $theme: default -->
<!-- page_number: true -->

twitter_thing
===

Qiming Zhang | qz718
Mingfei Huang | mh4925

---

# Features
1. User Profile (register / login / logout / unregister / get recent users)
2. Follow (follow / unfollow)
3. Tweet (new / delete / get user feed)

---

Architecture
![Architecture](Distributed_System.png)

---

# Implementation | VR

- Strong consistency model
- No single point of failure - can recover as long as quorum alive
- ~~code already there~~

---

# Implementation | DB

- Built on top of VR, run as a VR node
- Implement a set of generic key-value operations
	- get, set, del
- VR replica forwards DB requests to VR primary
- VR primary send write operations as VR command
	- set_command, del_command

---

# Implementation | API

- Host RESTful endpoint
- Implement all application logic
- Read / write from one pre-configured DB instance
- A stateless service

---

# Implementation | UI

- Some JS code
- Connect a user configured API endpoint

---

# DEMO