# 06: Audit Logging & SIEM Integration 📜

In enterprise security, blocking an attack is only half the battle. The other half is **recording** it. 

If you don't log attacks, you have no idea who is targeting you, what vulnerabilities they are looking for, or if they eventually succeeded.

## What is an Audit Log?
An audit log is an immutable record of security-relevant events. For a reverse proxy, this means logging:
- **Timestamp:** When did it happen?
- **Source IP:** Who did it?
- **Target Path:** What were they trying to access?
- **Violation Type:** Was it XSS? SQLi? A blocked IP?
- **Payload Snippet:** What exactly did they send? (Truncated to avoid massive logs).

## Structured Logging (JSON)
We will use JSON for our logs instead of plain text. 

**Why JSON?**
Because enterprise tools like **SIEMs** (Security Information and Event Management systems - e.g., Splunk, Datadog, ELK) are designed to ingest JSON. They can instantly parse fields like `{"violation": "XSS"}` and trigger alerts or draw graphs. Plain text logs require writing complex, brittle regex to extract data.

Next, we will build a dedicated Audit Logger package!
