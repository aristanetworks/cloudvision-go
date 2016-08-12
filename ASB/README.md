TODO TODO TODO TODO

# ASB (Automatic Server Bringup) State Diagram

    ------------------- reboot ---- [ reimage ]  ---------------- cmd ----------------------       
    |                                                                                      |
    |                                                                                      | 
    |               --- reboot ---- [ reimage_preserve ] -------- cmd ------------------   |
    |               |                                                                  |   |
    v               v                                                                  |   |
[ fresh ] ---> [ disk_ready ] ---> [ os_ready ] -- reboot --> [ config_ready ] ---> [ installed ] 



#             | donothing | unknown   | disk_rdy   | os_rdy | cfg_rdy | installed
# -------------------------------------------------------------------------------
# setupDisk   |           | V         |            |        |         |
# setupSW     |           | V         | V          |        |         |
# setupConfig |           | V         | V          | V      |         |
# setupVerify |           | V         | V          | V      | V       |
