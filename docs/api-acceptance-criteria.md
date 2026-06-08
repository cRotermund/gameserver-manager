# Acceptance Criteria

As an API consumer I must be able to...

* See the game servers managed by this control plane with their name, along with their basic info:
   * Basic hardware specs (cores/mem/bandwidth) 
   * What game they've been configured and built to run 
   * Their current top-level status (stopped, stopping, starting, running)
* Start and stop a server, and understand when that operation has completed
* Reboot the server without instance recycle (quick reboot)
* Retrieve the detailed operating state of a particular server
   * IP address 
   * Connected clients
   * CPU
   * Memory
   * Processes running
* Retrieve streamed, real-time observability telemetry for dynamic display
   * Server logs
   * Process logs
   * CPU/Memory/Disk/Network
   * Connected clients

As the API owner/system administrator I want to:

* Know that my API is rate-limited to protect from abuse or bad programming practices in the client
* Know that all control actions are associated with a user, and are logged for auditability
* Know that my infrastructure secrets are protected
* Know that my spend is protected from run-away costs
