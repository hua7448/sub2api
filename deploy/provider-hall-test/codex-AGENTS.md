# Server operations context

This host is the operator's test deployment server in California, USA.
Read SERVER-INVENTORY.md and the relevant deployments/<service>/README.md before deployments.
Check host services and listeners outside the execution sandbox when required; an empty sandbox result does not mean a service is absent.
Preserve existing Nginx virtual hosts and Sub2API services unless explicitly changing them.
Record each deployment's version, ports, domains, service units, data paths, backups and credential references.
Never put secret values in these instructions, deployment documentation, Git or command arguments.
Dujiao-Next credentials are in restricted runtime config with age-encrypted recovery archives; see its deployment README.
For Dujiao-Next deployment history and credential field mappings, read deployments/dujiao-next/DEPLOYMENT-RECORD.md. Read actual secrets only when needed for an authorized operation; do not echo them into logs or ordinary documentation.
When changing a deployment, update its operational README and append a dated entry to its deployment record. After credential rotation, update the encrypted credential record; do not assume bootstrap passwords remain current.
